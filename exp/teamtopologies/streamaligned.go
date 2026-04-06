package tt

import (
	"encoding/xml"
	"image"
	"image/color"
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
	return marshalTeamSVG(c.Team, e, canvas, 15, 15,
		color.RGBA{252, 237, 190, 229},
		color.RGBA{250, 216, 120, 229},
	)
}
