package wardley

import (
	"encoding/xml"
	"fmt"
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
	switch c.Type {
	case BuildComponent:
		drawing.DrawCircle(dst, 6, placement, color.RGBA{0x29, 0x80, 0xB9, 0xFF}, color.RGBA{0x29, 0x80, 0xB9, 0xFF})
	case BuyComponent:
		drawing.DrawCircle(dst, 6, placement, color.RGBA{0x27, 0xAE, 0x60, 0xFF}, color.RGBA{0x27, 0xAE, 0x60, 0xFF})
	case OutsourceComponent:
		drawing.DrawCircle(dst, 6, placement, color.RGBA{0x8E, 0x44, 0xAD, 0xFF}, color.RGBA{0x8E, 0x44, 0xAD, 0xFF})
	case DataProductComponent:
		drawing.DrawCircle(dst, 7, placement, color.RGBA{0x44, 0x44, 0x44, 0xff}, color.RGBA{246, 72, 22, 0xff})
	}
	drawing.DrawCircle(dst, 5, placement, color.Black, color.White)
}

func (c *Component) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	if c.PipelinedComponents != nil {
		pipelineGroup := xml.StartElement{
			Name: xml.Name{Local: "g"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "opacity"}, Value: "0.75"},
			},
		}
		if err := e.EncodeToken(pipelineGroup); err != nil {
			return err
		}
		for _, pipelineComp := range c.PipelinedComponents {
			err := pipelineComp.MarshalSVG(e, canvas)
			if err != nil {
				return err
			}
		}
		if err := e.EncodeToken(pipelineGroup.End()); err != nil {
			return err
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
		if err := encodePipelineTube(e, lowestBound, greaterBound); err != nil {
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
func encodePipelineTube(e *xml.Encoder, lowestBound, greaterBound image.Point) error {
	centerY := lowestBound.Y
	leftX := lowestBound.X
	rightX := greaterBound.X
	ry := 16
	rx := 10
	topY := centerY - ry
	bottomY := centerY + ry
	highlightH := ry * 2 / 5

	// Layer 1: inner face of left opening (dark depth)
	innerD := fmt.Sprintf(
		"M %d,%d A %d,%d 0 1,0 %d,%d A %d,%d 0 1,0 %d,%d Z",
		leftX-rx, centerY, rx, ry, leftX+rx, centerY,
		rx, ry, leftX-rx, centerY,
	)
	if err := e.Encode(svg.Path{D: innerD, Fill: svg.Color{Color: color.RGBA{130, 145, 180, 255}}}); err != nil {
		return err
	}

	// Layer 2: tube body (outer shell)
	bodyD := fmt.Sprintf(
		"M %d,%d L %d,%d A %d,%d 0 0 1 %d,%d L %d,%d A %d,%d 0 0 1 %d,%d Z",
		leftX, topY, rightX, topY,
		rx, ry, rightX, bottomY,
		leftX, bottomY,
		rx, ry, leftX, topY,
	)
	bodyEl := xml.StartElement{
		Name: xml.Name{Local: "path"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "d"}, Value: bodyD},
			{Name: xml.Name{Local: "style"}, Value: "fill:url(#pipelineTubeGradient);opacity:0.7"},
			{Name: xml.Name{Local: "stroke"}, Value: "rgb(140,155,190)"},
			{Name: xml.Name{Local: "stroke-width"}, Value: "1.5"},
		},
	}
	if err := e.EncodeToken(bodyEl); err != nil {
		return err
	}
	if err := e.EncodeToken(bodyEl.End()); err != nil {
		return err
	}

	// Layer 3: highlight strip along top
	highlightD := fmt.Sprintf(
		"M %d,%d L %d,%d A %d,%d 0 0 1 %d,%d L %d,%d Z",
		leftX, topY, rightX, topY,
		rx, highlightH, rightX, topY+highlightH,
		leftX, topY+highlightH,
	)
	highlightEl := xml.StartElement{
		Name: xml.Name{Local: "path"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "d"}, Value: highlightD},
			{Name: xml.Name{Local: "fill"}, Value: "white"},
			{Name: xml.Name{Local: "opacity"}, Value: "0.4"},
		},
	}
	if err := e.EncodeToken(highlightEl); err != nil {
		return err
	}
	if err := e.EncodeToken(highlightEl.End()); err != nil {
		return err
	}

	// Layer 4: left front rim (3D overlap)
	rimD := fmt.Sprintf(
		"M %d,%d A %d,%d 0 0 0 %d,%d",
		leftX, topY, rx, ry, leftX, bottomY,
	)
	rimEl := xml.StartElement{
		Name: xml.Name{Local: "path"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "d"}, Value: rimD},
			{Name: xml.Name{Local: "fill"}, Value: "none"},
			{Name: xml.Name{Local: "stroke"}, Value: "rgb(140,155,190)"},
			{Name: xml.Name{Local: "stroke-width"}, Value: "2"},
		},
	}
	if err := e.EncodeToken(rimEl); err != nil {
		return err
	}
	return e.EncodeToken(rimEl.End())
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

	hasNote := false
	for _, ann := range c.Annotations {
		if ann.Kind == "note" {
			hasNote = true
			break
		}
	}

	var iconElems []interface{}
	if elems := typeIndicator(c.Type); elems != nil {
		iconElems = elems
	} else {
		iconElems = []interface{}{baseCircle}
	}

	if hasNote {
		components = append(components, svg.Group{
			Attrs:    []xml.Attr{{Name: xml.Name{Local: "class"}, Value: "has-note"}},
			Children: iconElems,
		})
	} else {
		components = append(components, iconElems...)
	}

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
			sigChildren = append(sigChildren, signalIndicator(sig.Type, 20, offsetY)...)
		}
		components = append(components, svg.Group{
			Attrs: []xml.Attr{
				{Name: xml.Name{Local: "data-type"}, Value: "signal"},
				{Name: xml.Name{Local: "class"}, Value: "signal-animated"},
			},
			Children: sigChildren,
		})
	}

	if len(c.Annotations) > 0 {
		var warnChildren []interface{}
		for i, ann := range c.Annotations {
			if ann.Kind == "warning" {
				offsetX := -30 - i*28
				warnChildren = append(warnChildren, warningIndicator(offsetX, 0, ann.Text)...)
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
			gpChildren = append(gpChildren, gameplayBadge(gp.Type, 0, offsetY)...)
		}
		components = append(components, svg.Group{
			Attrs:    []xml.Attr{{Name: xml.Name{Local: "data-type"}, Value: "gameplay"}},
			Children: gpChildren,
		})
	}

	if len(c.Annotations) > 0 {
		var noteChildren []interface{}
		noteIdx := 0
		for _, ann := range c.Annotations {
			if ann.Kind == "note" {
				offsetX := 10 + noteIdx*12
				noteChildren = append(noteChildren, noteIndicator(offsetX, -10, ann.Text)...)
				noteIdx++
			}
		}
		if len(noteChildren) > 0 {
			components = append(components, svg.Group{
				Attrs:    []xml.Attr{{Name: xml.Name{Local: "data-type"}, Value: "note"}},
				Children: noteChildren,
			})
		}
	}

	return e.Encode(svg.Transform{
		Translate:  coords,
		Components: components,
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
	switch c.Type {
	case BuildComponent:
		drawing.DrawCircle(dst, 6, placement, color.RGBA{0x29, 0x80, 0xB9, 0xFF}, color.RGBA{0x29, 0x80, 0xB9, 0xFF})
	case BuyComponent:
		drawing.DrawCircle(dst, 6, placement, color.RGBA{0x27, 0xAE, 0x60, 0xFF}, color.RGBA{0x27, 0xAE, 0x60, 0xFF})
	case OutsourceComponent:
		drawing.DrawCircle(dst, 6, placement, color.RGBA{0x8E, 0x44, 0xAD, 0xFF}, color.RGBA{0x8E, 0x44, 0xAD, 0xFF})
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
