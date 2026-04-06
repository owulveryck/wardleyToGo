package tt

import (
	"encoding/xml"
	"image"
	"image/color"

	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/internal/svg"
	"github.com/owulveryck/wardleyToGo/internal/utils"
)

type Team struct {
	id             int64
	Area           image.Rectangle
	Label          string
	RenderingLayer int //The position of the element on the picture
}

func (t *Team) GetLayer() int {
	return t.RenderingLayer
}

func NewTeam(id int64) *Team {
	return &Team{
		id:             id,
		Area:           image.Rect(components.UndefinedCoord, components.UndefinedCoord, components.UndefinedCoord, components.UndefinedCoord),
		RenderingLayer: 1,
	}
}
func (t *Team) ID() int64 {
	return t.id
}

func (t *Team) String() string {
	return t.Label
}

func (t *Team) GetPosition() image.Point {
	return image.Pt((t.Area.Max.X-t.Area.Min.X)/2, (t.Area.Max.Y-t.Area.Min.Y)/2)
}
func (t *Team) GetArea() image.Rectangle {
	return t.Area
}

// marshalTeamSVG renders a team as an SVG rectangle with the given style parameters.
func marshalTeamSVG(team *Team, e *xml.Encoder, canvas image.Rectangle, rx, ry int, fill, stroke color.RGBA) error {
	placement := utils.CalcCoords(team.Area.Min, canvas)
	w, h := utils.Scale(team.Area.Dx(), team.Area.Dy(), canvas)
	return e.Encode(svg.Transform{
		Translate: placement,
		Components: []any{
			svg.Rectangle{
				R:           image.Rect(0, 0, w, h),
				Rx:          rx,
				Ry:          ry,
				Fill:        svg.Color{fill},
				Stroke:      svg.Color{stroke},
				StrokeWidth: "5px",
			},
		},
	})
}
