package wardleyToGo

import (
	"image"
	"strconv"

	svg "github.com/ajstarks/svgo"
	"gonum.org/v1/gonum/graph/simple"
)

// . An element is anything that have coordinates
type Element interface {
	// . GetPosition of the element wrt a 100x100 map
	GetPosition() image.Point
}

type Area interface {
	GetArea() image.Rectangle
}

// a Map is a DirectedGraph with a bunch of anotations
type Map struct {
	id    int64
	Title string
	area  image.Rectangle
	*simple.DirectedGraph
	Annotations          []*Annotation
	AnnotationsPlacement image.Point
}

// a Map can be a node of a map
func (m *Map) ID() int64 {
	return m.id
}

func (m *Map) GetPosition() image.Point {
	return image.Pt((m.area.Max.X-m.area.Min.X)/2, (m.area.Max.Y-m.area.Min.Y)/2)
}

// SVG representation, class is subMapElement and element
func (m *Map) SVG(s *svg.SVG, bounds image.Rectangle) {
	coords := calcCoords(m.GetPosition(), bounds)
	s.Gid(strconv.FormatInt(m.id, 10))
	s.Translate(coords.X, coords.Y)
	s.Text(10, 10, m.Title)
	s.Circle(0, 0, 5, `stroke-width="1"`, `stroke="black"`, `fill="black"`, `class="subMapElement, element"`)
	s.Gend()
	s.Gend()
}

// SVGer is any object that can represent itself on a map
type SVGer interface {
	// SVG is a method that represent the object on the svg mag with coordinates relatives to the bounds
	SVG(s *svg.SVG, bounds image.Rectangle)
}
