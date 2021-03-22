package plan

import (
	"strconv"

	svg "github.com/ajstarks/svgo"
)

type Annotation struct {
	Identifier int
	Coords     [][2]int
	Label      string
}

func NewAnnotation(identifier int) *Annotation {
	return &Annotation{
		Identifier: identifier,
		Coords:     make([][2]int, 0, 1),
	}
}

func (a *Annotation) SVG(s *svg.SVG, width, height, padLeft, padBottom int) {
	for _, coords := range a.Coords {
		s.Translate(coords[1]*(width-padLeft)/100+padLeft, (height-padLeft)-coords[0]*(height-padLeft)/100)
		s.Circle(0, 0, 15, `fill="white"`, `stroke="#595959"`, `stroke-width="2"`)
		s.Text(0, 5, strconv.Itoa(a.Identifier), `font-weight="14px"`, `font-size="14px"`, `text-anchor="middle"`)
		s.Gend()
	}

}

func (c *Annotation) String() string {
	return c.Label
}
