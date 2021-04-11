package wardley

import (
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components"
	"gonum.org/v1/gonum/graph"
)

const (
	RegularEdge wardleyToGo.EdgeType = iota | wardleyToGo.EdgeType(components.Wardley)
	EvolvedComponentEdge
	EvolvedEdge
)

type Collaboration struct {
	F, T  wardleyToGo.Component
	Label string
	Type  wardleyToGo.EdgeType
}

// From returns the from node of the edge.
func (c *Collaboration) From() graph.Node {
	return c.F
}

// To returns the to node of the edge.
func (c *Collaboration) To() graph.Node {
	return c.T
}

// ReversedEdge returns the edge reversal of the receiver
// if a reversal is valid for the data type.
// When a reversal is valid an edge of the same type as
// the receiver with nodes of the receiver swapped should
// be returned, otherwise the receiver should be returned
// unaltered.
func (c *Collaboration) ReversedEdge() graph.Edge {
	return &Collaboration{
		F:     c.T,
		T:     c.F,
		Label: c.Label,
		Type:  c.Type,
	}
}

func (c *Collaboration) GetType() wardleyToGo.EdgeType {
	return c.Type
}
