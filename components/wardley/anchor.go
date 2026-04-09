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
	LabelPlacement image.Point // LabelPlacement is relative to the placement
	Anchor         int         // text-anchor: AdjustStart, AdjustMiddle, AdjustEnd
	RenderingLayer int         //The position of the element on the picture
}

func NewAnchor(id int64) *Anchor {
	return &Anchor{
		id:             id,
		Placement:      image.Pt(components.UndefinedCoord, components.UndefinedCoord),
		LabelPlacement: image.Pt(components.UndefinedCoord, components.UndefinedCoord),
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
	labelP := a.LabelPlacement
	if labelP.X == components.UndefinedCoord {
		labelP.X = 0
	}
	if labelP.Y == components.UndefinedCoord {
		labelP.Y = 5
	}
	textAnchor := svg.TextAnchorMiddle
	switch a.Anchor {
	case AdjustStart:
		textAnchor = svg.TextAnchorStart
	case AdjustMiddle:
		textAnchor = svg.TextAnchorMiddle
	case AdjustEnd:
		textAnchor = svg.TextAnchorEnd
	}
	return e.Encode(svg.Transform{
		Translate: coords,
		Components: []interface{}{
			svg.Text{
				P:          labelP,
				Text:       []byte(a.Label),
				FontSize:   "14px",
				TextAnchor: textAnchor,
				TextAdjust: true,
			},
		},
	})
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
