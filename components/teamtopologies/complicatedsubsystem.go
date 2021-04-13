package tt

import (
	"image"

	svg "github.com/ajstarks/svgo"
	"github.com/owulveryck/wardleyToGo/internal/utils"
)

type ComplicatedSubsystemTeam struct {
	*Team
}

func NewComplicatedSubsystemTeam(id int64) *StreamAlignedTeam {
	return &StreamAlignedTeam{
		Team: NewTeam(id),
	}
}

func (sa *ComplicatedSubsystemTeam) SVGDraw(s *svg.SVG, bounds image.Rectangle) {
	sa.svg(s, bounds)
	w, h := utils.Scale(sa.Area.Dx(), sa.Area.Dy(), bounds)
	s.Roundrect(0, 0, w, h, 35, 35, `fill="rgb(236, 210, 177)"`, `opacity="0.9"`, `stroke="rgb(210,149,84)"`, `stroke-opacity="0.7"`, `stroke-width="5px"`)
	sa.svgEnd(s, bounds)
}
