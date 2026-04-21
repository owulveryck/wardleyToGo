# Plan : Menu "Import OWM" dans le playground UI

## Context

Le playground (`exp/ui/`) permet d'editer des cartes Wardley en WTG2. On vient de creer un convertisseur `exp/owm2wtg2` (Parse + Emit) qui transforme le format OWM en WTG2. L'objectif est d'ajouter un bouton "Import OWM" au burger menu qui ouvre un modal avec deux options : upload de fichier ou coller du texte.

## Approche

- **WASM** : ajouter une fonction `convertOWMToWTG2(text) → wtg2Text` dans `wasm/main.go` qui appelle le package Go existant
- **HTML** : ajouter un bouton "Import OWM" dans le burger menu (a cote de "Import WTG2")
- **JS** : creer un modal suivant le pattern existant (collab-dialog), puis appeler `loadWTG2IntoEditor()` avec le resultat converti

Pas de preview intermediaire : le flux existant (loadWTG2IntoEditor → render) montre le SVG immediatement. Layout empile (file upload en haut, textarea en bas, separateur "- ou -") — pas de tabs.

## 3 fichiers a modifier

### 1. `exp/ui/wasm/main.go`

**Ajout d'import** (ligne 13) :
```go
"github.com/owulveryck/wardleyToGo/exp/owm2wtg2"
```

**Nouvelle fonction** (apres `parseToState`, ~ligne 349) :
```go
func convertOWM(_ js.Value, args []js.Value) any {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()
    if len(args) < 1 {
        return "error: no input provided"
    }
    input := args[0].String()
    doc, err := owm2wtg2.Parse(bytes.NewBufferString(input))
    if err != nil {
        return fmt.Sprintf("error: %v", err)
    }
    var buf bytes.Buffer
    if err := owm2wtg2.Emit(doc, &buf); err != nil {
        return fmt.Sprintf("error: %v", err)
    }
    return buf.String()
}
```

**Registration** dans `main()` (ligne 18, apres `parseWTG2ToState`) :
```go
js.Global().Set("convertOWMToWTG2", js.FuncOf(convertOWM))
```

### 2. `exp/ui/index.html`

**Ajouter un bouton** apres la ligne 536 (`<input type="file" id="import-file-input" ...>`) :
```html
<button class="burger-item" id="import-owm"><span class="burger-item-icon"><svg ...upload icon...></svg></span> Import OWM</button>
```

Meme icone upload SVG que "Import WTG2" (ligne 535). Pas besoin de `<input type="file">` ici car le modal en creera un dynamiquement.

### 3. `exp/ui/app.js`

#### 3a. i18n — traductions FR (apres ligne 134 `toolbar.importWtg2`)

```javascript
'toolbar.importOwm': 'Importer un fichier OWM (Online Wardley Maps)',
'owm.title': 'Importer une carte OWM',
'owm.uploadLabel': 'Charger un fichier',
'owm.uploadHint': 'Fichier .owm ou .txt au format Online Wardley Maps',
'owm.orSeparator': '— ou —',
'owm.pasteLabel': 'Coller du texte OWM',
'owm.pastePlaceholder': 'title Ma Carte\ncomponent Client [0.95, 0.5]\n...',
'owm.import': 'Importer',
'owm.cancel': 'Annuler',
'owm.emptyError': 'Veuillez selectionner un fichier ou coller du texte OWM.',
'owm.convertError': 'Erreur de conversion : {message}',
```

#### 3b. i18n — traductions EN (apres ligne 269 `toolbar.importWtg2`)

```javascript
'toolbar.importOwm': 'Import OWM file (Online Wardley Maps)',
'owm.title': 'Import OWM Map',
'owm.uploadLabel': 'Upload a file',
'owm.uploadHint': '.owm or .txt file in Online Wardley Maps format',
'owm.orSeparator': '— or —',
'owm.pasteLabel': 'Paste OWM text',
'owm.pastePlaceholder': 'title My Map\ncomponent Customer [0.95, 0.5]\n...',
'owm.import': 'Import',
'owm.cancel': 'Cancel',
'owm.emptyError': 'Please select a file or paste OWM text.',
'owm.convertError': 'Conversion error: {message}',
```

#### 3c. Modal — fonctions (apres `handleImportFile`, ~ligne 1864)

3 fonctions :
- **`openOWMImportModal()`** : cree un overlay `collab-overlay` + dialog `collab-dialog` (CSS existant, lignes 261-270 de index.html). Contient :
  - `<input type="file" accept=".owm,.txt,.json">` pour upload
  - separateur texte "— ou —"
  - `<textarea>` pour coller du texte OWM
  - `<div id="owm-import-error">` pour erreurs inline
  - Boutons Cancel + Import
  - Events : Escape ferme, backdrop click ferme, file selection efface textarea et vice versa
  - Focus automatique sur le textarea apres ouverture

- **`closeOWMImportModal()`** : retire le DOM overlay

- **`doOWMImport()`** :
  1. Priorise le fichier si selectionne, sinon le texte colle
  2. Verifie que `convertOWMToWTG2` est disponible (WASM charge)
  3. Appelle `convertOWMToWTG2(text)` — retourne string WTG2 ou "error:..."
  4. Si erreur, affiche dans `#owm-import-error`
  5. Si succes, ferme le modal et appelle `loadWTG2IntoEditor(result)` (existant)

#### 3d. Event wiring (apres ligne 2321 `import-file-input.addEventListener`)

```javascript
document.getElementById('import-owm').addEventListener('click', function() {
    closeBurgerMenu();
    openOWMImportModal();
});
```

## Build

Apres les modifications :
```bash
cd exp/ui && make main.wasm
```

## Verification

1. `make run` puis ouvrir `http://localhost:8080`
2. Burger menu → "Import OWM" → modal s'ouvre
3. **Test coller** : coller le contenu de `exp/owm2wtg2/testdata/teashop.owm` dans le textarea → Import → le WTG2 apparait dans l'editeur, le SVG se rend
4. **Test fichier** : upload `teashop.owm` → meme resultat
5. **Test erreur** : coller du texte invalide → message d'erreur inline
6. **Test vide** : cliquer Import sans rien → message "Veuillez selectionner..."
7. **Test dismiss** : Escape, backdrop click, Cancel → modal se ferme
8. **Test i18n** : changer FR/EN → textes du modal traduits
9. **Test WASM non charge** : ouvrir le modal avant chargement WASM → message d'erreur
