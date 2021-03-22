package plan

import (
	"gonum.org/v1/gonum/graph/simple"
)

// . An element is anything that have coordinates
type Element interface {
	GetCoordinates() []int
}

// a Map is a DirectedGraph with a bunch of anotations
type Map struct {
	Title string
	*simple.DirectedGraph
	Annotations          []*Annotation
	AnnotationsPlacement [2]int
}
