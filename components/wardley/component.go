package wardley

import (
	"encoding/xml"
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	dotencoding "gonum.org/v1/gonum/graph/encoding"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/internal/drawing"
	"github.com/owulveryck/wardleyToGo/internal/svg"
	"github.com/owulveryck/wardleyToGo/internal/utils"
)

const (
	DefaultComponentRenderingLayer int = 10
)

// A Component is an element of the map
type Component struct {
	id                 int64
	Placement          image.Point // The placement of the component on a rectangle 100x100
	Label              string
	LabelPlacement     image.Point // LabelPlacement is relative to the placement
	Type               wardleyToGo.ComponentType
	RenderingLayer     int //The position of the element on the picture
	Configured         bool
	EvolutionPos       int
	Color              color.Color
	AbsoluteVisibility int
}

// GetAbsoluteVisibility returns the visibility of the component as seen from the anchor
func (c *Component) GetAbsoluteVisibility() int {
	return c.AbsoluteVisibility
}

func (c *Component) Attributes() []dotencoding.Attribute {
	return []dotencoding.Attribute{
		{
			Key:   "label",
			Value: c.Label,
		},
	}
}

// NewComponent with the corresponding id and default UndefinedCoords
func NewComponent(id int64) *Component {
	return &Component{
		id:             id,
		Placement:      image.Pt(components.UndefinedCoord, components.UndefinedCoord),
		LabelPlacement: image.Pt(components.UndefinedCoord, components.UndefinedCoord),
		RenderingLayer: 10,
		Color:          color.RGBA{R: 0, G: 0, B: 0, A: 255}, // black
	}
}

func (c *Component) GetLayer() int {
	return c.RenderingLayer
}

// Component fulfils the graph.Node interface
func (c *Component) ID() int64 {
	return c.id
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

func (c *Component) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	switch c.Type {
	case PipelineComponent:
		return c.marshalSVGPipeline(e, canvas, svg.Color{c.Color})
	default:
		return c.marshalSVG(e, canvas, svg.Color{c.Color})
	}
}

func (c *Component) marshalSVGPipeline(e *xml.Encoder, canvas image.Rectangle, col svg.Color) error {
	coords := components.CalcCoords(c.Placement, canvas)
	labelP := c.LabelPlacement
	if labelP.X == components.UndefinedCoord {
		labelP.X = 10
	}
	if labelP.Y == components.UndefinedCoord {
		labelP.Y = 10
	}
	fillColor := svg.White
	r, g, b, a := c.Color.RGBA()
	if r != 0 || g != 0 || b != 0 || a != 65535 {
		fillColor = svg.Color{col}
	}
	components := make([]interface{}, 0)
	components = append(components, svg.Rectangle{
		R: image.Rectangle{
			Min: image.Point{-5, -5},
			Max: image.Point{5, 5},
		},
		StrokeWidth: "3",
		Stroke:      col,
		Fill:        fillColor,
	})

	components = append(components, svg.Text{
		P:    labelP,
		Text: []byte(c.Label),
		Fill: col,
	})
	return e.Encode(svg.Transform{
		Translate:  coords,
		Components: components,
		//Classes:    []string{fmt.Sprintf("visibility%v", c.AbsoluteVisibility)},
	})
}
func (c *Component) marshalSVG(e *xml.Encoder, canvas image.Rectangle, col svg.Color) error {
	coords := components.CalcCoords(c.Placement, canvas)
	labelP := c.LabelPlacement
	if labelP.X == components.UndefinedCoord {
		labelP.X = 10
	}
	if labelP.Y == components.UndefinedCoord {
		labelP.Y = 10
	}
	fillColor := svg.White
	r, g, b, a := c.Color.RGBA()
	if r != 0 || g != 0 || b != 0 || a != 65535 {
		fillColor = svg.Color{col}
	}
	baseCircle := svg.Circle{
		R:           5,
		StrokeWidth: "1",
		Stroke:      col,
		Fill:        fillColor,
	}
	components := make([]interface{}, 0)
	switch c.Type {
	case BuildComponent:
		components = append(components, svg.Circle{
			R:           20,
			StrokeWidth: "1",
			Stroke:      svg.Black,
			Fill:        svg.Color{color.RGBA{0xd6, 0xd6, 0xd6, 0xff}},
		})
	case BuyComponent:
		components = append(components, svg.Circle{
			R:           20,
			StrokeWidth: "1",
			Fill:        svg.Color{color.RGBA{0xaa, 0xa5, 0xa9, 0xff}},
			Stroke:      svg.Color{color.RGBA{0xd6, 0xd6, 0xd6, 0xff}},
		})
	case OutsourceComponent:
		components = append(components, svg.Circle{
			R:           20,
			StrokeWidth: "1",
			Fill:        svg.Color{color.RGBA{0x44, 0x44, 0x44, 0xff}},
			Stroke:      svg.Color{color.RGBA{0x44, 0x44, 0x44, 0xff}},
		})
	case DataProductComponent:
		components = append(components, svg.Circle{
			R:           14,
			StrokeWidth: "1",
			Fill:        svg.Color{color.RGBA{246, 72, 22, 0xff}},
		})
	}
	components = append(components, baseCircle)
	components = append(components, svg.Text{
		P:    labelP,
		Text: []byte(c.Label),
		Fill: col,
	})
	return e.Encode(svg.Transform{
		Translate:  coords,
		Components: components,
		//Classes:    []string{fmt.Sprintf("visibility%v", c.AbsoluteVisibility)},
	})
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

func (c *EvolvedComponent) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	return c.marshalSVG(e, canvas, svg.Color{color.RGBA{255, 0, 0, 255}})
}
