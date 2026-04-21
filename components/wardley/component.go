package wardley

import (
	"encoding/xml"
	"image"
	"image/color"
	"image/draw"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/internal/drawing"
	"github.com/owulveryck/wardleyToGo/internal/svg"
	"github.com/owulveryck/wardleyToGo/internal/utils"
)

// Compile-time interface compliance checks.
var (
	_ wardleyToGo.Component = (*Component)(nil)
	_ wardleyToGo.Component = (*EvolvedComponent)(nil)
)

const (
	DefaultComponentRenderingLayer int = 10
	AdjustUndefined                int = iota
	AdjustStart
	AdjustMiddle
	AdjustEnd
)

// A Component is an element of the map
type Component struct {
	id                  int64
	Placement           image.Point // The placement of the component on a rectangle 100x100
	Label               string
	LabelPlacement      image.Point // LabelPlacement is relative to the placement
	Type                wardleyToGo.ComponentType
	RenderingLayer      int //The position of the element on the picture
	Configured          bool
	EvolutionPos        int
	Inertia             int
	InertiaKinds        []string // Qualified inertia: "tech", "financial", "human", "relational", "social"
	Color               color.Color
	AbsoluteVisibility  int
	Anchor              int
	PipelinedComponents []*Component
	PipelineReference   *Component
	Description         string
	Asset               string                // Capital type: "tech", "financial", "human", "relational", "social"
	Cost                string                // Free-text cost annotation
	Signals             []ComponentSignal     // Market signals attached to this component
	Annotations         []ComponentAnnotation // Notes/warnings attached to this component
	Gameplays           []ComponentGameplay   // Strategic maneuvers on this component
}

// ComponentSignal is a signal attached to a component for SVG rendering.
type ComponentSignal struct {
	Type string // "accelerating", "declining", "co-evolution", etc.
}

// ComponentAnnotation is a note or warning attached to a component.
type ComponentAnnotation struct {
	Kind string // "note" or "warning"
	Text string
}

// ComponentGameplay is a strategic maneuver annotation.
type ComponentGameplay struct {
	Type string // "ILC", "open-source", etc.
	Text string // optional description
}

// GetAbsoluteVisibility returns the visibility of the component as seen from the anchor
func (c *Component) GetAbsoluteVisibility() int {
	return c.AbsoluteVisibility
}

// NewComponent with the corresponding id and default UndefinedCoords
func NewComponent(id int64) *Component {
	return &Component{
		id:             id,
		Placement:      image.Pt(components.UndefinedCoord, components.UndefinedCoord),
		LabelPlacement: image.Pt(components.UndefinedCoord, components.UndefinedCoord),
		RenderingLayer: 10,
		Anchor:         AdjustUndefined,
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
	// Draw type-specific outer circle
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
}

func (c *Component) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	if c.PipelinedComponents != nil {
		for _, pipelineComp := range c.PipelinedComponents {
			err := pipelineComp.MarshalSVG(e, canvas)
			if err != nil {
				return err
			}
		}
	}
	switch c.Type {
	case PipelineComponent:
		return c.marshalSVGPipeline(e, canvas, svg.Color{Color: c.Color})
	default:
		return c.marshalSVG(e, canvas, svg.Color{Color: c.Color})
	}
}

func (c *Component) marshalSVGPipeline(e *xml.Encoder, canvas image.Rectangle, col svg.Color) error {
	// Draw the rectangle
	if len(c.PipelinedComponents) > 1 {
		rect := getBounds(c.PipelinedComponents)
		lowestBound := components.CalcCoords(rect.Min, canvas)
		greaterBound := components.CalcCoords(rect.Max, canvas)
		err := e.Encode(svg.Rectangle{
			R:           image.Rect(lowestBound.X, lowestBound.Y+10, greaterBound.X, greaterBound.Y-10),
			Rx:          0,
			Ry:          0,
			Fill:        svg.Transparent,
			Stroke:      svg.Black,
			StrokeWidth: "1",
			Style:       "",
		})
		if err != nil {
			return err
		}
	}

	coords := components.CalcCoords(c.Placement, canvas)
	labelP := c.LabelPlacement
	if labelP.X == components.UndefinedCoord {
		labelP.X = 10
	}
	if labelP.Y == components.UndefinedCoord {
		labelP.Y = 0
	}
	fillColor := svg.White
	r, g, b, a := c.Color.RGBA()
	if r != 0 || g != 0 || b != 0 || a != 65535 {
		fillColor = svg.Color{Color: col}
	}
	components := make([]interface{}, 0, 8)
	components = append(components, svg.Rectangle{
		R: image.Rectangle{
			Min: image.Point{-5, -15},
			Max: image.Point{5, -5},
		},
		StrokeWidth: "3",
		Stroke:      col,
		Fill:        fillColor,
	})
	anchor := svg.TextAnchorUndefined
	switch c.Anchor {
	case AdjustStart:
		anchor = svg.TextAnchorStart
	case AdjustMiddle:
		anchor = svg.TextAnchorMiddle
	case AdjustEnd:
		anchor = svg.TextAnchorEnd
	}

	components = append(components, svg.Text{
		P:          labelP,
		Text:       []byte(c.Label),
		TextAdjust: true,
		TextAnchor: anchor,
		Fill:       col,
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
		labelP.Y = 0
	}
	fillColor := svg.White
	r, g, b, a := c.Color.RGBA()
	if r != 0 || g != 0 || b != 0 || a != 65535 {
		fillColor = svg.Color{Color: col}
	}
	baseCircle := svg.Circle{
		R:           5,
		StrokeWidth: "1",
		Stroke:      col,
		Fill:        fillColor,
	}
	components := make([]interface{}, 0, 8)
	switch c.Type {
	case BuildComponent:
		components = append(components, svg.Circle{
			R:           20,
			StrokeWidth: "1",
			Stroke:      svg.Black,
			Fill:        svg.Color{Color: color.RGBA{0xd6, 0xd6, 0xd6, 0xff}},
		})
	case BuyComponent:
		components = append(components, svg.Circle{
			R:           20,
			StrokeWidth: "1",
			Fill:        svg.Color{Color: color.RGBA{0xaa, 0xa5, 0xa9, 0xff}},
			Stroke:      svg.Color{Color: color.RGBA{0xd6, 0xd6, 0xd6, 0xff}},
		})
	case OutsourceComponent:
		components = append(components, svg.Circle{
			R:           20,
			StrokeWidth: "1",
			Fill:        svg.Color{Color: color.RGBA{0x44, 0x44, 0x44, 0xff}},
			Stroke:      svg.Color{Color: color.RGBA{0x44, 0x44, 0x44, 0xff}},
		})
	case DataProductComponent:
		components = append(components, svg.Circle{
			R:           14,
			StrokeWidth: "1",
			Fill:        svg.Color{Color: color.RGBA{246, 72, 22, 0xff}},
		})
	}
	components = append(components, baseCircle)
	anchor := svg.TextAnchorUndefined
	switch c.Anchor {
	case AdjustStart:
		anchor = svg.TextAnchorStart
	case AdjustMiddle:
		anchor = svg.TextAnchorMiddle
	case AdjustEnd:
		anchor = svg.TextAnchorEnd
	}

	components = append(components, svg.Text{
		P:          labelP,
		Text:       []byte(c.Label),
		TextAdjust: true,
		TextAnchor: anchor,
		Fill:       col,
	})
	components = append(components, struct {
		XMLName xml.Name `xml:"title"`
		Text    string   `xml:",chardata"`
	}{
		Text: c.Description,
	})

	if len(c.Signals) > 0 {
		var sigChildren []interface{}
		for i, sig := range c.Signals {
			offsetY := -20 - i*28
			for _, el := range signalIndicator(sig.Type, 20, offsetY) {
				sigChildren = append(sigChildren, el)
			}
		}
		components = append(components, svg.Group{
			Attrs:    []xml.Attr{{Name: xml.Name{Local: "data-type"}, Value: "signal"}},
			Children: sigChildren,
		})
	}

	if len(c.Annotations) > 0 {
		var warnChildren []interface{}
		for i, ann := range c.Annotations {
			if ann.Kind == "warning" {
				offsetX := -30 - i*28
				for _, el := range warningIndicator(offsetX, 0, ann.Text) {
					warnChildren = append(warnChildren, el)
				}
			}
		}
		if len(warnChildren) > 0 {
			components = append(components, svg.Group{
				Attrs:    []xml.Attr{{Name: xml.Name{Local: "data-type"}, Value: "warning"}},
				Children: warnChildren,
			})
		}
	}

	if len(c.Gameplays) > 0 {
		var gpChildren []interface{}
		for i, gp := range c.Gameplays {
			offsetY := 18 + i*14
			for _, el := range gameplayBadge(gp.Type, 0, offsetY) {
				gpChildren = append(gpChildren, el)
			}
		}
		components = append(components, svg.Group{
			Attrs:    []xml.Attr{{Name: xml.Name{Local: "data-type"}, Value: "gameplay"}},
			Children: gpChildren,
		})
	}

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
	// Draw type-specific outer circle
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
}

// GetCoordinates fulfils the Element interface
func (e *EvolvedComponent) GetPosition() image.Point {
	return e.Component.GetPosition()
}

func (e *EvolvedComponent) String() string {
	return "[evolved]" + e.Label
}

func (c *EvolvedComponent) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	return c.marshalSVG(e, canvas, svg.Color{Color: color.RGBA{255, 0, 0, 255}})
}
