package wardleyToGo

import (
	"image"
	"strconv"

	svg "github.com/ajstarks/svgo"
	"gonum.org/v1/gonum/graph"
)

var (
	_ SVGer      = &Anchor{}
	_ graph.Node = &Anchor{}
	_ Element    = &Anchor{}
)

// An Anchor of the map
type Anchor struct {
	id        int64
	Placement image.Point
	Label     string
}

func NewAnchor(id int64) *Anchor {
	return &Anchor{
		id:        id,
		Placement: image.Pt(UndefinedCoord, UndefinedCoord),
	}
}

func (c *Anchor) ID() int64 {
	return c.id
}

func (c *Anchor) SVG(s *svg.SVG, bounds image.Rectangle) {
	coords := calcCoords(c.Placement, bounds)
	s.Gid(strconv.FormatInt(c.id, 10))
	s.Translate(coords.X, coords.Y)
	s.Text(0, 0, c.Label, `font-weight="14px"`, `font-size="14px"`, `text-anchor="middle"`)
	s.Gend()
	s.Gend()
}

func (c *Anchor) String() string {
	return c.Label
}

func (c *Anchor) GetPosition() image.Point {
	return c.Placement
}
