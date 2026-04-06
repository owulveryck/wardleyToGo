package wardley

import (
	"encoding/xml"
	"image"
	"image/color"
	"image/draw"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/internal/svg"
	"github.com/owulveryck/wardleyToGo/internal/utils"
)

// Compile-time interface compliance check.
var _ wardleyToGo.Component = (*Anchor)(nil)

// An Anchor of the map
type Anchor struct {
	id             int64
	Placement      image.Point
	Label          string
	RenderingLayer int //The position of the element on the picture
}

func NewAnchor(id int64) *Anchor {
	return &Anchor{
		id:        id,
		Placement: image.Pt(components.UndefinedCoord, components.UndefinedCoord),
	}
}

func (a *Anchor) GetLayer() int {
	return a.RenderingLayer
}

func (a *Anchor) ID() int64 {
	return a.id
}

func (a *Anchor) MarshalSVG(e *xml.Encoder, canvas image.Rectangle) error {
	coords := components.CalcCoords(a.Placement, canvas)
	//s.Gid(strconv.FormatInt(a.id, 10))
	e.Encode(svg.Transform{
		Translate: coords,
		Components: []interface{}{
			svg.Text{
				Text:       []byte(a.Label),
				FontSize:   "14px",
				TextAnchor: svg.TextAnchorMiddle,
			},
		},
	})
	return nil
}

func (a *Anchor) String() string {
	return a.Label
}

func (a *Anchor) GetPosition() image.Point {
	return a.Placement
}

func (a *Anchor) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	// Anchor bitmap rendering: place a single pixel marker at the anchor position.
	placement := utils.CalcCoords(a.Placement, r)
	dst.Set(placement.X, placement.Y, color.Black)
}
