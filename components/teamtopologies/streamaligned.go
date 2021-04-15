package tt

import (
	"encoding/xml"
	"image"
	"image/color"

	"github.com/owulveryck/wardleyToGo/internal/svg"
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

func (c *StreamAlignedTeam) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	placement := utils.CalcCoords(c.Area.Min, canvas)
	w, h := utils.Scale(c.Area.Dx(), c.Area.Dy(), canvas)
	return e.Encode(svg.Transform{
		Translate: placement,
		Components: []interface{}{
			svg.Rectangle{
				R:           image.Rect(0, 0, w, h),
				Rx:          15,
				Ry:          15,
				Fill:        svg.Color{color.RGBA{252, 237, 190, 229}},
				Stroke:      svg.Color{color.RGBA{250, 216, 120, 229}},
				StrokeWidth: "5px",
			},
		},
	})
}
