package tt

import (
	"image"

	svg "github.com/ajstarks/svgo"
)

type PlatformTeam struct {
	*Team
}

func NewPlatformTeam(id int64) *StreamAlignedTeam {
	return &StreamAlignedTeam{
		Team: newTeam(id),
	}
}

func (sa *PlatformTeam) SVG(s *svg.SVG, bounds image.Rectangle) {
	sa.svg(s, bounds)
	s.Rect(0, 0, sa.Area.Dx(), sa.Area.Dy(), `fill="rgb(170, 185, 215)"`, `opacity="0.95"`, `stroke="rgb(119,159,229)"`, `stroke-opacity="0.7"`, `stroke-width="5px"`)
	sa.svgEnd(s, bounds)
}
