package components

import (
	"github.com/owulveryck/wardleyToGo"
	"gonum.org/v1/gonum/graph"
)

// Collaboration is a generic collaboration type
type Collaboration struct {
	T         graph.Node
	F         graph.Node
	EdgeLabel string
	EdgeType  wardleyToGo.EdgeType
}

func (e Collaboration) From() graph.Node {
	return e.F
}

func (e Collaboration) ReversedEdge() graph.Edge {
	return Collaboration{
		F:         e.T,
		T:         e.F,
		EdgeLabel: e.EdgeLabel,
	}
}

func (e Collaboration) To() graph.Node {
	return e.T
}

func (e Collaboration) GetType() wardleyToGo.EdgeType {
	return e.EdgeType
}
