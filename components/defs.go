package components

import (
	"image"

	svg "github.com/ajstarks/svgo"
	"gonum.org/v1/gonum/graph"
)

//  A Component is anything that have coordinates
type Component interface {
	// . GetPosition of the element wrt a 100x100 map
	GetPosition() image.Point
	graph.Node
}

// An area is anything that covers a rectangle area
type Area interface {
	// GetArea should be expressed wrt a 100x100 map
	GetArea() image.Rectangle
	graph.Node
}

// SVGer is any object that can represent itself on a map
type SVGer interface {
	// SVG is a method that represent the object on the svg mag with coordinates relatives to the bounds
	SVG(s *svg.SVG, bounds image.Rectangle)
}
