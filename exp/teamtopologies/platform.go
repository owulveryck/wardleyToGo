package tt

import (
	"encoding/xml"
	"image"
	"image/color"
)

type PlatformTeam struct {
	*Team
}

func NewPlatformTeam(id int64) *PlatformTeam {
	return &PlatformTeam{
		Team: NewTeam(id),
	}
}

func (c *PlatformTeam) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	return marshalTeamSVG(c.Team, e, canvas, 0, 0,
		color.RGBA{170, 185, 215, 229},
		color.RGBA{119, 159, 229, 178},
	)
}
