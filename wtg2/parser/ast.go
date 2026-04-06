package parser

// Document is the intermediate representation of a parsed WTG2 file.
type Document struct {
	Title    string
	Date     string
	Author   string
	Scope    string
	Question string
	Stages   [4]string // Custom stage labels; default: "I", "II", "III", "IV"

	Nodes       []*NodeDecl
	Pipelines   []*PipelineDecl
	Edges       []*EdgeDecl
	Groups      []*GroupDecl
	Annotations []*AnnotationDecl
	Signals     []*SignalDecl
}

// NodeKind distinguishes anchors, components, and submaps.
type NodeKind int

const (
	KindComponent NodeKind = iota
	KindAnchor
	KindSubmap
)

// NodeDecl represents a node declaration in the source.
type NodeDecl struct {
	Name       string
	Kind       NodeKind
	Evolution  string  // Raw position string, e.g. "III.5"
	EvolvedTo  string  // Target position if ">>" present
	Inertia    int     // 0-3 (count of '!')
	Type       string  // "build", "buy", "outsource", or ""
	Color      string  // "#3498DB" or ""
	Visibility float64 // Explicit @visibility override, -1 if unset
	Note       string
}

// PipelineDecl represents a pipeline block.
type PipelineDecl struct {
	Name    string
	Members []*PipelineMemberDecl
}

// PipelineMemberDecl is a member within a pipeline.
type PipelineMemberDecl struct {
	Name     string
	Position string // Raw position string
}

// EdgeDecl represents a single edge (link) between two nodes.
type EdgeDecl struct {
	From          string // Node name, or "Pipeline:Member"
	To            string
	Label         string // Annotation text from -[text]->
	Bidirectional bool
}

// GroupDecl represents a visual group of components.
type GroupDecl struct {
	Name    string
	Members []string
}

// AnnotationDecl represents a note or warning on a node.
type AnnotationDecl struct {
	Kind   string // "note" or "warning"
	Text   string
	Target string // Node name
}

// SignalDecl represents a market signal on a node.
type SignalDecl struct {
	Type   string // "accelerating", "stagnating", "declining"
	Target string // Node name
}
