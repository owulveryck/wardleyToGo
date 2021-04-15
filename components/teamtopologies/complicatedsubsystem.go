package tt

import (
	"encoding/xml"
	"image"
	"image/color"

	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/internal/svg"
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

func (c *ComplicatedSubsystemTeam) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	placement := components.CalcCoords(c.Area.Min, canvas)
	w, h := utils.Scale(c.Area.Dx(), c.Area.Dy(), canvas)
	return e.Encode(svg.Transform{
		Translate: placement,
		Components: []interface{}{
			svg.Rectangle{
				R:           image.Rect(0, 0, w, h),
				Rx:          35,
				Ry:          35,
				Fill:        svg.Color{color.RGBA{236, 210, 177, 229}},
				Stroke:      svg.Color{color.RGBA{210, 149, 84, 178}},
				StrokeWidth: "5px",
			},
		},
	})
}
