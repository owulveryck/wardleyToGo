package wardleyToGo

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"sort"
	"strings"

	"github.com/owulveryck/wardleyToGo/internal/graph"
)

// collabEdge is an internal adapter that wraps a Collaboration to satisfy graph.Edge.
type collabEdge struct {
	c Collaboration
}

func (e collabEdge) From() graph.Node { return e.c.From().(graph.Node) }
func (e collabEdge) To() graph.Node   { return e.c.To().(graph.Node) }

// a Map is a directed graph whose components knows their own position wrt to an anchor.
// The anchor is the point A of a rectangle as defined by
//
//	A := image.Point{}
//	image.Rectangle{A, Pt(100, 100)}
type Map struct {
	id    int64
	Title string
	// Canvas is the function that will draw the initial map
	// allowing the placement of the axis, legend and so on
	Canvas               draw.Drawer
	Annotations          []*Annotation
	AnnotationsPlacement image.Point
	area                 image.Rectangle
	g                    *graph.DirectedGraph
	collabs              map[int64]map[int64]Collaboration
}

func (m *Map) String() string {
	var b strings.Builder
	b.WriteString("map {\n")
	for _, n := range m.Components() {
		if a, ok := n.(Area); ok {
			fmt.Fprintf(&b, "\t%v '%v' [%v,%v,%v,%v];\n", a.ID(), a,
				a.GetArea().Min.X, a.GetArea().Min.Y,
				a.GetArea().Max.X, a.GetArea().Max.Y)
		} else {
			fmt.Fprintf(&b, "\t%v '%v' [%v,%v];\n", n.ID(), n, n.GetPosition().X, n.GetPosition().Y)
		}
	}
	b.WriteString("\n")
	for _, e := range m.Collaborations() {
		fmt.Fprintf(&b, "\t%v -> %v [%v];\n", e.From().ID(), e.To().ID(), e.GetType())
	}
	b.WriteString("}\n")
	return b.String()
}

// NewMap with initial area of 100x100
func NewMap(id int64) *Map {
	return &Map{
		id:      id,
		area:    image.Rect(0, 0, 100, 100),
		g:       graph.NewDirectedGraph(),
		collabs: make(map[int64]map[int64]Collaboration),
	}
}

// a Map fulfills the graph.Node interface; therefore it can be part of a graph of maps
func (m *Map) ID() int64 {
	return m.id
}

// GetPosition fulfills the components.Component interface. Therefore a map can be a component of another map.
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
// If the Components and Collaboration elements of the maps are draw.Drawer, their methods
// are called accordingly
func (m *Map) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	if m.Canvas != nil {
		m.Canvas.Draw(dst, r, src, sp)
	}
	// Draw edges first
	for _, c := range m.Collaborations() {
		if e, ok := c.(draw.Drawer); ok {
			e.Draw(dst, r, src, sp)
		}
	}
	for _, n := range m.Components() {
		if d, ok := n.(draw.Drawer); ok {
			d.Draw(dst, r, src, sp)
		}
	}
}

// AddComponent add e to the graph. It returns an error if e is out-of-bounds,
// meaning its coordinates are less than 0 or more that 100.
// Area-type components skip bounds checking since their position is derived from their area.
func (m *Map) AddComponent(e Component) error {
	if _, isArea := e.(Area); !isArea {
		if !e.GetPosition().In(image.Rect(0, 0, 100, 100)) {
			return errors.New("component out of bounds")
		}
	}
	m.g.AddNode(e)
	return nil
}

func (m *Map) SetCollaboration(e Collaboration) error {
	fid := e.From().ID()
	tid := e.To().ID()
	if fid == tid {
		return fmt.Errorf("self-loop not allowed: node %d", fid)
	}
	// Store the Collaboration in our map
	if m.collabs[fid] == nil {
		m.collabs[fid] = make(map[int64]Collaboration)
	}
	m.collabs[fid][tid] = e
	// Store an adapter edge in the internal graph for traversal
	m.g.SetEdge(collabEdge{c: e})
	return nil
}

// RemoveEdge removes the edge from uid to vid.
func (m *Map) RemoveEdge(uid, vid int64) {
	m.g.RemoveEdge(uid, vid)
	if m.collabs[uid] != nil {
		delete(m.collabs[uid], vid)
	}
}

// Components returns all components in the map.
func (m *Map) Components() []Component {
	nodes := m.g.Nodes()
	result := make([]Component, 0, nodes.Len())
	for nodes.Next() {
		if c, ok := nodes.Node().(Component); ok {
			result = append(result, c)
		}
	}
	return result
}

// Collaborations returns all collaborations in the map in deterministic order,
// sorted by (from ID, to ID).
func (m *Map) Collaborations() []Collaboration {
	fromIDs := make([]int64, 0, len(m.collabs))
	for fid := range m.collabs {
		fromIDs = append(fromIDs, fid)
	}
	sort.Slice(fromIDs, func(i, j int) bool { return fromIDs[i] < fromIDs[j] })

	n := 0
	for _, fid := range fromIDs {
		n += len(m.collabs[fid])
	}
	result := make([]Collaboration, 0, n)
	for _, fid := range fromIDs {
		toIDs := make([]int64, 0, len(m.collabs[fid]))
		for tid := range m.collabs[fid] {
			toIDs = append(toIDs, tid)
		}
		sort.Slice(toIDs, func(i, j int) bool { return toIDs[i] < toIDs[j] })
		for _, tid := range toIDs {
			result = append(result, m.collabs[fid][tid])
		}
	}
	return result
}

// From returns all components reachable from the component with the given id.
func (m *Map) From(id int64) []Component {
	nodes := m.g.From(id)
	result := make([]Component, 0, nodes.Len())
	for nodes.Next() {
		if c, ok := nodes.Node().(Component); ok {
			result = append(result, c)
		}
	}
	return result
}

// To returns all components that have an edge to the component with the given id.
func (m *Map) To(id int64) []Component {
	nodes := m.g.To(id)
	result := make([]Component, 0, nodes.Len())
	for nodes.Next() {
		if c, ok := nodes.Node().(Component); ok {
			result = append(result, c)
		}
	}
	return result
}

// Node returns the component with the given id, or nil if not found.
func (m *Map) Node(id int64) Component {
	n := m.g.Node(id)
	if n == nil {
		return nil
	}
	if c, ok := n.(Component); ok {
		return c
	}
	return nil
}

// HasEdgeFromTo reports whether an edge exists from uid to vid.
func (m *Map) HasEdgeFromTo(uid, vid int64) bool {
	return m.g.HasEdgeFromTo(uid, vid)
}

// Chainer is a component that is part of a value chain
type Chainer interface {
	// GetAbsoluteVisibility returns the visibility of the component as seen from the anchor
	GetAbsoluteVisibility() int
}
