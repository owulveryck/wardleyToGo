package graph

// SimpleEdge is a basic directed edge between two nodes.
type SimpleEdge struct {
	F, T Node
}

// From returns the source node.
func (e SimpleEdge) From() Node { return e.F }

// To returns the destination node.
func (e SimpleEdge) To() Node { return e.T }
