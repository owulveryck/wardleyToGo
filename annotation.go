package wardleyToGo

import (
	"encoding/xml"
	"image"
	"image/color"
	"strconv"

	"github.com/owulveryck/wardleyToGo/internal/svg"
	"github.com/owulveryck/wardleyToGo/internal/utils"
)

// An annotation is a set of placements of a certain label
type Annotation struct {
	Identifier int
	Placements []image.Point
	Label      string
}

func NewAnnotation(identifier int) *Annotation {
	return &Annotation{
		Identifier: identifier,
		Placements: make([]image.Point, 0, 1),
	}
}

func (a *Annotation) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	for _, coords := range a.Placements {
		placement := utils.CalcCoords(coords, canvas)

		//s.Gid(strconv.FormatInt(a.id, 10))
		e.Encode(svg.Transform{
			Translate: placement,
			Components: []interface{}{
				svg.Circle{
					R:           15,
					Fill:        svg.White,
					Stroke:      svg.Color{color.RGBA{0x59, 0x59, 0x59, 0xff}},
					StrokeWidth: "2",
				},
				svg.Text{
					P:          image.Point{0, 5},
					Text:       []byte(strconv.Itoa(a.Identifier)),
					FontSize:   "14px",
					TextAnchor: svg.TextAnchorMiddle,
				},
			},
		})
	}
	return nil
}

func (a *Annotation) String() string {
	return a.Label
}
