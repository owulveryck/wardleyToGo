package wardleyToGo

// A Collaboration is an edge between two components with a certain type
type Collaboration interface {
	From() Component
	To() Component
	GetType() EdgeType
}
