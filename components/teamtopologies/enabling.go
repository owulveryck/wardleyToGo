package tt

import (
	"image"

	svg "github.com/ajstarks/svgo"
)

type EnablingTeam struct {
	*Team
}

func NewEnablingTeam(id int64) *StreamAlignedTeam {
	return &StreamAlignedTeam{
		Team: NewTeam(id),
	}
}

func (sa *EnablingTeam) SVG(s *svg.SVG, bounds image.Rectangle) {
	sa.svg(s, bounds)
	s.Roundrect(0, 0, sa.Area.Dx(), sa.Area.Dy(), 15, 15, `fill="rgb(170, 185, 215)"`, `opacity="0.95"`, `stroke="rgb(119,159,229)"`, `stroke-opacity="0.7"`, `stroke-width="5px"`)
	sa.svgEnd(s, bounds)
}
