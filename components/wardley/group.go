package wardley

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"
	"strings"

	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/internal/svg"
)

// Group represents a visual grouping of components rendered as an organic
// blob shape on the map.
type Group struct {
	id           int64
	Label        string
	MemberPoints []image.Point // positions in 100x100 space
	FillColor    color.Color   // RGBA with low alpha for translucency
	StrokeColor  color.Color   // same RGB, higher alpha for border
	TeamType     string        // EVT/PST team type: "explorer", "settler", "town-planner", "pioneer", "villager"
}

// NewGroup creates a Group with the given member positions and fill color.
// The stroke color is derived from the fill color with higher opacity.
func NewGroup(id int64, label string, members []image.Point, fillColor color.Color) *Group {
	r, g, b, _ := fillColor.RGBA()
	return &Group{
		id:           id,
		Label:        label,
		MemberPoints: members,
		FillColor:    fillColor,
		StrokeColor:  color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 0xA0},
	}
}

func (g *Group) ID() int64 { return g.id }

// GetPosition returns the centroid of all member points.
func (g *Group) GetPosition() image.Point {
	if len(g.MemberPoints) == 0 {
		return image.Point{}
	}
	var sx, sy int
	for _, p := range g.MemberPoints {
		sx += p.X
		sy += p.Y
	}
	n := len(g.MemberPoints)
	return image.Pt(sx/n, sy/n)
}

// GetArea returns the bounding rectangle of the group with padding.
// Implementing Area bypasses bounds checking in Map.AddComponent.
func (g *Group) GetArea() image.Rectangle {
	if len(g.MemberPoints) == 0 {
		return image.Rectangle{}
	}
	minX, minY := g.MemberPoints[0].X, g.MemberPoints[0].Y
	maxX, maxY := minX, minY
	for _, p := range g.MemberPoints[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	const pad = 5
	return image.Rect(minX-pad, minY-pad, maxX+pad, maxY+pad)
}

// GetLayer returns 1 so groups render between background (0) and components (10).
func (g *Group) GetLayer() int { return 1 }

func (g *Group) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	if len(g.MemberPoints) == 0 {
		return nil
	}

	// Convert member points to canvas coordinates
	canvasPoints := make([]fpoint, len(g.MemberPoints))
	for i, p := range g.MemberPoints {
		cp := components.CalcCoords(p, canvas)
		canvasPoints[i] = fpoint{float64(cp.X), float64(cp.Y)}
	}

	const padding = 50.0 // pixels of padding around members

	var pathD string

	switch len(canvasPoints) {
	case 1:
		pathD = ellipsePath(canvasPoints[0], padding, padding*0.7)
	default:
		pathD = blobPath(canvasPoints, padding)
	}

	// Render fill path
	fillPath := xml.StartElement{
		Name: xml.Name{Local: "path"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "d"}, Value: pathD},
		},
	}
	fc := svg.Color{Color: g.FillColor}
	fillPath.Attr = append(fillPath.Attr, must(fc.MarshalXMLAttr(xml.Name{Local: "fill"})))
	fillPath.Attr = append(fillPath.Attr, must(fc.MarshalXMLAttr(xml.Name{Local: "fill-opacity"})))
	fillPath.Attr = append(fillPath.Attr, xml.Attr{Name: xml.Name{Local: "stroke"}, Value: "none"})
	if err := e.EncodeToken(fillPath); err != nil {
		return err
	}
	if err := e.EncodeToken(fillPath.End()); err != nil {
		return err
	}

	// Render stroke path with team-type-specific styling
	strokeWidth := "2"
	if g.TeamType == "town-planner" {
		strokeWidth = "4"
	}
	strokePath := xml.StartElement{
		Name: xml.Name{Local: "path"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "d"}, Value: pathD},
			{Name: xml.Name{Local: "fill"}, Value: "none"},
			{Name: xml.Name{Local: "stroke-width"}, Value: strokeWidth},
		},
	}
	sc := svg.Color{Color: g.StrokeColor}
	strokePath.Attr = append(strokePath.Attr, must(sc.MarshalXMLAttr(xml.Name{Local: "stroke"})))
	strokePath.Attr = append(strokePath.Attr, must(sc.MarshalXMLAttr(xml.Name{Local: "stroke-opacity"})))

	// Apply team-type-specific border style
	if g.TeamType == "explorer" || g.TeamType == "pioneer" {
		strokePath.Attr = append(strokePath.Attr, xml.Attr{Name: xml.Name{Local: "stroke-dasharray"}, Value: "8,4"})
	}

	if err := e.EncodeToken(strokePath); err != nil {
		return err
	}
	if err := e.EncodeToken(strokePath.End()); err != nil {
		return err
	}

	// For town-planner, render a second inner stroke for double-line effect
	if g.TeamType == "town-planner" {
		innerPath := xml.StartElement{
			Name: xml.Name{Local: "path"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "d"}, Value: pathD},
				{Name: xml.Name{Local: "fill"}, Value: "none"},
				{Name: xml.Name{Local: "stroke-width"}, Value: "1"},
			},
		}
		innerPath.Attr = append(innerPath.Attr, must(sc.MarshalXMLAttr(xml.Name{Local: "stroke"})))
		innerPath.Attr = append(innerPath.Attr, must(sc.MarshalXMLAttr(xml.Name{Local: "stroke-opacity"})))
		if err := e.EncodeToken(innerPath); err != nil {
			return err
		}
		if err := e.EncodeToken(innerPath.End()); err != nil {
			return err
		}
	}

	return nil
}

func must(a xml.Attr, err error) xml.Attr {
	if err != nil {
		panic(err)
	}
	return a
}

// fpoint is a floating-point 2D point for intermediate calculations.
type fpoint struct {
	X, Y float64
}

// ellipsePath generates an SVG path for an ellipse using two arc commands.
func ellipsePath(center fpoint, rx, ry float64) string {
	return fmt.Sprintf(
		"M %.1f,%.1f A %.1f,%.1f 0 1,0 %.1f,%.1f A %.1f,%.1f 0 1,0 %.1f,%.1f Z",
		center.X-rx, center.Y,
		rx, ry, center.X+rx, center.Y,
		rx, ry, center.X-rx, center.Y,
	)
}

// blobPath generates a smooth organic blob around the given points.
func blobPath(points []fpoint, padding float64) string {
	// Generate padded point cloud: 8 points per member on a circle
	padded := make([]fpoint, 0, len(points)*8)
	for _, p := range points {
		for j := range 8 {
			angle := float64(j) * (2 * math.Pi / 8)
			padded = append(padded, fpoint{
				X: p.X + padding*math.Cos(angle),
				Y: p.Y + padding*math.Sin(angle),
			})
		}
	}

	// Compute convex hull
	hull := convexHull(padded)
	if len(hull) < 3 {
		center := points[0]
		return ellipsePath(center, padding, padding*0.7)
	}

	// Smooth corners with quadratic Bezier curves
	const roundFactor = 0.3
	n := len(hull)

	// For each vertex, compute approach and departure points
	type cornerPoints struct {
		approach, vertex, departure fpoint
	}
	corners := make([]cornerPoints, n)
	for i := range n {
		prev := hull[(i-1+n)%n]
		curr := hull[i]
		next := hull[(i+1)%n]

		corners[i] = cornerPoints{
			approach: fpoint{
				X: curr.X + roundFactor*(prev.X-curr.X),
				Y: curr.Y + roundFactor*(prev.Y-curr.Y),
			},
			vertex: curr,
			departure: fpoint{
				X: curr.X + roundFactor*(next.X-curr.X),
				Y: curr.Y + roundFactor*(next.Y-curr.Y),
			},
		}
	}

	// Build SVG path: start at departure of vertex 0
	var b strings.Builder
	fmt.Fprintf(&b, "M %.1f,%.1f", corners[0].departure.X, corners[0].departure.Y)
	for i := 1; i < n; i++ {
		// Line to approach of vertex i
		fmt.Fprintf(&b, " L %.1f,%.1f", corners[i].approach.X, corners[i].approach.Y)
		// Quadratic Bezier around vertex i
		fmt.Fprintf(&b, " Q %.1f,%.1f %.1f,%.1f",
			corners[i].vertex.X, corners[i].vertex.Y,
			corners[i].departure.X, corners[i].departure.Y)
	}
	// Close: line to approach of vertex 0, curve around vertex 0
	fmt.Fprintf(&b, " L %.1f,%.1f", corners[0].approach.X, corners[0].approach.Y)
	fmt.Fprintf(&b, " Q %.1f,%.1f %.1f,%.1f",
		corners[0].vertex.X, corners[0].vertex.Y,
		corners[0].departure.X, corners[0].departure.Y)
	b.WriteString(" Z")

	return b.String()
}

// convexHull computes the convex hull of a set of points using Andrew's
// monotone chain algorithm. Returns vertices in counter-clockwise order.
func convexHull(points []fpoint) []fpoint {
	n := len(points)
	if n < 3 {
		return points
	}

	// Sort by X, then by Y
	sort.Slice(points, func(i, j int) bool {
		if points[i].X != points[j].X {
			return points[i].X < points[j].X
		}
		return points[i].Y < points[j].Y
	})

	// Remove near-duplicates
	const eps = 0.5
	unique := []fpoint{points[0]}
	for i := 1; i < n; i++ {
		dx := points[i].X - unique[len(unique)-1].X
		dy := points[i].Y - unique[len(unique)-1].Y
		if dx*dx+dy*dy > eps*eps {
			unique = append(unique, points[i])
		}
	}
	points = unique
	n = len(points)
	if n < 3 {
		return points
	}

	hull := make([]fpoint, 0, 2*n)

	// Lower hull
	for _, p := range points {
		for len(hull) >= 2 && fcross(hull[len(hull)-2], hull[len(hull)-1], p) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p)
	}

	// Upper hull
	lower := len(hull) + 1
	for i := n - 2; i >= 0; i-- {
		p := points[i]
		for len(hull) >= lower && fcross(hull[len(hull)-2], hull[len(hull)-1], p) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p)
	}

	return hull[:len(hull)-1] // Remove last point (duplicate of first)
}

// fcross returns the cross product of vectors (A-O) and (B-O).
func fcross(o, a, b fpoint) float64 {
	return (a.X-o.X)*(b.Y-o.Y) - (a.Y-o.Y)*(b.X-o.X)
}
