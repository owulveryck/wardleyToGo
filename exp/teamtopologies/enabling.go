package tt

import (
	"encoding/xml"
	"image"
	"image/color"
)

type EnablingTeam struct {
	*Team
}

func NewEnablingTeam(id int64) *EnablingTeam {
	return &EnablingTeam{
		Team: NewTeam(id),
	}
}

func (c *EnablingTeam) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	return marshalTeamSVG(c.Team, e, canvas, 15, 15,
		color.RGBA{170, 185, 215, 229},
		color.RGBA{119, 159, 229, 178},
	)
}
