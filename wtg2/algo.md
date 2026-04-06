
## 🗺️ Vue d'ensemble du pipeline de rendu

```
Document source
     │
     ▼
┌─────────────┐
│   Parsing   │ ← BNF → AST
└──────┬──────┘
       ▼
┌─────────────┐
│ Validation  │ ← Erreurs sémantiques
│ sémantique  │
└──────┬──────┘
       ▼
┌─────────────┐
│  Layout Y   │ ← Algo 1 : Visibilité verticale
└──────┬──────┘
       ▼
┌─────────────┐
│  Layout X   │ ← Algo 2 : Position horizontale (trivial, depuis l'évolution)
└──────┬──────┘
       ▼
┌─────────────┐
│ Résolution  │ ← Algo 3 : Anti-chevauchement des nœuds
│ collisions  │
└──────┬──────┘
       ▼
┌─────────────┐
│  Routage    │ ← Algo 4 : Tracé des liens
│  des liens  │
└──────┬──────┘
       ▼
┌─────────────┐
│ Placement   │ ← Algo 5 : Labels des nœuds
│  labels     │
└──────┬──────┘
       ▼
┌─────────────┐
│ Placement   │ ← Algo 6 : Annotations des liens
│ link labels │
└──────┬──────┘
       ▼
┌─────────────┐
│ Placement   │ ← Algo 7 : Notes, warnings, signaux
│ annotations │
└──────┬──────┘
       ▼
┌─────────────┐
│  Pipelines  │ ← Algo 8 : Rendu des pipelines
└──────┬──────┘
       ▼
┌─────────────┐
│  Groupes    │ ← Algo 9 : Enveloppes visuelles
└──────┬──────┘
       ▼
┌─────────────┐
│  Styling    │ ← Algo 10 : Couleurs, épaisseurs, inertie, mouvement
└──────┬──────┘
       ▼
   Rendu final (SVG / Canvas)
```

---

## Algo 1 — Placement vertical (visibilité)

### Problème
Assigner une coordonnée Y à chaque nœud en respectant la sémantique Wardley : les ancres en haut, les dépendances profondes en bas.

### Entrée
- Graphe orienté des dépendances (edges `->`)
- Overrides explicites (`@0.7`)
- Nœuds `anchor` (forcés en haut)

### Algorithme

```
PLACEMENT_VERTICAL(graphe, overrides):

  1. Tri topologique
     ─────────────
     topo = tri_topologique(graphe)
     // Gère les cycles : détection et erreur sémantique

  2. Calcul du rang (profondeur depuis les ancres)
     ─────────────
     pour chaque noeud n dans topo:
       si n est anchor:
         rang[n] = 0
       sinon:
         rang[n] = max(rang[parent] + 1 pour parent dans parents(n))

  3. Normalisation en [0, 1]
     ─────────────
     rang_max = max(rang.values())
     pour chaque noeud n:
       si n dans overrides:
         Y[n] = overrides[n]              // override explicite
       sinon:
         Y[n] = 1.0 - (rang[n] / rang_max)  // 1.0 = haut

  4. Marge minimale entre rangs
     ─────────────
     // Garantir un espacement minimum entre niveaux
     espacement_min = 1.0 / (rang_max + 2)
     pour chaque rang r:
       Y_rang[r] = 1.0 - (r + 1) * espacement_min

  5. Distribution intra-rang (nœuds au même rang)
     ─────────────
     // Voir Algo 3 pour l'anti-chevauchement
```

### Complexité
`O(V + E)` — tri topologique + une passe linéaire.

### Cas limites
| Cas | Traitement |
|-----|-----------|
| Cycle dans le graphe | Erreur sémantique, refuser le rendu |
| Nœud orphelin (pas dans un edge) | Placer en bas par défaut (`Y = 0.2`) |
| Override `@` qui viole l'ordre topo | Accepter (l'utilisateur sait ce qu'il fait) |
| Ancres multiples | Toutes au rang 0, distribuées horizontalement |

---

## Algo 2 — Position horizontale (évolution)

### Problème
Trivial : convertir la notation romaine en coordonnée X.

### Algorithme

```
POSITION_X(position):

  // Chaque phase occupe 25% de l'axe
  base = { I: 0.0, II: 0.25, III: 0.50, IV: 0.75 }

  si position = roman:
    retourner base[roman] + 0.125    // centre de la phase

  si position = roman.decimal:
    retourner base[roman] + (decimal / 10.0) * 0.25
    // Ex: II.7 = 0.25 + 0.7 * 0.25 = 0.425
```

### Complexité
`O(1)` par nœud.

---

## Algo 3 — Anti-chevauchement des nœuds

### Problème
Deux nœuds au même rang avec des positions X proches se chevauchent visuellement.

### Entrée
- Positions (X, Y) calculées par Algo 1 et 2
- Dimensions de chaque nœud (largeur du texte + padding)

### Algorithme

```
ANTI_CHEVAUCHEMENT(noeuds):

  1. Grouper par bande Y
     ─────────────
     // Deux nœuds "se chevauchent en Y" si |Y1 - Y2| < hauteur_noeud
     bandes = regrouper_par_proximite_Y(noeuds, seuil=hauteur_noeud)

  2. Pour chaque bande, résoudre en X
     ─────────────
     pour chaque bande:
       trier les noeuds par X
       pour i = 1 à len(bande) - 1:
         chevauchement = (bande[i-1].X + bande[i-1].largeur/2 + marge)
                       - (bande[i].X - bande[i].largeur/2)
         si chevauchement > 0:
           // Écarter symétriquement
           bande[i-1].X -= chevauchement / 2
           bande[i].X   += chevauchement / 2

  3. Si écartement en X impossible (bord de carte)
     ─────────────
     // Décaler en Y légèrement (micro-ajustement)
     bande[i].Y += micro_offset
     // Le micro_offset ne doit pas changer le rang visuel
```

### Complexité
`O(n log n)` — tri par X dans chaque bande.

---

## Algo 4 — Routage des liens

### Problème
Tracer les liens entre nœuds sans traverser d'autres nœuds, avec un rendu lisible.

### Algorithme

```
ROUTAGE_LIENS(liens, noeuds):

  pour chaque lien (source, cible, type, annotation):

    1. Points d'ancrage
       ─────────────
       // Choisir le meilleur point de connexion sur chaque nœud
       // Préférer : bas du source → haut de la cible (flux descendant)
       p_source = point_ancrage(source, direction_vers=cible)
       p_cible  = point_ancrage(cible, direction_vers=source)

    2. Détection d'obstacles
       ─────────────
       obstacles = noeuds_entre(p_source, p_cible, tous_les_noeuds)

    3. Tracé
       ─────────────
       si pas d'obstacles:
         // Ligne directe avec légère courbe de Bézier
         chemin = bezier_quadratique(p_source, point_controle, p_cible)
       sinon:
         // Contournement : ajouter des waypoints
         waypoints = contourner_obstacles(p_source, p_cible, obstacles)
         chemin = spline_catmull_rom(p_source, waypoints, p_cible)

    4. Flèche directionnelle
       ─────────────
       si type == "->":
         ajouter_fleche(chemin.fin, angle=tangente_fin(chemin))
       si type == "<->":
         ajouter_fleche(chemin.debut, ...)
         ajouter_fleche(chemin.fin, ...)
```

### Point d'ancrage intelligent

```
POINT_ANCRAGE(noeud, direction_vers):
  
  angle = atan2(cible.Y - noeud.Y, cible.X - noeud.X)
  
  // 8 points d'ancrage possibles (N, NE, E, SE, S, SW, W, NW)
  // Choisir celui dont l'angle est le plus proche de `angle`
  // En privilégiant S (bas) pour les sources et N (haut) pour les cibles
  
  si noeud est source:
    preference = [S, SE, SW, E, W, NE, NW, N]  // favoriser le bas
  sinon:
    preference = [N, NE, NW, E, W, SE, SW, S]  // favoriser le haut
  
  retourner meilleur_ancrage(noeud, angle, preference)
```

### Complexité
`O(E × V)` naïf. Optimisable avec un index spatial (quadtree).

---

## Algo 5 — Placement des labels de nœuds

### Problème
Placer le texte de chaque nœud sans chevaucher d'autres nœuds ni d'autres labels.

### Algorithme : Eight-Position Model

```
PLACEMENT_LABELS(noeuds, liens):

  // 8 positions candidates par nœud
  positions = [N, NE, E, SE, S, SW, W, NW]

  1. Score de chaque position candidate
     ─────────────
     pour chaque noeud n:
       pour chaque pos dans positions:
         rect = calculer_rect_label(n, pos)
         score[n][pos] = 0

         // Pénalité : chevauchement avec d'autres nœuds
         pour chaque autre noeud m:
           si rect chevauche m.rect:
             score[n][pos] -= 100

         // Pénalité : chevauchement avec labels déjà placés
         pour chaque label déjà placé:
           si rect chevauche label.rect:
             score[n][pos] -= 80

         // Pénalité : chevauchement avec des liens
         pour chaque lien:
           si rect intersecte lien.chemin:
             score[n][pos] -= 50

         // Bonus : position préférée (NE par défaut Wardley)
         si pos == NE:
           score[n][pos] += 10

         // Pénalité : label hors carte
         si rect sort des limites:
           score[n][pos] -= 200

  2. Assignation gloutonne (greedy)
     ─────────────
     trier les noeuds par nombre de voisins décroissant
     // Les nœuds les plus contraints en premier
     pour chaque noeud n:
       label_pos[n] = pos avec max(score[n][pos])
       marquer rect comme occupé

  3. Fallback
     ─────────────
     si aucune position sans chevauchement:
       utiliser la position avec le score le moins négatif
       réduire la taille de police si nécessaire
```

### Complexité
`O(n² × 8)` — acceptable pour des cartes < 200 nœuds.

### Alternative avancée
Pour des cartes très denses : reformuler comme un problème de **satisfaction de contraintes** (CSP) ou utiliser l'algorithme de **Simulated Annealing** pour le placement de labels (algorithme de Christensen/Marks/Shieber).

---

## Algo 6 — Labels sur les liens (annotations)

### Problème
Placer le texte `[flux critique]` sur le lien correspondant, lisiblement.

### Algorithme

```
PLACEMENT_LABEL_LIEN(lien, texte):

  1. Point médian du chemin
     ─────────────
     point = milieu(lien.chemin)
     angle = tangente(lien.chemin, au_milieu)

  2. Orientation du texte
     ─────────────
     // Toujours lisible gauche→droite
     si angle > 90° ou angle < -90°:
       angle += 180°

  3. Décalage perpendiculaire
     ─────────────
     // Décaler le texte au-dessus du lien
     offset = normale(angle) * 8px
     point_label = point + offset

  4. Anti-chevauchement
     ─────────────
     rect_label = calculer_rect(texte, point_label, angle)
     si rect_label chevauche un noeud ou un autre label:
       // Déplacer le label au 1/3 ou 2/3 du chemin
       point_label = position_alternative(lien.chemin, [0.33, 0.66])
```

---

## Algo 7 — Placement des annotations (notes, warnings, signaux)

### Problème
Afficher les annotations liées à un nœud sans encombrer la carte.

### Algorithme

```
PLACEMENT_ANNOTATIONS(annotations, noeuds, zone_carte):

  1. Stratégie : panneau latéral vs inline
     ─────────────
     si nombre(annotations) > 5:
       mode = PANNEAU_LATERAL     // Numérotées avec lignes de rappel
     sinon:
       mode = INLINE              // Directement sur la carte

  2. Mode INLINE
     ─────────────
     pour chaque annotation a:
       noeud_cible = a.noeud_ref

       // Bulle d'annotation
       rect = calculer_bulle(a.texte)

       // Placement : chercher un espace libre autour du nœud
       position = chercher_espace_libre(
         centre = noeud_cible.position,
         rayon_initial = 60px,
         rayon_max = 200px,
         obstacles = tous_les_rects_occupés
       )

       // Ligne de rappel (connecteur)
       tracer_ligne_pointillee(position, noeud_cible.position)

  3. Mode PANNEAU LATÉRAL
     ─────────────
     // Colonne à droite de la carte
     panneau_x = carte.largeur + marge
     pour chaque annotation a (triées par Y du nœud cible):
       afficher numéro circled sur le noeud
       afficher dans le panneau : "① texte de l'annotation"
       tracer_ligne_pointillee(noeud, panneau)

  4. Styling par type
     ─────────────
     note    → fond bleu pâle, bordure bleue
     warning → fond jaune pâle, bordure orange, icône ⚠
     signal  → icône directionnelle (↗ ↘ →) sur le nœud
```

---

## Algo 8 — Rendu des pipelines

### Problème
Afficher un pipeline comme une barre horizontale contenant ses membres positionnés.

### Algorithme

```
RENDU_PIPELINE(pipeline, Y_pipeline):

  membres = pipeline.membres
  x_min = min(POSITION_X(m.position) pour m dans membres) - padding
  x_max = max(POSITION_X(m.position) pour m dans membres) + padding

  1. Barre de pipeline
     ─────────────
     tracer_rectangle_arrondi(
       x = x_min, y = Y_pipeline - hauteur/2,
       largeur = x_max - x_min,
       hauteur = hauteur_pipeline,
       style = fond_transparent, bordure_pointillee
     )

  2. Membres positionnés dans la barre
     ─────────────
     pour chaque membre m:
       m.X = POSITION_X(m.position)
       m.Y = Y_pipeline
       tracer_noeud(m, style=petit_cercle)
       tracer_label(m.nom, position=au_dessus)

  3. Liens entrants/sortants
     ─────────────
     // Les liens vers le pipeline sans ancre précise
     // arrivent au centre de la barre
     // Les liens vers pipeline:membre arrivent au membre
```

---

## Algo 9 — Enveloppes de groupes

### Problème
Dessiner un contour englobant autour des composants d'un groupe.

### Algorithme : Convex Hull + Padding

```
ENVELOPPE_GROUPE(groupe, noeuds):

  1. Collecter les positions
     ─────────────
     points = []
     pour chaque noeud dans groupe.membres:
       points.ajouter(noeud.rect.coins())  // 4 coins du rectangle

  2. Enveloppe convexe
     ─────────────
     hull = convex_hull(points)           // Graham scan, O(n log n)

  3. Padding et lissage
     ─────────────
     hull_padded = offset_polygone(hull, padding=20px)
     hull_lisse = arrondir_coins(hull_padded, rayon=10px)

  4. Rendu
     ─────────────
     tracer_polygone(hull_lisse,
       fond = couleur_groupe avec opacité 0.08,
       bordure = couleur_groupe avec opacité 0.3,
       style = pointillé
     )
     placer_label(groupe.nom,
       position = coin_superieur_gauche(hull_lisse),
       style = italique, couleur_groupe
     )
```

### Alternative
Pour des groupes non convexes (nœuds dispersés) : utiliser un **alpha shape** ou **bubble set** (Algorithme de Collins et al., 2009) qui génère des formes organiques englobantes.

---

## Algo 10 — Styling (couleurs, épaisseurs, marqueurs visuels)

### 10a. Couleurs des liens

```
STYLE_LIEN(lien):

  // Couleur de base
  si lien.annotation existe:
    couleur = #2980B9          // bleu — lien informatif
  sinon si lien.type == "<->":
    couleur = #8E44AD          // violet — bidirectionnel
  sinon:
    couleur = #7F8C8D          // gris — dépendance standard

  // Épaisseur selon le nombre de dépendants
  fan_out = nombre_de_liens_sortants(lien.source)
  epaisseur = 1.5 + min(fan_out * 0.3, 2.0)  // 1.5px à 3.5px

  // Opacité selon la profondeur
  profondeur = max(rang[source], rang[cible])
  opacite = 1.0 - profondeur * 0.1            // les liens profonds s'estompent
```

### 10b. Marqueurs d'inertie

```
RENDU_INERTIE(noeud):

  si noeud.inertie == "!":
    // Petite barre verticale sur le nœud
    tracer_barre(noeud.X, noeud.Y, hauteur=15px, couleur=#E74C3C)

  si noeud.inertie == "!!":
    // Double barre
    tracer_barre(noeud.X - 3, ...) 
    tracer_barre(noeud.X + 3, ...)

  si noeud.inertie == "!!!":
    // Triple barre + halo rouge
    tracer_barres(3, ...)
    tracer_halo(noeud, couleur=#E74C3C, opacite=0.15)
```

### 10c. Marqueurs de mouvement (évolution projetée)

```
RENDU_MOUVEMENT(noeud):

  si noeud.evolution a un mouvement (>> position_cible):
    x_actuel = POSITION_X(noeud.position_actuelle)
    x_cible  = POSITION_X(noeud.position_cible)

    // Flèche pointillée horizontale
    tracer_fleche_pointillee(
      de = (x_actuel + rayon_noeud, noeud.Y),
      vers = (x_cible, noeud.Y),
      couleur = #E74C3C,
      style = pointillé
    )

    // Nœud fantôme à la position cible
    tracer_cercle(x_cible, noeud.Y,
      style = pointillé, opacite = 0.3
    )
```

### 10d. Marqueurs de signaux

```
RENDU_SIGNAL(signal):

  noeud = signal.noeud_ref

  icone = {
    accelerating: "↗",    // flèche montante verte
    stagnating:   "→",    // flèche horizontale orange
    declining:    "↘"     // flèche descendante rouge
  }

  couleur = {
    accelerating: #27AE60,
    stagnating:   #F39C12,
    declining:    #E74C3C
  }

  placer_icone(
    icone[signal.type],
    position = noeud.coin_superieur_droit,
    couleur = couleur[signal.type],
    taille = 14px
  )
```

### 10e. Styling des types (build / buy / outsource)

```
STYLE_NOEUD(noeud):

  // Forme de base
  forme = cercle_plein      // défaut Wardley

  // Couleur selon le type
  si noeud.type == "build":
    contour = #2C3E50        // gris foncé
    remplissage = #2C3E50
  si noeud.type == "buy":
    contour = #2C3E50
    remplissage = #FFFFFF    // cercle vide (convention Wardley)
  si noeud.type == "outsource":
    contour = #2C3E50
    remplissage = #FFFFFF
    forme = carre             // carré pour outsource

  // Override couleur explicite
  si noeud.color existe:
    remplissage = noeud.color
    contour = assombrir(noeud.color, 20%)
```

---

## 📋 Tableau récapitulatif

| # | Algorithme | Complexité | Difficulté | Dépendances |
|---|-----------|------------|------------|-------------|
| 1 | Placement vertical | O(V+E) | ⭐⭐ | Graphe de dépendances |
| 2 | Position horizontale | O(V) | ⭐ | Notation romaine |
| 3 | Anti-chevauchement nœuds | O(n log n) | ⭐⭐ | Algo 1, 2 |
| 4 | Routage des liens | O(E×V) | ⭐⭐⭐ | Algo 1, 2, 3 |
| 5 | Labels des nœuds | O(n²) | ⭐⭐⭐ | Algo 1, 2, 3, 4 |
| 6 | Labels des liens | O(E) | ⭐⭐ | Algo 4 |
| 7 | Annotations | O(A×n) | ⭐⭐ | Algo 1-5 |
| 8 | Pipelines | O(P×M) | ⭐⭐ | Algo 1, 2 |
| 9 | Enveloppes groupes | O(n log n) | ⭐⭐ | Algo 1, 2, 3 |
| 10 | Styling | O(V+E) | ⭐ | Tous |

### Ordre d'implémentation recommandé

```
Phase 1 — MVP fonctionnel
  Algo 2 (X trivial) → Algo 1 (Y) → Algo 3 (collisions)
  → Algo 10e (styling basique)

Phase 2 — Lisibilité
  Algo 4 (liens) → Algo 5 (labels nœuds) → Algo 6 (labels liens)

Phase 3 — Richesse
  Algo 8 (pipelines) → Algo 10a-d (styling avancé)

Phase 4 — Polish
  Algo 7 (annotations) → Algo 9 (groupes)
```

