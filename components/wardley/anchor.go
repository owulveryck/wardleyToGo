package wardley

import (
	"image"
	"strconv"

	svg "github.com/ajstarks/svgo"
	"github.com/owulveryck/wardleyToGo/components"
)

// An Anchor of the map
type Anchor struct {
	id             int64
	Placement      image.Point
	Label          string
	RenderingLayer int //The position of the element on the picture
}

func NewAnchor(id int64) *Anchor {
	return &Anchor{
		id:        id,
		Placement: image.Pt(components.UndefinedCoord, components.UndefinedCoord),
	}
}

func (a *Anchor) GetLayer() int {
	return a.RenderingLayer
}

func (a *Anchor) ID() int64 {
	return a.id
}

func (a *Anchor) SVGDraw(s *svg.SVG, bounds image.Rectangle) {
	coords := components.CalcCoords(a.Placement, bounds)
	s.Gid(strconv.FormatInt(a.id, 10))
	s.Translate(coords.X, coords.Y)
	s.Text(0, 0, a.Label, `font-weight="14px"`, `font-size="14px"`, `text-anchor="middle"`)
	s.Gend()
	s.Gend()
}

func (a *Anchor) String() string {
	return a.Label
}

func (a *Anchor) GetPosition() image.Point {
	return a.Placement
}
