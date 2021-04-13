package wardley

import (
	"image"
	"image/color"
	"image/draw"
	"strconv"

	svg "github.com/ajstarks/svgo"
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/internal/drawing"
	"github.com/owulveryck/wardleyToGo/internal/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	DefaultComponentRenderingLayer int = 10
)

// A Component is an element of the map
type Component struct {
	id             int64
	Placement      image.Point // The placement of the component on a rectangle 100x100
	Label          string
	LabelPlacement image.Point // LabelPlacement is relative to the placement
	Type           wardleyToGo.ComponentType
	RenderingLayer int //The position of the element on the picture
}

// NewComponent with the corresponding id and default UndefinedCoords
func NewComponent(id int64) *Component {
	return &Component{
		id:             id,
		Placement:      image.Pt(components.UndefinedCoord, components.UndefinedCoord),
		LabelPlacement: image.Pt(components.UndefinedCoord, components.UndefinedCoord),
		RenderingLayer: 10,
	}
}

func (c *Component) GetLayer() int {
	return c.RenderingLayer
}

// Component fulfils the graph.Node interface
func (c *Component) ID() int64 {
	return c.id
}

func drawCircle(dst draw.Image, r int, p image.Point, stroke, fill color.Color) {

}

// Draw aligns r.Min in dst with sp in src and then replaces the
// rectangle r in dst with the result of drawing src on dst.
func (c *Component) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	placement := utils.CalcCoords(c.Placement, r)
	//coords := components.CalcCoords(c.Placement, r)
	labelP := c.LabelPlacement
	if labelP.X == components.UndefinedCoord {
		labelP.X = 18
	}
	if labelP.Y == components.UndefinedCoord {
		labelP.Y = 10
	}
	labelP = labelP.Add(placement)
	// First create the circle with a correct resolution
	switch c.Type {
	case BuildComponent:
		drawing.DrawCircle(dst, 10, placement, color.Black, color.RGBA{0xd6, 0xd6, 0xd6, 0xff})
	case BuyComponent:
		drawing.DrawCircle(dst, 10, placement, color.RGBA{0xAA, 0xA5, 0xa9, 0xff}, color.RGBA{0xd6, 0xd6, 0xd6, 0xff})
	case OutsourceComponent:
		drawing.DrawCircle(dst, 10, placement, color.RGBA{0x44, 0x44, 0x44, 0xff}, color.RGBA{0x44, 0x44, 0x44, 0xff})
	case DataProductComponent:
		drawing.DrawCircle(dst, 7, placement, color.RGBA{0x44, 0x44, 0x44, 0xff}, color.RGBA{246, 72, 22, 0xff})
	}
	drawing.DrawCircle(dst, 5, placement, color.Black, color.White)

	// Create the circle with the correct
	dot := fixed.P(labelP.X, labelP.Y)
	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}),
		Face: basicfont.Face7x13,
		Dot:  dot,
	}
	d.DrawString(c.Label)
}

// SVGDraw is a representation of the component
func (c *Component) SVGDraw(s *svg.SVG, bounds image.Rectangle) {
	coords := components.CalcCoords(c.Placement, bounds)
	labelP := c.LabelPlacement
	if labelP.X == components.UndefinedCoord {
		labelP.X = 10
	}
	if labelP.Y == components.UndefinedCoord {
		labelP.Y = 10
	}
	s.Gid(strconv.FormatInt(c.id, 10))
	s.Translate(coords.X, coords.Y)
	s.Text(labelP.X, labelP.Y, c.Label)
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
	c := NewComponent(id)
	return &EvolvedComponent{c}
}

func (e *EvolvedComponent) SVGDraw(s *svg.SVG, bounds image.Rectangle) {
	coords := components.CalcCoords(e.Placement, bounds)
	labelP := e.LabelPlacement
	if labelP.X == components.UndefinedCoord {
		labelP.X = 10
	}
	if labelP.Y == components.UndefinedCoord {
		labelP.Y = 10
	}
	s.Gid(strconv.FormatInt(e.id, 10))
	s.Translate(coords.X, coords.Y)
	s.Text(labelP.X, labelP.Y, e.Label, `fill="red"`)
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

// Draw aligns r.Min in dst with sp in src and then replaces the
// rectangle r in dst with the result of drawing src on dst.
func (c *EvolvedComponent) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	placement := utils.CalcCoords(c.Placement, r)
	//coords := components.CalcCoords(c.Placement, r)
	labelP := c.LabelPlacement
	if labelP.X == components.UndefinedCoord {
		labelP.X = 18
	}
	if labelP.Y == components.UndefinedCoord {
		labelP.Y = 10
	}
	labelP = labelP.Add(placement)
	// First create the circle with a correct resolution
	switch c.Type {
	case BuildComponent:
		drawing.DrawCircle(dst, 10, placement, color.Black, color.RGBA{0xd6, 0xd6, 0xd6, 0xff})
	case BuyComponent:
		drawing.DrawCircle(dst, 10, placement, color.RGBA{0xAA, 0xA5, 0xa9, 0xff}, color.RGBA{0xd6, 0xd6, 0xd6, 0xff})
	case OutsourceComponent:
		drawing.DrawCircle(dst, 10, placement, color.RGBA{0x44, 0x44, 0x44, 0xff}, color.RGBA{0x44, 0x44, 0x44, 0xff})
	case DataProductComponent:
		drawing.DrawCircle(dst, 7, placement, color.RGBA{0x44, 0x44, 0x44, 0xff}, color.RGBA{246, 72, 22, 0xff})
	}
	drawing.DrawCircle(dst, 5, placement, color.RGBA{0xff, 0, 0, 0xff}, color.White)

	// Create the circle with the correct
	dot := fixed.P(labelP.X, labelP.Y)
	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xFF}),
		Face: basicfont.Face7x13,
		Dot:  dot,
	}
	d.DrawString(c.Label)
}

// GetCoordinates fulfils the Element interface
func (e *EvolvedComponent) GetPosition() image.Point {
	return e.Component.GetPosition()
}

func (e *EvolvedComponent) String() string {
	return "[evolved]" + e.Label
}
