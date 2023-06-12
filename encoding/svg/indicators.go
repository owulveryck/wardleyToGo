package svgmap

import (
	"encoding/xml"
	"image"
	"image/color"
	"math"

	"github.com/owulveryck/wardleyToGo/internal/svg"
)

type EvolutionIndication struct {
	Y      int
	Labels []string
}

func (e *EvolutionIndication) MarshalSVG(enc *xml.Encoder, box image.Rectangle, canvas image.Rectangle) {
	canvasWidth := canvas.Max.X - canvas.Min.X
	enc.Encode(svg.Transform{
		Translate: image.Point{canvas.Min.X + 7, canvas.Min.Y + e.Y},
		Rotate:    0,
		Scale:     0,
		Classes:   []string{},
		Components: []interface{}{
			svg.Text{
				P:          image.Point{},
				FontWeight: "bold",
				FontSize:   "11px",
				Text:       []byte(e.Labels[0]),
				TextAnchor: svg.TextAnchorStart,
				Fill:       svg.Color{color.RGBA{19, 36, 84, 80}},
				FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
				TextAdjust: true,
				MaxChars:   20,
			}},
	})
	enc.Encode(svg.Transform{
		Translate: image.Point{position(canvasWidth, stage12) + 47, canvas.Min.Y + e.Y},
		Rotate:    0,
		Scale:     0,
		Classes:   []string{},
		Components: []interface{}{
			svg.Text{
				P:          image.Point{},
				FontWeight: "bold",
				FontSize:   "11px",
				Text:       []byte(e.Labels[1]),
				TextAnchor: svg.TextAnchorStart,
				Fill:       svg.Color{color.RGBA{19, 36, 84, 80}},
				FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
				TextAdjust: true,
				MaxChars:   25,
			}},
	})
	enc.Encode(svg.Transform{
		Translate: image.Point{position(canvasWidth, stage23) + 47, canvas.Min.Y + e.Y},
		Rotate:    0,
		Scale:     0,
		Classes:   []string{},
		Components: []interface{}{
			svg.Text{
				P:          image.Point{},
				FontWeight: "bold",
				FontSize:   "11px",
				Text:       []byte(e.Labels[2]),
				TextAnchor: svg.TextAnchorStart,
				Fill:       svg.Color{color.RGBA{19, 36, 84, 80}},
				FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
				TextAdjust: true,
				MaxChars:   35,
			}},
	})
	enc.Encode(svg.Transform{
		Translate: image.Point{position(canvasWidth, stage34) + 47, canvas.Min.Y + e.Y},
		Rotate:    0,
		Scale:     0,
		Classes:   []string{},
		Components: []interface{}{
			svg.Text{
				P:          image.Point{},
				FontWeight: "bold",
				FontSize:   "11px",
				Text:       []byte(e.Labels[3]),
				TextAnchor: svg.TextAnchorStart,
				Fill:       svg.Color{color.RGBA{19, 36, 84, 80}},
				FontFamily: "Century Gothic,CenturyGothic,AppleGothic,sans-serif",
				TextAdjust: true,
				MaxChars:   35,
			}},
	})
}

func position(width int, stage float64) int {
	return int(math.Round(float64(width) * stage))
}

func AllEvolutionIndications() []Annotator {
	return []Annotator{
		&EvolutionIndication{Y: 235, Labels: []string{"Market: Undefined", "Market: Forming", "Market: Growing", "Market: Mature"}},
		&EvolutionIndication{Y: 335, Labels: []string{"Failure: High / tolerated / assumed", "Failure: Moderate / unsurprising but disappointed", "Failure: Not tolerated, focus on constant improvement", "Failure: Operational efficiency and surprised by failure"}},
		&EvolutionIndication{Y: 435, Labels: []string{"Focus of value: High future worth", "Focus of value: Seeking profis / ROI?", "Focus of value: High profitability", "Focus of value: High volume / reducing margin"}},
		&EvolutionIndication{Y: 535, Labels: []string{"Comparison: Constantly changing / a differential / unstable", "Comparison: Learning from others / testing the water / some evidential support", "Comparison: Feature difference", "Comparison: Essential / operational advantage"}},
	}
}
