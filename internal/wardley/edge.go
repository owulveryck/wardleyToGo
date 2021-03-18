package wardley

type EdgeType int

const (
	RegularEdge EdgeType = iota
	EvolvedComponentEdge
	EvolvedEdge
)
