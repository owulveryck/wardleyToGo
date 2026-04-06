package tt

import (
	"encoding/xml"
	"image"
	"image/color"
)

type ComplicatedSubsystemTeam struct {
	*Team
}

func NewComplicatedSubsystemTeam(id int64) *ComplicatedSubsystemTeam {
	return &ComplicatedSubsystemTeam{
		Team: NewTeam(id),
	}
}

func (c *ComplicatedSubsystemTeam) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	return marshalTeamSVG(c.Team, e, canvas, 35, 35,
		color.RGBA{236, 210, 177, 229},
		color.RGBA{210, 149, 84, 178},
	)
}
