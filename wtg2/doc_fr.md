# WTG2 -- Langage de description de cartes Wardley

## Table des matieres

1. [Introduction](#introduction)
2. [Structure d'un document](#structure-dun-document)
3. [Commentaires](#commentaires)
4. [Metadonnees](#metadonnees)
5. [Labels des phases d'evolution](#labels-des-phases-devolution)
6. [Legende](#legende)
7. [Noeuds](#noeuds)
   - [Ancres](#ancres)
   - [Composants](#composants)
   - [Sous-cartes](#sous-cartes)
   - [Declaration en bloc](#declaration-en-bloc)
8. [Positionnement sur l'axe d'evolution](#positionnement-sur-laxe-devolution)
9. [Mouvement d'evolution](#mouvement-devolution)
10. [Inertie](#inertie)
11. [Visibilite](#visibilite)
12. [Pipelines](#pipelines)
13. [Chaine de valeur (edges)](#chaine-de-valeur-edges)
14. [Groupes](#groupes)
15. [Annotations](#annotations)
16. [Signaux](#signaux)
17. [Manouvres strategiques (gameplays)](#manoeuvres-strategiques-gameplays)
18. [Focus](#focus)
19. [Regles sur les identifiants](#regles-sur-les-identifiants)
20. [Regles semantiques](#regles-semantiques)
21. [Exemple complet](#exemple-complet)

---

## Introduction

**WTG2** est un langage declaratif concu pour decrire des [cartes Wardley](https://learnwardleymapping.com/). Une carte Wardley est un outil strategique qui represente une **chaine de valeur** (axe vertical) en fonction du **stade d'evolution** de chaque composant (axe horizontal).

### Les deux axes d'une carte Wardley

**Axe vertical -- Visibilite (chaine de valeur) :**
Les composants visibles par l'utilisateur sont en haut de la carte. Plus on descend, plus on s'approche de l'infrastructure invisible.

**Axe horizontal -- Evolution :**
Les composants se deplacent de gauche a droite au fur et a mesure qu'ils murissent, a travers quatre phases :

| Phase | Chiffre romain | Description |
|-------|---------------|-------------|
| Genese | I | Nouveau, mal compris, forte incertitude |
| Sur-mesure | II | Compris mais specifique, demande de l'expertise |
| Produit | III | De plus en plus standardise, disponible comme produit |
| Commodite | IV | Hautement standardise, paiement a l'usage, invisible |

Un fichier `.wtg2` est un fichier texte brut qui, une fois parse, genere une structure de carte rendable en SVG.

---

## Structure d'un document

Un document WTG2 est compose de **declarations** (statements), dans un ordre libre. Toutefois, l'ordre canonique recommande est :

```
1. Metadonnees      (title, date, author, scope, question, doctrine)
2. Configuration    (stages, legend)
3. Noeuds           (anchor, component, submap, pipeline)
4. Chaine de valeur (edges / dependances)
5. Groupes          (group)
6. Annotations      (note, warning, signal, gameplay)
```

Toutes les sections sont optionnelles. Les commentaires et lignes vides peuvent apparaitre n'importe ou.

---

## Commentaires

WTG2 supporte deux styles de commentaires :

```
// Commentaire sur une ligne

/* Commentaire
   sur plusieurs lignes */
```

Les commentaires sont ignores par le parser et servent a documenter la carte.

---

## Metadonnees

Les metadonnees decrivent le contexte de la carte. Chaque champ est optionnel.

```
title: Plateforme de Navigation -- Strategie 2026
date: 2026-01-15
author: Cellule Strategie Produit
scope: Application mobile B2C, marche europeen
question: "Ou investir pour se differencier face a Google Maps ?"
doctrine: context
```

| Champ | Description |
|-------|-------------|
| `title` | Titre de la carte |
| `date` | Date au format ISO `AAAA-MM-JJ` |
| `author` | Auteur ou equipe |
| `scope` | Perimetre de la carte |
| `question` | Question strategique a laquelle la carte repond (entre guillemets) |
| `doctrine` | Phase de maturite organisationnelle : `hygiene`, `context`, `excellence` ou `evolution` |

---

## Labels des phases d'evolution

Par defaut, les phases sont marquees par les chiffres romains I, II, III, IV. Vous pouvez les personnaliser :

```
stages: Genese, Sur-mesure, Produit, Commodite
```

Exactement quatre labels, separes par des virgules. Ils apparaitront au bas de la carte comme intitules des zones.

---

## Legende

Le mot-cle `legend` active l'affichage d'un panneau de legende a droite de la carte. Le contenu est auto-detecte a partir des elements presents.

```
legend
```

Pas de configuration necessaire. Le viewBox SVG est automatiquement elargi pour accueillir la legende sans deformer la carte.

---

## Noeuds

Un noeud represente un element sur la carte. Il existe trois types de noeuds.

### Ancres

Une **ancre** represente un utilisateur ou un acteur. Elle est toujours placee en haut de la carte et ne necessite pas de position d'evolution.

```
anchor Automobiliste
anchor Collectivite Locale
```

Une ancre peut optionnellement avoir une position (pour un placement horizontal precis) :

```
anchor Automobiliste : III.5
```

### Composants

Un **composant** est l'element fondamental d'une carte Wardley. Le mot-cle `component` est optionnel -- un nom suivi d'une position est automatiquement traite comme un composant.

```
// Ces deux ecritures sont equivalentes :
component Application Mobile : III.5
Application Mobile : III.5
```

Un composant peut etre qualifie par un **type de sourcing** :

| Type | Signification |
|------|--------------|
| `build` | Construit en interne |
| `buy` | Achete comme produit |
| `outsource` | Externalise a un tiers |

```
Infrastructure Cloud : IV.3 (buy)
Reseau Mobile : IV.7 (outsource)
Moteur de Calcul : II.7 (build)
```

### Sous-cartes

Une **sous-carte** (submap) represente un composant encapsulant une carte Wardley distincte. Elle est rendue comme un composant avec un marqueur visuel specifique.

```
submap Systeme de Paiement : III.6
```

### Declaration en bloc

Pour ajouter des proprietes supplementaires, utilisez la declaration en bloc avec accolades :

```
Moteur de Calcul : II.7 {
  type: build
  color: #3498DB
  asset: tech
  cost: "1.2M EUR/an, 12 ETP"
  note: "Notre differenciateur cle"
}
```

| Champ | Valeurs | Description |
|-------|---------|-------------|
| `type` | `build`, `buy`, `outsource` | Type de sourcing |
| `asset` | `tech`, `financial`, `human`, `relational`, `social` | Nature du capital (voir ci-dessous) |
| `color` | `#RRGGBB`, `#RGB`, nom CSS | Couleur de rendu |
| `visibility` | `0.0` a `1.0` | Position verticale forcee |
| `cost` | Texte libre | Annotation de cout |
| `note` | Texte entre guillemets | Description du composant |

#### Classification des actifs (asset)

Chaque composant represente un type de capital organisationnel :

| Type | Description | Exemple |
|------|------------|---------|
| `tech` | Capital technologique : code, infra, brevets | Un moteur de routage, un pipeline de donnees |
| `financial` | Capital financier : modeles de revenus, pouvoir tarifaire | Un systeme de facturation |
| `human` | Capital humain : expertise, competences, savoir tacite | Une equipe ML, des experts metier |
| `relational` | Capital relationnel : partenariats, marque, contrats | Une API partenaire, un portefeuille de brevets |
| `social` | Capital social/environnemental : communaute, reglementation | Communaute open-source, conformite reglementaire |

---

## Positionnement sur l'axe d'evolution

La position horizontale utilise la **notation en chiffres romains** avec une subdivision decimale optionnelle :

```
<romain>.<chiffre>
```

Ou `<romain>` est `I`, `II`, `III` ou `IV`, et `<chiffre>` est de `0` a `9`.

Chaque phase occupe 25% de l'axe. Le chiffre decimal subdivise la phase (0 = debut, 9 = fin). Sans decimal, le centre de la phase est utilise.

### Table de correspondance

| Position | Coordonnee (0-100) | Signification |
|----------|-------------------|---------------|
| `I.0` | 0 | Debut de Genese |
| `I.5` | 12 | Milieu de Genese |
| `I.9` | 22 | Fin de Genese |
| `II.0` | 25 | Debut de Sur-mesure |
| `II.5` | 37 | Milieu de Sur-mesure |
| `II.9` | 47 | Fin de Sur-mesure |
| `III.0` | 50 | Debut de Produit |
| `III.5` | 62 | Milieu de Produit |
| `III.9` | 72 | Fin de Produit |
| `IV.0` | 75 | Debut de Commodite |
| `IV.5` | 87 | Milieu de Commodite |
| `IV.9` | 97 | Fin de Commodite |

**Formule :** `floor((base + chiffre/10 * 0.25) * 100)` ou `base` vaut `I=0.00, II=0.25, III=0.50, IV=0.75`.

Sans decimal (ex: `III` seul), le composant est place au centre de la phase (12, 37, 62 ou 87).

---

## Mouvement d'evolution

L'operateur `>>` indique qu'un composant est en train d'evoluer d'une position vers une autre :

```
Composant : II.7 >> III.5
```

Cela genere une fleche pointillee de la position actuelle (II.7) vers la position cible (III.5), avec un composant "fantome" a la position d'arrivee.

---

## Inertie

L'inertie represente la resistance au changement. Elle se note avec des points d'exclamation (`!`) entre la position actuelle et l'operateur de mouvement `>>`.

### Niveaux d'inertie

```
Composant : II.7 ! >> III.5      // inertie moderee
Composant : II.7 !! >> III.5     // inertie forte
Composant : II.7 !!! >> III.5    // inertie bloquante
```

### Inertie qualifiee

L'inertie peut etre qualifiee par sa nature grace a une liste de types entre parentheses :

```
Composant : II.7 !!(tech,human) >> III.5     // inertie technique + humaine
Composant : II.7 !(financial) >> III.5        // inertie financiere seule
Composant : II.7 !!!(tech,financial,human) >> III.5  // trois types de resistance
```

| Type | Signification | Symptome typique |
|------|--------------|-----------------|
| `tech` | Dette technique, verrouillage technologique | "On a toujours utilise Java" |
| `financial` | Couts irrecuperables, modeles de revenus etablis | "On a investi 5M la-dedans" |
| `human` | Manque de competences, obsolescence de l'expertise | "L'equipe ne connait pas le cloud-native" |
| `relational` | Obligations contractuelles, dependances partenaires | "On a un contrat fournisseur de 3 ans" |
| `social` | Resistance culturelle, inertie reglementaire | "Ce n'est pas comme ca qu'on fonctionne ici" |

L'inertie qualifiee est rendue visuellement avec des barres colorees cote a cote sur la fleche d'evolution.

---

## Visibilite

Par defaut, la position verticale d'un composant est **calculee automatiquement** a partir du graphe de dependances : les ancres sont en haut, les composants d'infrastructure en bas.

Vous pouvez forcer la position verticale avec l'operateur `@` :

```
Composant : III.5 @0.9    // pres du haut de la carte
Composant : III.5 @0.1    // pres du bas de la carte
```

Les valeurs vont de `0.0` (bas) a `1.0` (haut). Cette override est rarement necessaire -- le layout automatique produit generalement un bon resultat.

La visibilite peut aussi etre definie dans un bloc :

```
Composant : III.5 {
  visibility: 0.8
}
```

---

## Pipelines

Un **pipeline** modelise un composant strategique existant sous **plusieurs formes simultanees** a differents stades d'evolution. C'est l'un des concepts les plus puissants des cartes Wardley.

### Concept

Un pipeline represente par exemple un moteur de calcul qui existe en trois implementations :
- Un algorithme classique (produit mature)
- Un algorithme IA predictif (sur-mesure, en developpement)
- Un algorithme quantique (recherche, genese)

Ces trois implementations remplissent la meme fonction mais se situent a des niveaux d'evolution differents.

### Syntaxe

```
// D'abord, declarer le composant parent
Moteur de Calcul : II.7

// Puis, declarer le pipeline
pipeline Moteur de Calcul {
  Algo Dijkstra Classique : III.5
  Algo Predictif IA : II.3
  Algo Quantique : I.2
}
```

### Regles

- Le nom du pipeline **doit correspondre** a un composant deja declare.
- Les membres sont positionnes uniquement sur l'axe d'evolution ; leur position verticale est celle du composant parent.
- La plage horizontale du pipeline couvre de son membre le plus a gauche a son membre le plus a droite.
- Un pipeline **ne peut pas contenir** un autre pipeline.
- Chaque membre est declare avec un nom et une position : `Nom : Position`.

### Rendu visuel

Le pipeline est affiche comme un rectangle horizontal contenant ses membres, positionnes sur l'axe d'evolution. Le composant parent est le point d'ancrage pour les liens (edges) de la chaine de valeur.

### Reference a un membre de pipeline

Dans les edges, vous pouvez cibler un membre specifique d'un pipeline avec la syntaxe `:` :

```
Alertes Trafic -> Moteur de Calcul:Algo Predictif IA
```

---

## Chaine de valeur (edges)

Les edges definissent les dependances entre composants et constituent la chaine de valeur.

### Syntaxe de base

```
A -> B                              // A depend de B
A <-> B                             // relation bidirectionnelle
A -[texte du label]-> B             // dependance annotee
A <-[texte du label]-> B            // bidirectionnelle annotee
```

### Chaines

Les edges peuvent etre enchaines sur une seule ligne :

```
Utilisateur -> Application -> API -> Base de Donnees -> Cloud
```

Cela cree quatre edges distincts : Utilisateur->Application, Application->API, API->Base de Donnees, Base de Donnees->Cloud.

### Labels

Les labels permettent de documenter la nature de la relation :

```
Collectivite -> API Partenaires -[Open Data, licence annuelle]-> Modele de Donnees
```

Le texte entre crochets est affiche le long de l'edge sur la carte.

### Reference a un membre de pipeline

```
Alertes Trafic -> Moteur de Calcul:Algo Predictif IA
```

La syntaxe `Pipeline:Membre` permet de creer un lien vers un membre specifique d'un pipeline plutot que vers le composant parent.

---

## Groupes

Les groupes permettent de regrouper visuellement des composants, typiquement par equipe ou par domaine. Ils sont **purement visuels** et ne creent pas de namespace.

### Syntaxe de base

```
group Equipe Core Navigation {
  Moteur de Calcul
  Algo Predictif IA
  Modele de Donnees Carto
}
```

### Directives de groupe

| Directive | Valeurs | Description |
|-----------|---------|-------------|
| `color:` | `#RRGGBB`, `#RGB`, nom CSS | Couleur de l'enveloppe du groupe |
| `team:` | `explorer`, `settler`, `town-planner`, `pioneer`, `villager` | Type d'equipe EVT/PST |

```
group Equipe R&D {
  team: explorer
  color: #E74C3C
  Algo Quantique
  Cache Experimental
}
```

### Alignement EVT/PST

Le modele Explorer-Villager-Town-Planner aligne les types d'equipe aux phases d'evolution :

| Type d'equipe | Phase d'evolution | Etat d'esprit |
|--------------|-------------------|---------------|
| `explorer` / `pioneer` | Genese (I) | Decouverte, intuition, tolerance a l'echec |
| `settler` / `villager` | Sur-mesure/Produit (II-III) | Productisation, standards, analyse |
| `town-planner` | Commodite (IV) | Industrialisation, optimisation des couts |

Un desalignement entre le type d'equipe et la phase d'evolution des composants est un signal strategique a examiner.

---

## Annotations

Les annotations attachent des informations textuelles a un composant.

### Notes

```
note "Candidat a l'externalisation Q4 2026" on Systeme de Paiement
note "Partenariat signe avec 12 metropoles" on API Partenaires B2G
```

### Avertissements

```
warning "SPOF -- aucun fallback si indisponible" on Moteur de Calcul
warning "Vendor lock-in AWS, cout en hausse de 30%/an" on Infrastructure Cloud
```

Les notes sont rendues avec un style neutre (bleu) tandis que les warnings utilisent un style d'alerte (orange/rouge).

---

## Signaux

Les signaux marquent les dynamiques de marche et les forces climatiques sur un composant.

```
signal accelerating on Algo Predictif IA
signal stagnating on Algo Dijkstra Classique
signal declining on Donnees OSM
```

### Types de signaux

| Signal | Signification |
|--------|--------------|
| `accelerating` | Evolution rapide vers la commodite |
| `stagnating` | L'evolution a atteint un plateau |
| `declining` | Regression en pertinence ou en usage |
| `co-evolution` | Technologie et pratique evoluant ensemble |
| `red-queen` | Doit evoluer constamment juste pour maintenir sa position |
| `commoditization` | Force gravitationnelle vers l'utilitaire |
| `network-effects` | La valeur croit avec le nombre d'utilisateurs |
| `economies-of-scale` | Avantage de cout par le volume |

```
signal co-evolution on Pratiques DevOps
signal commoditization on Infrastructure Cloud
signal network-effects on Plateforme Sociale
signal economies-of-scale on Cloud Compute
signal red-queen on Application
```

---

## Manoeuvres strategiques (gameplays)

Les gameplays annotent les manoeuvres strategiques deliberees appliquees a un composant.

### Syntaxe

```
gameplay <type> on <Composant>
gameplay <type> "description optionnelle" on <Composant>
```

### Types de gameplays

| Gameplay | Description | Contexte typique |
|----------|-------------|-----------------|
| `ILC` | Innovate-Leverage-Commoditize : fournir l'infra, observer, absorber | Plateforme avec ecosysteme |
| `open-source` | Commoditiser une couche pour capturer la valeur adjacente | Concurrent avec une rente proprietaire |
| `land-grab` | Sacrifier la rentabilite pour la part de marche | Nouveau marche avec effets de reseau |
| `embrace-extend` | Adopter un standard, ajouter des extensions proprietaires | Standard a controler |
| `tower-moat` | Eriger des barrieres : brevets, lock-in, protocoles fermes | Proteger une rente existante |
| `FUD` | Semer le doute pour ralentir l'adoption d'un concurrent | Concurrent en forte traction |
| `strangler-fig` | Remplacer progressivement un systeme legacy | Systeme legacy bloquant l'evolution |
| `signal-distortion` | Induire les concurrents en erreur sur l'intention strategique | Desinformation concurrentielle |

### Exemples

```
gameplay ILC on API Plateforme
gameplay open-source "Commoditiser le compute pour capturer la couche IA" on Infrastructure Cloud
gameplay strangler-fig on Legacy CRM
gameplay land-grab on Nouveau Marche
gameplay FUD "Ralentir l'adoption du concurrent" on Solution Rivale
```

La description entre guillemets est optionnelle et fournit du contexte strategique supplementaire.

---

## Focus

Le mot-cle `focus` met en valeur un composant et tous ses descendants dans la chaine de valeur. Les elements hors focus passent en opacite reduite.

```
focus Moteur de Calcul
```

Plusieurs focus peuvent etre declares -- leurs sous-arbres sont alors reunis :

```
focus Moteur de Calcul
focus Flux Temps Reel Capteurs
```

Le focus est utile pour presenter une partie specifique de la carte lors d'une discussion ou d'une presentation.

---

## Regles sur les identifiants

Les identifiants (noms de composants, de groupes, etc.) :

- Commencent par une lettre ou un chiffre
- Peuvent contenir des lettres, chiffres, `.`, `-`, `'`, `_` et des **espaces**
- Les espaces sont autorises a l'interieur (ex: `Application Mobile`)
- Ne peuvent pas etre un mot reserve utilise seul

### Mots reserves

`anchor`, `component`, `submap`, `pipeline`, `group`, `note`, `warning`, `signal`, `gameplay`, `legend`, `focus`, `title`, `date`, `author`, `scope`, `question`, `stages`, `doctrine`, `evolution`, `type`, `asset`, `color`, `visibility`, `cost`, `team`, `build`, `buy`, `outsource`, `accelerating`, `stagnating`, `declining`, `co-evolution`, `red-queen`, `commoditization`, `network-effects`, `economies-of-scale`, `ILC`, `open-source`, `land-grab`, `embrace-extend`, `tower-moat`, `FUD`, `strangler-fig`, `signal-distortion`, `explorer`, `settler`, `town-planner`, `pioneer`, `villager`, `tech`, `financial`, `human`, `relational`, `social`, `on`, `hygiene`, `context`, `excellence`, `legend`

---

## Regles semantiques

1. Tout noeud reference dans un edge, une annotation, un signal ou un pipeline **doit etre declare** avec une position quelque part dans le document.
2. Le nom d'un pipeline **doit correspondre** a un composant deja declare.
3. Les pipelines **ne peuvent pas etre imbriques**.
4. Les groupes **ne creent pas de namespace** -- les composants restent globaux.
5. Les ancres **ne necessitent pas** de position d'evolution (elles sont placees automatiquement en haut).
6. Le graphe de dependances **doit etre acyclique** (pas de dependances circulaires).
7. Un noeud orphelin (non reference dans un edge) est place en bas de la carte par defaut.

---

## Exemple complet

```wtg2
// Carte Wardley -- Plateforme de Navigation GPS

title: Plateforme de Navigation -- Strategie 2026
date: 2026-01-15
author: Cellule Strategie Produit
scope: Application mobile B2C, marche europeen
question: "Ou investir pour se differencier face a Google Maps ?"
doctrine: context

stages: Genese, Sur-mesure, Produit, Commodite
legend

// Ancres
anchor Automobiliste
anchor Collectivite Locale

// Couche visible
Application Mobile : III.5
Itineraire Affiche : III.2
Alertes Trafic en Temps Reel : II.3

// Coeur metier avec inertie qualifiee
Moteur de Calcul d'Itineraire : II.7 !!(tech,human) >> III.5 {
  type: build
  asset: tech
  color: #3498DB
  cost: "1.2M EUR/an, 12 ETP"
  note: "Notre differenciateur cle"
}

// Pipeline : le moteur existe sous plusieurs formes
pipeline Moteur de Calcul d'Itineraire {
  Algo Dijkstra Classique : III.5
  Algo Predictif IA : II.3
  Algo Quantique : I.2
}

Modele de Donnees Carto : III.1 (buy)
API Partenaires B2G : II.1 {
  asset: relational
}

// Infrastructure
Donnees OSM : III.8 (buy)
Flux Temps Reel Capteurs : I.8 !(relational) >> II.5 {
  type: build
  asset: tech
  color: #E67E22
  cost: "300k EUR/an"
  note: "Partenariat en cours avec Waze/TomTom"
}
Infrastructure Cloud : IV.3 (buy) {
  cost: "500k EUR/an, en hausse de 30%"
}
CDN : IV.5 (buy)
Reseau Mobile : IV.7 (outsource)

submap Systeme de Paiement : III.6

// Chaine de valeur
Automobiliste -> Application Mobile -> Itineraire Affiche -> Moteur de Calcul d'Itineraire
Application Mobile -> Alertes Trafic en Temps Reel -> Flux Temps Reel Capteurs
Moteur de Calcul d'Itineraire -> Modele de Donnees Carto -> Donnees OSM
Moteur de Calcul d'Itineraire -> Infrastructure Cloud
Flux Temps Reel Capteurs -> Infrastructure Cloud

Collectivite Locale -> API Partenaires B2G -[Open Data, licence annuelle]-> Modele de Donnees Carto
Collectivite Locale -> Alertes Trafic en Temps Reel

Application Mobile -> CDN -> Infrastructure Cloud
Infrastructure Cloud -> Reseau Mobile
Application Mobile -> Systeme de Paiement

// Lien vers un membre specifique du pipeline
Alertes Trafic en Temps Reel -> Moteur de Calcul d'Itineraire:Algo Predictif IA

// Groupes avec types d'equipe
group Equipe Core Navigation {
  team: settler
  Moteur de Calcul d'Itineraire
  Algo Predictif IA
  Modele de Donnees Carto
}

group Equipe Plateforme {
  team: town-planner
  Infrastructure Cloud
  CDN
  Systeme de Paiement
}

group Equipe Data {
  team: explorer
  Flux Temps Reel Capteurs
  Donnees OSM
  Algo Quantique
}

// Avertissements
warning "SPOF -- aucun fallback si indisponible" on Moteur de Calcul d'Itineraire
warning "Vendor lock-in AWS, cout en hausse de 30%/an" on Infrastructure Cloud
warning "Dependance critique a un seul fournisseur" on Donnees OSM

// Notes strategiques
note "Candidat a l'externalisation Q4 2026" on Systeme de Paiement
note "Partenariat signe avec 12 metropoles" on API Partenaires B2G
note "Budget R&D 400k EUR, horizon 2028" on Algo Quantique

// Signaux de marche
signal accelerating on Algo Predictif IA
signal co-evolution on Flux Temps Reel Capteurs
signal stagnating on Algo Dijkstra Classique
signal commoditization on Infrastructure Cloud
signal declining on Donnees OSM

// Manoeuvres strategiques
gameplay strangler-fig "Remplacer Dijkstra par IA Predictive" on Moteur de Calcul d'Itineraire
gameplay open-source "Commoditiser les donnees carto pour reduire la dependance" on Donnees OSM
gameplay ILC on API Partenaires B2G

// Focus (optionnel)
focus Moteur de Calcul d'Itineraire
```

---

## Resume rapide de la syntaxe

| Element | Syntaxe | Exemple |
|---------|---------|---------|
| Metadonnee | `cle: valeur` | `title: Ma Carte` |
| Ancre | `anchor Nom` | `anchor Client` |
| Composant | `Nom : Position` | `API : III.5` |
| Composant type | `Nom : Position (type)` | `DB : III.8 (buy)` |
| Sous-carte | `submap Nom : Position` | `submap Paiement : III.6` |
| Bloc | `Nom : Position { ... }` | voir ci-dessus |
| Evolution | `Nom : Pos >> Pos` | `API : II.7 >> III.5` |
| Inertie | `! !! !!!` | `API : II.7 !! >> III.5` |
| Inertie qualifiee | `!!(types)` | `API : II.7 !!(tech,human) >> III.5` |
| Visibilite | `@0.0-1.0` | `API : III.5 @0.8` |
| Pipeline | `pipeline Nom { ... }` | voir ci-dessus |
| Edge | `A -> B` | `Client -> App` |
| Edge annote | `A -[texte]-> B` | `A -[flux critique]-> B` |
| Edge bidi | `A <-> B` | `A <-> B` |
| Edge chaine | `A -> B -> C` | `Client -> App -> DB` |
| Ref pipeline | `Pipeline:Membre` | `Moteur:Algo IA` |
| Groupe | `group Nom { ... }` | voir ci-dessus |
| Note | `note "texte" on Nom` | `note "Important" on API` |
| Warning | `warning "texte" on Nom` | `warning "Risque" on DB` |
| Signal | `signal type on Nom` | `signal accelerating on IA` |
| Gameplay | `gameplay type on Nom` | `gameplay ILC on API` |
| Focus | `focus Nom` | `focus Moteur` |
| Legende | `legend` | `legend` |
| Stages | `stages: A, B, C, D` | `stages: Genese, Sur-mesure, Produit, Commodite` |
