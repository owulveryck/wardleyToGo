package wardleyToGo

import (
	"image"
	"strconv"

	svg "github.com/ajstarks/svgo"
	"github.com/owulveryck/wardleyToGo/internal/utils"
)

// An annotation is a set of placements of a certain label
type Annotation struct {
	Identifier int
	Placements []image.Point
	Label      string
}

func NewAnnotation(identifier int) *Annotation {
	return &Annotation{
		Identifier: identifier,
		Placements: make([]image.Point, 0, 1),
	}
}

// Annotation fulfils the svgmap.SVGer interface
func (a *Annotation) SVGDraw(s *svg.SVG, bounds image.Rectangle) {
	for _, coords := range a.Placements {
		placement := utils.CalcCoords(coords, bounds)
		s.Translate(placement.X, placement.Y)
		s.Circle(0, 0, 15, `fill="white"`, `stroke="#595959"`, `stroke-width="2"`)
		s.Text(0, 5, strconv.Itoa(a.Identifier), `font-weight="14px"`, `font-size="14px"`, `text-anchor="middle"`)
		s.Gend()
	}

}

func (a *Annotation) String() string {
	return a.Label
}
