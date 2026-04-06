package wardleyToGo

import (
	"image"
)

// A Component is a node of a graph that have coordinates.
// A Component can represent itself on a 100x100 map
type Component interface {
	// GetPosition of the element wrt a 100x100 map
	GetPosition() image.Point
	// ID returns a unique identifier for this component
	ID() int64
}

// An area is anything that covers a rectangle area on a map
type Area interface {
	// GetArea should be expressed wrt a 100x100 map
	GetArea() image.Rectangle
	// ID returns a unique identifier for this area
	ID() int64
}
