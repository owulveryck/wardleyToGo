package wardleyToGo

import (
	"gonum.org/v1/gonum/graph"
)

// A Collaboration is an edge between two components with a certain type
type Collaboration interface {
	graph.Edge
	GetType() EdgeType
}
