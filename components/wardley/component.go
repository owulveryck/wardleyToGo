package wardley

import (
	"image"
	"strconv"

	svg "github.com/ajstarks/svgo"
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components"
)

const (
	// This is a RegularComponent
	RegularComponent wardleyToGo.ComponentType = iota | components.Wardley
	// BuildComponent ...
	BuildComponent
	// Off the shelf element
	BuyComponent
	// OutsourceComponent ...
	OutsourceComponent
	// DataProductComponent ...
	DataProductComponent
)

// A Component is an element of the map
type Component struct {
	id            int64
	Placement     image.Point // The placement of the component on a rectangle 100x100
	Label         string
	LabePlacement image.Point // LabelPlacement is relative to the placement
	Type          wardleyToGo.ComponentType
}

// NewComponent with the corresponding id and default UndefinedCoords
func NewComponent(id int64) *Component {
	return &Component{
		id:        id,
		Placement: image.Pt(components.UndefinedCoord, components.UndefinedCoord),
	}
}

// Component fulfils the graph.Node interface
func (c *Component) ID() int64 {
	return c.id
}

// SVG is a representation of the component
func (c *Component) SVG(s *svg.SVG, bounds image.Rectangle) {
	coords := components.CalcCoords(c.Placement, bounds)
	s.Gid(strconv.FormatInt(c.id, 10))
	s.Translate(coords.X, coords.Y)
	s.Text(c.LabePlacement.X, c.Placement.Y, c.Label)
	switch c.Type {
	case BuildComponent:
		s.Circle(0, 0, 20, `fill="#D6D6D6"`, `stroke="#000000"`, `class="element, buildComponent"`)
	case BuyComponent:
		s.Circle(0, 0, 20, `fill="#AAA5A9"`, `stroke="#D6D6D6"`, `class="element, buyComponent"`)
	case OutsourceComponent:
		s.Circle(0, 0, 20, `fill="#444444"`, `stroke="#444444"`, `class="element, outsourceComponent"`)
	case DataProductComponent:
		s.Circle(0, 0, 14, `fill="rgb(246,72,22)"`, `class="element, dataProductComponent"`)
	}
	s.Circle(0, 0, 5, `stroke-width="1"`, `stroke="black"`, `fill="white"`, `class="element, component"`)
	s.Gend()
	s.Gend()
}

func (c *Component) String() string {
	return c.Label
}

func (c *Component) GetPosition() image.Point {
	return c.Placement
}

type EvolvedComponent struct {
	*Component
}

func (e *EvolvedComponent) ID() int64 {
	return e.id
}

func NewEvolvedComponent(id int64) *EvolvedComponent {
	return &EvolvedComponent{NewComponent(id)}
}

func (e *EvolvedComponent) SVG(s *svg.SVG, bounds image.Rectangle) {
	coords := components.CalcCoords(e.Placement, bounds)
	s.Gid(strconv.FormatInt(e.id, 10))
	s.Translate(coords.X, coords.Y)
	s.Text(e.LabePlacement.X, e.Placement.Y, e.Label, `fill="red"`)
	switch e.Type {
	case BuildComponent:
		s.Circle(0, 0, 20, `fill="#D6D6D6"`, `stroke="#000000"`, `class="element, buildComponent"`)
	case BuyComponent:
		s.Circle(0, 0, 20, `fill="#AAA5A9"`, `stroke="#D6D6D6"`, `class="element, buyComponent"`)
	case OutsourceComponent:
		s.Circle(0, 0, 20, `fill="#444444"`, `stroke="#444444"`, `class="element, outsourceComponent"`)
	case DataProductComponent:
		s.Circle(0, 0, 14, `fill="rgb(246,72,22)"`, `class="element, dataProductComponent"`)
	}
	s.Circle(0, 0, 5, `stroke-width="1"`, `stroke="red"`, `fill="white"`, `class="element, component"`)
	s.Gend()
	s.Gend()
}

// GetCoordinates fulfils the Element interface
func (e *EvolvedComponent) GetPosition() image.Point {
	return e.Component.GetPosition()
}

func (e *EvolvedComponent) String() string {
	return "[evolved]" + e.Label
}