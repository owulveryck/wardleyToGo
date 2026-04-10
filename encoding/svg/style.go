package svgmap

import (
	"encoding/xml"
	"image"
	"image/color"

	"github.com/owulveryck/wardleyToGo/internal/svg"
)

type WardleyStyle struct {
	evolutionSteps []Evolution
}

type Evolution struct {
	Position  float64
	Label     string // Phase indicator (I, II, III, IV) displayed at zone boundary
	ZoneLabel string // Descriptive label displayed centered within the zone
}

const (
	stage12 float64 = 0.174
	stage23 float64 = 0.4
	stage34 float64 = 0.7
)

//	DataEvolution
//
// https://learnwardleymapping.com/2020/01/22/visualizing-the-interaction-of-evolution-and-data-measurement/
var DataEvolution = []Evolution{
	{
		Position:  0,
		Label:     "I",
		ZoneLabel: "Unmodeled",
	},
	{
		Position:  stage12,
		Label:     "II",
		ZoneLabel: "Divergent",
	},
	{
		Position:  stage23,
		Label:     "III",
		ZoneLabel: "Convergent",
	},
	{
		Position:  stage34,
		Label:     "IV",
		ZoneLabel: "Modeled",
	},
}

var DefaultEvolution = []Evolution{
	{
		Position:  0,
		Label:     "I",
		ZoneLabel: "🧪 Genesis",
	},
	{
		Position:  stage12,
		Label:     "II",
		ZoneLabel: "⚒️  Custom-Built",
	},
	{
		Position:  stage23,
		Label:     "III",
		ZoneLabel: "🛒 Product\n(+rental)",
	},
	{
		Position:  stage34,
		Label:     "IV",
		ZoneLabel: "⛽ Commodity\n(+utility)",
	},
}

func NewWardleyStyle(evolutionSteps []Evolution) *WardleyStyle {
	return &WardleyStyle{
		evolutionSteps: evolutionSteps,
	}

}
func (w *WardleyStyle) MarshalStyleSVG(enc *xml.Encoder, box, canvas image.Rectangle) {
	_ = enc.Encode(svg.Rectangle{
		R:    box,
		Fill: svg.Gray(128),
	})
	_ = enc.Encode(svg.Defs{
		Gradient: svg.LinearGradient{
			ID: "wardleyGradient",
			X1: "0%", Y1: "0%", X2: "100%", Y2: "0%",
			Stops: []svg.Stop{
				{
					Offset: "0%",
					StopColor: svg.Color{
						Color: color.RGBA{196, 196, 196, 255},
					},
				},
				{
					Offset:    "30%",
					StopColor: svg.White,
				},
				{
					Offset:    "70%",
					StopColor: svg.White,
				},
				{
					Offset: "100%",
					StopColor: svg.Color{
						Color: color.RGBA{196, 196, 196, 255},
					},
				},
			},
		},
		Markers: []svg.Marker{
			{
				ID:           "arrow",
				RefX:         15,
				RefY:         0,
				MarkerWidth:  12,
				MarkerHeight: 12,
				ViewBox:      "0 -5 10 10",
				Path: &svg.Path{
					D:    "M0,-5L10,0L0,5",
					Fill: svg.Red,
				},
			},
			{
				ID:           "graphArrow",
				RefX:         9,
				RefY:         0,
				MarkerWidth:  12,
				MarkerHeight: 12,
				ViewBox:      "0 -5 10 10",
				Path: &svg.Path{
					D:    "M0,-5L10,0L0,5",
					Fill: svg.Black,
				},
			},
		},
	})
	_ = enc.Encode(svg.Rectangle{
		R:     canvas,
		Style: "fill:url(#wardleyGradient)",
	})

	verticals := make([]interface{}, 0)
	verticals = append(verticals, svg.Line{
		F:           image.Point{0, 0},
		T:           image.Point{canvas.Dy(), 0},
		Stroke:      svg.Black,
		StrokeWidth: "1",
		MarkerEnd:   "url(#graphArrow)",
	})
	for i := 1; i < len(w.evolutionSteps); i++ {
		position := w.evolutionSteps[i].Position
		verticals = append(verticals, svg.Line{
			F:               image.Point{0, int(float64(canvas.Dx()) * position)},
			T:               image.Point{canvas.Dy(), int(float64(canvas.Dx()) * position)},
			Stroke:          svg.Gray(0xb8),
			StrokeWidth:     "1",
			StrokeDashArray: []int{2, 2},
		})
	}

	verticals = append(verticals, svg.Text{
		P:          image.Point{5, -10},
		Text:       []byte(`Invisible`),
		TextAnchor: svg.TextAnchorStart,
	})
	verticals = append(verticals, svg.Text{
		P:          image.Point{canvas.Dy() - 5, -10},
		Text:       []byte(`Visible`),
		TextAnchor: svg.TextAnchorEnd,
	})
	verticals = append(verticals, svg.Text{
		P:          image.Point{canvas.Dy() / 2, -10},
		Text:       []byte(`Value Chain`),
		TextAnchor: svg.TextAnchorMiddle,
		FontWeight: "bold",
	})
	_ = enc.Encode(svg.Transform{
		Rotate:     270,
		Translate:  image.Point{canvas.Min.X, canvas.Max.Y},
		Components: verticals,
	})
	_ = enc.Encode(svg.Line{
		F:         image.Point{canvas.Min.X, canvas.Max.Y},
		T:         canvas.Max,
		Stroke:    svg.Black,
		MarkerEnd: "url(#graphArrow)",
	})
	_ = enc.Encode(svg.Circle{
		P: image.Point{canvas.Min.X + 7, canvas.Min.Y + 5},
		R: 5,
	})
	_ = enc.Encode(svg.Text{
		P:          image.Point{canvas.Min.X + 7, canvas.Min.Y + 15},
		FontWeight: "bold",
		FontSize:   "11px",
		Text:       []byte(`Uncharted`),
		TextAnchor: svg.TextAnchorStart,
	})
	_ = enc.Encode(svg.Text{
		P:          image.Point{canvas.Max.X - 5, canvas.Min.Y + 15},
		FontWeight: "bold",
		FontSize:   "11px",
		Text:       []byte(`Industrialised`),
		TextAnchor: svg.TextAnchorEnd,
	})
	// Phase indicators (I, II, III, IV) at zone boundaries
	for i := 0; i < len(w.evolutionSteps); i++ {
		axis := w.evolutionSteps[i]
		anchor := svg.TextAnchorMiddle
		if i == 0 {
			anchor = svg.TextAnchorStart
		}
		_ = enc.Encode(svg.Text{
			P:          image.Point{int(float64(canvas.Dx())*axis.Position) + canvas.Min.X, canvas.Max.Y + 20},
			Text:       []byte(axis.Label),
			TextAnchor: anchor,
			FontWeight: "bold",
			FontSize:   "14px",
		})
	}
	// Zone labels centered within each zone
	for i := 0; i < len(w.evolutionSteps); i++ {
		if w.evolutionSteps[i].ZoneLabel == "" {
			continue
		}
		leftPos := w.evolutionSteps[i].Position
		rightPos := 1.0
		if i+1 < len(w.evolutionSteps) {
			rightPos = w.evolutionSteps[i+1].Position
		}
		centerX := int(float64(canvas.Dx())*(leftPos+rightPos)/2) + canvas.Min.X
		zoneWidth := rightPos - leftPos
		maxChars := int(zoneWidth * 40)
		if maxChars < 8 {
			maxChars = 8
		}
		_ = enc.Encode(svg.Text{
			P:          image.Point{centerX, canvas.Max.Y + 20},
			Text:       []byte(w.evolutionSteps[i].ZoneLabel),
			TextAnchor: svg.TextAnchorMiddle,
			FontSize:   "12px",
			TextAdjust: true,
			MaxChars:   maxChars,
		})
	}
	_ = enc.Encode(svg.Text{
		P:          image.Point{canvas.Max.X, canvas.Max.Y + 20},
		Text:       []byte(`Evolution`),
		TextAnchor: svg.TextAnchorEnd,
		FontWeight: "bold",
	})
}
