package wardleyToGo

import (
	"errors"
	"image"
	"strconv"

	svg "github.com/ajstarks/svgo"
	"github.com/owulveryck/wardleyToGo/components"
	"gonum.org/v1/gonum/graph/simple"
)

// a Map is a directed graph whose components knows their own position wrt to an anchor.
// The anchor is the point A of a rectangle as defined by
//  A := image.Point{}
//  image.Rectangle{A, Pt(100, 100)}
type Map struct {
	id                   int64
	Title                string
	Annotations          []*Annotation
	AnnotationsPlacement image.Point
	area                 image.Rectangle
	g                    *simple.DirectedGraph
}

// NewMap with initial area of 100x100
func NewMap(id int64) *Map {
	return &Map{
		id:   id,
		area: image.Rect(0, 0, 100, 100),
		g:    simple.NewDirectedGraph(),
	}
}

// a Map fulfills the graph.Node interface; thererfore if can be part of a graph of maps
func (m *Map) ID() int64 {
	return m.id
}

// GetPosition fulfills the componnts.Component interface. Therefore a map can be a component of another map.
// This allows doing submaping.
// The position is the  center of the area of the map
func (m *Map) GetPosition() image.Point {
	return image.Pt((m.area.Max.X-m.area.Min.X)/2, (m.area.Max.Y-m.area.Min.Y)/2)
}
func (m *Map) GetArea() image.Rectangle {
	return m.area
}

// SVG representation, class is subMapElement and element
func (m *Map) SVG(s *svg.SVG, bounds image.Rectangle) {
	coords := components.CalcCoords(m.GetPosition(), bounds)
	s.Gid(strconv.FormatInt(m.id, 10))
	s.Translate(coords.X, coords.Y)
	s.Text(10, 10, m.Title)
	s.Circle(0, 0, 5, `stroke-width="1"`, `stroke="black"`, `fill="black"`, `class="subMapElement, element"`)
	s.Gend()
	s.Gend()
}

// AddComponent add e to the graph. It returns an error if e is out-of-bounds,
// meaning its coordinates are less than 0 or more that 100
func (m *Map) AddComponent(e components.Component) error {
	if !e.GetPosition().In(image.Rect(0, 0, 100, 100)) {
		return errors.New("component out of bounds")
	}
	m.g.AddNode(e)
	return nil
}

func (m *Map) SetEdge(e *Edge) error {
	m.g.SetEdge(e)
	return nil
}
