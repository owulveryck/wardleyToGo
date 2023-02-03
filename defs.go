package wardleyToGo

import (
	"image"

	"gonum.org/v1/gonum/graph"
)

// A Component is a node of a graph that have coordinates.
// A Component can represent iself on a 100x100 map
type Component interface {
	// GetPosition of the element wrt a 100x100 map
	GetPosition() image.Point
	graph.Node
}

// An area is anything that covers a rectangle area on a map
type Area interface {
	// GetArea should be expressed wrt a 100x100 map
	GetArea() image.Rectangle
	graph.Node
}

type Positioner interface {
	SetPosition(image.Point)
}
