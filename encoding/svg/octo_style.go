package svgmap

import (
	"encoding/xml"
	"image"
	"image/color"

	"github.com/owulveryck/wardleyToGo/v2/internal/svg"
)

type OctoStyle struct {
	evolutionSteps []Evolution
	WithValueChain bool
	WithSpace      bool
	WithControls   bool
	Annotators     []Annotator
}

func NewOctoStyle(evolutionSteps []Evolution, annotators ...Annotator) *OctoStyle {
	svg.UpdateDefaultFont("'Outfit', sans-serif")
	return &OctoStyle{
		evolutionSteps: evolutionSteps,
		WithValueChain: true,
		WithSpace:      true,
		WithControls:   false,
		Annotators:     annotators,
	}

}

type Annotator interface {
	MarshalSVG(enc *xml.Encoder, box, canvas image.Rectangle)
}

func (w *OctoStyle) MarshalStyleSVG(enc *xml.Encoder, box, canvas image.Rectangle) {
	_ = enc.Encode(svg.Defs{
		Gradients: []svg.LinearGradient{
			{
				ID: "wardleyGradient",
				X1: "0%", Y1: "0%", X2: "100%", Y2: "0%",
				Stops: []svg.Stop{
					{
						Offset: "0%",
						StopColor: svg.Color{
							Color: color.RGBA{236, 237, 243, 255},
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
							Color: color.RGBA{236, 237, 243, 255},
						},
					},
				},
			},
			{
				ID: "pipelineTubeGradient",
				X1: "0%", Y1: "0%", X2: "0%", Y2: "100%",
				Stops: []svg.Stop{
					{
						Offset:    "0%",
						StopColor: svg.Color{Color: color.RGBA{180, 195, 225, 255}},
					},
					{
						Offset:    "25%",
						StopColor: svg.Color{Color: color.RGBA{220, 230, 248, 255}},
					},
					{
						Offset:    "50%",
						StopColor: svg.Color{Color: color.RGBA{205, 218, 242, 255}},
					},
					{
						Offset:    "80%",
						StopColor: svg.Color{Color: color.RGBA{185, 200, 230, 255}},
					},
					{
						Offset:    "100%",
						StopColor: svg.Color{Color: color.RGBA{160, 175, 210, 255}},
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
	if w.WithValueChain {
		verticals = append(verticals, svg.Line{
			F:           image.Point{0, 0},
			T:           image.Point{canvas.Dy(), 0},
			Stroke:      svg.Color{Color: color.RGBA{19, 36, 84, 255}},
			StrokeWidth: "1",
			MarkerEnd:   "url(#graphArrow)",
		})
	}
	for i := 1; i < len(w.evolutionSteps); i++ {
		position := w.evolutionSteps[i].Position
		verticals = append(verticals, svg.Line{
			F:               image.Point{0, int(float64(canvas.Dx()) * position)},
			T:               image.Point{canvas.Dy(), int(float64(canvas.Dx()) * position)},
			Stroke:          svg.Color{Color: color.RGBA{19, 36, 84, 255}},
			StrokeWidth:     "1",
			StrokeDashArray: []int{2, 2},
		})
	}

	if w.WithValueChain {
		verticals = append(verticals, svg.Text{
			P:          image.Point{5, -10},
			Text:       []byte(`Invisible`),
			Fill:       svg.Color{Color: color.RGBA{19, 36, 84, 255}},
			FontFamily: "'Outfit', sans-serif",
			TextAnchor: svg.TextAnchorStart,
		})
		verticals = append(verticals, svg.Text{
			P:          image.Point{canvas.Dy() - 5, -10},
			Text:       []byte(`Visible`),
			Fill:       svg.Color{Color: color.RGBA{19, 36, 84, 255}},
			FontFamily: "'Outfit', sans-serif",
			TextAnchor: svg.TextAnchorEnd,
		})
		verticals = append(verticals, svg.Text{
			P:          image.Point{canvas.Dy() / 2, -10},
			Text:       []byte(`Value Chain`),
			Fill:       svg.Color{Color: color.RGBA{19, 36, 84, 255}},
			TextAnchor: svg.TextAnchorMiddle,
			FontFamily: "'Outfit', sans-serif",
			FontWeight: "500",
		})
	}
	_ = enc.Encode(svg.Transform{
		Rotate:     270,
		Translate:  image.Point{canvas.Min.X, canvas.Max.Y},
		Components: verticals,
	})
	_ = enc.Encode(svg.Line{
		F:         image.Point{canvas.Min.X, canvas.Max.Y},
		T:         canvas.Max,
		Stroke:    svg.Color{Color: color.RGBA{19, 36, 84, 255}},
		MarkerEnd: "url(#graphArrow)",
	})
	if w.WithControls {
		displayControls(enc, box, canvas)
	}
	if w.WithSpace {
		_ = enc.Encode(svg.Text{
			P:          image.Point{canvas.Min.X + 7, canvas.Min.Y + 15},
			FontWeight: "500",
			FontSize:   "11px",
			Text:       []byte(`Uncharted`),
			TextAnchor: svg.TextAnchorStart,
			Fill:       svg.Color{Color: color.RGBA{19, 36, 84, 255}},
			FontFamily: "'Outfit', sans-serif",
		})
		_ = enc.Encode(svg.Text{
			P:          image.Point{canvas.Max.X - 5, canvas.Min.Y + 15},
			FontWeight: "500",
			FontSize:   "11px",
			Fill:       svg.Color{Color: color.RGBA{19, 36, 84, 255}},
			Text:       []byte(`Industrialised`),
			TextAnchor: svg.TextAnchorEnd,
			FontFamily: "'Outfit', sans-serif",
		})
	}
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
			FontWeight: "600",
			FontSize:   "14px",
			Fill:       svg.Color{Color: color.RGBA{19, 36, 84, 255}},
			FontFamily: "'Outfit', sans-serif",
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
			Fill:       svg.Color{Color: color.RGBA{19, 36, 84, 180}},
			FontFamily: "'Outfit', sans-serif",
		})
	}
	_ = enc.Encode(svg.Text{
		P:          image.Point{canvas.Max.X, canvas.Max.Y + 20},
		Text:       []byte(`Evolution`),
		TextAnchor: svg.TextAnchorEnd,
		Fill:       svg.Color{Color: color.RGBA{19, 36, 84, 255}},
		FontFamily: "'Outfit', sans-serif",
		FontWeight: "500",
	})
	for _, a := range w.Annotators {
		a.MarshalSVG(enc, box, canvas)
	}
}

func displayControls(enc *xml.Encoder, _, canvas image.Rectangle) {
	visibilityGroup := makeGroup("visibilitytoggle", 0)
	visibilityGroup.Attr = append(visibilityGroup.Attr, xml.Attr{
		Name:  xml.Name{Local: "onclick"},
		Value: "toggleVisibility()",
	},
	)
	_ = enc.EncodeToken(visibilityGroup.StartElement)

	_ = enc.Encode(svg.Circle{
		P: image.Point{canvas.Min.X + 105, canvas.Max.Y + 35},
		R: 5,
	})
	_ = enc.Encode(svg.Text{
		P:          image.Point{canvas.Min.X + 112, canvas.Max.Y + 39},
		FontWeight: "bold",
		FontSize:   "11px",
		Text:       []byte(`Toggle visibility`),
		TextAnchor: svg.TextAnchorStart,
		Fill:       svg.Color{Color: color.RGBA{19, 36, 84, 255}},
		FontFamily: "'Outfit', sans-serif",
	})
	_ = enc.EncodeToken(visibilityGroup.End())
	linkGroup := makeGroup("linktoggle", 0)
	linkGroup.Attr = append(linkGroup.Attr, xml.Attr{
		Name:  xml.Name{Local: "onclick"},
		Value: "toggleLinks()",
	},
	)
	_ = enc.EncodeToken(linkGroup.StartElement)

	_ = enc.Encode(svg.Circle{
		P: image.Point{canvas.Min.X + 5, canvas.Max.Y + 35},
		R: 5,
	})
	_ = enc.Encode(svg.Text{
		P:          image.Point{canvas.Min.X + 12, canvas.Max.Y + 39},
		FontWeight: "bold",
		FontSize:   "11px",
		Text:       []byte(`Toggle links`),
		TextAnchor: svg.TextAnchorStart,
		Fill:       svg.Color{Color: color.RGBA{19, 36, 84, 255}},
		FontFamily: "'Outfit', sans-serif",
	})
	_ = enc.EncodeToken(linkGroup.End())

}
