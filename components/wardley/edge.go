package wardley

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/internal/drawing"
	"github.com/owulveryck/wardleyToGo/internal/svg"
	"github.com/owulveryck/wardleyToGo/internal/utils"
)

// Compile-time interface compliance check.
var _ wardleyToGo.Collaboration = (*Collaboration)(nil)

type Collaboration struct {
	F, T               wardleyToGo.Component
	Label              string
	Type               wardleyToGo.EdgeType
	Inertia            image.Point
	InertiaKinds       []string // Qualified inertia types for colored bar rendering
	CurveOffset        int      // perpendicular offset in pixels for Bézier control point; 0 = straight line
	RenderingLayer     int
	Visibility         int
	AbsoluteVisibility int
}

// GetAbsoluteVisibility returns the visibility of the component as seen from the anchor
func (c *Collaboration) GetAbsoluteVisibility() int {
	return c.AbsoluteVisibility
}

// From returns the source component of the edge.
func (c *Collaboration) From() wardleyToGo.Component {
	return c.F
}

// To returns the destination component of the edge.
func (c *Collaboration) To() wardleyToGo.Component {
	return c.T
}

func (c *Collaboration) GetLayer() int {
	return c.RenderingLayer
}

func (c *Collaboration) GetType() wardleyToGo.EdgeType {
	return c.Type
}

func (c *Collaboration) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	fromCoord := c.F.GetPosition()
	toCoord := c.T.GetPosition()
	coordsF := utils.CalcCoords(fromCoord, canvas)
	coordsT := utils.CalcCoords(toCoord, canvas)

	var stroke svg.Color
	var markerEnd string
	var dashArray []int
	classes := []string{}

	switch c.Type {
	case RegularEdge:
		stroke = svg.Gray(128)
	case EvolvedComponentEdge:
		markerEnd = "url(#arrow)"
		dashArray = []int{5, 5}
		stroke = svg.Red
		classes = append(classes, "evolutionEdge")
		if c.Inertia.X != 0 {
			inertiaPosition := utils.CalcCoords(c.Inertia, canvas)
			if len(c.InertiaKinds) > 0 {
				// Qualified inertia: render colored bars side by side
				barWidth := 10 / len(c.InertiaKinds)
				if barWidth < 3 {
					barWidth = 3
				}
				totalWidth := barWidth * len(c.InertiaKinds)
				startX := inertiaPosition.X - totalWidth/2
				for i, kind := range c.InertiaKinds {
					kindColor := inertiaKindColor(kind)
					bar := svg.Rectangle{
						R: image.Rectangle{
							Min: image.Point{
								X: startX + i*barWidth,
								Y: coordsF.Y - 15,
							},
							Max: image.Point{
								X: startX + (i+1)*barWidth,
								Y: coordsF.Y + 15,
							},
						},
						Fill:   svg.Color{Color: kindColor},
						Stroke: svg.Color{Color: kindColor},
					}
					if err := e.Encode(bar); err != nil {
						return err
					}
				}
			} else {
				// Unqualified inertia: single black bar (backward compatible)
				inertia := svg.Rectangle{
					R: image.Rectangle{
						Min: image.Point{
							X: inertiaPosition.X - 5,
							Y: coordsF.Y - 15,
						},
						Max: image.Point{
							X: inertiaPosition.X + 5,
							Y: coordsF.Y + 15,
						},
					},
					Fill:   svg.Black,
					Stroke: svg.Black,
				}
				if err := e.Encode(inertia); err != nil {
					return err
				}
			}
		}
	case EvolvedEdge:
		stroke = svg.Red
	}

	if c.CurveOffset != 0 {
		// Render as a quadratic Bézier curve.
		cx, cy := curveControlPoint(coordsF, coordsT, c.CurveOffset)
		d := fmt.Sprintf("M %d,%d Q %d,%d %d,%d",
			coordsF.X, coordsF.Y, cx, cy, coordsT.X, coordsT.Y)
		if err := e.Encode(svg.Path{
			D:               d,
			Stroke:          stroke,
			StrokeWidth:     "1",
			StrokeDashArray: dashArray,
			MarkerEnd:       markerEnd,
			Class:           classes,
		}); err != nil {
			return err
		}
		if c.Label != "" {
			// Place label near the control point (midpoint of curve).
			return e.Encode(edgeLabelText(cx, cy-8, c.Label))
		}
		return nil
	}

	if err := e.Encode(svg.Line{
		F:               coordsF,
		T:               coordsT,
		StrokeWidth:     "1",
		Stroke:          stroke,
		StrokeDashArray: dashArray,
		MarkerEnd:       markerEnd,
		Class:           classes,
	}); err != nil {
		return err
	}
	if c.Label != "" {
		// Place label at midpoint of straight edge.
		mx := (coordsF.X + coordsT.X) / 2
		my := (coordsF.Y + coordsT.Y) / 2
		return e.Encode(edgeLabelText(mx, my-8, c.Label))
	}
	return nil
}

// edgeLabelText creates SVG elements for an edge label at (x, y).
// It wraps the text in a transform group so multi-line labels render
// with proper line spacing.
func edgeLabelText(x, y int, label string) svg.Transform {
	return svg.Transform{
		Translate: image.Pt(x, y),
		Components: []any{
			svg.Text{
				P:          image.Pt(0, 0),
				Text:       []byte(label),
				TextAnchor: svg.TextAnchorMiddle,
				TextAdjust: true,
				Fill:       svg.Gray(100),
				FontSize:   "11px",
			},
		},
	}
}

// curveControlPoint computes the control point for a quadratic Bézier
// curve by offsetting the midpoint perpendicular to the line direction.
func curveControlPoint(from, to image.Point, offset int) (int, int) {
	mx := float64(from.X+to.X) / 2
	my := float64(from.Y+to.Y) / 2
	dx := float64(to.X - from.X)
	dy := float64(to.Y - from.Y)
	length := math.Sqrt(dx*dx + dy*dy)
	if length < 1 {
		return int(mx), int(my)
	}
	// Perpendicular unit vector
	px := -dy / length
	py := dx / length
	return int(mx + px*float64(offset)), int(my + py*float64(offset))
}

// Draw aligns r.Min in dst with sp in src and then replaces the
// rectangle r in dst with the result of drawing src on dst.
func (c *Collaboration) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	dash := [2]int{0, 0}
	coordsF := utils.CalcCoords(c.F.GetPosition(), r)
	coordsT := utils.CalcCoords(c.T.GetPosition(), r)
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
