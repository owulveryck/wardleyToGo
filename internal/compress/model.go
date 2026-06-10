package compress

// Statement type indices for DStmtType distribution.
const (
	stmtMeta       = iota // 0: meta-decl (title, date, author, scope, question, doctrine)
	stmtStages            // 1: stage-decl
	stmtNode              // 2: node-decl (anchor, component, submap)
	stmtPipeline          // 3: pipeline-decl
	stmtEdge              // 4: edge
	stmtGroup             // 5: group-decl
	stmtAnnotation        // 6: annotation (note, warning)
	stmtSignal            // 7: signal
	stmtGameplay          // 8: gameplay
	stmtFocus             // 9: focus-decl
	stmtLegend            // 10: legend-decl
	stmtCount             // 11: sentinel
)

// Meta field indices for DMetaField distribution.
const (
	metaTitle    = iota // 0
	metaDate           // 1
	metaAuthor         // 2
	metaScope          // 3
	metaQuestion       // 4
	metaDoctrine       // 5
	metaCount          // 6
)

// Node kind indices for DNodeKind distribution.
const (
	nodeComponent = iota // 0
	nodeAnchor           // 1
	nodeSubmap           // 2
	nodeCount            // 3
)

// Node format indices for DNodeFormat distribution.
const (
	fmtBare      = iota // 0: identifier only (no evolution, no block)
	fmtShorthand        // 1: Name : evolution [type] [visibility]
	fmtBlock            // 2: Name [: evolution] { config... }
	fmtCount            // 3
)

// Shorthand variant indices for DShorthandVar distribution.
const (
	shEvoOnly       = iota // 0: evolution only
	shEvoType              // 1: evolution + (type)
	shEvoVis               // 2: evolution + @visibility
	shEvoTypeVis           // 3: evolution + (type) + @visibility
	shCount                // 4
)

// Evolution form indices for DEvoForm distribution.
const (
	evoPositionOnly = iota // 0: just position
	evoMove                // 1: position >> position
	evoInertiaMove         // 2: position ! >> position
	evoCount               // 3
)

// Roman numeral indices for DRoman distribution.
const (
	romanI    = iota // 0
	romanII          // 1
	romanIII         // 2
	romanIV          // 3
	romanCount       // 4
)

// Type value indices for DTypeValue distribution.
const (
	typeBuild     = iota // 0
	typeBuy              // 1
	typeOutsource        // 2
	typeCount            // 3
)

// Link type indices for DLinkType distribution.
const (
	linkArrow     = iota // 0: ->
	linkBidir            // 1: <->
	linkLabeled          // 2: -[text]->
	linkCount            // 3
)

// Signal type indices for DSignalType distribution.
const (
	sigAccelerating   = iota // 0
	sigStagnating            // 1
	sigDeclining             // 2
	sigCoEvolution           // 3
	sigRedQueen              // 4
	sigCommoditization       // 5
	sigNetworkEffects        // 6
	sigEconomiesOfScale      // 7
	sigCount                 // 8
)

// Gameplay type indices for DGameplayType distribution.
const (
	gpILC              = iota // 0
	gpOpenSource              // 1
	gpLandGrab                // 2
	gpEmbraceExtend           // 3
	gpTowerMoat               // 4
	gpFUD                     // 5
	gpStranglerFig            // 6
	gpSignalDistortion        // 7
	gpCount                   // 8
)

// Inertia level indices for DInertiaLevel distribution.
const (
	inertia1    = iota // 0: !
	inertia2           // 1: !!
	inertia3           // 2: !!!
	inertiaCount       // 3
)

// Inertia kind indices for DInertiaKind distribution.
const (
	ikTech       = iota // 0
	ikFinancial         // 1
	ikHuman             // 2
	ikRelational        // 3
	ikSocial            // 4
	ikCount             // 5
)

// Doctrine value indices for DDoctrineValue distribution.
const (
	docHygiene    = iota // 0
	docContext           // 1
	docExcellence        // 2
	docEvolution         // 3
	docCount             // 4
)

// Team type indices for DTeamType distribution.
const (
	teamExplorer    = iota // 0
	teamSettler            // 1
	teamTownPlanner        // 2
	teamPioneer            // 3
	teamVillager           // 4
	teamCount              // 5
)

// Annotation kind indices for DAnnotationKind distribution.
const (
	annoNote    = iota // 0
	annoWarning        // 1
	annoCount          // 2
)

// Asset value indices for DAssetValue distribution.
const (
	assetTech       = iota // 0
	assetFinancial         // 1
	assetHuman             // 2
	assetRelational        // 3
	assetSocial            // 4
	assetCount             // 5
)

// Boolean decision indices (yes/no).
const (
	boolNo  = 0
	boolYes = 1
)

// String tables mapping indices to WTG2 keywords.
var signalTypes = [sigCount]string{
	"accelerating", "stagnating", "declining", "co-evolution",
	"red-queen", "commoditization", "network-effects", "economies-of-scale",
}

var gameplayTypes = [gpCount]string{
	"ILC", "open-source", "land-grab", "embrace-extend",
	"tower-moat", "FUD", "strangler-fig", "signal-distortion",
}

var typeValues = [typeCount]string{"build", "buy", "outsource"}

var doctrineValues = [docCount]string{"hygiene", "context", "excellence", "evolution"}

var teamTypes = [teamCount]string{"explorer", "settler", "town-planner", "pioneer", "villager"}

var inertiaKinds = [ikCount]string{"tech", "financial", "human", "relational", "social"}

var assetValues = [assetCount]string{"tech", "financial", "human", "relational", "social"}

// Lookup functions for encoding.

func signalTypeIndex(s string) int {
	for i, v := range signalTypes {
		if v == s {
			return i
		}
	}
	return 0
}

func gameplayTypeIndex(s string) int {
	for i, v := range gameplayTypes {
		if v == s {
			return i
		}
	}
	return 0
}

func typeValueIndex(s string) int {
	for i, v := range typeValues {
		if v == s {
			return i
		}
	}
	return 0
}

func doctrineValueIndex(s string) int {
	for i, v := range doctrineValues {
		if v == s {
			return i
		}
	}
	return 0
}

func teamTypeIndex(s string) int {
	for i, v := range teamTypes {
		if v == s {
			return i
		}
	}
	return 0
}

func inertiaKindIndex(s string) int {
	for i, v := range inertiaKinds {
		if v == s {
			return i
		}
	}
	return 0
}

func assetValueIndex(s string) int {
	for i, v := range assetValues {
		if v == s {
			return i
		}
	}
	return 0
}

// Static probability distributions (Level 2 - empirical frequencies).
// Frequencies are scaled to sum to ~1000 for good arithmetic coding precision.
var (
	distStmtType = NewDistribution([]uint32{
		100, // meta
		10,  // stages
		300, // node
		30,  // pipeline
		350, // edge
		30,  // group
		60,  // annotation
		40,  // signal
		30,  // gameplay
		20,  // focus
		10,  // legend
	})

	distMetaField = NewDistribution([]uint32{
		200, // title
		180, // date
		200, // author
		150, // scope
		150, // question
		120, // doctrine
	})

	distNodeKind = NewDistribution([]uint32{
		750, // component
		180, // anchor
		70,  // submap
	})

	distNodeFormat = NewDistribution([]uint32{
		100, // bare (no evolution, no block)
		600, // shorthand
		300, // block
	})

	distShorthandVar = NewDistribution([]uint32{
		400, // evo only
		350, // evo + type
		100, // evo + visibility
		150, // evo + type + visibility
	})

	distEvoForm = NewDistribution([]uint32{
		700, // position only
		200, // move (>>)
		100, // inertia + move
	})

	distRoman = NewDistribution([]uint32{
		150, // I
		350, // II
		350, // III
		150, // IV
	})

	distHasDecimal = NewDistribution([]uint32{
		150, // no
		850, // yes
	})

	distDecimalDigit = Uniform(10) // 0-9

	distTypeValue = NewDistribution([]uint32{
		300, // build
		450, // buy
		250, // outsource
	})

	distLinkType = NewDistribution([]uint32{
		800, // ->
		100, // <->
		100, // -[text]->
	})

	distSignalType = NewDistribution([]uint32{
		200, // accelerating
		150, // stagnating
		150, // declining
		100, // co-evolution
		50,  // red-queen
		200, // commoditization
		80,  // network-effects
		70,  // economies-of-scale
	})

	distGameplayType = Uniform(int(gpCount))

	distBool = NewDistribution([]uint32{
		500, // no
		500, // yes
	})

	distHasNext = NewDistribution([]uint32{
		70,  // no (end)
		930, // yes (continue)
	})

	distIdentKnown = NewDistribution([]uint32{
		300, // new
		700, // known
	})

	distInertiaLevel = NewDistribution([]uint32{
		400, // !
		400, // !!
		200, // !!!
	})

	distInertiaKind = Uniform(int(ikCount))

	distDoctrineValue = Uniform(int(docCount))

	distTeamType = Uniform(int(teamCount))

	distAnnotationKind = NewDistribution([]uint32{
		600, // note
		400, // warning
	})

	distAssetValue = Uniform(int(assetCount))

	// Byte-level distribution for all text/identifier encoding.
	// 256 symbols — one per byte value. Covers ASCII and UTF-8 continuation bytes.
	distByte *Distribution
)

func init() {
	freqs := make([]uint32, 256)
	for i := range freqs {
		freqs[i] = 1
	}
	// Boost printable ASCII by typical English/French frequency
	for c := byte('a'); c <= 'z'; c++ {
		freqs[c] = 80
	}
	for c := byte('A'); c <= 'Z'; c++ {
		freqs[c] = 60
	}
	for c := byte('0'); c <= '9'; c++ {
		freqs[c] = 30
	}
	freqs[' '] = 200
	freqs['-'] = 40
	freqs['.'] = 40
	freqs[','] = 40
	freqs['\''] = 40
	freqs[':'] = 20
	freqs['#'] = 15
	freqs['"'] = 15
	// UTF-8 continuation bytes (0x80-0xBF) and lead bytes (0xC0-0xF7)
	for c := 0x80; c <= 0xBF; c++ {
		freqs[c] = 10
	}
	for c := 0xC0; c <= 0xF7; c++ {
		freqs[c] = 5
	}
	distByte = NewDistribution(freqs)
}
