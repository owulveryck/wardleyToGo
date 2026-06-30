package wardley

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/owulveryck/wardleyToGo/v2"
	"github.com/owulveryck/wardleyToGo/v2/internal/drawing"
	"github.com/owulveryck/wardleyToGo/v2/internal/svg"
	"github.com/owulveryck/wardleyToGo/v2/internal/utils"
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
			wallCX := inertiaPosition.X
			wallCY := coordsF.Y
			halfThick := inertiaWallHalfThick(c.InertiaKinds)

			// Segment 1: source to front of the wall
			if err := e.Encode(svg.Line{
				F:               coordsF,
				T:               image.Point{X: wallCX - halfThick, Y: wallCY},
				StrokeWidth:     "1",
				Stroke:          stroke,
				StrokeDashArray: dashArray,
				Class:           classes,
			}); err != nil {
				return err
			}

			// 3D perspective wall
			if err := encodeInertiaWall(e, wallCX, wallCY, c.InertiaKinds); err != nil {
				return err
			}

			// Segment 2: back of wall to target (faded)
			if err := e.Encode(svg.Line{
				F:               image.Point{X: wallCX + halfThick, Y: wallCY},
				T:               coordsT,
				StrokeWidth:     "1",
				Stroke:          stroke,
				StrokeDashArray: dashArray,
				MarkerEnd:       markerEnd,
				Class:           append(append([]string{}, classes...), "evolutionEdgeFaded"),
			}); err != nil {
				return err
			}

			if c.Label != "" {
				mx := (coordsF.X + coordsT.X) / 2
				my := (coordsF.Y + coordsT.Y) / 2
				if err := e.Encode(edgeLabelText(mx, my-8, c.Label)); err != nil {
					return err
				}
			}
			return nil
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
				FontWeight: "300",
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

func inertiaWallHalfThick(kinds []string) int {
	halfThick := 6
	if n := len(kinds); n > 0 {
		if needed := n * 3; needed > halfThick {
			halfThick = needed
		}
	}
	return halfThick
}

func encodeInertiaWall(e *xml.Encoder, cx, cy int, kinds []string) error {
	halfThick := inertiaWallHalfThick(kinds)
	halfH := 18
	shearX := 12
	shearY := -10

	// Points
	ax, ay := cx-halfThick, cy-halfH   // A: front top-left
	bx, by := cx+halfThick, cy-halfH   // B: front top-right
	cx2, cy2 := cx+halfThick, cy+halfH // C: front bottom-right
	dx, dy := cx-halfThick, cy+halfH   // D: front bottom-left
	ex, ey := ax+shearX, ay+shearY     // E: back top-left
	fx, fy := bx+shearX, by+shearY     // F: back top-right
	gx, gy := cx2+shearX, cy2+shearY   // G: back bottom-right

	// Face 1: front face (wall thickness visible face-on)
	if len(kinds) > 0 {
		stripeW := (2 * halfThick) / len(kinds)
		if stripeW < 2 {
			stripeW = 2
		}
		for i, kind := range kinds {
			col := inertiaKindColor(kind)
			sx := ax + i*stripeW
			sxEnd := ax + (i+1)*stripeW
			if i == len(kinds)-1 {
				sxEnd = bx
			}
			d := fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Z",
				sx, ay, sxEnd, by, sxEnd, cy2, sx, dy)
			el := xml.StartElement{
				Name: xml.Name{Local: "path"},
				Attr: []xml.Attr{
					{Name: xml.Name{Local: "d"}, Value: d},
					{Name: xml.Name{Local: "style"}, Value: fmt.Sprintf("fill:rgba(%d,%d,%d,0.5)", col.R, col.G, col.B)},
					{Name: xml.Name{Local: "stroke"}, Value: fmt.Sprintf("rgba(%d,%d,%d,0.6)", col.R, col.G, col.B)},
					{Name: xml.Name{Local: "stroke-width"}, Value: "0.5"},
				},
			}
			if err := e.EncodeToken(el); err != nil {
				return err
			}
			if err := e.EncodeToken(el.End()); err != nil {
				return err
			}
		}
	} else {
		d := fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Z",
			ax, ay, bx, by, cx2, cy2, dx, dy)
		el := xml.StartElement{
			Name: xml.Name{Local: "path"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "d"}, Value: d},
				{Name: xml.Name{Local: "style"}, Value: "fill:rgba(60,60,80,0.5)"},
				{Name: xml.Name{Local: "stroke"}, Value: "rgba(40,40,60,0.6)"},
				{Name: xml.Name{Local: "stroke-width"}, Value: "0.5"},
			},
		}
		if err := e.EncodeToken(el); err != nil {
			return err
		}
		if err := e.EncodeToken(el.End()); err != nil {
			return err
		}
	}

	// Face 2: right side face (back wall visible in perspective)
	sideD := fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Z",
		bx, by, fx, fy, gx, gy, cx2, cy2)
	sideEl := xml.StartElement{
		Name: xml.Name{Local: "path"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "d"}, Value: sideD},
			{Name: xml.Name{Local: "style"}, Value: "fill:rgba(40,40,60,0.5)"},
			{Name: xml.Name{Local: "stroke"}, Value: "rgba(30,30,50,0.6)"},
			{Name: xml.Name{Local: "stroke-width"}, Value: "0.5"},
		},
	}
	if err := e.EncodeToken(sideEl); err != nil {
		return err
	}
	if err := e.EncodeToken(sideEl.End()); err != nil {
		return err
	}

	// Face 3: top face (perspective depth)
	topD := fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Z",
		ax, ay, bx, by, fx, fy, ex, ey)
	topEl := xml.StartElement{
		Name: xml.Name{Local: "path"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "d"}, Value: topD},
			{Name: xml.Name{Local: "style"}, Value: "fill:rgba(100,100,130,0.4)"},
			{Name: xml.Name{Local: "stroke"}, Value: "rgba(80,80,110,0.5)"},
			{Name: xml.Name{Local: "stroke-width"}, Value: "0.5"},
		},
	}
	if err := e.EncodeToken(topEl); err != nil {
		return err
	}
	return e.EncodeToken(topEl.End())
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
