package owm2wtg2

// OWMDocument is the intermediate representation of a parsed OWM file.
type OWMDocument struct {
	Title     string
	Evolution [4]string // custom stage labels (from "evolution A->B->C->D")
	Style     string    // "wardley", "handwritten", etc.

	Components []*OWMComponent
	Evolves    []*OWMEvolve
	Anchors    []*OWMAnchor
	Edges      []*OWMEdge
	Pipelines  []*OWMPipeline
	Notes      []*OWMNote
	Submaps    []*OWMSubmap
	URLs       map[string]string // urlName -> URL
	Areas      []*OWMArea        // pioneers, settlers, townplanners
	Signals    []*OWMSignal      // accelerator, deaccelerator
}

// OWMComponent represents a component declaration.
type OWMComponent struct {
	Name       string
	Visibility float64 // 0.0–1.0
	Maturity   float64 // 0.0–1.0
	Type       string  // "build", "buy", "outsource", "market", "dataProduct", ""
	Inertia    bool
	LabelX     int
	LabelY     int
	HasLabel   bool
}

// OWMEvolve represents an evolve directive.
type OWMEvolve struct {
	Name     string
	NewName  string  // non-empty for "evolve Name->NewName maturity"
	Maturity float64 // target maturity 0.0–1.0
	Type     string  // "build", "buy", "outsource", ""
	LabelX   int
	LabelY   int
	HasLabel bool
}

// OWMAnchor represents an anchor declaration.
type OWMAnchor struct {
	Name       string
	Visibility float64
	Maturity   float64
}

// OWMEdge represents a link between two components.
type OWMEdge struct {
	From     string
	To       string
	Label    string // label text for +'label'> flow
	FlowType string // "regular", "bidirectional", "past", "future", "labeled"
}

// OWMPipeline represents a pipeline declaration.
type OWMPipeline struct {
	Name     string
	StartMat float64 // 0.0–1.0
	EndMat   float64 // 0.0–1.0
	HasRange bool
}

// OWMNote represents a positioned note.
type OWMNote struct {
	Text       string
	Visibility float64
	Maturity   float64
}

// OWMSubmap represents a submap component.
type OWMSubmap struct {
	Name       string
	Visibility float64
	Maturity   float64
	URLName    string // reference to URLs map
}

// OWMArea represents a pioneers/settlers/townplanners area.
type OWMArea struct {
	Kind string // "pioneers", "settlers", "townplanners"
	Vis1 float64
	Mat1 float64
	Vis2 float64
	Mat2 float64
}

// OWMSignal represents an accelerator or deaccelerator.
type OWMSignal struct {
	Kind       string // "accelerator", "deaccelerator"
	Name       string
	Visibility float64
	Maturity   float64
}
