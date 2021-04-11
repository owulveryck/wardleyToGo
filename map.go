package wardleyToGo

import (
	"errors"
	"image"
	"image/draw"

	"gonum.org/v1/gonum/graph/simple"
)

// a Map is a directed graph whose components knows their own position wrt to an anchor.
// The anchor is the point A of a rectangle as defined by
//  A := image.Point{}
//  image.Rectangle{A, Pt(100, 100)}
type Map struct {
	id    int64
	Title string
	// Canvas is the function that will draw the initial map
	// allowing the placement of the axis, legend and so on
	Canvas               draw.Drawer
	Annotations          []*Annotation
	AnnotationsPlacement image.Point
	area                 image.Rectangle
	*simple.DirectedGraph
}

// NewMap with initial area of 100x100
func NewMap(id int64) *Map {
	return &Map{
		id:            id,
		area:          image.Rect(0, 0, 100, 100),
		DirectedGraph: simple.NewDirectedGraph(),
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

// Draw aligns r.Min in dst with sp in src and then replaces the
// rectangle r in dst with the result of drawing src on dst.
// If the Components and Collaboration elemts of the maps are draw.Drawer, their methods
// are called accordingly
func (m *Map) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	if m.Canvas != nil {
		m.Canvas.Draw(dst, r, src, sp)
	}
	nodes := m.Nodes()
	for nodes.Next() {
		if n, ok := nodes.Node().(draw.Drawer); ok {
			n.Draw(dst, r, src, sp)
		}
	}
	edges := m.Edges()
	for edges.Next() {
		if e, ok := edges.Edge().(draw.Drawer); ok {
			e.Draw(dst, r, src, sp)
		}
	}
}

// SVG representation, class is subMapElement and element
/*
func (m *Map) SVG(s *svg.SVG, bounds image.Rectangle) {
	coords := utils.CalcCoords(m.GetPosition(), bounds)
	s.Gid(strconv.FormatInt(m.id, 10))
	s.Translate(coords.X, coords.Y)
	s.Text(10, 10, m.Title)
	s.Circle(0, 0, 5, `stroke-width="1"`, `stroke="black"`, `fill="black"`, `class="subMapElement, element"`)
	s.Gend()
	s.Gend()
}
*/

// AddComponent add e to the graph. It returns an error if e is out-of-bounds,
// meaning its coordinates are less than 0 or more that 100
func (m *Map) AddComponent(e Component) error {
	if !e.GetPosition().In(image.Rect(0, 0, 100, 100)) {
		return errors.New("component out of bounds")
	}
	m.DirectedGraph.AddNode(e)
	return nil
}

func (m *Map) SetCollaboration(e Collaboration) error {
	m.DirectedGraph.SetEdge(e)
	return nil
}
