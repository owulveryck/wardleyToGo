package wardley

import "gonum.org/v1/gonum/graph/simple"

// . An element is anything that have coordinates
type Element interface {
	GetCoordinates() [2]int
}

// a Map is a DirectedGraph with a bunch of anotations
type Map struct {
	*simple.DirectedGraph
	Annotations []Element
}
