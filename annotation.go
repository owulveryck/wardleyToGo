package wardleyToGo

import (
	"image"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

var (
	_ SVGer = &Annotation{}
)

type Annotation struct {
	Identifier int
	Placement  []image.Point
	Label      string
}

func NewAnnotation(identifier int) *Annotation {
	return &Annotation{
		Identifier: identifier,
		Placement:  make([]image.Point, 0, 1),
	}
}

func (a *Annotation) SVG(s *svg.SVG, bounds image.Rectangle) {
	for _, coords := range a.Placement {
		placement := calcCoords(coords, bounds)
		s.Translate(placement.X, placement.Y)
		s.Circle(0, 0, 15, `fill="white"`, `stroke="#595959"`, `stroke-width="2"`)
		s.Text(0, 5, strconv.Itoa(a.Identifier), `font-weight="14px"`, `font-size="14px"`, `text-anchor="middle"`)
		s.Gend()
	}

}

func (a *Annotation) String() string {
	return a.Label
}
