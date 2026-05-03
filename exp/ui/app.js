// ================================================================
// 1. CodeMirror WTG2 syntax mode
// ================================================================
CodeMirror.defineSimpleMode("wtg2", {
    start: [
        { regex: /\/\/.*/, token: "comment" },
        { regex: /\/\*.*?\*\//, token: "comment" },
        { regex: /^(title|date|author|scope|question|stages)\s*:/, token: "keyword" },
        { regex: /^(anchor|submap|pipeline|group|component)\b/, token: "keyword" },
        { regex: /^(note|warning)\b/, token: "tag" },
        { regex: /^(signal)\s+(accelerating|stagnating|declining)\b/, token: ["tag", "variable-2"] },
        { regex: /->|<->/, token: "operator" },
        { regex: /-\[/, token: "operator", push: "edgeLabel" },
        { regex: />>/, token: "operator" },
        { regex: /!!?/, token: "variable-2" },
        { regex: /\(buy\)|\(outsource\)|\(build\)/, token: "atom" },
        { regex: /\{/, token: "bracket", indent: true },
        { regex: /\}/, token: "bracket", dedent: true },
        { regex: /"[^"]*"/, token: "string" },
        { regex: /\b[IVX]+\.\d+\b/, token: "number" },
        { regex: /#[0-9a-fA-F]{3,8}\b/, token: "number" },
        { regex: /\b(type|color|note)\s*:/, token: "property" },
    ],
    edgeLabel: [
        { regex: /\]->/, token: "operator", pop: true },
        { regex: /./, token: "string" },
    ],
    meta: { lineComment: "//" }
});

// ================================================================
// 2. Templates
// ================================================================
const TEMPLATES = {
    teashop: {
        meta: { title: "Tea Shop", author: "", question: "" },
        stages: ["Genesis", "Custom-Built", "Product", "Commodity"],
        chains: "Customer -> Cup of Tea -> Tea\nCup of Tea -> Hot Water -> Water\nHot Water -> Kettle -> Power\nCup of Tea -> Cup",
        components: [
            { name: "Customer", kind: "anchor", evolution: 50, type: "", evolving: false, evolvedTo: 75, inertia: 0 },
            { name: "Cup of Tea", kind: "component", evolution: 62, type: "", evolving: false, evolvedTo: 75, inertia: 0 },
            { name: "Tea", kind: "component", evolution: 90, type: "", evolving: false, evolvedTo: 95, inertia: 0 },
            { name: "Hot Water", kind: "component", evolution: 80, type: "", evolving: false, evolvedTo: 85, inertia: 0 },
            { name: "Water", kind: "component", evolution: 90, type: "", evolving: false, evolvedTo: 95, inertia: 0 },
            { name: "Kettle", kind: "component", evolution: 40, type: "build", evolving: true, evolvedTo: 82, inertia: 1, inertiaKinds: ["tech"] },
            { name: "Power", kind: "component", evolution: 92, type: "outsource", evolving: false, evolvedTo: 95, inertia: 0 },
            { name: "Cup", kind: "component", evolution: 95, type: "", evolving: false, evolvedTo: 97, inertia: 0 },
        ],
        edges: [
            { from: "Customer", to: "Cup of Tea" },
            { from: "Cup of Tea", to: "Tea" },
            { from: "Cup of Tea", to: "Hot Water" },
            { from: "Hot Water", to: "Water" },
            { from: "Hot Water", to: "Kettle" },
            { from: "Kettle", to: "Power" },
            { from: "Cup of Tea", to: "Cup" },
        ],
        groups: [],
        annotations: [{ kind: "note", text: "Considering outsourcing", target: "Kettle" }],
        signals: [],
        legend: true
    },
    navigation: {
        meta: { title: "Plateforme de Navigation", author: "Cellule Strategie Produit", question: "Ou investir pour se differencier ?" },
        stages: ["Genese", "Sur-mesure", "Produit", "Commodite"],
        chains: "Automobiliste -> Application Mobile -> Itineraire Affiche -> Moteur de Calcul\nApplication Mobile -> Alertes Trafic -> Flux Temps Reel\nMoteur de Calcul -> Modele de Donnees -> Donnees OSM\nMoteur de Calcul -> Infrastructure Cloud\nApplication Mobile -> CDN -> Infrastructure Cloud\nInfrastructure Cloud -> Reseau Mobile",
        components: [
            { name: "Automobiliste", kind: "anchor", evolution: 50, type: "", evolving: false, evolvedTo: 75, inertia: 0 },
            { name: "Application Mobile", kind: "component", evolution: 62, type: "", evolving: false, evolvedTo: 75, inertia: 0 },
            { name: "Itineraire Affiche", kind: "component", evolution: 55, type: "", evolving: false, evolvedTo: 75, inertia: 0 },
            { name: "Moteur de Calcul", kind: "component", evolution: 42, type: "build", evolving: false, evolvedTo: 62, inertia: 0, isPipeline: true, pipelineMembers: [{ name: "Algo Dijkstra", evolution: 82 }, { name: "Algo Predictif IA", evolution: 37 }, { name: "Algo Quantique", evolution: 5 }] },
            { name: "Alertes Trafic", kind: "component", evolution: 37, type: "", evolving: false, evolvedTo: 75, inertia: 0 },
            { name: "Flux Temps Reel", kind: "component", evolution: 27, type: "build", evolving: true, evolvedTo: 42, inertia: 1 },
            { name: "Modele de Donnees", kind: "component", evolution: 52, type: "buy", evolving: false, evolvedTo: 75, inertia: 0 },
            { name: "Donnees OSM", kind: "component", evolution: 70, type: "buy", evolving: false, evolvedTo: 75, inertia: 0 },
            { name: "Infrastructure Cloud", kind: "component", evolution: 82, type: "buy", evolving: false, evolvedTo: 75, inertia: 0 },
            { name: "CDN", kind: "component", evolution: 87, type: "buy", evolving: false, evolvedTo: 75, inertia: 0 },
            { name: "Reseau Mobile", kind: "component", evolution: 92, type: "outsource", evolving: false, evolvedTo: 75, inertia: 0 },
        ],
        edges: [
            { from: "Automobiliste", to: "Application Mobile" },
            { from: "Application Mobile", to: "Itineraire Affiche" },
            { from: "Itineraire Affiche", to: "Moteur de Calcul" },
            { from: "Application Mobile", to: "Alertes Trafic" },
            { from: "Alertes Trafic", to: "Flux Temps Reel" },
            { from: "Moteur de Calcul", to: "Modele de Donnees" },
            { from: "Modele de Donnees", to: "Donnees OSM" },
            { from: "Moteur de Calcul", to: "Infrastructure Cloud" },
            { from: "Application Mobile", to: "CDN" },
            { from: "CDN", to: "Infrastructure Cloud" },
            { from: "Infrastructure Cloud", to: "Reseau Mobile" },
        ],
        groups: [
            { name: "Equipe Core Navigation", members: ["Moteur de Calcul", "Modele de Donnees"], color: "#3498DB" },
            { name: "Equipe Plateforme", members: ["Infrastructure Cloud", "CDN"], color: "#E67E22" },
        ],
        annotations: [
            { kind: "warning", text: "SPOF - aucun fallback", target: "Moteur de Calcul" },
            { kind: "note", text: "Partenariat en cours", target: "Flux Temps Reel" },
        ],
        signals: [
            { type: "accelerating", target: "Flux Temps Reel" },
            { type: "declining", target: "Donnees OSM" },
        ],
        legend: true
    },
    blank: {
        meta: { title: "Ma Carte Wardley", author: "", question: "" },
        stages: ["", "", "", ""],
        chains: "",
        components: [],
        edges: [],
        groups: [],
        annotations: [],
        signals: [],
        legend: false
    }
};

// ================================================================
// 2b. Internationalization (i18n)
// ================================================================
const translations = {
    fr: {
        'toolbar.title': 'WTG2 Playground',
        'mode.guided': 'Guide',
        'mode.editor': 'Editeur',
        'toolbar.ratio': 'Ratio',
        'toolbar.resolution': 'Def.',
        'toolbar.dlSvg': 'Telecharger SVG',
        'toolbar.dlPng': 'Telecharger PNG',
        'toolbar.detach': 'Ouvrir dans une fenetre separee',
        'toolbar.dlWtg2': 'Exporter le source WTG2',
        'toolbar.importWtg2': 'Importer un fichier WTG2',
        'toolbar.importOwm': 'Importer un fichier OWM (Online Wardley Maps)',
        'owm.title': 'Importer une carte OWM',
        'owm.uploadLabel': 'Charger un fichier',
        'owm.uploadHint': 'Fichier .owm ou .txt au format Online Wardley Maps',
        'owm.orSeparator': '— ou —',
        'owm.pasteLabel': 'Coller du texte OWM',
        'owm.pastePlaceholder': 'title Ma Carte\\ncomponent Client [0.95, 0.5]\\n...',
        'owm.import': 'Importer',
        'owm.cancel': 'Annuler',
        'owm.emptyError': 'Veuillez selectionner un fichier ou coller du texte OWM.',
        'owm.convertError': 'Erreur de conversion : {message}',
        'toolbar.share': 'Copier le lien de partage',
        'toolbar.shareLabel': 'Partager',
        'status.loading': 'Chargement WASM...',
        'status.ready': 'Pret',
        'status.ok': 'OK',
        'status.wasmNotLoaded': 'WASM pas encore charge',
        'mobile.edit': 'Editer',
        'mobile.preview': 'Apercu',
        'template.label': 'Template :',
        'template.choose': '-- Choisir --',
        'template.teashop': 'Tea Shop (simple)',
        'template.navigation': 'Plateforme Navigation (riche)',
        'template.blank': 'Vide',
        'step.mapInfo': 'Carte',
        'step.valueChain': 'Chaine de valeur',
        'step.evolution': 'Evolution',
        'step.enrichment': 'Enrichissement',
        'step.navigation': 'Navigation',
        'section.mapInfo': 'Infos carte',
        'field.title': 'Titre',
        'field.titlePlaceholder': 'Ma Carte Wardley',
        'field.author': 'Auteur',
        'field.authorPlaceholder': 'Votre nom',
        'field.question': 'Question strategique',
        'field.questionPlaceholder': 'Ou investir ?',
        'section.valueChain': 'Chaine de valeur',
        'field.dependencies': 'Dependances (une chaine par ligne)',
        'field.dependenciesPlaceholder': 'Client -> Application -> Base de donnees\nApplication -> API -> Service externe',
        'field.dependenciesHint': 'Tapez des chaines avec -> pour lier les composants. Les composants sans dependance entrante deviennent automatiquement des ancres.',
        'comp.detected': 'Composants detectes',
        'comp.count': '{count} composant{s}',
        'comp.anchor': 'ancre',
        'comp.component': 'composant',
        'comp.delete': 'Supprimer',
        'comp.focus': 'Focus',
        'comp.typeToggle': 'Changer le type (build / buy / outsource)',
        'evo.type': 'Type:',
        'evo.evolution': 'Evolution',
        'evo.towards': 'vers',
        'evo.inertia0': "Pas d'inertie",
        'evo.inertia1': '! Moderee',
        'evo.inertia2': '!! Forte',
        'evo.inertia3': '!!! Bloquante',
        'evo.inertiaKinds': 'Types :',
        'evo.inertia.tech': 'Tech',
        'evo.inertia.financial': 'Financier',
        'evo.inertia.human': 'Humain',
        'evo.inertia.relational': 'Relationnel',
        'evo.inertia.social': 'Social',
        'enrich.groups': 'Groupes',
        'enrich.addGroup': '+ Ajouter un groupe',
        'enrich.annotations': 'Annotations',
        'enrich.addAnnotation': '+ Ajouter une annotation',
        'enrich.signals': 'Signaux',
        'enrich.addSignal': '+ Ajouter un signal',
        'enrich.showLegend': 'Afficher la legende',
        'evo.pipeline': 'Pipeline',
        'pipeline.addMember': '+ Membre',
        'pipeline.newMember': 'Nouveau membre',
        'group.defaultName': 'Nouveau groupe',
        'group.namePlaceholder': 'Nom du groupe',
        'annotation.note': 'Note',
        'annotation.warning': 'Warning',
        'annotation.textPlaceholder': 'Texte...',
        'annotation.selectComp': '-- Composant --',
        'nav.prev': 'Precedent',
        'nav.next': 'Suivant',
        'nav.finish': 'Terminer',
        'nav.animationControls': 'Controles d\'animation',
        'nav.mode': 'Mode',
        'nav.modeDepth': 'Par profondeur de dependance',
        'nav.modeYRank': 'Par position Y',
        'nav.reset': 'Reinitialiser',
        'nav.animPrev': 'Precedent',
        'nav.animNext': 'Suivant',
        'nav.showAll': 'Tout afficher',
        'nav.keyboardHint': 'Clavier : Fleches pour naviguer, Debut pour reinitialiser, Fin pour tout afficher',
        'nav.stepCounter': 'Etape {current} / {total}',
        'empty.message': 'Commencez a definir votre chaine de valeur pour voir l\'apercu',
        'share.tooLong': 'Attention : lien tres long ({length} car.). Certains navigateurs pourraient le tronquer.',
        'share.copied': 'Lien copie dans le presse-papier !',
        'share.manual': "Lien mis a jour dans la barre d'adresse (copie manuelle)",
        'share.error': 'Erreur : {message}',
        'status.urlError': 'Erreur URL: {message}',
        'toolbar.collab': 'Rejoindre une session collaborative',
        'toolbar.collabLabel': 'Collaborer',
        'menu.newProject': 'Nouveau Projet',
        'menu.documentation': 'Documentation',
        'menu.licenses': 'Licences tierces',
        'menu.exportPng': 'Export PNG',
        'menu.exportSvg': 'Export SVG',
        'menu.exportGif': 'Export GIF animé',
        'menu.exportApng': 'Export APNG animé',
        'export.progress': 'Export en cours... {percent}%',
        'export.encoding': 'Encodage...',
        'export.cancel': 'Annuler',
        'export.noAnimations': 'Aucune animation détectée, export statique.',
        'export.error': 'Erreur d\'export : {message}',
        'menu.confirmNew': 'Creer un nouveau projet ? Les modifications non sauvegardees seront perdues.',
        'menu.myMaps': 'Mes Cartes',
        'menu.templates': 'Templates',
        'menu.confirmDelete': 'Supprimer cette carte ? Cette action est irreversible.',
        'maps.title': 'Mes Cartes',
        'maps.empty': 'Aucune carte sauvegardee.',
        'maps.open': 'Ouvrir',
        'maps.delete': 'Supprimer',
        'maps.close': 'Fermer',
        'maps.created': 'Creee le',
        'maps.modified': 'Modifiee le',
        'maps.current': '(en cours)',
        'maps.untitled': 'Sans titre',
        'toolbar.llmHint': 'Syntaxe WTG2 :',
        'toolbar.llmHintTitle': 'Fichier skill utilisable avec ChatGPT, Claude, etc.',
        'toolbar.llmHintLink': '\uD83E\uDD16 Skill LLM',
        'collab.title': 'Session collaborative',
        'collab.urlLabel': 'URL de la session',
        'collab.urlPlaceholder': 'ws://serveur:8081/s/session-id/access-id',
        'collab.urlHint': 'Collez l\'URL fournie par le serveur de collaboration',
        'collab.nameLabel': 'Votre nom',
        'collab.namePlaceholder': 'Alice',
        'collab.connect': 'Se connecter',
        'collab.cancel': 'Annuler',
        'collab.disconnect': 'Quitter',
        'collab.connected': 'Session collaborative',
        'collab.spectator': 'Spectateur',
        'collab.invalidUrl': 'URL invalide (doit commencer par ws:// ou wss://)',
        'collab.connectionError': 'Erreur de connexion : {message}',
        'collab.reconnecting': 'Reconnexion...',
        'onboarding.welcome.title': 'Bienvenue dans le Playground',
        'onboarding.welcome.body': 'Créez votre carte Wardley en 4 étapes simples. Ce guide rapide vous montre les bases.',
        'onboarding.step1.title': 'Étape 1 : Chaîne de valeur',
        'onboarding.step1.body': 'Définissez vos composants et leurs dépendances. Tapez des chaînes avec -> pour les relier.',
        'onboarding.step1.tryit': 'Essayez : tapez <strong>Client -> Application -> Base de données</strong>',
        'onboarding.step2.title': 'Étape 2 : Évolution',
        'onboarding.step2.body': 'Positionnez chaque composant sur l\'axe d\'évolution. Cliquez sur une carte pour l\'ouvrir et ajuster le curseur.',
        'onboarding.step3.title': 'Étape 3 : Enrichissement',
        'onboarding.step3.body': 'Ajoutez des groupes, annotations et signaux pour donner du contexte à votre carte.',
        'onboarding.preview.title': 'Aperçu en temps réel',
        'onboarding.preview.body': 'Vos modifications apparaissent instantanément dans l\'aperçu SVG à droite.',
        'onboarding.skip': 'Passer',
        'onboarding.next': 'Suivant',
        'onboarding.finish': 'Commencer',
    },
    en: {
        'toolbar.title': 'WTG2 Playground',
        'mode.guided': 'Guided',
        'mode.editor': 'Editor',
        'toolbar.ratio': 'Ratio',
        'toolbar.resolution': 'Res.',
        'toolbar.dlSvg': 'Download SVG',
        'toolbar.dlPng': 'Download PNG',
        'toolbar.detach': 'Open in separate window',
        'toolbar.dlWtg2': 'Export WTG2 source',
        'toolbar.importWtg2': 'Import WTG2 file',
        'toolbar.importOwm': 'Import OWM file (Online Wardley Maps)',
        'owm.title': 'Import OWM Map',
        'owm.uploadLabel': 'Upload a file',
        'owm.uploadHint': '.owm or .txt file in Online Wardley Maps format',
        'owm.orSeparator': '— or —',
        'owm.pasteLabel': 'Paste OWM text',
        'owm.pastePlaceholder': 'title My Map\\ncomponent Customer [0.95, 0.5]\\n...',
        'owm.import': 'Import',
        'owm.cancel': 'Cancel',
        'owm.emptyError': 'Please select a file or paste OWM text.',
        'owm.convertError': 'Conversion error: {message}',
        'toolbar.share': 'Copy share link',
        'toolbar.shareLabel': 'Share',
        'status.loading': 'Loading WASM...',
        'status.ready': 'Ready',
        'status.ok': 'OK',
        'status.wasmNotLoaded': 'WASM not loaded yet',
        'mobile.edit': 'Edit',
        'mobile.preview': 'Preview',
        'template.label': 'Template:',
        'template.choose': '-- Choose --',
        'template.teashop': 'Tea Shop (simple)',
        'template.navigation': 'Navigation Platform (rich)',
        'template.blank': 'Blank',
        'step.mapInfo': 'Map',
        'step.valueChain': 'Value Chain',
        'step.evolution': 'Evolution',
        'step.enrichment': 'Enrichment',
        'step.navigation': 'Navigation',
        'section.mapInfo': 'Map Info',
        'field.title': 'Title',
        'field.titlePlaceholder': 'My Wardley Map',
        'field.author': 'Author',
        'field.authorPlaceholder': 'Your name',
        'field.question': 'Strategic question',
        'field.questionPlaceholder': 'Where to invest?',
        'section.valueChain': 'Value Chain',
        'field.dependencies': 'Dependencies (one chain per line)',
        'field.dependenciesPlaceholder': 'Customer -> Application -> Database\nApplication -> API -> External Service',
        'field.dependenciesHint': 'Type chains with -> to link components. Components with no inbound dependency become anchors automatically.',
        'comp.detected': 'Detected components',
        'comp.count': '{count} component{s}',
        'comp.anchor': 'anchor',
        'comp.component': 'component',
        'comp.delete': 'Delete',
        'comp.focus': 'Focus',
        'comp.typeToggle': 'Change type (build / buy / outsource)',
        'evo.type': 'Type:',
        'evo.evolution': 'Evolution',
        'evo.towards': 'to',
        'evo.inertia0': 'No inertia',
        'evo.inertia1': '! Moderate',
        'evo.inertia2': '!! Strong',
        'evo.inertia3': '!!! Blocking',
        'evo.inertiaKinds': 'Types:',
        'evo.inertia.tech': 'Tech',
        'evo.inertia.financial': 'Financial',
        'evo.inertia.human': 'Human',
        'evo.inertia.relational': 'Relational',
        'evo.inertia.social': 'Social',
        'enrich.groups': 'Groups',
        'enrich.addGroup': '+ Add group',
        'enrich.annotations': 'Annotations',
        'enrich.addAnnotation': '+ Add annotation',
        'enrich.signals': 'Signals',
        'enrich.addSignal': '+ Add signal',
        'enrich.showLegend': 'Show legend',
        'evo.pipeline': 'Pipeline',
        'pipeline.addMember': '+ Member',
        'pipeline.newMember': 'New member',
        'group.defaultName': 'New group',
        'group.namePlaceholder': 'Group name',
        'annotation.note': 'Note',
        'annotation.warning': 'Warning',
        'annotation.textPlaceholder': 'Text...',
        'annotation.selectComp': '-- Component --',
        'nav.prev': 'Previous',
        'nav.next': 'Next',
        'nav.finish': 'Finish',
        'nav.animationControls': 'Animation Controls',
        'nav.mode': 'Mode',
        'nav.modeDepth': 'By dependency depth',
        'nav.modeYRank': 'By Y position',
        'nav.reset': 'Reset',
        'nav.animPrev': 'Prev',
        'nav.animNext': 'Next',
        'nav.showAll': 'Show All',
        'nav.keyboardHint': 'Keyboard: Arrow keys to navigate, Home to reset, End to show all',
        'nav.stepCounter': 'Step {current} / {total}',
        'empty.message': 'Start defining your value chain to see the preview',
        'share.tooLong': 'Warning: very long link ({length} chars). Some browsers may truncate it.',
        'share.copied': 'Link copied to clipboard!',
        'share.manual': 'Link updated in address bar (manual copy)',
        'share.error': 'Error: {message}',
        'status.urlError': 'URL Error: {message}',
        'toolbar.collab': 'Join a collaborative session',
        'toolbar.collabLabel': 'Collaborate',
        'menu.newProject': 'New Project',
        'menu.documentation': 'Documentation',
        'menu.licenses': 'Third-party licenses',
        'menu.exportPng': 'Export PNG',
        'menu.exportSvg': 'Export SVG',
        'menu.exportGif': 'Export animated GIF',
        'menu.exportApng': 'Export animated APNG',
        'export.progress': 'Exporting... {percent}%',
        'export.encoding': 'Encoding...',
        'export.cancel': 'Cancel',
        'export.noAnimations': 'No animations detected, exporting static image.',
        'export.error': 'Export error: {message}',
        'menu.confirmNew': 'Create a new project? Unsaved changes will be lost.',
        'menu.myMaps': 'My Maps',
        'menu.templates': 'Templates',
        'menu.confirmDelete': 'Delete this map? This action cannot be undone.',
        'maps.title': 'My Maps',
        'maps.empty': 'No saved maps.',
        'maps.open': 'Open',
        'maps.delete': 'Delete',
        'maps.close': 'Close',
        'maps.created': 'Created',
        'maps.modified': 'Modified',
        'maps.current': '(current)',
        'maps.untitled': 'Untitled',
        'toolbar.llmHint': 'WTG2 syntax:',
        'toolbar.llmHintTitle': 'Skill file usable with ChatGPT, Claude, etc.',
        'toolbar.llmHintLink': '\uD83E\uDD16 LLM Skill',
        'collab.title': 'Collaborative session',
        'collab.urlLabel': 'Session URL',
        'collab.urlPlaceholder': 'ws://server:8081/s/session-id/access-id',
        'collab.urlHint': 'Paste the URL provided by the collaboration server',
        'collab.nameLabel': 'Your name',
        'collab.namePlaceholder': 'Alice',
        'collab.connect': 'Connect',
        'collab.cancel': 'Cancel',
        'collab.disconnect': 'Leave',
        'collab.connected': 'Collaborative session',
        'collab.spectator': 'Spectator',
        'collab.invalidUrl': 'Invalid URL (must start with ws:// or wss://)',
        'collab.connectionError': 'Connection error: {message}',
        'collab.reconnecting': 'Reconnecting...',
        'onboarding.welcome.title': 'Welcome to the Playground',
        'onboarding.welcome.body': 'Create your Wardley Map in 4 simple steps. This quick guide shows you the basics.',
        'onboarding.step1.title': 'Step 1: Value Chain',
        'onboarding.step1.body': 'Define your components and their dependencies. Type chains with -> to link them.',
        'onboarding.step1.tryit': 'Try it: type <strong>Customer -> Application -> Database</strong>',
        'onboarding.step2.title': 'Step 2: Evolution',
        'onboarding.step2.body': 'Position each component on the evolution axis. Click a card to open it and adjust the slider.',
        'onboarding.step3.title': 'Step 3: Enrichment',
        'onboarding.step3.body': 'Add groups, annotations and signals to give context to your map.',
        'onboarding.preview.title': 'Live Preview',
        'onboarding.preview.body': 'Your changes appear instantly in the SVG preview on the right.',
        'onboarding.skip': 'Skip',
        'onboarding.next': 'Next',
        'onboarding.finish': 'Get Started',
    }
};

let currentLang = localStorage.getItem('wtg2-lang') || 'fr';

function t(key, params) {
    const str = (translations[currentLang] && translations[currentLang][key])
             || translations['fr'][key] || key;
    if (!params) return str;
    return str.replace(/\{(\w+)\}/g, (_, k) => params[k] !== undefined ? params[k] : '{' + k + '}');
}

function applyStaticTranslations() {
    document.querySelectorAll('[data-i18n]').forEach(el => {
        el.textContent = t(el.dataset.i18n);
    });
    document.querySelectorAll('[data-i18n-title]').forEach(el => {
        el.title = t(el.dataset.i18nTitle);
    });
    document.querySelectorAll('[data-i18n-placeholder]').forEach(el => {
        el.placeholder = t(el.dataset.i18nPlaceholder);
    });
}

// ================================================================
// 3. Wizard State
// ================================================================
let wizardState = {
    meta: { title: "", author: "", question: "" },
    stages: ["", "", "", ""],
    components: [],
    edges: [],
    groups: [],
    annotations: [],
    signals: [],
    legend: false,
    focus: ""
};

let currentStep = 0;
let currentMode = 'guided'; // 'guided' or 'editor'
let editor = null;
let renderTimer = null;
let detachedWindow = null;
let expandedEvoIdx = -1;
let currentMapId = null;
let mapsIndex = [];

// ================================================================
// 3b. Undo/Redo
// ================================================================
let undoStack = [];
let redoStack = [];
let undoTimer = null;
const UNDO_MAX = 50;

function scheduleUndoPush() {
    clearTimeout(undoTimer);
    undoTimer = setTimeout(pushUndoState, 600);
}

function pushUndoState() {
    undoStack.push(JSON.parse(JSON.stringify(wizardState)));
    if (undoStack.length > UNDO_MAX) undoStack.shift();
    redoStack = [];
    updateUndoButtons();
}

function undo() {
    if (currentMode !== 'guided' || undoStack.length === 0) return;
    redoStack.push(JSON.parse(JSON.stringify(wizardState)));
    wizardState = undoStack.pop();
    expandedEvoIdx = -1;
    syncWizardUI();
    scheduleRender();
    updateUndoButtons();
}

function redo() {
    if (currentMode !== 'guided' || redoStack.length === 0) return;
    undoStack.push(JSON.parse(JSON.stringify(wizardState)));
    wizardState = redoStack.pop();
    expandedEvoIdx = -1;
    syncWizardUI();
    scheduleRender();
    updateUndoButtons();
}

function updateUndoButtons() {
    const undoBtn = document.getElementById('undo-btn');
    const redoBtn = document.getElementById('redo-btn');
    if (undoBtn) undoBtn.disabled = undoStack.length === 0;
    if (redoBtn) redoBtn.disabled = redoStack.length === 0;
}

document.addEventListener('keydown', function(e) {
    if (currentMode !== 'guided') return;
    if ((e.ctrlKey || e.metaKey) && e.key === 'z' && !e.shiftKey) { e.preventDefault(); undo(); }
    if ((e.ctrlKey || e.metaKey) && e.key === 'z' && e.shiftKey) { e.preventDefault(); redo(); }
    if ((e.ctrlKey || e.metaKey) && e.key === 'y') { e.preventDefault(); redo(); }
});

// ================================================================
// 4. Evolution helpers
// ================================================================
function evoToRoman(val) {
    // val: 0-99 integer -> "I.0" to "IV.9"
    val = Math.max(0, Math.min(99, Math.round(val)));
    const phases = ['I', 'II', 'III', 'IV'];
    const phaseIdx = Math.min(3, Math.floor(val / 25));
    const sub = Math.min(9, Math.round((val - phaseIdx * 25) / 25 * 10));
    return phases[phaseIdx] + '.' + sub;
}

function romanToEvo(roman) {
    // "III.5" -> 62
    const match = roman.match(/^(I{1,3}|IV)\s*\.\s*(\d)$/);
    if (!match) return 50;
    const phases = { 'I': 0, 'II': 1, 'III': 2, 'IV': 3 };
    const phaseIdx = phases[match[1]];
    if (phaseIdx === undefined) return 50;
    const sub = parseInt(match[2]);
    return phaseIdx * 25 + Math.round(sub / 10 * 25);
}

// ================================================================
// 5. Chain parser
// ================================================================
function parseChains() {
    const text = chainEditor ? chainEditor.getValue() : document.getElementById('chain-input').value;
    const componentMap = new Map();
    const edges = [];
    const inbound = new Set();

    for (const rawLine of text.split('\n')) {
        const line = rawLine.trim();
        if (!line || line.startsWith('//')) continue;

        const parts = line.split(/\s*->\s*/).map(s => s.trim()).filter(Boolean);
        if (parts.length === 0) continue;

        for (const name of parts) {
            if (!componentMap.has(name)) {
                componentMap.set(name, { order: componentMap.size });
            }
        }
        for (let i = 0; i < parts.length - 1; i++) {
            edges.push({ from: parts[i], to: parts[i + 1] });
            inbound.add(parts[i + 1]);
        }
    }

    // Preserve existing component evolution positions when possible
    const existingMap = new Map();
    for (const c of wizardState.components) {
        existingMap.set(c.name, c);
    }

    const total = componentMap.size;
    const components = [];
    for (const [name, data] of componentMap) {
        const existing = existingMap.get(name);
        const isAnchor = !inbound.has(name);
        if (existing) {
            existing.kind = isAnchor ? 'anchor' : existing.kind;
            components.push(existing);
        } else {
            // Assign default evolution spread
            const fraction = total > 1 ? data.order / (total - 1) : 0.5;
            const evo = Math.round(20 + fraction * 70); // Spread from 20 to 90
            components.push({
                name,
                kind: isAnchor ? 'anchor' : 'component',
                evolution: evo,
                type: '',
                evolving: false,
                evolvedTo: Math.min(99, evo + 15),
                inertia: 0,
                isPipeline: false,
                pipelineMembers: []
            });
        }
    }

    wizardState.components = components;
    wizardState.edges = deduplicateEdges(edges);

    // Clean up groups/annotations/signals referencing removed components
    const names = new Set(components.map(c => c.name));
    wizardState.groups = wizardState.groups.map(g => ({
        ...g,
        members: g.members.filter(m => names.has(m))
    }));
    wizardState.annotations = wizardState.annotations.filter(a => names.has(a.target));
    wizardState.signals = wizardState.signals.filter(s => names.has(s.target));

    renderCompList();
    onWizardChange();
}

function deduplicateEdges(edges) {
    const seen = new Set();
    return edges.filter(e => {
        const key = e.from + '|' + e.to;
        if (seen.has(key)) return false;
        seen.add(key);
        return true;
    });
}

// ================================================================
// 6. WTG2 Code Generator
// ================================================================
function generateWTG2(state) {
    const lines = [];

    // Metadata
    if (state.meta.title) lines.push('title: ' + state.meta.title);
    if (state.meta.author) lines.push('author: ' + state.meta.author);
    if (state.meta.question) lines.push('question: "' + state.meta.question + '"');
    if (lines.length > 0) lines.push('');

    // Stages (evolution zone labels)
    const hasStages = state.stages && state.stages.some(s => s.trim() !== '');
    if (hasStages) {
        lines.push('stages: ' + state.stages.join(', '));
        lines.push('');
    }

    // Anchors
    const anchors = state.components.filter(c => c.kind === 'anchor');
    for (const a of anchors) {
        if (a.evolution !== undefined && a.evolution !== 50) {
            lines.push('anchor ' + a.name + ' : ' + evoToRoman(a.evolution));
        } else {
            lines.push('anchor ' + a.name);
        }
    }
    if (anchors.length) lines.push('');

    // Components
    const comps = state.components.filter(c => c.kind !== 'anchor');
    for (const c of comps) {
        let line = c.name + ' : ' + evoToRoman(c.evolution);
        if (c.inertia > 0) {
            line += ' ' + '!'.repeat(c.inertia);
            if (c.inertiaKinds && c.inertiaKinds.length > 0) line += '(' + c.inertiaKinds.join(',') + ')';
        }
        if (c.evolving && c.evolvedTo !== undefined) line += ' >> ' + evoToRoman(c.evolvedTo);
        if (c.type) line += ' (' + c.type + ')';
        lines.push(line);
    }
    if (comps.length) lines.push('');

    // Pipelines
    const pipelineComps = state.components.filter(c => c.isPipeline && c.pipelineMembers && c.pipelineMembers.length > 0);
    for (const pc of pipelineComps) {
        const valid = pc.pipelineMembers.filter(m => m.name && m.name.trim());
        if (!valid.length) continue;
        lines.push('pipeline ' + pc.name + ' {');
        for (const m of valid) lines.push('  ' + m.name + ' : ' + evoToRoman(m.evolution));
        lines.push('}');
        lines.push('');
    }

    // Edges
    for (const e of state.edges) {
        lines.push(e.from + ' -> ' + e.to);
    }
    if (state.edges.length) lines.push('');

    // Groups
    for (const g of state.groups) {
        lines.push('group ' + g.name + ' {');
        for (const m of g.members) lines.push('  ' + m);
        if (g.color) lines.push('  color: ' + g.color);
        lines.push('}');
        lines.push('');
    }

    // Annotations
    for (const a of state.annotations) {
        if (a.text && a.target) {
            lines.push(a.kind + ' "' + a.text + '" on ' + a.target);
        }
    }

    // Signals
    for (const s of state.signals) {
        if (s.type && s.target) {
            lines.push('signal ' + s.type + ' on ' + s.target);
        }
    }

    // Legend
    if (state.legend) {
        lines.push('');
        lines.push('legend');
    }

    // Focus
    if (state.focus) {
        lines.push('');
        lines.push('focus ' + state.focus);
    }

    return lines.join('\n');
}

// ================================================================
// 7. UI Renderers
// ================================================================

function renderCompList() {
    const container = document.getElementById('comp-list');
    const count = document.getElementById('comp-count');
    const n = wizardState.components.length;
    count.textContent = t('comp.count', { count: n, s: n !== 1 ? 's' : '' });

    container.innerHTML = wizardState.components.map((c, i) => `
        <div class="comp-item">
            <span class="comp-name">${escapeHtml(c.name)}</span>
            <span class="comp-badge ${c.kind === 'anchor' ? 'anchor' : ''}">${c.kind === 'anchor' ? t('comp.anchor') : t('comp.component')}</span>
            ${c.kind !== 'anchor' ? `<button class="comp-type-toggle ${c.type || ''}" onclick="cycleType(${i})" title="${t('comp.typeToggle')}">${c.type || '\u2014'}</button>` : ''}
            <button class="comp-focus ${wizardState.focus === c.name ? 'active' : ''}" onclick="toggleFocus(${i})" title="${t('comp.focus')}"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="width:13px;height:13px"><circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="3"/></svg></button>
            <button class="comp-remove" onclick="removeComponent(${i})" title="${t('comp.delete')}"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="width:13px;height:13px"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
        </div>
    `).join('');
}

function evoSummaryHtml(c) {
    let parts = [];
    if (c.type) parts.push(`<span class="summary-type ${c.type}">${c.type}</span>`);
    if (c.evolving) parts.push(`<span class="summary-evolving">&rarr; ${evoToRoman(c.evolvedTo)}</span>`);
    if (c.inertia > 0) parts.push(`<span class="summary-inertia">${'!'.repeat(c.inertia)}</span>`);
    return parts.join('');
}

function renderEvoList() {
    const container = document.getElementById('evo-list');
    container.innerHTML = wizardState.components.map((c, compIdx) => {
        const isAnchor = c.kind === 'anchor';
        if (!isAnchor && c.isPipeline) return renderPipelineItem(c, compIdx);
        const isCollapsed = compIdx !== expandedEvoIdx;
        return `
        <div class="evo-item${isCollapsed ? ' collapsed' : ''}" data-evo-idx="${compIdx}">
            <div class="evo-item-header" onclick="toggleEvoCard(${compIdx})">
                <span class="evo-item-name">${escapeHtml(c.name)}</span>
                ${isAnchor ? '<span class="comp-badge anchor">' + t('comp.anchor') + '</span>' : ''}
                <span class="evo-item-summary">${evoSummaryHtml(c)}</span>
                <span class="evo-item-pos" id="evo-pos-${compIdx}">${evoToRoman(c.evolution)}</span>
                <span class="evo-item-chevron"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="width:10px;height:10px"><polyline points="6 9 12 15 18 9"/></svg></span>
            </div>
            <div class="evo-item-body">
                <div class="evo-slider-track">
                    <div class="evo-slider-phases">
                        <div>I</div><div>II</div><div>III</div><div>IV</div>
                    </div>
                    <input type="range" min="0" max="99" value="${c.evolution}"
                        oninput="updateEvolution(${compIdx}, this.value)"
                        onchange="onWizardChange()">
                </div>
                ${isAnchor ? '' : `
                <div class="evo-item-options">
                    <label>${t('evo.type')}</label>
                    <div class="pill-group">
                        <button class="${!c.type ? 'active' : ''}" onclick="setType(${compIdx},'')">&mdash;</button>
                        <button class="${c.type === 'build' ? 'active' : ''}" onclick="setType(${compIdx},'build')">build</button>
                        <button class="${c.type === 'buy' ? 'active' : ''}" onclick="setType(${compIdx},'buy')">buy</button>
                        <button class="${c.type === 'outsource' ? 'active' : ''}" onclick="setType(${compIdx},'outsource')">outsource</button>
                    </div>
                    <label style="margin-left:8px">
                        <input type="checkbox" ${c.evolving ? 'checked' : ''} onchange="toggleEvolving(${compIdx}, this.checked)"> ${t('evo.evolution')}
                    </label>
                    ${c.evolving ? `
                    <span style="font-size:11px;color:#888">${t('evo.towards')}</span>
                    <input type="range" min="0" max="99" value="${c.evolvedTo}" style="width:60px;accent-color:#6ab04c"
                        oninput="updateEvolvedTo(${compIdx}, this.value)"
                        onchange="onWizardChange()">
                    <span style="font-size:11px;color:#6ab04c;font-family:monospace" id="evo-to-${compIdx}">${evoToRoman(c.evolvedTo)}</span>
                    <select onchange="setInertia(${compIdx}, this.value)" style="margin-left:4px">
                        <option value="0" ${c.inertia === 0 ? 'selected' : ''}>${t('evo.inertia0')}</option>
                        <option value="1" ${c.inertia === 1 ? 'selected' : ''}>${t('evo.inertia1')}</option>
                        <option value="2" ${c.inertia === 2 ? 'selected' : ''}>${t('evo.inertia2')}</option>
                        <option value="3" ${c.inertia === 3 ? 'selected' : ''}>${t('evo.inertia3')}</option>
                    </select>
                    ${c.inertia > 0 ? `
                    <span style="font-size:11px;color:#888;margin-left:6px">${t('evo.inertiaKinds')}</span>
                    ${['tech','financial','human','relational','social'].map(k =>
                        `<label style="font-size:11px;margin-left:4px;white-space:nowrap">
                            <input type="checkbox" ${(c.inertiaKinds || []).includes(k) ? 'checked' : ''}
                                onchange="toggleInertiaKind(${compIdx}, '${k}', this.checked)">
                            ${t('evo.inertia.' + k)}
                        </label>`
                    ).join('')}
                    ` : ''}
                    ` : ''}
                    <label style="margin-left:8px">
                        <input type="checkbox" onchange="togglePipeline(${compIdx}, this.checked)"> ${t('evo.pipeline')}
                    </label>
                </div>`}
            </div>
        </div>`;
    }).join('');
    initPipelineDrag();
    initPipelineLabelClicks();
}

function toggleEvoCard(idx) {
    const items = document.querySelectorAll('.evo-item[data-evo-idx]');
    if (expandedEvoIdx === idx) {
        expandedEvoIdx = -1;
        items.forEach(el => el.classList.add('collapsed'));
    } else {
        expandedEvoIdx = idx;
        items.forEach(el => {
            const elIdx = parseInt(el.dataset.evoIdx);
            el.classList.toggle('collapsed', elIdx !== idx);
        });
    }
}

function renderPipelineItem(c, compIdx) {
    const members = (c.pipelineMembers || []).slice().sort((a, b) => a.evolution - b.evolution);
    const minEvo = members.length ? members[0].evolution : 0;
    const maxEvo = members.length ? members[members.length - 1].evolution : 99;
    const spanLeft = (minEvo / 99 * 100).toFixed(2);
    const spanWidth = ((maxEvo - minEvo) / 99 * 100).toFixed(2);

    const handles = members.map((m, sortedIdx) => {
        const origIdx = c.pipelineMembers.indexOf(m);
        const leftPct = (m.evolution / 99 * 100).toFixed(2);
        const belowClass = sortedIdx % 2 !== 0 ? ' pipeline-handle--below' : '';
        return `<div class="pipeline-handle${belowClass}" data-member-idx="${origIdx}" style="left:${leftPct}%">
            <div class="pipeline-handle-grip" data-comp-idx="${compIdx}" data-member-idx="${origIdx}"></div>
            <span class="pipeline-handle-label" data-comp-idx="${compIdx}" data-member-idx="${origIdx}">${escapeHtml(m.name)}</span>
            <span class="pipeline-handle-evo">${evoToRoman(m.evolution)}</span>
            <button class="pipeline-handle-remove" onclick="removePipelineMember(${compIdx},${origIdx})"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="width:10px;height:10px"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
        </div>`;
    }).join('');

    const isCollapsed = compIdx !== expandedEvoIdx;
    return `
    <div class="evo-item evo-item--pipeline${isCollapsed ? ' collapsed' : ''}" data-evo-idx="${compIdx}">
        <div class="evo-item-header" onclick="toggleEvoCard(${compIdx})">
            <span class="evo-item-name">${escapeHtml(c.name)}</span>
            <span class="evo-item-summary">${evoSummaryHtml(c)}<span style="font-size:10px;color:#8691AB;margin-left:2px">pipeline</span></span>
            <span class="evo-item-pos" id="evo-pos-${compIdx}">${evoToRoman(c.evolution)}</span>
            <span class="evo-item-chevron"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="width:10px;height:10px"><polyline points="6 9 12 15 18 9"/></svg></span>
        </div>
        <div class="evo-item-body">
            <div class="evo-slider-track evo-slider-track--secondary">
                <div class="evo-slider-phases">
                    <div>I</div><div>II</div><div>III</div><div>IV</div>
                </div>
                <input type="range" min="0" max="99" value="${c.evolution}"
                    oninput="updateEvolution(${compIdx}, this.value)"
                    onchange="onWizardChange()">
            </div>
            <div class="pipeline-track" data-comp-idx="${compIdx}">
                <div class="evo-slider-phases">
                    <div>I</div><div>II</div><div>III</div><div>IV</div>
                </div>
                ${members.length > 1 ? `<div class="pipeline-span" style="left:${spanLeft}%;width:${spanWidth}%"></div>` : ''}
                ${handles}
                <button class="pipeline-add-btn" onclick="addPipelineMember(${compIdx})" title="${t('pipeline.addMember')}">+</button>
            </div>
            <div class="evo-item-options">
                <label>${t('evo.type')}</label>
                <div class="pill-group">
                    <button class="${!c.type ? 'active' : ''}" onclick="setType(${compIdx},'')">&mdash;</button>
                    <button class="${c.type === 'build' ? 'active' : ''}" onclick="setType(${compIdx},'build')">build</button>
                    <button class="${c.type === 'buy' ? 'active' : ''}" onclick="setType(${compIdx},'buy')">buy</button>
                    <button class="${c.type === 'outsource' ? 'active' : ''}" onclick="setType(${compIdx},'outsource')">outsource</button>
                </div>
                <label style="margin-left:8px">
                    <input type="checkbox" checked onchange="togglePipeline(${compIdx}, this.checked)"> ${t('evo.pipeline')}
                </label>
            </div>
        </div>
    </div>`;
}

function renderGroups() {
    const container = document.getElementById('groups-container');
    const compNames = wizardState.components.filter(c => c.kind !== 'anchor').map(c => c.name);

    container.innerHTML = wizardState.groups.map((g, gi) => `
        <div class="group-card">
            <div class="group-card-header">
                <input type="text" value="${escapeHtml(g.name)}" placeholder="${t('group.namePlaceholder')}"
                    oninput="wizardState.groups[${gi}].name=this.value;onWizardChange()">
                <input type="color" value="${g.color || '#5e9ed6'}"
                    oninput="wizardState.groups[${gi}].color=this.value;onWizardChange()">
                <button onclick="removeGroup(${gi})" title="${t('comp.delete')}"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="width:12px;height:12px"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
            </div>
            <div class="group-members">
                ${compNames.map(n => `
                    <label>
                        <input type="checkbox" ${g.members.includes(n) ? 'checked' : ''}
                            onchange="toggleGroupMember(${gi},'${escapeAttr(n)}',this.checked)">
                        ${escapeHtml(n)}
                    </label>
                `).join('')}
            </div>
        </div>
    `).join('');
    const gc = document.getElementById('groups-count');
    if (gc) gc.textContent = wizardState.groups.length ? ' (' + wizardState.groups.length + ')' : '';
}

function toggleEnrichSection(id) {
    document.getElementById(id).classList.toggle('enrich-collapsed');
}

function renderAnnotations() {
    const container = document.getElementById('annotations-container');
    const compNames = wizardState.components.map(c => c.name);

    container.innerHTML = wizardState.annotations.map((a, ai) => `
        <div class="annotation-row">
            <select onchange="wizardState.annotations[${ai}].kind=this.value;onWizardChange()">
                <option value="note" ${a.kind === 'note' ? 'selected' : ''}>${t('annotation.note')}</option>
                <option value="warning" ${a.kind === 'warning' ? 'selected' : ''}>${t('annotation.warning')}</option>
            </select>
            <input type="text" value="${escapeHtml(a.text)}" placeholder="${t('annotation.textPlaceholder')}"
                oninput="wizardState.annotations[${ai}].text=this.value;onWizardChange()">
            <select onchange="wizardState.annotations[${ai}].target=this.value;onWizardChange()">
                <option value="">${t('annotation.selectComp')}</option>
                ${compNames.map(n => `<option value="${escapeAttr(n)}" ${a.target === n ? 'selected' : ''}>${escapeHtml(n)}</option>`).join('')}
            </select>
            <button onclick="removeAnnotation(${ai})"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="width:12px;height:12px"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
        </div>
    `).join('');
    const ac = document.getElementById('annotations-count');
    if (ac) ac.textContent = wizardState.annotations.length ? ' (' + wizardState.annotations.length + ')' : '';
}

function renderSignals() {
    const container = document.getElementById('signals-container');
    const compNames = wizardState.components.map(c => c.name);

    container.innerHTML = wizardState.signals.map((s, si) => `
        <div class="annotation-row">
            <select onchange="wizardState.signals[${si}].type=this.value;onWizardChange()">
                <option value="accelerating" ${s.type === 'accelerating' ? 'selected' : ''}>accelerating</option>
                <option value="stagnating" ${s.type === 'stagnating' ? 'selected' : ''}>stagnating</option>
                <option value="declining" ${s.type === 'declining' ? 'selected' : ''}>declining</option>
            </select>
            <select onchange="wizardState.signals[${si}].target=this.value;onWizardChange()" style="flex:1">
                <option value="">${t('annotation.selectComp')}</option>
                ${compNames.map(n => `<option value="${escapeAttr(n)}" ${s.target === n ? 'selected' : ''}>${escapeHtml(n)}</option>`).join('')}
            </select>
            <button onclick="removeSignal(${si})"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="width:12px;height:12px"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg></button>
        </div>
    `).join('');
    const sc = document.getElementById('signals-count');
    if (sc) sc.textContent = wizardState.signals.length ? ' (' + wizardState.signals.length + ')' : '';
}

// ================================================================
// 8. Wizard Actions
// ================================================================

function removeComponent(idx) {
    const name = wizardState.components[idx].name;
    wizardState.components.splice(idx, 1);
    wizardState.edges = wizardState.edges.filter(e => e.from !== name && e.to !== name);
    wizardState.groups = wizardState.groups.map(g => ({ ...g, members: g.members.filter(m => m !== name) }));
    wizardState.annotations = wizardState.annotations.filter(a => a.target !== name);
    wizardState.signals = wizardState.signals.filter(s => s.target !== name);
    if (wizardState.focus === name) wizardState.focus = "";
    // Update chain input to reflect removal
    updateChainInputFromState();
    renderCompList();
    onWizardChange();
}

function toggleFocus(idx) {
    const name = wizardState.components[idx].name;
    wizardState.focus = (wizardState.focus === name) ? "" : name;
    renderCompList();
    onWizardChange();
}

function updateChainInputFromState() {
    const lines = wizardState.edges.map(e => e.from + ' -> ' + e.to);
    if (chainEditor) {
        chainEditor.setValue(lines.join('\n'));
        setTimeout(() => chainEditor.refresh(), 10);
    } else document.getElementById('chain-input').value = lines.join('\n');
}

function updateMeta() {
    wizardState.meta.title = document.getElementById('meta-title').value;
    wizardState.meta.author = document.getElementById('meta-author').value;
    wizardState.meta.question = document.getElementById('meta-question').value;
    onWizardChange();
}

function updateEvolution(idx, val) {
    wizardState.components[idx].evolution = parseInt(val);
    const posEl = document.getElementById('evo-pos-' + idx);
    if (posEl) posEl.textContent = evoToRoman(parseInt(val));
    scheduleRender();
}

function updateEvolvedTo(idx, val) {
    wizardState.components[idx].evolvedTo = parseInt(val);
    const el = document.getElementById('evo-to-' + idx);
    if (el) el.textContent = evoToRoman(parseInt(val));
    scheduleRender();
}

function setType(idx, type) {
    wizardState.components[idx].type = type;
    renderCompList();
    renderEvoList();
    onWizardChange();
}

function cycleType(idx) {
    const order = ['', 'build', 'buy', 'outsource'];
    const current = wizardState.components[idx].type || '';
    const next = order[(order.indexOf(current) + 1) % order.length];
    wizardState.components[idx].type = next;
    renderCompList();
    renderEvoList();
    onWizardChange();
}

function toggleEvolving(idx, checked) {
    wizardState.components[idx].evolving = checked;
    if (checked && wizardState.components[idx].evolvedTo <= wizardState.components[idx].evolution) {
        wizardState.components[idx].evolvedTo = Math.min(99, wizardState.components[idx].evolution + 15);
    }
    renderEvoList();
    onWizardChange();
}

function setInertia(idx, val) {
    wizardState.components[idx].inertia = parseInt(val);
    if (parseInt(val) === 0) {
        wizardState.components[idx].inertiaKinds = [];
    }
    renderEvoList();
    onWizardChange();
}

function toggleInertiaKind(idx, kind, checked) {
    const c = wizardState.components[idx];
    if (!c.inertiaKinds) c.inertiaKinds = [];
    if (checked) {
        if (!c.inertiaKinds.includes(kind)) c.inertiaKinds.push(kind);
    } else {
        c.inertiaKinds = c.inertiaKinds.filter(k => k !== kind);
    }
    onWizardChange();
}

// ================================================================
// 8b. Pipeline Actions
// ================================================================

function togglePipeline(compIdx, checked) {
    const c = wizardState.components[compIdx];
    c.isPipeline = checked;
    if (checked && (!c.pipelineMembers || c.pipelineMembers.length === 0)) {
        c.pipelineMembers = [
            { name: c.name + ' (v1)', evolution: Math.max(0, c.evolution - 15) },
            { name: c.name + ' (v2)', evolution: Math.min(99, c.evolution + 15) }
        ];
    }
    if (!checked) c.pipelineMembers = [];
    renderEvoList();
    onWizardChange();
}

function addPipelineMember(compIdx) {
    const c = wizardState.components[compIdx];
    const members = c.pipelineMembers || [];
    let midEvo = c.evolution;
    if (members.length >= 2) {
        const evos = members.map(m => m.evolution);
        midEvo = Math.round((Math.min(...evos) + Math.max(...evos)) / 2);
    }
    c.pipelineMembers.push({ name: t('pipeline.newMember'), evolution: midEvo });
    renderEvoList();
    onWizardChange();
}

function removePipelineMember(compIdx, memberIdx) {
    const c = wizardState.components[compIdx];
    c.pipelineMembers.splice(memberIdx, 1);
    if (c.pipelineMembers.length === 0) c.isPipeline = false;
    renderEvoList();
    onWizardChange();
}

// ================================================================
// 8c. Pipeline Drag Handler (global listeners, grip listeners re-attached)
// ================================================================

let _pipelineDragState = { handle: null, compIdx: -1, memberIdx: -1, trackRect: null };

function initPipelineDrag() {
    // Only re-attach per-grip listeners; document-level listeners are set once below
    document.querySelectorAll('.pipeline-handle-grip').forEach(grip => {
        grip.addEventListener('mousedown', _pipelineDragStart);
        grip.addEventListener('touchstart', _pipelineDragStart, { passive: false });
    });
}

function _pipelineDragStart(e) {
    const grip = e.target.closest('.pipeline-handle-grip');
    if (!grip) return;
    e.preventDefault();
    const s = _pipelineDragState;
    s.handle = grip.closest('.pipeline-handle');
    s.compIdx = parseInt(grip.dataset.compIdx);
    s.memberIdx = parseInt(grip.dataset.memberIdx);
    s.trackRect = grip.closest('.pipeline-track').getBoundingClientRect();
    s.handle.classList.add('dragging');
    document.body.style.cursor = 'grabbing';
    document.body.style.userSelect = 'none';
}

(function() {
    function onMove(e) {
        const s = _pipelineDragState;
        if (!s.handle) return;
        e.preventDefault();
        const clientX = e.touches ? e.touches[0].clientX : e.clientX;
        let pct = (clientX - s.trackRect.left) / s.trackRect.width;
        pct = Math.max(0, Math.min(1, pct));
        const evo = Math.round(pct * 99);
        s.handle.style.left = (pct * 100) + '%';
        const evoLabel = s.handle.querySelector('.pipeline-handle-evo');
        if (evoLabel) evoLabel.textContent = evoToRoman(evo);
        wizardState.components[s.compIdx].pipelineMembers[s.memberIdx].evolution = evo;
        const track = s.handle.closest('.pipeline-track');
        const span = track ? track.querySelector('.pipeline-span') : null;
        if (span) {
            const evos = wizardState.components[s.compIdx].pipelineMembers.map(m => m.evolution);
            const minE = Math.min(...evos), maxE = Math.max(...evos);
            span.style.left = (minE / 99 * 100) + '%';
            span.style.width = ((maxE - minE) / 99 * 100) + '%';
        }
        scheduleRender();
    }
    function onEnd() {
        const s = _pipelineDragState;
        if (!s.handle) return;
        s.handle.classList.remove('dragging');
        s.handle = null;
        document.body.style.cursor = '';
        document.body.style.userSelect = '';
        onWizardChange();
    }
    document.addEventListener('mousemove', onMove);
    document.addEventListener('mouseup', onEnd);
    document.addEventListener('touchmove', onMove, { passive: false });
    document.addEventListener('touchend', onEnd);
})();

// ================================================================
// 8d. Pipeline Label Inline Editing
// ================================================================

function initPipelineLabelClicks() {
    document.querySelectorAll('.pipeline-handle-label').forEach(label => {
        label.addEventListener('click', function(e) {
            e.stopPropagation();
            const compIdx = parseInt(this.dataset.compIdx);
            const memberIdx = parseInt(this.dataset.memberIdx);
            const currentName = wizardState.components[compIdx].pipelineMembers[memberIdx].name;
            const input = document.createElement('input');
            input.type = 'text';
            input.value = currentName;
            input.style.cssText = 'font-size:10px;width:80px;padding:1px 4px;border:1px solid #00D2DD;border-radius:3px;outline:none;background:#fff;color:#0E2356;font-family:inherit;';
            this.replaceWith(input);
            input.focus();
            input.select();
            function commit() {
                const val = input.value.trim().replace(/:/g, '');
                if (val) wizardState.components[compIdx].pipelineMembers[memberIdx].name = val;
                renderEvoList();
                onWizardChange();
            }
            input.addEventListener('blur', commit);
            input.addEventListener('keydown', function(ev) {
                if (ev.key === 'Enter') { ev.preventDefault(); commit(); }
                if (ev.key === 'Escape') { ev.preventDefault(); renderEvoList(); }
            });
        });
    });
}

function addGroup() {
    wizardState.groups.push({ name: t('group.defaultName'), members: [], color: '#5e9ed6' });
    renderGroups();
}

function removeGroup(idx) {
    wizardState.groups.splice(idx, 1);
    renderGroups();
    onWizardChange();
}

function toggleGroupMember(gi, name, checked) {
    const g = wizardState.groups[gi];
    if (checked && !g.members.includes(name)) {
        g.members.push(name);
    } else if (!checked) {
        g.members = g.members.filter(m => m !== name);
    }
    onWizardChange();
}

function addAnnotation() {
    wizardState.annotations.push({ kind: 'note', text: '', target: '' });
    renderAnnotations();
}

function removeAnnotation(idx) {
    wizardState.annotations.splice(idx, 1);
    renderAnnotations();
    onWizardChange();
}

function addSignal() {
    wizardState.signals.push({ type: 'accelerating', target: '' });
    renderSignals();
}

function removeSignal(idx) {
    wizardState.signals.splice(idx, 1);
    renderSignals();
    onWizardChange();
}

// ================================================================
// 9. Steps Navigation
// ================================================================

function goToStep(step) {
    currentStep = step;
    document.querySelectorAll('.wizard-step').forEach((el, i) => {
        el.classList.toggle('active', i === step);
    });
    document.querySelectorAll('.step-indicator').forEach((el, i) => {
        el.classList.toggle('active', i === step);
        if (i < step && wizardState.components.length > 0) {
            el.classList.add('completed');
        }
    });

    document.getElementById('btn-prev').style.visibility = step === 0 ? 'hidden' : 'visible';
    const nextBtn = document.getElementById('btn-next');
    nextBtn.innerHTML = step === 4 ? t('nav.finish') : t('nav.next') + ' \u25B6';

    if (step === 2) renderEvoList();
    if (step === 3) {
        renderGroups();
        renderAnnotations();
        renderSignals();
    }
    if (step === 4) {
        initAnimationEngine();
    } else {
        teardownAnimationEngine();
    }
}

function nextStep() {
    if (currentStep < 4) goToStep(currentStep + 1);
}

function prevStep() {
    if (currentStep > 0) goToStep(currentStep - 1);
}

// ================================================================
// 9b. Animation Engine (Step 4 - Navigation)
// ================================================================

let animSteps = [];
let animCurrentStep = -1;
let animVisibleIDs = new Set();
let animActive = false;
let animKeyHandler = null;

function initAnimationEngine() {
    const svgEl = document.querySelector('#output-pane svg');
    if (!svgEl) return;
    animActive = true;
    buildAnimSteps(svgEl);
    animResetState(svgEl);
    updateAnimCounter();

    if (animKeyHandler) document.removeEventListener('keydown', animKeyHandler);
    animKeyHandler = function(e) {
        if (!animActive) return;
        if (['INPUT', 'SELECT', 'TEXTAREA'].includes(document.activeElement.tagName)) return;
        switch (e.key) {
            case 'ArrowRight': e.preventDefault(); animNext(); break;
            case 'ArrowLeft':  e.preventDefault(); animPrev(); break;
            case 'Home':       e.preventDefault(); animReset(); break;
            case 'End':        e.preventDefault(); animShowAll(); break;
        }
    };
    document.addEventListener('keydown', animKeyHandler);
}

function teardownAnimationEngine() {
    if (!animActive) return;
    var svgEl = document.querySelector('#output-pane svg');
    if (svgEl) animShowAllImmediate(svgEl);
    animActive = false;
    animSteps = [];
    animCurrentStep = -1;
    animVisibleIDs = new Set();
    if (animKeyHandler) {
        document.removeEventListener('keydown', animKeyHandler);
        animKeyHandler = null;
    }
}

function buildAnimSteps(svgEl) {
    animSteps = [];
    var mode = document.getElementById('anim-mode').value;
    var attrKey = mode === 'yrank' ? 'data-y-rank' : 'data-depth';

    var components = svgEl.querySelectorAll('[data-depth]');
    var edges = svgEl.querySelectorAll('[data-child-id]');
    var signals = svgEl.querySelectorAll("[data-type='signal']");
    var warnings = svgEl.querySelectorAll("[data-type='warning']");
    var gameplays = svgEl.querySelectorAll("[data-type='gameplay']");
    var groups = svgEl.querySelectorAll("[data-type='group']");

    components.forEach(function(el) { el.style.opacity = '0'; el.style.pointerEvents = 'none'; });
    edges.forEach(function(el) { el.style.opacity = '0'; el.style.pointerEvents = 'none'; });
    signals.forEach(function(el) { el.style.opacity = '0'; });
    warnings.forEach(function(el) { el.style.opacity = '0'; });
    gameplays.forEach(function(el) { el.style.opacity = '0'; });
    groups.forEach(function(el) { el.style.opacity = '0'; el.style.pointerEvents = 'none'; });

    var levelMap = new Map();
    components.forEach(function(el) {
        var type = el.getAttribute('data-type');
        if (type === 'group') return;
        var level = parseInt(el.getAttribute(attrKey) || '0');
        if (!levelMap.has(level)) levelMap.set(level, []);
        levelMap.get(level).push(el);
    });

    var sortedLevels = Array.from(levelMap.keys()).filter(function(l) { return l >= 0; }).sort(function(a, b) { return a - b; });

    var depthOf = new Map();
    components.forEach(function(el) {
        var level = parseInt(el.getAttribute(attrKey) || '0');
        depthOf.set(el.id, level);
    });

    var edgesByLevel = new Map();
    edges.forEach(function(edge) {
        var childId = edge.getAttribute('data-child-id');
        var edgeId = edge.id || '';
        var match = edgeId.match(/^edge_(\d+)_(\d+)$/);
        var sourceId = match ? 'element_' + match[1] : null;
        var sourceDepth = sourceId && depthOf.has(sourceId) ? depthOf.get(sourceId) : 0;
        var targetDepth = depthOf.has(childId) ? depthOf.get(childId) : 0;
        var showAt = Math.max(sourceDepth, targetDepth);
        if (!edgesByLevel.has(showAt)) edgesByLevel.set(showAt, []);
        edgesByLevel.get(showAt).push(edge);
    });

    for (var li = 0; li < sortedLevels.length; li++) {
        var level = sortedLevels[li];
        var nodesAtLevel = levelMap.get(level);
        var edgesAtLevel = edgesByLevel.get(level) || [];
        animSteps.push({ nodes: nodesAtLevel, edges: edgesAtLevel, groups: groups, extras: [] });
    }

    var allExtras = Array.from(signals).concat(Array.from(warnings));
    if (allExtras.length > 0) {
        animSteps.push({ nodes: [], edges: [], groups: groups, extras: allExtras });
    }

    var allGameplays = Array.from(gameplays);
    if (allGameplays.length > 0) {
        animSteps.push({ nodes: [], edges: [], groups: groups, extras: allGameplays });
    }
}

function animResetState(svgEl) {
    animCurrentStep = -1;
    animVisibleIDs = new Set();
    if (!svgEl) svgEl = document.querySelector('#output-pane svg');
    if (!svgEl) return;

    svgEl.querySelectorAll('[data-depth]').forEach(function(el) {
        el.style.opacity = '0'; el.style.pointerEvents = 'none'; el.style.transition = 'none';
        var inner = el.querySelector('g[transform]');
        if (inner && inner._finalTransform) {
            inner.setAttribute('transform', inner._finalTransform);
            inner.style.transition = 'none';
        }
    });
    svgEl.querySelectorAll('[data-child-id]').forEach(function(el) {
        el.style.opacity = '0'; el.style.pointerEvents = 'none'; el.style.transition = 'none';
        animResetEdgeDrawing(el);
    });
    svgEl.querySelectorAll("[data-type='signal']").forEach(function(el) { el.style.opacity = '0'; el.style.transition = 'none'; });
    svgEl.querySelectorAll("[data-type='warning']").forEach(function(el) { el.style.opacity = '0'; el.style.transition = 'none'; });
    svgEl.querySelectorAll("[data-type='gameplay']").forEach(function(el) { el.style.opacity = '0'; el.style.transition = 'none'; });
    svgEl.querySelectorAll("[data-type='group']").forEach(function(el) { el.style.opacity = '0'; el.style.pointerEvents = 'none'; el.style.transition = 'none'; });

    updateAnimCounter();
}

function animApplyStep(stepIndex, withAnimation) {
    var step = animSteps[stepIndex];
    if (!step) return;
    var duration = withAnimation ? 0.5 : 0;

    for (var ni = 0; ni < step.nodes.length; ni++) {
        var node = step.nodes[ni];
        if (withAnimation) {
            animDeployFromParent(node, duration);
        } else {
            node.style.transition = 'none';
            node.style.opacity = '1';
            node.style.pointerEvents = 'auto';
        }
        animVisibleIDs.add(node.id);
    }

    for (var ei = 0; ei < step.edges.length; ei++) {
        var edge = step.edges[ei];
        if (withAnimation) {
            animDrawEdge(edge, duration);
        } else {
            edge.style.transition = 'none';
            edge.style.opacity = '1';
            edge.style.pointerEvents = 'auto';
        }
    }

    for (var xi = 0; xi < step.extras.length; xi++) {
        var extra = step.extras[xi];
        extra.style.transition = withAnimation ? 'opacity ' + duration + 's ease' : 'none';
        extra.style.opacity = '1';
    }

    if (step.groups) {
        step.groups.forEach(function(g) {
            var membersAttr = g.getAttribute('data-members');
            if (!membersAttr) return;
            var members = membersAttr.split(',');
            var allVisible = members.every(function(id) { return animVisibleIDs.has(id.trim()); });
            if (allVisible && g.style.opacity !== '1') {
                g.style.transition = withAnimation ? 'opacity 0.4s ease' : 'none';
                g.style.opacity = '1';
                g.style.pointerEvents = 'auto';
            }
        });
    }
}

function animGoToStep(target) {
    if (target < 0 || target >= animSteps.length) return;
    var svgEl = document.querySelector('#output-pane svg');
    if (!svgEl) return;

    if (target <= animCurrentStep) {
        animResetState(svgEl);
        void svgEl.getBoundingClientRect();
        for (var i = 0; i <= target; i++) {
            animApplyStep(i, i === target);
        }
    } else {
        for (var i = animCurrentStep + 1; i <= target; i++) {
            animApplyStep(i, i === target);
        }
    }
    animCurrentStep = target;
    updateAnimCounter();
}

function animGetTranslateCoords(el) {
    var g = el.querySelector('g[transform]');
    if (!g) return null;
    var attr = g.getAttribute('transform');
    var match = attr && attr.match(/translate\(\s*([^,\s]+)\s*,\s*([^)\s]+)\s*\)/);
    if (!match) return null;
    return { x: parseFloat(match[1]), y: parseFloat(match[2]), g: g, attr: attr };
}

function animDeployFromParent(node, duration) {
    var parentId = node.getAttribute('data-parent-id');
    var coords = animGetTranslateCoords(node);
    var svgEl = document.querySelector('#output-pane svg');

    if (parentId && coords && svgEl) {
        var parentEl = svgEl.getElementById(parentId);
        var parentCoords = parentEl && animGetTranslateCoords(parentEl);
        if (parentCoords) {
            coords.g._finalTransform = coords.attr;
            coords.g.setAttribute('transform', 'translate(' + parentCoords.x + ',' + parentCoords.y + ')');
            coords.g.style.transition = 'none';
            node.style.transition = 'none';
            node.style.opacity = '1';
            node.style.pointerEvents = 'auto';
            void coords.g.getBoundingClientRect();
            coords.g.style.transition = 'transform ' + duration + 's ease';
            coords.g.setAttribute('transform', coords.attr);
            return;
        }
    }
    node.style.transition = 'opacity ' + duration + 's ease';
    node.style.opacity = '1';
    node.style.pointerEvents = 'auto';
}

function animIsEvolutionEdge(line) {
    var cls = line.getAttribute('class') || '';
    return cls.indexOf('evolutionEdge') >= 0;
}

function animDrawEdge(edgeEl, duration) {
    var line = edgeEl.querySelector('line, path');
    if (line && typeof line.getTotalLength === 'function') {
        if (animIsEvolutionEdge(line)) {
            edgeEl.style.transition = 'opacity ' + duration + 's ease';
            edgeEl.style.opacity = '1';
            edgeEl.style.pointerEvents = 'auto';
            return;
        }
        var length = line.getTotalLength();
        line.style.strokeDasharray = length;
        line.style.strokeDashoffset = length;
        line.style.transition = 'none';
        edgeEl.style.opacity = '1';
        edgeEl.style.pointerEvents = 'auto';
        void line.getBoundingClientRect();
        line.style.transition = 'stroke-dashoffset ' + duration + 's ease';
        line.style.strokeDashoffset = '0';
    } else {
        edgeEl.style.transition = 'opacity ' + duration + 's ease';
        edgeEl.style.opacity = '1';
        edgeEl.style.pointerEvents = 'auto';
    }
}

function animResetEdgeDrawing(edgeEl) {
    var line = edgeEl.querySelector('line, path');
    if (line && !animIsEvolutionEdge(line)) {
        line.style.strokeDasharray = '';
        line.style.strokeDashoffset = '';
        line.style.transition = 'none';
    }
}

function animShowAllImmediate(svgEl) {
    if (!svgEl) svgEl = document.querySelector('#output-pane svg');
    if (!svgEl) return;
    svgEl.querySelectorAll('[data-depth]').forEach(function(el) {
        el.style.opacity = '1'; el.style.pointerEvents = 'auto'; el.style.transition = 'none';
    });
    svgEl.querySelectorAll('[data-child-id]').forEach(function(el) {
        el.style.opacity = '1'; el.style.pointerEvents = 'auto'; el.style.transition = 'none';
        animResetEdgeDrawing(el);
    });
    svgEl.querySelectorAll("[data-type='signal']").forEach(function(el) { el.style.opacity = '1'; el.style.transition = 'none'; });
    svgEl.querySelectorAll("[data-type='warning']").forEach(function(el) { el.style.opacity = '1'; el.style.transition = 'none'; });
    svgEl.querySelectorAll("[data-type='gameplay']").forEach(function(el) { el.style.opacity = '1'; el.style.transition = 'none'; });
    svgEl.querySelectorAll("[data-type='group']").forEach(function(el) { el.style.opacity = '1'; el.style.pointerEvents = 'auto'; el.style.transition = 'none'; });
}

function animReset() {
    if (!animActive) return;
    animResetState();
}

function animPrev() {
    if (!animActive) return;
    animGoToStep(animCurrentStep - 1);
}

function animNext() {
    if (!animActive) return;
    animGoToStep(animCurrentStep + 1);
}

function animShowAll() {
    if (!animActive) return;
    if (animSteps.length === 0) return;
    animResetState();
    for (var i = 0; i < animSteps.length; i++) {
        animApplyStep(i, false);
    }
    animCurrentStep = animSteps.length - 1;
    updateAnimCounter();
}

function onAnimModeChange() {
    if (!animActive) return;
    var svgEl = document.querySelector('#output-pane svg');
    if (!svgEl) return;
    buildAnimSteps(svgEl);
    animResetState(svgEl);
}

function updateAnimCounter() {
    var el = document.getElementById('anim-counter');
    if (!el) return;
    if (animSteps.length === 0) {
        el.textContent = '-';
    } else {
        el.textContent = t('nav.stepCounter', { current: animCurrentStep + 1, total: animSteps.length });
    }
    var prevBtn = document.getElementById('anim-prev');
    var nextBtn = document.getElementById('anim-next');
    if (prevBtn) prevBtn.disabled = animCurrentStep < 0;
    if (nextBtn) nextBtn.disabled = animCurrentStep >= animSteps.length - 1;
}

// ================================================================
// 10. Mode Switching
// ================================================================

function setMode(mode) {
    // When switching from editor to guided, sync editor text back to wizardState.
    if (mode === 'guided' && currentMode === 'editor' && editor) {
        const text = editor.getValue();
        if (text.trim() && typeof parseWTG2ToState === 'function') {
            const result = parseWTG2ToState(text);
            if (result.startsWith('error:')) {
                // Parse failed — stay in editor mode and show the error.
                document.getElementById('status').textContent = result;
                return;
            }
            try {
                const parsed = JSON.parse(result);
                wizardState.meta = parsed.meta || { title: "", author: "", question: "" };
                wizardState.stages = parsed.stages || ["", "", "", ""];
                wizardState.components = parsed.components || [];
                wizardState.edges = parsed.edges || [];
                wizardState.groups = parsed.groups || [];
                wizardState.annotations = parsed.annotations || [];
                wizardState.signals = parsed.signals || [];
                wizardState.legend = parsed.legend || false;
                wizardState.focus = parsed.focus || "";
                syncWizardUI();
            } catch (e) {
                document.getElementById('status').textContent = 'error: ' + e.message;
                return;
            }
        }
    }

    currentMode = mode;
    document.getElementById('mode-guided').classList.toggle('active', mode === 'guided');
    document.getElementById('mode-editor').classList.toggle('active', mode === 'editor');
    var slider = document.querySelector('.mode-toggle .toggle-slider');
    if (slider) slider.classList.toggle('editor', mode === 'editor');
    document.getElementById('steps-bar').classList.toggle('visible', mode === 'guided');
    document.getElementById('wizard-pane').style.display = mode === 'guided' ? 'flex' : 'none';
    document.getElementById('editor-pane').style.display = mode === 'editor' ? 'flex' : 'none';
    document.getElementById('undo-redo-group').style.display = mode === 'guided' ? '' : 'none';
    undoStack = []; redoStack = []; updateUndoButtons();

    if (mode === 'editor') {
        if (!editor) initEditor();
        // Load generated WTG2 into editor
        const wtg2 = generateWTG2(wizardState);
        editor.setValue(wtg2);
        setTimeout(() => editor.refresh(), 10);
    }

    render();
}

// Rebuild wizard UI elements from the current wizardState.
function syncWizardUI() {
    document.getElementById('meta-title').value = wizardState.meta.title;
    document.getElementById('meta-author').value = wizardState.meta.author;
    document.getElementById('meta-question').value = wizardState.meta.question;
    for (let i = 0; i < 4; i++) {
        document.getElementById('stage-' + i).value = wizardState.stages[i];
    }
    document.getElementById('legend-toggle').checked = wizardState.legend;

    // Reconstruct chain-input from edges
    const lines = wizardState.edges.map(e => e.from + ' -> ' + e.to);
    if (chainEditor) {
        chainEditor.setValue(lines.join('\n'));
        setTimeout(() => chainEditor.refresh(), 10);
    } else document.getElementById('chain-input').value = lines.join('\n');

    renderCompList();
    goToStep(currentStep);
    saveState();
}

// ================================================================
// 11. Template Loading
// ================================================================

function loadTemplate(key) {
    if (!key || !TEMPLATES[key]) return;
    const tpl = TEMPLATES[key];

    wizardState.meta = { ...tpl.meta };
    wizardState.stages = tpl.stages ? [...tpl.stages] : ["", "", "", ""];
    wizardState.components = tpl.components.map(c => ({ ...c }));
    wizardState.edges = tpl.edges.map(e => ({ ...e }));
    wizardState.groups = tpl.groups.map(g => ({ ...g, members: [...g.members] }));
    wizardState.annotations = tpl.annotations.map(a => ({ ...a }));
    wizardState.signals = tpl.signals.map(s => ({ ...s }));
    wizardState.legend = tpl.legend;
    wizardState.focus = "";

    // Update UI
    document.getElementById('meta-title').value = wizardState.meta.title;
    document.getElementById('meta-author').value = wizardState.meta.author;
    document.getElementById('meta-question').value = wizardState.meta.question;
    if (chainEditor) {
        chainEditor.setValue(tpl.chains);
        setTimeout(() => chainEditor.refresh(), 10);
    } else document.getElementById('chain-input').value = tpl.chains;
    document.getElementById('legend-toggle').checked = wizardState.legend;
    for (let i = 0; i < 4; i++) {
        document.getElementById('stage-' + i).value = wizardState.stages[i];
    }

    renderCompList();
    goToStep(1);
    onWizardChange();
}

// ================================================================
// 12. Section Toggle
// ================================================================

function toggleSection(header) {
    header.parentElement.classList.toggle('collapsed');
}

// ================================================================
// 13. Render Pipeline
// ================================================================

function onWizardChange() {
    saveState();
    scheduleRender();
    scheduleUndoPush();
}

function scheduleRender() {
    clearTimeout(renderTimer);
    renderTimer = setTimeout(render, 400);
}

function execInlineScripts(container) {
    container.querySelectorAll('script').forEach(function(old) {
        var s = document.createElement('script');
        var code = old.textContent;
        code = code.replace(/\bconst\s+/g, 'var ');
        code = code.replace(/\blet\s+/g, 'var ');
        s.textContent = code;
        old.parentNode.replaceChild(s, old);
    });
}

function render() {
    let text;
    if (currentMode === 'guided') {
        text = generateWTG2(wizardState);
    } else if (editor) {
        text = editor.getValue();
    } else {
        return;
    }

    if (!text.trim()) {
        document.getElementById('output-pane').innerHTML = '<div style="color:#888;padding:40px;text-align:center">' + escapeHtml(t('empty.message')) + '</div>';
        return;
    }

    if (typeof generateSVG !== 'function') {
        document.getElementById('status').textContent = t('status.wasmNotLoaded');
        return;
    }

    const isStatic = false;
    const [baseW, baseH] = getBaseDimensions();
    const result = generateSVG(text, isStatic, getResolution(), baseW, baseH);
    const output = document.getElementById('output-pane');

    if (result.startsWith('error:')) {
        output.innerHTML = '<div class="error">' + escapeHtml(result) + '</div>';
    } else {
        output.innerHTML = result;
        execInlineScripts(output);
    }
    document.getElementById('status').textContent = t('status.ok');
    applyCanvasTransform();
    updateDetached();
    if (animActive && currentStep === 4) {
        initAnimationEngine();
    }
}

function getResolution() {
    return parseInt(document.getElementById('resolution').value, 10);
}

function getBaseDimensions() {
    const parts = document.getElementById('ratio').value.split(',');
    return [parseInt(parts[0], 10), parseInt(parts[1], 10)];
}

// ================================================================
// 13b. Canvas Zoom/Pan
// ================================================================
let canvasScale = 1;
let canvasPanX = 0;
let canvasPanY = 0;
let _panActive = false;
let _panStartX = 0;
let _panStartY = 0;
let _panStartPanX = 0;
let _panStartPanY = 0;

function applyCanvasTransform() {
    const pane = document.getElementById('output-pane');
    if (!pane) return;
    if (canvasScale === 1 && canvasPanX === 0 && canvasPanY === 0) {
        pane.style.transform = '';
    } else {
        pane.style.transform = `translate(${canvasPanX}px, ${canvasPanY}px) scale(${canvasScale})`;
    }
    const label = document.getElementById('canvas-zoom-label');
    if (label) label.textContent = Math.round(canvasScale * 100) + '%';
}

function canvasZoomIn() { canvasZoomTo(canvasScale + 0.15); }
function canvasZoomOut() { canvasZoomTo(canvasScale - 0.15); }
function canvasZoomTo(s) {
    canvasScale = Math.max(0.25, Math.min(3, s));
    applyCanvasTransform();
}
function canvasReset() {
    canvasScale = 1; canvasPanX = 0; canvasPanY = 0;
    applyCanvasTransform();
}

(function initCanvasControls() {
    document.addEventListener('DOMContentLoaded', function() {
        const vp = document.getElementById('canvas-viewport');
        if (!vp) return;

        vp.addEventListener('wheel', function(e) {
            if (e.ctrlKey || e.metaKey) {
                e.preventDefault();
                const rect = vp.getBoundingClientRect();
                const cx = e.clientX - rect.left;
                const cy = e.clientY - rect.top;
                const oldScale = canvasScale;
                const delta = e.deltaY > 0 ? -0.1 : 0.1;
                canvasScale = Math.max(0.25, Math.min(3, canvasScale + delta));
                const ratio = canvasScale / oldScale;
                canvasPanX = cx - ratio * (cx - canvasPanX);
                canvasPanY = cy - ratio * (cy - canvasPanY);
            } else {
                e.preventDefault();
                canvasPanX -= e.deltaX;
                canvasPanY -= e.deltaY;
            }
            applyCanvasTransform();
        }, { passive: false });

        vp.addEventListener('mousedown', function(e) {
            if (e.target.closest('svg') || e.target.closest('#canvas-controls')) return;
            _panActive = true;
            _panStartX = e.clientX;
            _panStartY = e.clientY;
            _panStartPanX = canvasPanX;
            _panStartPanY = canvasPanY;
            document.getElementById('output-pane').classList.add('panning');
            e.preventDefault();
        });

        document.addEventListener('mousemove', function(e) {
            if (!_panActive) return;
            canvasPanX = _panStartPanX + (e.clientX - _panStartX);
            canvasPanY = _panStartPanY + (e.clientY - _panStartY);
            applyCanvasTransform();
        });

        document.addEventListener('mouseup', function() {
            if (!_panActive) return;
            _panActive = false;
            document.getElementById('output-pane').classList.remove('panning');
        });
    });
})();

// ================================================================
// 13c. Autocompletion
// ================================================================

function extractComponentNames(text) {
    const names = new Set();
    const lines = text.split('\n');
    const metaRe = /^(title|date|author|scope|question|stages|legend|focus)\b/;
    const defRe = /^(?:component\s+|anchor\s+|submap\s+)?(.+?)\s+:\s+[IVX]+\.\d+/;
    const anchorRe = /^anchor\s+(.+?)\s*$/;
    const pipelineRe = /^pipeline\s+(.+?)\s*\{/;

    for (const raw of lines) {
        const line = raw.trim();
        if (!line || line.startsWith('//') || line.startsWith('/*')) continue;
        if (metaRe.test(line)) continue;
        if (/^(note|warning|signal)\b/.test(line)) continue;

        let m;
        if ((m = line.match(pipelineRe))) { names.add(m[1].trim()); continue; }
        if ((m = line.match(defRe))) { names.add(m[1].trim()); continue; }
        if ((m = line.match(anchorRe))) { names.add(m[1].trim()); continue; }

        if (line.includes('->')) {
            var edgeParts = line.split(/\s*(?:-\[.*?\])?->\s*/);
            for (var ep of edgeParts) {
                var ename = ep.trim();
                if (ename && !/^(group|note|warning|signal)\b/.test(ename)) {
                    names.add(ename);
                }
            }
            continue;
        }

        if (/^group\s+/.test(line)) continue;
    }
    return Array.from(names).sort();
}

function isInsideGroupBlock(cm, lineNo) {
    for (var i = lineNo - 1; i >= 0; i--) {
        var l = cm.getLine(i).trim();
        if (l === '}') return false;
        if (/^group\s+.+\{/.test(l)) return true;
        if (/^pipeline\s+.+\{/.test(l)) return false;
    }
    return false;
}

function wtg2Hint(cm) {
    var cursor = cm.getCursor();
    var line = cm.getLine(cursor.line);
    var textUpTo = line.substring(0, cursor.ch);

    var partial = '';
    var startCh = 0;

    var arrowMatch = textUpTo.match(/(?:->|]\->)\s*(.*)$/);
    if (arrowMatch) {
        partial = arrowMatch[1];
        startCh = cursor.ch - partial.length;
    } else if (isInsideGroupBlock(cm, cursor.line)) {
        partial = textUpTo.trimStart();
        startCh = textUpTo.length - partial.length;
    } else {
        var keywords = /^(title|date|author|scope|question|stages|legend|focus|anchor|component|submap|pipeline|group|note|warning|signal)\b/;
        if (!keywords.test(textUpTo.trim())) {
            if (textUpTo.indexOf(' : ') >= 0) return null;
            partial = textUpTo.trimStart();
            startCh = textUpTo.length - partial.length;
        } else {
            return null;
        }
    }

    if (!partial) return null;

    var allNames = extractComponentNames(cm.getValue());
    var lowerPartial = partial.toLowerCase();
    var matches = allNames.filter(function(n) {
        return n.toLowerCase().startsWith(lowerPartial) && n.toLowerCase() !== lowerPartial;
    });

    if (matches.length === 0) return null;

    return {
        list: matches,
        from: CodeMirror.Pos(cursor.line, startCh),
        to: cursor
    };
}

// --- Guided mode: chain-input as CodeMirror editor ---

var chainEditor = null;

function chainHint(cm) {
    var cursor = cm.getCursor();
    var line = cm.getLine(cursor.line);
    var textUpTo = line.substring(0, cursor.ch);

    var partial = '';
    var startCh = 0;

    var arrowMatch = textUpTo.match(/(?:->)\s*(.*)$/);
    if (arrowMatch) {
        partial = arrowMatch[1];
        startCh = cursor.ch - partial.length;
    } else {
        partial = textUpTo.replace(/^\s+/, '');
        startCh = textUpTo.length - partial.length;
    }

    if (!partial) return null;

    var names = (typeof wizardState !== 'undefined' && wizardState.components)
        ? wizardState.components.map(function(c) { return c.name; })
        : [];
    var lowerPartial = partial.toLowerCase();
    var matches = names.filter(function(n) {
        return n.toLowerCase().startsWith(lowerPartial) && n.toLowerCase() !== lowerPartial;
    });

    if (matches.length === 0) return null;

    return {
        list: matches,
        from: CodeMirror.Pos(cursor.line, startCh),
        to: cursor
    };
}

function initChainEditor() {
    var textarea = document.getElementById('chain-input');
    if (!textarea) return;

    chainEditor = CodeMirror.fromTextArea(textarea, {
        mode: null,
        theme: 'default',
        lineNumbers: false,
        tabSize: 2,
        lineWrapping: true,
        placeholder: textarea.getAttribute('placeholder'),
        hintOptions: { completeSingle: false },
        extraKeys: {
            'Ctrl-Space': function(cm) { cm.showHint({ hint: chainHint }); }
        }
    });
    chainEditor.on('change', function(cm, changeObj) {
        if (changeObj.origin === 'setValue') return;
        parseChains();
    });
    chainEditor.on('inputRead', function(cm, changeObj) {
        if (changeObj.origin !== '+input') return;
        var cursor = cm.getCursor();
        var line = cm.getLine(cursor.line);
        var textUpTo = line.substring(0, cursor.ch);

        var partial = null;
        var arrowMatch = textUpTo.match(/(?:->)\s*(.+)$/);
        if (arrowMatch) partial = arrowMatch[1];
        else partial = textUpTo.replace(/^\s+/, '');

        if (partial && partial.trim().length >= 1) {
            setTimeout(function() {
                cm.showHint({ hint: chainHint, completeSingle: false });
            }, 0);
        }
    });
}

// ================================================================
// 14. Editor Init
// ================================================================

function initEditor() {
    editor = CodeMirror.fromTextArea(document.getElementById('source'), {
        mode: 'wtg2',
        theme: 'default',
        lineNumbers: true,
        tabSize: 2,
        hintOptions: { completeSingle: false },
        extraKeys: {
            'Ctrl-Space': function(cm) { cm.showHint({ hint: wtg2Hint }); }
        }
    });
    editor.on('change', function() {
        localStorage.setItem('wtg2-source', editor.getValue());
        scheduleRender();
    });
    editor.on('inputRead', function(cm, changeObj) {
        if (changeObj.origin !== '+input') return;
        var cursor = cm.getCursor();
        var line = cm.getLine(cursor.line);
        var textUpTo = line.substring(0, cursor.ch);

        var partial = null;
        var arrowMatch = textUpTo.match(/(?:->|]\->)\s*(.+)$/);
        if (arrowMatch) {
            partial = arrowMatch[1];
        } else if (isInsideGroupBlock(cm, cursor.line)) {
            partial = textUpTo.trimStart();
        } else {
            var kw = /^(title|date|author|scope|question|stages|legend|focus|anchor|component|submap|pipeline|group|note|warning|signal)\b/;
            if (!kw.test(textUpTo.trim()) && textUpTo.indexOf(' : ') < 0) {
                partial = textUpTo.trimStart();
            }
        }

        if (partial && partial.trim().length >= 1) {
            setTimeout(function() {
                cm.showHint({ hint: wtg2Hint, completeSingle: false });
            }, 0);
        }
    });
}

// ================================================================
// 15. Export / Import / Share Functions
// ================================================================

function getCurrentWTG2() {
    if (currentMode === 'guided') return generateWTG2(wizardState);
    if (editor) return editor.getValue();
    return '';
}

function getExportSVG() {
    const text = getCurrentWTG2();
    if (!text.trim() || typeof generateSVG !== 'function') return null;
    const [baseW, baseH] = getBaseDimensions();
    const result = generateSVG(text, true, getResolution(), baseW, baseH);
    if (result.startsWith('error:')) return null;
    return result;
}

function getMapTitle() {
    const title = wizardState.meta.title || 'wardley-map';
    return title.replace(/[^a-zA-Z0-9_\-\s]/g, '').replace(/\s+/g, '_');
}

function downloadSVG() {
    const svgStr = getExportSVG();
    if (!svgStr) return;
    const blob = new Blob([svgStr], { type: 'image/svg+xml' });
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = getMapTitle() + '.svg';
    a.click();
    URL.revokeObjectURL(a.href);
}

function downloadPNG() {
    const svgStr = getExportSVG();
    if (!svgStr) return;
    const resPct = getResolution();
    const [baseW, baseH] = getBaseDimensions();
    const w = Math.round(baseW * resPct / 100);
    const h = Math.round(baseH * resPct / 100);

    const canvas = document.createElement('canvas');
    canvas.width = w;
    canvas.height = h;
    const ctx = canvas.getContext('2d');
    ctx.fillStyle = '#fff';
    ctx.fillRect(0, 0, w, h);

    const img = new Image();
    const blob = new Blob([svgStr], { type: 'image/svg+xml' });
    const url = URL.createObjectURL(blob);
    img.onload = () => {
        ctx.drawImage(img, 0, 0, w, h);
        URL.revokeObjectURL(url);
        canvas.toBlob(pngBlob => {
            const a = document.createElement('a');
            a.href = URL.createObjectURL(pngBlob);
            a.download = getMapTitle() + '_' + resPct + 'pct.png';
            a.click();
            URL.revokeObjectURL(a.href);
        }, 'image/png');
    };
    img.src = url;
}

// ================================================================
// 15a. Animated Export (GIF / APNG)
// ================================================================

var _exportCancelled = false;

function getAnimatedSVG() {
    var text = getCurrentWTG2();
    if (!text.trim() || typeof generateSVG !== 'function') return null;
    var dims = getBaseDimensions();
    var result = generateSVG(text, false, getResolution(), dims[0], dims[1]);
    if (result.startsWith('error:')) return null;
    return result;
}

function showExportProgress(onCancel) {
    var overlay = document.createElement('div');
    overlay.className = 'export-progress-overlay';
    overlay.id = 'export-progress-overlay';
    overlay.innerHTML =
        '<div class="export-progress-dialog">' +
        '<div class="export-progress-label" id="export-progress-label">' + t('export.progress', { percent: 0 }) + '</div>' +
        '<div class="export-progress-bar-track"><div class="export-progress-bar-fill" id="export-progress-bar-fill"></div></div>' +
        '<button class="btn" id="export-cancel-btn">' + t('export.cancel') + '</button>' +
        '</div>';
    document.body.appendChild(overlay);
    document.getElementById('export-cancel-btn').addEventListener('click', function() {
        _exportCancelled = true;
        if (onCancel) onCancel();
        hideExportProgress();
    });
}

function updateExportProgress(percent, label) {
    var el = document.getElementById('export-progress-label');
    var bar = document.getElementById('export-progress-bar-fill');
    if (el) el.textContent = label || t('export.progress', { percent: Math.round(percent) });
    if (bar) bar.style.width = percent + '%';
}

function hideExportProgress() {
    var overlay = document.getElementById('export-progress-overlay');
    if (overlay) overlay.remove();
}

function captureAnimationFrames(svgString, w, h, onProgress, fps) {
    return new Promise(function(resolve, reject) {
        var DURATION_MS = 3000;
        var FPS = fps || 15;
        var FRAME_COUNT = Math.round(DURATION_MS * FPS / 1000);
        var FRAME_INTERVAL = DURATION_MS / FRAME_COUNT;

        var container = document.createElement('div');
        container.style.cssText = 'position:fixed;left:-9999px;top:-9999px;width:' + w + 'px;height:' + h + 'px;overflow:hidden;';
        container.innerHTML = svgString;
        document.body.appendChild(container);

        var svgEl = container.querySelector('svg');
        if (!svgEl) {
            container.remove();
            reject(new Error('No SVG element found'));
            return;
        }

        svgEl.querySelectorAll('script').forEach(function(s) { s.remove(); });
        svgEl.querySelectorAll('foreignObject').forEach(function(f) { f.remove(); });

        svgEl.setAttribute('width', w);
        svgEl.setAttribute('height', h);

        var canvas = document.createElement('canvas');
        canvas.width = w;
        canvas.height = h;
        var ctx = canvas.getContext('2d', { willReadFrequently: true });

        var animations = svgEl.getAnimations({ subtree: true });
        var hasAnimations = animations.length > 0;
        var totalFrames = hasAnimations ? FRAME_COUNT : 1;

        if (hasAnimations) {
            animations.forEach(function(a) { a.pause(); });
        }

        var frames = [];
        var frameIndex = 0;
        var serializer = new XMLSerializer();

        function processNextFrame() {
            if (_exportCancelled) {
                container.remove();
                reject(new Error('cancelled'));
                return;
            }

            if (frameIndex >= totalFrames) {
                container.remove();
                resolve({ frames: frames, delay: Math.round(1000 / FPS), hasAnimations: hasAnimations });
                return;
            }

            var timeMs = frameIndex * FRAME_INTERVAL;

            if (hasAnimations) {
                animations.forEach(function(a) {
                    a.currentTime = timeMs;
                });
            }

            requestAnimationFrame(function() {
                var savedStyles = [];
                if (hasAnimations) {
                    animations.forEach(function(a) {
                        var el = a.effect.target;
                        if (!el) return;
                        var cs = window.getComputedStyle(el);
                        savedStyles.push({
                            el: el,
                            origStyle: el.getAttribute('style') || '',
                            transform: cs.transform,
                            opacity: cs.opacity,
                            strokeDashoffset: cs.strokeDashoffset
                        });
                    });

                    savedStyles.forEach(function(info) {
                        if (info.transform && info.transform !== 'none') {
                            info.el.style.transform = info.transform;
                        }
                        info.el.style.opacity = info.opacity;
                        if (info.strokeDashoffset) {
                            info.el.style.strokeDashoffset = info.strokeDashoffset;
                        }
                        info.el.style.animation = 'none';
                    });
                }

                var svgStr = serializer.serializeToString(svgEl);

                if (hasAnimations) {
                    savedStyles.forEach(function(info) {
                        if (info.origStyle) {
                            info.el.setAttribute('style', info.origStyle);
                        } else {
                            info.el.removeAttribute('style');
                        }
                    });
                }

                var blob = new Blob([svgStr], { type: 'image/svg+xml;charset=utf-8' });
                var url = URL.createObjectURL(blob);
                var img = new Image();
                img.onload = function() {
                    ctx.fillStyle = '#fff';
                    ctx.fillRect(0, 0, w, h);
                    ctx.drawImage(img, 0, 0, w, h);
                    URL.revokeObjectURL(url);

                    frames.push(ctx.getImageData(0, 0, w, h));

                    frameIndex++;
                    if (onProgress) onProgress(frameIndex / totalFrames);
                    setTimeout(processNextFrame, 0);
                };
                img.onerror = function(e) {
                    console.error('SVG frame render failed at frame', frameIndex, e);
                    URL.revokeObjectURL(url);
                    container.remove();
                    reject(new Error('Failed to render SVG frame'));
                };
                img.src = url;
            });
        }

        setTimeout(processNextFrame, 50);
    });
}

function downloadGIF() {
    closeBurgerMenu();
    if (typeof GIF === 'undefined') {
        showToast(t('export.error', { message: 'gif.js not loaded' }));
        return;
    }
    var svgStr = getAnimatedSVG();
    if (!svgStr) return;

    var resPct = getResolution();
    var dims = getBaseDimensions();
    var GIF_FPS = 10;
    var w = Math.round(dims[0] * resPct / 100);
    var h = Math.round(dims[1] * resPct / 100);

    _exportCancelled = false;
    var gifEncoder = null;

    showExportProgress(function() {
        if (gifEncoder) gifEncoder.abort();
    });

    captureAnimationFrames(svgStr, w, h, function(p) {
        updateExportProgress(p * 50, t('export.progress', { percent: Math.round(p * 50) }));
    }, GIF_FPS).then(function(result) {
        if (_exportCancelled) return;

        if (!result.hasAnimations) {
            showToast(t('export.noAnimations'));
        }

        updateExportProgress(50, t('export.encoding'));

        gifEncoder = new GIF({
            workers: 2,
            quality: 10,
            width: w,
            height: h,
            workerScript: 'gif.worker.js',
            globalPalette: true
        });

        var tempCanvas = document.createElement('canvas');
        tempCanvas.width = w;
        tempCanvas.height = h;
        var tempCtx = tempCanvas.getContext('2d');

        result.frames.forEach(function(imageData) {
            tempCtx.putImageData(imageData, 0, 0);
            gifEncoder.addFrame(tempCanvas, { delay: result.delay, copy: true });
        });

        gifEncoder.on('progress', function(p) {
            updateExportProgress(50 + p * 50, t('export.encoding'));
        });

        gifEncoder.on('finished', function(blob) {
            hideExportProgress();
            var a = document.createElement('a');
            a.href = URL.createObjectURL(blob);
            a.download = getMapTitle() + '.gif';
            a.click();
            URL.revokeObjectURL(a.href);
        });

        gifEncoder.render();
    }).catch(function(err) {
        console.error('GIF export error:', err);
        hideExportProgress();
        if (err.message !== 'cancelled') {
            showToast(t('export.error', { message: err.message }));
        }
    });
}

function downloadAPNG() {
    closeBurgerMenu();
    if (typeof UPNG === 'undefined') {
        showToast(t('export.error', { message: 'UPNG.js not loaded' }));
        return;
    }
    var svgStr = getAnimatedSVG();
    if (!svgStr) return;

    var resPct = getResolution();
    var dims = getBaseDimensions();
    var w = Math.round(dims[0] * resPct / 100);
    var h = Math.round(dims[1] * resPct / 100);

    _exportCancelled = false;

    showExportProgress(null);

    captureAnimationFrames(svgStr, w, h, function(p) {
        updateExportProgress(p * 70, t('export.progress', { percent: Math.round(p * 70) }));
    }).then(function(result) {
        if (_exportCancelled) return;

        if (!result.hasAnimations) {
            showToast(t('export.noAnimations'));
        }

        updateExportProgress(70, t('export.encoding'));

        setTimeout(function() {
            try {
                var bufs = result.frames.map(function(f) { return f.data.buffer; });
                var delays = result.frames.map(function() { return result.delay; });
                var apngData = UPNG.encode(bufs, w, h, 0, delays);

                hideExportProgress();

                var blob = new Blob([apngData], { type: 'image/png' });
                var a = document.createElement('a');
                a.href = URL.createObjectURL(blob);
                a.download = getMapTitle() + '.apng';
                a.click();
                URL.revokeObjectURL(a.href);
            } catch (err) {
                hideExportProgress();
                showToast(t('export.error', { message: err.message }));
            }
        }, 50);
    }).catch(function(err) {
        console.error('APNG export error:', err);
        hideExportProgress();
        if (err.message !== 'cancelled') {
            showToast(t('export.error', { message: err.message }));
        }
    });
}

function detachMap() {
    if (detachedWindow && !detachedWindow.closed) {
        detachedWindow.focus();
        updateDetached();
        return;
    }
    detachedWindow = window.open('', 'WardleyMap', 'width=1100,height=900,resizable=yes');
    if (!detachedWindow) return;
    detachedWindow.document.write('<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Wardley Map</title><style>*{margin:0;padding:0}body{width:100vw;height:100vh;background:#fff;overflow:hidden}#map{width:100%;height:100%;display:flex;align-items:center;justify-content:center}#map svg{width:100%;height:100%}</style></head><body><div id="map"></div></body></html>');
    detachedWindow.document.close();
    updateDetached();
}

function updateDetached() {
    if (!detachedWindow || detachedWindow.closed) { detachedWindow = null; return; }
    const svgStr = getExportSVG();
    if (!svgStr) return;
    const container = detachedWindow.document.getElementById('map');
    if (container) {
        container.innerHTML = svgStr;
        container.querySelectorAll('script').forEach(function(old) {
            var s = detachedWindow.document.createElement('script');
            s.textContent = old.textContent;
            old.parentNode.replaceChild(s, old);
        });
    }
}

// ================================================================
// 15b. WTG2 Export / Import
// ================================================================

function downloadWTG2() {
    const text = getCurrentWTG2();
    if (!text.trim()) return;
    const blob = new Blob([text], { type: 'text/plain' });
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = getMapTitle() + '.wtg2';
    a.click();
    URL.revokeObjectURL(a.href);
}

function importWTG2() {
    document.getElementById('import-file-input').value = '';
    document.getElementById('import-file-input').click();
}

function handleImportFile(event) {
    const file = event.target.files && event.target.files[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = function(e) {
        loadWTG2IntoEditor(e.target.result);
    };
    reader.readAsText(file);
}

function loadWTG2IntoEditor(text) {
    setMode('editor');
    if (editor) {
        editor.setValue(text);
        editor.setCursor(0, 0);
    }
    render();
}

function openOWMImportModal() {
    var overlay = document.createElement('div');
    overlay.className = 'collab-overlay';
    overlay.id = 'owm-import-overlay';
    overlay.innerHTML =
        '<div class="collab-dialog" style="width:500px">' +
        '  <h2>' + t('owm.title') + '</h2>' +
        '  <div class="field">' +
        '    <label>' + t('owm.uploadLabel') + '</label>' +
        '    <input type="file" id="owm-file-input" accept=".owm,.txt,.json" style="width:100%;background:rgba(255,255,255,0.7);border:1px solid rgba(14,35,86,0.10);border-radius:8px;padding:8px 10px;color:#0E2356;font-size:13px;">' +
        '    <div class="hint">' + t('owm.uploadHint') + '</div>' +
        '  </div>' +
        '  <div style="text-align:center;color:#8691AB;font-size:12px;margin:12px 0;">' + t('owm.orSeparator') + '</div>' +
        '  <div class="field">' +
        '    <label>' + t('owm.pasteLabel') + '</label>' +
        '    <textarea id="owm-text-input" placeholder="' + t('owm.pastePlaceholder').replace(/"/g, '&quot;') + '" style="width:100%;min-height:140px;background:rgba(255,255,255,0.7);border:1px solid rgba(14,35,86,0.10);border-radius:8px;padding:10px 12px;color:#0E2356;font-size:13px;font-family:\'SF Mono\',\'Fira Code\',monospace;resize:vertical;"></textarea>' +
        '  </div>' +
        '  <div id="owm-import-error" style="display:none;color:#e74c3c;font-size:12px;margin-top:8px;"></div>' +
        '  <div class="btn-row">' +
        '    <button class="btn" id="owm-cancel-btn">' + t('owm.cancel') + '</button>' +
        '    <button class="btn btn-primary" id="owm-import-btn">' + t('owm.import') + '</button>' +
        '  </div>' +
        '</div>';
    document.body.appendChild(overlay);

    setTimeout(function() { document.getElementById('owm-text-input').focus(); }, 50);

    document.getElementById('owm-cancel-btn').addEventListener('click', closeOWMImportModal);
    document.getElementById('owm-import-btn').addEventListener('click', doOWMImport);

    overlay.addEventListener('click', function(e) {
        if (e.target === overlay) closeOWMImportModal();
    });
    overlay.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') closeOWMImportModal();
    });

    document.getElementById('owm-file-input').addEventListener('change', function() {
        document.getElementById('owm-text-input').value = '';
    });
    document.getElementById('owm-text-input').addEventListener('input', function() {
        if (this.value.trim()) document.getElementById('owm-file-input').value = '';
    });
}

function closeOWMImportModal() {
    var overlay = document.getElementById('owm-import-overlay');
    if (overlay) overlay.remove();
}

function doOWMImport() {
    var errorEl = document.getElementById('owm-import-error');
    errorEl.style.display = 'none';

    var fileInput = document.getElementById('owm-file-input');
    var textInput = document.getElementById('owm-text-input');
    var pastedText = textInput.value.trim();

    if (fileInput.files && fileInput.files[0]) {
        var reader = new FileReader();
        reader.onload = function(e) {
            processOWMText(e.target.result);
        };
        reader.readAsText(fileInput.files[0]);
    } else if (pastedText) {
        processOWMText(pastedText);
    } else {
        errorEl.textContent = t('owm.emptyError');
        errorEl.style.display = 'block';
    }
}

function processOWMText(owmText) {
    var errorEl = document.getElementById('owm-import-error');

    if (typeof convertOWMToWTG2 !== 'function') {
        errorEl.textContent = t('status.wasmNotLoaded');
        errorEl.style.display = 'block';
        return;
    }

    var result = convertOWMToWTG2(owmText);
    if (typeof result === 'string' && result.indexOf('error:') === 0) {
        errorEl.textContent = t('owm.convertError', { message: result.substring(7) });
        errorEl.style.display = 'block';
        return;
    }

    closeOWMImportModal();
    loadWTG2IntoEditor(result);
}

// ================================================================
// 15c. Compression & URL Sharing
// ================================================================

function compressText(text) {
    const byteArray = new TextEncoder().encode(text);
    const cs = new CompressionStream('gzip');
    const writer = cs.writable.getWriter();
    writer.write(byteArray);
    writer.close();
    return new Response(cs.readable).arrayBuffer();
}

function decompressBuffer(buffer) {
    const cs = new DecompressionStream('gzip');
    const writer = cs.writable.getWriter();
    writer.write(buffer);
    writer.close();
    return new Response(cs.readable).arrayBuffer().then(function(arrayBuffer) {
        return new TextDecoder().decode(arrayBuffer);
    });
}

function arrayBufferToUrlSafeBase64(buffer) {
    let binary = '';
    const bytes = new Uint8Array(buffer);
    for (let i = 0; i < bytes.byteLength; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

function urlSafeBase64ToArrayBuffer(base64) {
    let std = base64.replace(/-/g, '+').replace(/_/g, '/');
    while (std.length % 4) std += '=';
    const binary = window.atob(std);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
        bytes[i] = binary.charCodeAt(i);
    }
    return bytes.buffer;
}

async function shareURL() {
    const text = getCurrentWTG2();
    if (!text.trim()) return;
    try {
        const compressed = await compressText(text);
        const encoded = arrayBufferToUrlSafeBase64(compressed);
        const url = location.origin + location.pathname + '?wtg2=' + encoded;
        if (url.length > 8000) {
            showToast(t('share.tooLong', { length: url.length }));
        }
        window.history.replaceState({}, '', '?wtg2=' + encoded);
        try {
            await navigator.clipboard.writeText(url);
            showToast(t('share.copied'));
        } catch (e) {
            showToast(t('share.manual'));
        }
    } catch (e) {
        showToast(t('share.error', { message: e.message }));
    }
}

function showToast(message) {
    const toast = document.getElementById('share-toast');
    toast.textContent = message;
    toast.classList.add('visible');
    clearTimeout(toast._timer);
    toast._timer = setTimeout(function() { toast.classList.remove('visible'); }, 3000);
}

// ================================================================
// 16. Persistence
// ================================================================

function saveState() {
    saveCurrentMapToStorage();
}

function loadState() {
    loadMapsIndex();
    var savedId = localStorage.getItem('wtg2-current-map-id');
    if (savedId && mapsIndex.find(function(m) { return m.id === savedId; })) {
        return openMap(savedId);
    }
    if (mapsIndex.length > 0) {
        var sorted = mapsIndex.slice().sort(function(a, b) {
            return new Date(b.modifiedAt) - new Date(a.modifiedAt);
        });
        return openMap(sorted[0].id);
    }
    return false;
}

function restoreUI() {
    document.getElementById('meta-title').value = wizardState.meta.title || '';
    document.getElementById('meta-author').value = wizardState.meta.author || '';
    document.getElementById('meta-question').value = wizardState.meta.question || '';
    document.getElementById('legend-toggle').checked = wizardState.legend;
    for (let i = 0; i < 4; i++) {
        document.getElementById('stage-' + i).value = (wizardState.stages && wizardState.stages[i]) || '';
    }
    updateChainInputFromState();
    renderCompList();
}

// ================================================================
// 16b. Multi-Map Management
// ================================================================

function generateMapId() {
    if (crypto.randomUUID) return crypto.randomUUID();
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = Math.random() * 16 | 0;
        return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
    });
}

function formatDateShort(isoString) {
    try {
        var d = new Date(isoString);
        return d.toLocaleDateString(currentLang === 'fr' ? 'fr-FR' : 'en-US', {
            day: 'numeric', month: 'short', year: 'numeric',
            hour: '2-digit', minute: '2-digit'
        });
    } catch (e) { return isoString; }
}

function loadMapsIndex() {
    try {
        var raw = localStorage.getItem('wtg2-maps-index');
        mapsIndex = raw ? JSON.parse(raw) : [];
    } catch (e) { mapsIndex = []; }
    return mapsIndex;
}

function saveMapsIndex() {
    try { localStorage.setItem('wtg2-maps-index', JSON.stringify(mapsIndex)); }
    catch (e) { /* quota exceeded */ }
}

function saveCurrentMapToStorage() {
    if (!currentMapId) return;
    try {
        localStorage.setItem('wtg2-map-' + currentMapId, JSON.stringify(wizardState));
        var entry = mapsIndex.find(function(m) { return m.id === currentMapId; });
        if (entry) {
            entry.modifiedAt = new Date().toISOString();
            entry.title = (wizardState.meta && wizardState.meta.title) || '';
            saveMapsIndex();
        }
    } catch (e) { /* quota exceeded */ }
}

function parseStateFromJSON(parsed) {
    return {
        meta: parsed.meta || { title: '', author: '', question: '' },
        stages: parsed.stages || ['', '', '', ''],
        components: (parsed.components || []).map(function(c) {
            return Object.assign({}, c, {
                isPipeline: c.isPipeline || false,
                pipelineMembers: c.pipelineMembers || []
            });
        }),
        edges: parsed.edges || [],
        groups: parsed.groups || [],
        annotations: parsed.annotations || [],
        signals: parsed.signals || [],
        legend: parsed.legend || false,
        focus: parsed.focus || ''
    };
}

function createNewMap() {
    if (currentMapId) saveCurrentMapToStorage();

    var id = generateMapId();
    var now = new Date().toISOString();
    mapsIndex.push({ id: id, title: '', createdAt: now, modifiedAt: now });
    saveMapsIndex();

    wizardState = { meta: { title: '', author: '', question: '' }, stages: ['', '', '', ''], components: [], edges: [], groups: [], annotations: [], signals: [], legend: false, focus: '' };
    currentMapId = id;
    localStorage.setItem('wtg2-current-map-id', id);
    saveCurrentMapToStorage();

    undoStack = [];
    redoStack = [];
    updateUndoButtons();
    currentStep = 0;
    if (editor) editor.setValue(generateWTG2(wizardState));
    canvasReset();
    syncWizardUI();
    render();
}

function openMap(id) {
    if (currentMapId && currentMapId !== id) saveCurrentMapToStorage();

    try {
        var raw = localStorage.getItem('wtg2-map-' + id);
        if (!raw) return false;
        wizardState = parseStateFromJSON(JSON.parse(raw));
        currentMapId = id;
        localStorage.setItem('wtg2-current-map-id', id);

        undoStack = [];
        redoStack = [];
        updateUndoButtons();
        restoreUI();
        render();
        return true;
    } catch (e) { return false; }
}

function deleteMap(id) {
    localStorage.removeItem('wtg2-map-' + id);
    mapsIndex = mapsIndex.filter(function(m) { return m.id !== id; });
    saveMapsIndex();

    if (currentMapId === id) {
        currentMapId = null;
        if (mapsIndex.length > 0) {
            var sorted = mapsIndex.slice().sort(function(a, b) {
                return new Date(b.modifiedAt) - new Date(a.modifiedAt);
            });
            openMap(sorted[0].id);
        } else {
            createNewMap();
        }
    }
}

function openMapsModal() {
    loadMapsIndex();

    var overlay = document.createElement('div');
    overlay.className = 'maps-overlay';
    overlay.id = 'maps-overlay';

    var listHTML = '';
    if (mapsIndex.length === 0) {
        listHTML = '<div class="maps-list-empty">' + t('maps.empty') + '</div>';
    } else {
        var sorted = mapsIndex.slice().sort(function(a, b) {
            return new Date(b.modifiedAt) - new Date(a.modifiedAt);
        });
        listHTML = sorted.map(function(m) {
            var isCurrent = m.id === currentMapId;
            var title = m.title || t('maps.untitled');
            var created = formatDateShort(m.createdAt);
            var modified = formatDateShort(m.modifiedAt);
            return '<div class="map-entry' + (isCurrent ? ' current' : '') + '" data-map-id="' + m.id + '">' +
                '<div class="map-entry-info">' +
                    '<div class="map-entry-title">' + escapeHtml(title) +
                        (isCurrent ? ' <span class="current-badge">' + t('maps.current') + '</span>' : '') +
                    '</div>' +
                    '<div class="map-entry-dates">' +
                        t('maps.created') + ' ' + created + ' &middot; ' +
                        t('maps.modified') + ' ' + modified +
                    '</div>' +
                '</div>' +
                '<div class="map-entry-actions">' +
                    (isCurrent ? '' : '<button class="btn btn-primary btn-open-map" data-id="' + m.id + '">' + t('maps.open') + '</button>') +
                    '<button class="btn btn-danger btn-delete-map" data-id="' + m.id + '">' + t('maps.delete') + '</button>' +
                '</div>' +
            '</div>';
        }).join('');
    }

    overlay.innerHTML =
        '<div class="maps-dialog">' +
        '<h2>' + t('maps.title') + '</h2>' +
        '<div class="maps-list">' + listHTML + '</div>' +
        '<div class="btn-row">' +
        '<button class="btn" id="maps-close-btn">' + t('maps.close') + '</button>' +
        '</div>' +
        '</div>';

    document.body.appendChild(overlay);

    overlay.addEventListener('click', function(e) {
        var openBtn = e.target.closest('.btn-open-map');
        var deleteBtn = e.target.closest('.btn-delete-map');
        if (openBtn) {
            openMap(openBtn.getAttribute('data-id'));
            closeMapsModal();
        } else if (deleteBtn) {
            if (confirm(t('menu.confirmDelete'))) {
                deleteMap(deleteBtn.getAttribute('data-id'));
                closeMapsModal();
                if (mapsIndex.length > 0) openMapsModal();
            }
        } else if (e.target === overlay) {
            closeMapsModal();
        }
    });

    document.getElementById('maps-close-btn').addEventListener('click', closeMapsModal);
    document.addEventListener('keydown', mapsModalEscHandler);
}

function mapsModalEscHandler(e) {
    if (e.key === 'Escape' && document.getElementById('maps-overlay')) {
        closeMapsModal();
    }
}

function closeMapsModal() {
    var overlay = document.getElementById('maps-overlay');
    if (overlay) overlay.remove();
    document.removeEventListener('keydown', mapsModalEscHandler);
}

function migrateToMultiMap() {
    if (localStorage.getItem('wtg2-maps-index')) return;
    var oldState = localStorage.getItem('wtg2-guided-state');
    if (!oldState) return;
    try {
        var parsed = JSON.parse(oldState);
        var id = generateMapId();
        var now = new Date().toISOString();
        var index = [{ id: id, title: (parsed.meta && parsed.meta.title) || '', createdAt: now, modifiedAt: now }];
        localStorage.setItem('wtg2-map-' + id, oldState);
        localStorage.setItem('wtg2-maps-index', JSON.stringify(index));
        localStorage.setItem('wtg2-current-map-id', id);
        localStorage.removeItem('wtg2-guided-state');
    } catch (e) { /* migration failed */ }
}

// ================================================================
// 17. Responsive: Mobile Tabs
// ================================================================

function setMobileTab(tab) {
    const btns = document.querySelectorAll('#mobile-tabs button');
    btns[0].classList.toggle('active', tab === 'edit');
    btns[1].classList.toggle('active', tab === 'preview');

    document.getElementById('left-pane').classList.toggle('hidden', tab !== 'edit');
    document.getElementById('right-pane').classList.toggle('visible', tab === 'preview');

    if (tab === 'preview') render();
    if (tab === 'edit' && editor) setTimeout(() => editor.refresh(), 10);
}

// ================================================================
// 18. Split Pane Divider
// ================================================================

(function initDivider() {
    const divider = document.getElementById('divider');
    const leftPane = document.getElementById('left-pane');
    const main = document.getElementById('main');
    let isDragging = false;

    function startDrag(e) {
        isDragging = true;
        divider.classList.add('dragging');
        document.body.style.cursor = 'col-resize';
        document.body.style.userSelect = 'none';
        e.preventDefault();
    }

    function doDrag(e) {
        if (!isDragging) return;
        const clientX = e.touches ? e.touches[0].clientX : e.clientX;
        const rect = main.getBoundingClientRect();
        const pct = ((clientX - rect.left) / rect.width) * 100;
        const clamped = Math.max(20, Math.min(70, pct));
        leftPane.style.width = clamped + '%';
        if (editor) editor.refresh();
    }

    function stopDrag() {
        if (!isDragging) return;
        isDragging = false;
        divider.classList.remove('dragging');
        document.body.style.cursor = '';
        document.body.style.userSelect = '';
    }

    divider.addEventListener('mousedown', startDrag);
    document.addEventListener('mousemove', doDrag);
    document.addEventListener('mouseup', stopDrag);
    divider.addEventListener('touchstart', startDrag, { passive: false });
    document.addEventListener('touchmove', doDrag, { passive: false });
    document.addEventListener('touchend', stopDrag);
})();

// ================================================================
// 19. Burger Menu
// ================================================================

function openBurgerMenu() {
    document.getElementById('burger-panel').classList.add('open');
    document.getElementById('burger-backdrop').classList.add('open');
}
function closeBurgerMenu() {
    document.getElementById('burger-panel').classList.remove('open');
    document.getElementById('burger-backdrop').classList.remove('open');
}
document.getElementById('burger-btn').addEventListener('click', openBurgerMenu);
document.getElementById('burger-backdrop').addEventListener('click', closeBurgerMenu);
document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape' && document.getElementById('burger-panel').classList.contains('open')) {
        closeBurgerMenu();
    }
});

document.getElementById('burger-new').addEventListener('click', function() {
    if (confirm(t('menu.confirmNew'))) {
        createNewMap();
    }
    closeBurgerMenu();
});

document.getElementById('burger-my-maps').addEventListener('click', function() {
    closeBurgerMenu();
    openMapsModal();
});

document.querySelectorAll('.burger-template').forEach(function(btn) {
    btn.addEventListener('click', function() {
        var key = btn.getAttribute('data-template');
        createNewMap();
        loadTemplate(key);
        closeBurgerMenu();
    });
});

// ================================================================
// 20. Toolbar Events
// ================================================================

document.getElementById('ratio').addEventListener('change', render);
const resSlider = document.getElementById('resolution');
const resLabel = document.getElementById('res-label');
resSlider.addEventListener('input', () => {
    resLabel.textContent = resSlider.value + '%';
    scheduleRender();
});
document.getElementById('dl-svg').addEventListener('click', downloadSVG);
document.getElementById('dl-png').addEventListener('click', downloadPNG);
document.getElementById('dl-gif').addEventListener('click', downloadGIF);
document.getElementById('dl-apng').addEventListener('click', downloadAPNG);
document.getElementById('detach').addEventListener('click', detachMap);
document.getElementById('dl-wtg2').addEventListener('click', downloadWTG2);
document.getElementById('import-wtg2').addEventListener('click', importWTG2);
document.getElementById('import-file-input').addEventListener('change', handleImportFile);
document.getElementById('import-owm').addEventListener('click', function() {
    closeBurgerMenu();
    openOWMImportModal();
});
document.getElementById('share-url').addEventListener('click', shareURL);

// Language selector
document.getElementById('lang-select').addEventListener('change', function() {
    currentLang = this.value;
    localStorage.setItem('wtg2-lang', currentLang);
    document.documentElement.lang = currentLang;
    applyStaticTranslations();
    renderCompList();
    if (currentStep === 2) renderEvoList();
    if (currentStep === 3) { renderGroups(); renderAnnotations(); renderSignals(); }
    goToStep(currentStep);
});

// Restore saved language on load
(function initLang() {
    const savedLang = localStorage.getItem('wtg2-lang');
    if (savedLang && translations[savedLang]) currentLang = savedLang;
    document.documentElement.lang = currentLang;
    document.getElementById('lang-select').value = currentLang;
    applyStaticTranslations();
})();

initChainEditor();

// ================================================================
// 20. Utilities
// ================================================================

function escapeHtml(s) {
    const d = document.createElement('div');
    d.textContent = s;
    return d.innerHTML;
}

function escapeAttr(s) {
    return s.replace(/'/g, "\\'").replace(/"/g, '&quot;');
}

// ================================================================
// 20b. Onboarding / Coach Marks
// ================================================================

const ONBOARDING_KEY = 'wtg2-onboarding-done';

const onboardingSteps = [
    {
        target: '#wizard-pane',
        titleKey: 'onboarding.welcome.title',
        bodyKey: 'onboarding.welcome.body',
        tryItKey: null,
        arrow: 'left',
        padding: 8,
        beforeShow: null,
    },
    {
        target: '#chain-editor-wrapper',
        titleKey: 'onboarding.step1.title',
        bodyKey: 'onboarding.step1.body',
        tryItKey: 'onboarding.step1.tryit',
        arrow: 'left',
        padding: 6,
        beforeShow: function() { goToStep(1); },
    },
    {
        target: '.steps-bar .step-indicator:nth-child(5)',
        titleKey: 'onboarding.step2.title',
        bodyKey: 'onboarding.step2.body',
        tryItKey: null,
        arrow: 'top',
        padding: 4,
        beforeShow: function() { goToStep(2); },
    },
    {
        target: '.steps-bar .step-indicator:nth-child(7)',
        titleKey: 'onboarding.step3.title',
        bodyKey: 'onboarding.step3.body',
        tryItKey: null,
        arrow: 'top',
        padding: 4,
        beforeShow: function() { goToStep(3); },
    },
    {
        target: '#right-pane',
        titleKey: 'onboarding.preview.title',
        bodyKey: 'onboarding.preview.body',
        tryItKey: null,
        arrow: 'right',
        padding: 12,
        beforeShow: function() { goToStep(1); },
    },
];

var onboardingCurrentStep = 0;
var onboardingActive = false;

function shouldShowOnboarding() {
    if (localStorage.getItem(ONBOARDING_KEY)) return false;
    if (new URLSearchParams(window.location.search).has('wtg2')) return false;
    if (window.innerWidth < 768) return false;
    return true;
}

function startOnboarding() {
    onboardingActive = true;
    onboardingCurrentStep = 0;
    loadTemplate('blank');
    var overlay = document.getElementById('onboarding-overlay');
    overlay.classList.add('active');
    showOnboardingStep(0);
}

function showOnboardingStep(stepIdx) {
    onboardingCurrentStep = stepIdx;
    var step = onboardingSteps[stepIdx];
    var overlay = document.getElementById('onboarding-overlay');
    var total = onboardingSteps.length;

    if (step.beforeShow) step.beforeShow();

    requestAnimationFrame(function() {
        var targetEl = document.querySelector(step.target);
        if (!targetEl) { finishOnboarding(); return; }

        var rect = targetEl.getBoundingClientRect();
        var pad = step.padding || 8;

        var dotsHtml = '';
        for (var i = 0; i < total; i++) {
            dotsHtml += '<span class="onboarding-dot' + (i === stepIdx ? ' active' : '') + '"></span>';
        }

        var isLast = stepIdx === total - 1;
        var nextLabel = isLast ? t('onboarding.finish') : t('onboarding.next');
        var tryItHtml = step.tryItKey ? '<div class="try-it">' + t(step.tryItKey) + '</div>' : '';

        var tooltipStyle = '';
        var arrowClass = 'arrow-' + step.arrow;

        switch (step.arrow) {
            case 'left':
                tooltipStyle = 'top:' + Math.max(8, rect.top + rect.height / 2 - 60) + 'px;left:' + (rect.right + pad + 16) + 'px;';
                break;
            case 'right':
                tooltipStyle = 'top:' + Math.max(8, rect.top + rect.height / 2 - 60) + 'px;right:' + (window.innerWidth - rect.left + pad + 16) + 'px;';
                break;
            case 'top':
                tooltipStyle = 'top:' + (rect.bottom + pad + 16) + 'px;left:' + rect.left + 'px;';
                break;
            case 'bottom':
                tooltipStyle = 'bottom:' + (window.innerHeight - rect.top + pad + 16) + 'px;left:' + rect.left + 'px;';
                break;
        }

        overlay.innerHTML =
            '<div class="onboarding-backdrop"></div>' +
            '<div class="onboarding-spotlight pulse" style="' +
                'top:' + (rect.top - pad) + 'px;' +
                'left:' + (rect.left - pad) + 'px;' +
                'width:' + (rect.width + pad * 2) + 'px;' +
                'height:' + (rect.height + pad * 2) + 'px;' +
            '"></div>' +
            '<div class="onboarding-tooltip ' + arrowClass + '" style="' + tooltipStyle + '" ' +
                 'role="dialog" tabindex="-1">' +
                '<h3>' + t(step.titleKey) + '</h3>' +
                '<p>' + t(step.bodyKey) + '</p>' +
                tryItHtml +
                '<div class="onboarding-nav">' +
                    '<button class="onboarding-skip" onclick="finishOnboarding()">' + t('onboarding.skip') + '</button>' +
                    '<div class="onboarding-dots">' + dotsHtml + '</div>' +
                    '<button class="onboarding-next" onclick="nextOnboardingStep()">' + nextLabel + '</button>' +
                '</div>' +
            '</div>';

        var tooltip = overlay.querySelector('.onboarding-tooltip');
        if (tooltip) {
            tooltip.addEventListener('keydown', function(e) {
                if (e.key === 'Escape') finishOnboarding();
                if (e.key === 'ArrowRight' || e.key === 'Enter') nextOnboardingStep();
                if (e.key === 'ArrowLeft' && onboardingCurrentStep > 0) showOnboardingStep(onboardingCurrentStep - 1);
            });
            tooltip.focus();
        }
    });
}

function nextOnboardingStep() {
    if (onboardingCurrentStep < onboardingSteps.length - 1) {
        showOnboardingStep(onboardingCurrentStep + 1);
    } else {
        finishOnboarding();
    }
}

function finishOnboarding() {
    onboardingActive = false;
    localStorage.setItem(ONBOARDING_KEY, '1');
    var overlay = document.getElementById('onboarding-overlay');
    overlay.classList.remove('active');
    overlay.innerHTML = '';
    loadTemplate('teashop');
    goToStep(0);
}

// ================================================================
// 21. WASM Init & Startup
// ================================================================

const go = new Go();

(async function initWasm() {
    const statusEl = document.getElementById('status');
    const progressContainer = document.getElementById('wasm-progress');
    const progressBar = document.getElementById('wasm-progress-bar');

    try {
        let result;
        const response = await fetch('main.wasm');
        const contentLength = response.headers.get('Content-Length');

        if (contentLength && response.body) {
            if (progressContainer) progressContainer.style.display = '';
            const total = parseInt(contentLength, 10);
            const reader = response.body.getReader();
            const chunks = [];
            let loaded = 0;

            while (true) {
                const { done, value } = await reader.read();
                if (done) break;
                chunks.push(value);
                loaded += value.length;
                const pct = Math.min(100, Math.round(loaded / total * 100));
                if (statusEl) statusEl.textContent = t('status.loading') + ' ' + pct + '%';
                if (progressBar) progressBar.style.width = pct + '%';
            }

            if (progressContainer) progressContainer.style.display = 'none';
            const wasmBytes = new Uint8Array(loaded);
            let offset = 0;
            for (const chunk of chunks) {
                wasmBytes.set(chunk, offset);
                offset += chunk.length;
            }
            result = await WebAssembly.instantiate(wasmBytes, go.importObject);
        } else {
            if (progressContainer) {
                progressContainer.style.display = '';
                if (progressBar) {
                    progressBar.style.width = '100%';
                    progressBar.style.animation = 'indeterminate 1.5s ease-in-out infinite';
                }
            }
            result = await WebAssembly.instantiateStreaming(response, go.importObject);
            if (progressContainer) progressContainer.style.display = 'none';
        }

        go.run(result.instance);
        if (statusEl) statusEl.textContent = t('status.ready');

        const params = new URLSearchParams(window.location.search);
        if (params.has('wtg2')) {
            try {
                const buffer = urlSafeBase64ToArrayBuffer(params.get('wtg2'));
                const text = await decompressBuffer(buffer);
                loadWTG2IntoEditor(text);
                return;
            } catch (e) {
                if (statusEl) statusEl.textContent = t('status.urlError', { message: e.message });
            }
        }

        migrateToMultiMap();
        if (loadState()) {
            // openMap already calls restoreUI and render
        } else {
            createNewMap();
            loadTemplate('teashop');
        }
        render();

        if (shouldShowOnboarding()) {
            setTimeout(startOnboarding, 600);
        }
    } catch (err) {
        if (progressContainer) progressContainer.style.display = 'none';
        if (statusEl) statusEl.textContent = 'WASM: ' + err;
        migrateToMultiMap();
        if (loadState()) {
            // loaded
        } else {
            createNewMap();
            loadTemplate('teashop');
        }
    }
})();

// ================================================================
// 22. COLLABORATION
// ================================================================

var collabWS = null;
var collabClientId = null;
var collabMode = null;
var collabUsers = {};
var collabVersion = 0;
var collabRemoteUpdate = false;
var collabReconnectTimer = null;
var collabReconnectDelay = 1000;
var collabLastUrl = null;
var collabLastName = null;
var collabRemoteCursors = {};

function openCollabDialog() {
    if (collabWS) {
        // Already connected — show disconnect option
        disconnectSession();
        return;
    }

    var savedName = localStorage.getItem('wtg2-collab-name') || '';
    var overlay = document.createElement('div');
    overlay.className = 'collab-overlay';
    overlay.id = 'collab-overlay';
    overlay.innerHTML =
        '<div class="collab-dialog">' +
        '  <h2>' + t('collab.title') + '</h2>' +
        '  <div class="field">' +
        '    <label>' + t('collab.urlLabel') + '</label>' +
        '    <input type="text" id="collab-url-input" placeholder="' + t('collab.urlPlaceholder') + '">' +
        '    <div class="hint">' + t('collab.urlHint') + '</div>' +
        '    <div class="error-hint" id="collab-url-error"></div>' +
        '  </div>' +
        '  <div class="field">' +
        '    <label>' + t('collab.nameLabel') + '</label>' +
        '    <input type="text" id="collab-name-input" placeholder="' + t('collab.namePlaceholder') + '" value="' + savedName.replace(/"/g, '&quot;') + '">' +
        '  </div>' +
        '  <div class="btn-row">' +
        '    <button class="btn" onclick="closeCollabDialog()">' + t('collab.cancel') + '</button>' +
        '    <button class="btn btn-primary" onclick="doConnect()">' + t('collab.connect') + '</button>' +
        '  </div>' +
        '</div>';
    document.body.appendChild(overlay);

    // Focus URL input
    setTimeout(function() { document.getElementById('collab-url-input').focus(); }, 50);

    // Allow Enter to connect
    overlay.addEventListener('keydown', function(e) {
        if (e.key === 'Enter') doConnect();
        if (e.key === 'Escape') closeCollabDialog();
    });
}

function closeCollabDialog() {
    var overlay = document.getElementById('collab-overlay');
    if (overlay) overlay.remove();
}

function doConnect() {
    var urlInput = document.getElementById('collab-url-input');
    var nameInput = document.getElementById('collab-name-input');
    var errorEl = document.getElementById('collab-url-error');
    var wsUrl = urlInput.value.trim();
    var userName = nameInput.value.trim() || 'Anonymous';

    if (!wsUrl.match(/^wss?:\/\//)) {
        errorEl.textContent = t('collab.invalidUrl');
        errorEl.style.display = 'block';
        return;
    }

    localStorage.setItem('wtg2-collab-name', userName);
    closeCollabDialog();
    connectToSession(wsUrl, userName);
}

function connectToSession(wsUrl, userName) {
    collabLastUrl = wsUrl;
    collabLastName = userName;

    // If in guided mode, switch to editor with current WTG2
    if (currentMode === 'guided') {
        setMode('editor');
    }

    var separator = wsUrl.indexOf('?') >= 0 ? '&' : '?';
    var fullUrl = wsUrl + separator + 'name=' + encodeURIComponent(userName);

    var ws = new WebSocket(fullUrl);
    collabWS = ws;

    ws.onopen = function() {
        collabReconnectDelay = 1000;
        // Send hello
        ws.send(JSON.stringify({ type: 'hello', payload: { name: userName } }));
    };

    ws.onmessage = function(evt) {
        var msg = JSON.parse(evt.data);
        switch (msg.type) {
            case 'welcome': onWelcome(msg.payload); break;
            case 'op': onRemoteOp(msg.payload); break;
            case 'ack': onAck(msg.payload); break;
            case 'user_joined': onUserJoined(msg.payload); break;
            case 'user_left': onUserLeft(msg.payload); break;
            case 'cursor': onRemoteCursor(msg.payload); break;
            case 'full_sync': onFullSync(msg.payload); break;
            case 'error': console.error('collab error:', msg.payload.message); break;
        }
    };

    ws.onclose = function() {
        if (collabWS === ws) {
            // Unexpected close — try to reconnect
            updateCollabBar();
            document.getElementById('collab-status').textContent = t('collab.reconnecting');
            collabReconnectTimer = setTimeout(function() {
                if (collabWS === ws) {
                    connectToSession(collabLastUrl, collabLastName);
                }
            }, collabReconnectDelay);
            collabReconnectDelay = Math.min(collabReconnectDelay * 2, 30000);
        }
    };

    ws.onerror = function() {
        // onclose will handle reconnection
    };
}

function disconnectSession() {
    if (collabReconnectTimer) {
        clearTimeout(collabReconnectTimer);
        collabReconnectTimer = null;
    }
    var ws = collabWS;
    collabWS = null;
    collabClientId = null;
    collabMode = null;
    collabUsers = {};
    collabVersion = 0;
    if (ws) ws.close();

    // Remove remote cursors
    for (var id in collabRemoteCursors) {
        if (collabRemoteCursors[id].bookmark) collabRemoteCursors[id].bookmark.clear();
    }
    collabRemoteCursors = {};

    // Restore editor
    if (editor) editor.setOption('readOnly', false);

    // Update UI
    document.getElementById('collab-bar').classList.remove('active');
    document.getElementById('collab-btn').innerHTML = '<span class="burger-item-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg></span> <span data-i18n="toolbar.collabLabel">' + t('toolbar.collabLabel') + '</span>';
    var dot = document.getElementById('collab-dot');
    if (dot) dot.classList.remove('active');
}

function onWelcome(payload) {
    collabClientId = payload.clientId;
    collabMode = payload.mode;
    collabVersion = payload.version;
    collabUsers = {};
    (payload.users || []).forEach(function(u) {
        collabUsers[u.id] = u;
    });

    // Load document into editor
    if (!editor) initEditor();
    collabRemoteUpdate = true;
    editor.setValue(payload.document.join('\n'));
    collabRemoteUpdate = false;

    // Set read-only for spectators
    if (collabMode === 'ro') {
        editor.setOption('readOnly', true);
    } else {
        editor.setOption('readOnly', false);
    }

    updateCollabBar();
    render();
}

function onRemoteOp(payload) {
    if (!editor) return;
    collabVersion = payload.version;

    collabRemoteUpdate = true;
    var lastLine = editor.lastLine();
    switch (payload.type) {
        case 'insert': {
            if (payload.lineStart > lastLine) {
                // Inserting past the last line — append with preceding newline
                var text = '\n' + payload.lines.join('\n');
                var endPos = { line: lastLine, ch: editor.getLine(lastLine).length };
                editor.replaceRange(text, endPos, endPos);
            } else {
                var pos = { line: payload.lineStart, ch: 0 };
                var text = payload.lines.join('\n') + '\n';
                editor.replaceRange(text, pos, pos);
            }
            break;
        }
        case 'delete': {
            var from = { line: payload.lineStart, ch: 0 };
            var endLine = payload.lineStart + payload.lineCount;
            if (endLine > lastLine) {
                var to = { line: lastLine, ch: editor.getLine(lastLine).length };
                if (payload.lineStart > 0) {
                    from = { line: payload.lineStart - 1, ch: editor.getLine(payload.lineStart - 1).length };
                }
                editor.replaceRange('', from, to);
            } else {
                var to = { line: endLine, ch: 0 };
                editor.replaceRange('', from, to);
            }
            break;
        }
        case 'replace': {
            var from = { line: payload.lineStart, ch: 0 };
            var endLine = payload.lineStart + payload.lineCount;
            if (endLine > lastLine) {
                // Range reaches past last line — don't add trailing \n
                var to = { line: lastLine, ch: editor.getLine(lastLine).length };
                var text = payload.lines.join('\n');
                editor.replaceRange(text, from, to);
            } else {
                var to = { line: endLine, ch: 0 };
                var text = payload.lines.join('\n') + '\n';
                editor.replaceRange(text, from, to);
            }
            break;
        }
    }
    collabRemoteUpdate = false;
    scheduleRender();
}

function onAck(payload) {
    collabVersion = payload.version;
}

function onFullSync(payload) {
    collabVersion = payload.version;
    if (!editor) return;
    collabRemoteUpdate = true;
    editor.setValue(payload.document.join('\n'));
    collabRemoteUpdate = false;
    scheduleRender();
}

function onUserJoined(payload) {
    collabUsers[payload.id] = payload;
    updateCollabBar();
    showToast(payload.name + ' joined');
}

function onUserLeft(payload) {
    var user = collabUsers[payload.id];
    delete collabUsers[payload.id];
    // Remove cursor
    if (collabRemoteCursors[payload.id]) {
        if (collabRemoteCursors[payload.id].bookmark) collabRemoteCursors[payload.id].bookmark.clear();
        delete collabRemoteCursors[payload.id];
    }
    updateCollabBar();
    if (user) showToast(user.name + ' left');
}

function onRemoteCursor(payload) {
    if (!editor) return;
    var prev = collabRemoteCursors[payload.clientId];
    if (prev && prev.bookmark) prev.bookmark.clear();

    var user = collabUsers[payload.clientId];
    if (!user) return;

    var el = document.createElement('span');
    el.className = 'remote-cursor';
    el.style.borderLeftColor = user.color;
    var label = document.createElement('span');
    label.className = 'remote-cursor-label';
    label.style.background = user.color;
    label.textContent = user.name;
    el.appendChild(label);

    var bookmark = editor.setBookmark({ line: payload.line, ch: payload.ch }, { widget: el, insertLeft: true });
    collabRemoteCursors[payload.clientId] = { bookmark: bookmark };
}

function updateCollabBar() {
    var bar = document.getElementById('collab-bar');
    var statusEl = document.getElementById('collab-status');
    var usersEl = document.getElementById('collab-users');

    if (!collabWS || !collabClientId) {
        bar.classList.remove('active');
        return;
    }

    bar.classList.add('active');
    var userCount = Object.keys(collabUsers).length + 1;
    statusEl.textContent = t('collab.connected') + ' (' + userCount + ')';

    var html = '';
    if (collabMode === 'ro') {
        html += '<span class="spectator-badge">' + t('collab.spectator') + '</span>';
    }
    for (var id in collabUsers) {
        var u = collabUsers[id];
        html += '<span class="collab-user"><span class="collab-dot" style="background:' + u.color + '"></span>' + escapeHtml(u.name) + '</span>';
    }
    usersEl.innerHTML = html;

    // Update button text and collab indicator
    document.getElementById('collab-btn').innerHTML = '<span class="burger-item-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg></span> <span data-i18n="toolbar.collabLabel">' + t('toolbar.collabLabel') + '</span> (' + userCount + ')';
    var dot = document.getElementById('collab-dot');
    if (dot) dot.classList.add('active');
}

function localChangeToOp(change) {
    // CodeMirror change: { from: {line, ch}, to: {line, ch}, text: [...], removed: [...] }
    var fromLine = change.from.line;
    var toLine = change.to.line;
    var removedCount = change.removed.length;
    var addedLines = change.text;

    // Determine if lines were added/removed/replaced
    if (removedCount <= 1 && addedLines.length <= 1) {
        // Single-line edit — send as replace of that line
        // Get the full content of the affected line after the edit
        var lineContent = editor.getLine(fromLine);
        if (lineContent === undefined) return null;
        return { type: 'replace', lineStart: fromLine, lineCount: 1, lines: [lineContent], version: collabVersion };
    }

    // Multi-line change: deleted lines replaced by new lines
    if (removedCount > 1 && addedLines.length >= 1) {
        // It's a replace
        var lines = [];
        for (var i = fromLine; i < fromLine + addedLines.length; i++) {
            var l = editor.getLine(i);
            if (l !== undefined) lines.push(l);
        }
        return { type: 'replace', lineStart: fromLine, lineCount: removedCount, lines: lines, version: collabVersion };
    }

    if (addedLines.length > 1) {
        // Lines were inserted (pressing Enter or pasting multi-line)
        var lines = [];
        for (var i = fromLine + 1; i < fromLine + addedLines.length; i++) {
            var l = editor.getLine(i);
            if (l !== undefined) lines.push(l);
        }
        if (lines.length > 0) {
            // Also need to update the current line if it was split
            var currentLine = editor.getLine(fromLine);
            return { type: 'replace', lineStart: fromLine, lineCount: 1, lines: [currentLine].concat(lines), version: collabVersion };
        }
    }

    if (removedCount > 1) {
        // Lines were deleted
        var currentLine = editor.getLine(fromLine);
        return { type: 'replace', lineStart: fromLine, lineCount: removedCount, lines: [currentLine], version: collabVersion };
    }

    return null;
}

// Hook into CodeMirror changes for collaboration
var origInitEditor = initEditor;
initEditor = function() {
    origInitEditor();
    editor.on('change', function(cm, change) {
        if (collabRemoteUpdate) return;
        if (!collabWS || collabMode !== 'rw') return;

        var op = localChangeToOp(change);
        if (op) {
            collabWS.send(JSON.stringify({ type: 'op', payload: op }));
        }
    });

    // Send cursor position
    editor.on('cursorActivity', function(cm) {
        if (!collabWS) return;
        var cursor = cm.getCursor();
        collabWS.send(JSON.stringify({ type: 'cursor', payload: { line: cursor.line, ch: cursor.ch } }));
    });
};
