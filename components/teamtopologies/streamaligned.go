package tt

import (
	"image"

	svg "github.com/ajstarks/svgo"
	"github.com/owulveryck/wardleyToGo/internal/utils"
)

type StreamAlignedTeam struct {
	*Team
}

func NewStreamAlignedTeam(id int64) *StreamAlignedTeam {
	return &StreamAlignedTeam{
		Team: NewTeam(id),
	}
}

func (sa *StreamAlignedTeam) SVGDraw(s *svg.SVG, bounds image.Rectangle) {
	sa.svg(s, bounds)
	w, h := utils.Scale(sa.Area.Dx(), sa.Area.Dy(), bounds)
	s.Roundrect(0, 0, w, h, 15, 15, `fill="rgb(252, 237, 190)"`, `opacity="0.9"`, `stroke="rgb(250,216,120)"`, `stroke-opacity="0.9"`, `stroke-width="5px"`)
	sa.svgEnd(s, bounds)
}
