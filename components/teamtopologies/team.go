package tt

import (
	"image"

	svg "github.com/ajstarks/svgo"
	"github.com/owulveryck/wardleyToGo/components"
)

type Team struct {
	id    int64
	Area  image.Rectangle
	Label string
}

func newTeam(id int64) *Team {
	return &Team{
		id:   id,
		Area: image.Rect(components.UndefinedCoord, components.UndefinedCoord, components.UndefinedCoord, components.UndefinedCoord),
	}
}

func (t *Team) String() string {
	return t.Label
}

func (t *Team) GetPosition() image.Point {
	return image.Pt((t.Area.Max.X-t.Area.Min.X)/2, (t.Area.Max.Y-t.Area.Min.Y)/2)
}
func (t *Team) GetArea() image.Rectangle {
	return t.Area
}

func (t *Team) svg(s *svg.SVG, bounds image.Rectangle) {
	placement := components.CalcCoords(t.Area.Min, bounds)
	s.Translate(placement.X, placement.Y)
}
func (t *Team) svgEnd(s *svg.SVG, _ image.Rectangle) {
	s.Gend()
}
