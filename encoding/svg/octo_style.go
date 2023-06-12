package svgmap

import (
	"encoding/xml"
	"image"
	"image/color"

	"github.com/owulveryck/wardleyToGo/internal/svg"
)

type OctoStyle struct {
	evolutionSteps []Evolution
	WithValueChain bool
	WithSpace      bool
	WithControls   bool
	Annotators     []Annotator
}

func NewOctoStyle(evolutionSteps []Evolution, annotators ...Annotator) *OctoStyle {
	svg.UpdateDefaultFont("Century Gothic,CenturyGothic,AppleGothic,sans-serif")
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
	enc.Encode(svg.Defs{
		Gradient: svg.LinearGradient{
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
	enc.Encode(svg.Rectangle{
		R:     canvas,
		Style: "fill:url(#wardleyGradient)",
	})

	verticals := make([]interface{}, 0)
	if w.WithValueChain {
		verticals = append(verticals, svg.Line{
			F:           image.Point{0, 0},
			T:           image.Point{canvas.Dy(), 0},
			Stroke:      svg.Color{color.RGBA{19, 36, 84, 255}},
			StrokeWidth: "1",
			MarkerEnd:   "url(#graphArrow)",
		})
	}
	for i := 1; i < len(w.evolutionSteps); i++ {
		position := w.evolutionSteps[i].Position
		verticals = append(verticals, svg.Line{
			F:               image.Point{0, int(float64(canvas.Dx()) * position)},
			T:               image.Point{canvas.Dy(), int(float64(canvas.Dx()) * position)},
			Stroke:          svg.Color{color.RGBA{19, 36, 84, 255}},
			StrokeWidth:     "1",
			StrokeDashArray: []int{2, 2},
		})
	}

	if w.WithValueChain {
		verticals = append(verticals, svg.Text{
			P:          image.Point{5, -10},
			Text:       []byte(`Invisible`),
			Fill:       svg.Color{color.RGBA{19, 36, 84, 255}},
			FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
			TextAnchor: svg.TextAnchorStart,
		})
		verticals = append(verticals, svg.Text{
			P:          image.Point{canvas.Dy() - 5, -10},
			Text:       []byte(`Visible`),
			Fill:       svg.Color{color.RGBA{19, 36, 84, 255}},
			FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
			TextAnchor: svg.TextAnchorEnd,
		})
		verticals = append(verticals, svg.Text{
			P:          image.Point{canvas.Dy() / 2, -10},
			Text:       []byte(`Value Chain`),
			Fill:       svg.Color{color.RGBA{19, 36, 84, 255}},
			TextAnchor: svg.TextAnchorMiddle,
			FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
			FontWeight: "bold",
		})
	}
	enc.Encode(svg.Transform{
		Rotate:     270,
		Translate:  image.Point{canvas.Min.X, canvas.Max.Y},
		Components: verticals,
	})
	enc.Encode(svg.Line{
		F:         image.Point{canvas.Min.X, canvas.Max.Y},
		T:         canvas.Max,
		Stroke:    svg.Color{color.RGBA{19, 36, 84, 255}},
		MarkerEnd: "url(#graphArrow)",
	})
	if w.WithControls {
		displayControls(enc, box, canvas)
	}
	if w.WithSpace {
		enc.Encode(svg.Text{
			P:          image.Point{canvas.Min.X + 7, canvas.Min.Y + 15},
			FontWeight: "bold",
			FontSize:   "11px",
			Text:       []byte(`Uncharted`),
			TextAnchor: svg.TextAnchorStart,
			Fill:       svg.Color{color.RGBA{19, 36, 84, 255}},
			FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
		})
		enc.Encode(svg.Text{
			P:          image.Point{canvas.Max.X - 5, canvas.Min.Y + 15},
			FontWeight: "bold",
			FontSize:   "11px",
			Fill:       svg.Color{color.RGBA{19, 36, 84, 255}},
			Text:       []byte(`Industrialised`),
			TextAnchor: svg.TextAnchorEnd,
			FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
		})
	}
	for i := 0; i < len(w.evolutionSteps); i++ {
		axis := w.evolutionSteps[i]
		enc.Encode(svg.Text{
			P:          image.Point{int(float64(canvas.Dx())*axis.Position) + canvas.Min.X, canvas.Max.Y + 15},
			Text:       []byte(axis.Label),
			Fill:       svg.Color{color.RGBA{19, 36, 84, 255}},
			FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
		})
	}
	enc.Encode(svg.Text{
		P:          image.Point{canvas.Max.X, canvas.Max.Y + 15},
		Text:       []byte(`Evolution`),
		TextAnchor: svg.TextAnchorEnd,
		Fill:       svg.Color{color.RGBA{19, 36, 84, 255}},
		FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
		FontWeight: "bold",
	})
	for _, a := range w.Annotators {
		a.MarshalSVG(enc, box, canvas)
	}
}

func displayControls(enc *xml.Encoder, _, canvas image.Rectangle) {
	visibilityGroup := makeGroup("visibilitytoggle", 0)
	visibilityGroup.StartElement.Attr = append(visibilityGroup.StartElement.Attr, xml.Attr{
		Name:  xml.Name{Local: "onclick"},
		Value: "toggleVisibility()",
	},
	)
	enc.EncodeToken(visibilityGroup.StartElement)

	enc.Encode(svg.Circle{
		P: image.Point{canvas.Min.X + 105, canvas.Max.Y + 35},
		R: 5,
	})
	enc.Encode(svg.Text{
		P:          image.Point{canvas.Min.X + 112, canvas.Max.Y + 39},
		FontWeight: "bold",
		FontSize:   "11px",
		Text:       []byte(`Toggle visibility`),
		TextAnchor: svg.TextAnchorStart,
		Fill:       svg.Color{color.RGBA{19, 36, 84, 255}},
		FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
	})
	enc.EncodeToken(visibilityGroup.End())
	linkGroup := makeGroup("linktoggle", 0)
	linkGroup.StartElement.Attr = append(linkGroup.StartElement.Attr, xml.Attr{
		Name:  xml.Name{Local: "onclick"},
		Value: "toggleLinks()",
	},
	)
	enc.EncodeToken(linkGroup.StartElement)

	enc.Encode(svg.Circle{
		P: image.Point{canvas.Min.X + 5, canvas.Max.Y + 35},
		R: 5,
	})
	enc.Encode(svg.Text{
		P:          image.Point{canvas.Min.X + 12, canvas.Max.Y + 39},
		FontWeight: "bold",
		FontSize:   "11px",
		Text:       []byte(`Toggle links`),
		TextAnchor: svg.TextAnchorStart,
		Fill:       svg.Color{color.RGBA{19, 36, 84, 255}},
		FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
	})
	enc.EncodeToken(linkGroup.End())

}
