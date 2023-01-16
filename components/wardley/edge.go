package wardley

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/internal/drawing"
	"github.com/owulveryck/wardleyToGo/internal/svg"
	"github.com/owulveryck/wardleyToGo/internal/utils"
	"gonum.org/v1/gonum/graph"
	dotencoding "gonum.org/v1/gonum/graph/encoding"
)

type Collaboration struct {
	F, T               wardleyToGo.Component
	Label              string
	Type               wardleyToGo.EdgeType
	RenderingLayer     int
	Visibility         int
	AbsoluteVisibility int
}

// GetAbsoluteVisibility returns the visibility of the component as seen from the anchor
func (c *Collaboration) GetAbsoluteVisibility() int {
	return c.AbsoluteVisibility
}

func (c *Collaboration) Attributes() []dotencoding.Attribute {
	return []dotencoding.Attribute{
		{
			Key:   "minlen",
			Value: strconv.Itoa(c.Visibility),
		},
	}
}

// From returns the from node of the edge.
func (c *Collaboration) From() graph.Node {
	return c.F
}

// To returns the to node of the edge.
func (c *Collaboration) To() graph.Node {
	return c.T
}

// ReversedEdge returns the edge reversal of the receiver
// if a reversal is valid for the data type.
// When a reversal is valid an edge of the same type as
// the receiver with nodes of the receiver swapped should
// be returned, otherwise the receiver should be returned
// unaltered.
func (c *Collaboration) ReversedEdge() graph.Edge {
	return &Collaboration{
		F:     c.T,
		T:     c.F,
		Label: c.Label,
		Type:  c.Type,
	}
}

func (c *Collaboration) GetLayer() int {
	return c.RenderingLayer
}

func (c *Collaboration) GetType() wardleyToGo.EdgeType {
	return c.Type
}

func (c *Collaboration) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	fromCoord := c.F.(wardleyToGo.Component).GetPosition()
	toCoord := c.T.(wardleyToGo.Component).GetPosition()
	coordsF := utils.CalcCoords(fromCoord, canvas)
	coordsT := utils.CalcCoords(toCoord, canvas)
	line := svg.Line{
		F:           coordsF,
		T:           coordsT,
		StrokeWidth: "1",
		Class:       []string{fmt.Sprintf("visibility%v", c.AbsoluteVisibility)},
	}
	switch c.Type {
	case RegularEdge:
		line.Stroke = svg.Gray(128)
	case EvolvedComponentEdge:
		line.MarkerEnd = "url(#arrow)"
		line.StrokeDashArray = []int{5, 5}
		line.Stroke = svg.Red
		line.Class = append(line.Class, "evolutionEdge")
	case EvolvedEdge:
		line.Stroke = svg.Red
	}
	return e.Encode(line)
}

// Draw aligns r.Min in dst with sp in src and then replaces the
// rectangle r in dst with the result of drawing src on dst.
func (c *Collaboration) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	dash := [2]int{0, 0}
	coordsF := utils.CalcCoords(c.F.(wardleyToGo.Component).GetPosition(), r)
	coordsT := utils.CalcCoords(c.T.(wardleyToGo.Component).GetPosition(), r)
	var col color.Color
	switch c.Type {
	case EvolvedComponentEdge:
		col = color.RGBA{0xff, 0x00, 0x00, 0xff}
		dash = [2]int{5, 5}
		drawing.Arrow(dst, coordsF.X, coordsF.Y, coordsT.X, coordsT.Y, col, dash)
	case EvolvedEdge:
		col = color.RGBA{0xff, 0x00, 0x00, 0xff}
		drawing.Line(dst, coordsF.X, coordsF.Y, coordsT.X, coordsT.Y, col, dash)
	default:
		col = color.Gray{Y: 128}
		drawing.Line(dst, coordsF.X, coordsF.Y, coordsT.X, coordsT.Y, col, dash)
	}
}
