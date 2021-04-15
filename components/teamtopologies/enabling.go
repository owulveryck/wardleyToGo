package tt

import (
	"encoding/xml"
	"image"
	"image/color"

	"github.com/owulveryck/wardleyToGo/internal/svg"
	"github.com/owulveryck/wardleyToGo/internal/utils"
)

type EnablingTeam struct {
	*Team
}

func NewEnablingTeam(id int64) *StreamAlignedTeam {
	return &StreamAlignedTeam{
		Team: NewTeam(id),
	}
}

func (c *EnablingTeam) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	placement := utils.CalcCoords(c.Area.Min, canvas)
	w, h := utils.Scale(c.Area.Dx(), c.Area.Dy(), canvas)
	return e.Encode(svg.Transform{
		Translate: placement,
		Components: []interface{}{
			svg.Rectangle{
				R:           image.Rect(0, 0, w, h),
				Rx:          15,
				Ry:          15,
				Fill:        svg.Color{color.RGBA{170, 185, 215, 229}},
				Stroke:      svg.Color{color.RGBA{119, 159, 229, 178}},
				StrokeWidth: "5px",
			},
		},
	})
}
