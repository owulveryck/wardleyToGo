package svgmap

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/owulveryck/wardleyToGo/internal/svg"
)

// LegendItem describes a single entry in the legend.
type LegendItem struct {
	Category string      // "Components", "Edges", "Signals", "Gameplays", "Groups"
	Label    string      // Human-readable label
	Type     string      // Internal type key for rendering the correct icon
	Color    color.Color // Optional color override (used for groups)
}

// Legend implements Annotator and renders a legend panel next to the map.
type Legend struct {
	Items []LegendItem
}

const (
	legendFontFamily = "Century Gothic,CenturyGothic,AppleGothic,sans-serif"
	legendFontSize   = "11px"
	legendLineHeight = 24
	legendIconOffset = 15 // X offset for icon center from left edge
	legendTextOffset = 35 // X offset for label text from left edge
	legendPadding    = 15
)

// LegendWidth is the width added to the viewBox when the legend is active.
const LegendWidth = 220

var legendTextColor = svg.Color{Color: color.RGBA{19, 36, 84, 255}}

func (l *Legend) MarshalSVG(enc *xml.Encoder, box, canvas image.Rectangle) {
	if len(l.Items) == 0 {
		return
	}

	// Legend area: right of canvas, within box
	x0 := canvas.Max.X + 20
	y0 := canvas.Min.Y
	areaWidth := box.Max.X - x0 - 10

	// Background rectangle
	_ = enc.Encode(svg.Rectangle{
		R:           image.Rect(x0, y0, x0+areaWidth, box.Max.Y-50),
		Rx:          6,
		Ry:          6,
		Fill:        svg.Color{Color: color.RGBA{250, 250, 252, 255}},
		Stroke:      svg.Color{Color: color.RGBA{200, 200, 210, 255}},
		StrokeWidth: "1",
	})

	// Title
	y := y0 + legendPadding + 14
	_ = enc.Encode(svg.Transform{
		Translate: image.Point{x0 + areaWidth/2, y},
		Components: []any{
			svg.Text{
				P:          image.Point{0, 0},
				Text:       []byte("Legend"),
				FontWeight: "bold",
				FontSize:   "13px",
				FontFamily: legendFontFamily,
				Fill:       legendTextColor,
				TextAnchor: svg.TextAnchorMiddle,
			},
		},
	})

	// Separator line under title
	y += 8
	_ = enc.Encode(svg.Line{
		F:           image.Point{x0 + 10, y},
		T:           image.Point{x0 + areaWidth - 10, y},
		Stroke:      svg.Color{Color: color.RGBA{200, 200, 210, 255}},
		StrokeWidth: "1",
	})

	y += legendLineHeight / 2

	// Group items by category
	categories := []string{"Components", "Edges", "Groups", "Signals", "Gameplays", "Other"}
	grouped := make(map[string][]LegendItem)
	for _, item := range l.Items {
		grouped[item.Category] = append(grouped[item.Category], item)
	}

	for _, cat := range categories {
		items, ok := grouped[cat]
		if !ok || len(items) == 0 {
			continue
		}

		// Category header
		y += legendLineHeight
		_ = enc.Encode(svg.Transform{
			Translate: image.Point{x0 + legendPadding, y},
			Components: []any{
				svg.Text{
					P:          image.Point{0, 0},
					Text:       []byte(cat),
					FontWeight: "bold",
					FontSize:   "13px",
					FontFamily: legendFontFamily,
					Fill:       legendTextColor,
					TextAnchor: svg.TextAnchorStart,
				},
			},
		})

		// Items
		for _, item := range items {
			y += legendLineHeight
			marshalLegendIcon(enc, item, image.Point{x0 + legendIconOffset, y})
			_ = enc.Encode(svg.Transform{
				Translate: image.Point{x0 + legendTextOffset, y - 2},
				Components: []any{
					svg.Text{
						P:          image.Point{0, 0},
						Text:       []byte(item.Label),
						FontSize:   legendFontSize,
						FontFamily: legendFontFamily,
						Fill:       legendTextColor,
						TextAnchor: svg.TextAnchorStart,
					},
				},
			})
		}
	}
}

func marshalLegendIcon(enc *xml.Encoder, item LegendItem, p image.Point) {
	switch item.Type {
	case "component":
		// Simple circle (regular component)
		_ = enc.Encode(svg.Circle{
			P:           p,
			R:           5,
			Fill:        svg.White,
			Stroke:      svg.Black,
			StrokeWidth: "1",
		})
	case "build":
		// Double circle: outer gray, inner white
		_ = enc.Encode(svg.Circle{
			P:           p,
			R:           8,
			Fill:        svg.Gray(196),
			Stroke:      svg.Gray(196),
			StrokeWidth: "1",
		})
		_ = enc.Encode(svg.Circle{
			P:           p,
			R:           4,
			Fill:        svg.White,
			Stroke:      svg.Black,
			StrokeWidth: "1",
		})
	case "buy":
		_ = enc.Encode(svg.Circle{
			P:           p,
			R:           8,
			Fill:        svg.Color{Color: color.RGBA{170, 165, 169, 255}},
			Stroke:      svg.Color{Color: color.RGBA{170, 165, 169, 255}},
			StrokeWidth: "1",
		})
		_ = enc.Encode(svg.Circle{
			P:           p,
			R:           4,
			Fill:        svg.White,
			Stroke:      svg.Black,
			StrokeWidth: "1",
		})
	case "outsource":
		_ = enc.Encode(svg.Circle{
			P:           p,
			R:           8,
			Fill:        svg.Color{Color: color.RGBA{68, 68, 68, 255}},
			Stroke:      svg.Color{Color: color.RGBA{68, 68, 68, 255}},
			StrokeWidth: "1",
		})
		_ = enc.Encode(svg.Circle{
			P:           p,
			R:           4,
			Fill:        svg.White,
			Stroke:      svg.Black,
			StrokeWidth: "1",
		})
	case "evolved":
		// Red circle for evolved component
		_ = enc.Encode(svg.Circle{
			P:           p,
			R:           5,
			Fill:        svg.White,
			Stroke:      svg.Red,
			StrokeWidth: "1.5",
		})
	case "pipeline":
		// Small rectangle
		_ = enc.Encode(svg.Rectangle{
			R:           image.Rect(p.X-10, p.Y-4, p.X+10, p.Y+4),
			Fill:        svg.Transparent,
			Stroke:      svg.Black,
			StrokeWidth: "1",
		})
	case "edge":
		// Gray line
		_ = enc.Encode(svg.Line{
			F:           image.Point{p.X - 10, p.Y},
			T:           image.Point{p.X + 10, p.Y},
			Stroke:      svg.Gray(128),
			StrokeWidth: "1",
		})
	case "evolved_edge":
		// Red dashed line
		_ = enc.Encode(svg.Line{
			F:               image.Point{p.X - 10, p.Y},
			T:               image.Point{p.X + 10, p.Y},
			Stroke:          svg.Red,
			StrokeWidth:     "1",
			StrokeDashArray: []int{3, 3},
		})
	case "inertia":
		// Black bar
		_ = enc.Encode(svg.Rectangle{
			R:    image.Rect(p.X-2, p.Y-6, p.X+2, p.Y+6),
			Fill: svg.Black,
		})
	case "group":
		// Colored ellipse using the group's actual color
		fillColor := color.RGBA{52, 152, 219, 40}
		strokeColor := color.RGBA{52, 152, 219, 255}
		if item.Color != nil {
			r, g, b, _ := item.Color.RGBA()
			fillColor = color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 0x30}
			strokeColor = color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 0xFF}
		}
		_ = enc.Encode(svg.Circle{
			P:           p,
			R:           7,
			Fill:        svg.Color{Color: fillColor},
			Stroke:      svg.Color{Color: strokeColor},
			StrokeWidth: "1.5",
		})
	case "signal":
		// Generic signal: small diamond
		marshalDiamond(enc, p, svg.Color{Color: color.RGBA{230, 126, 34, 255}})
	case "gameplay":
		// Small rounded badge
		_ = enc.Encode(svg.Rectangle{
			R:           image.Rect(p.X-8, p.Y-5, p.X+8, p.Y+5),
			Rx:          3,
			Ry:          3,
			Fill:        svg.Color{Color: color.RGBA{155, 89, 182, 255}},
			Stroke:      svg.Color{Color: color.RGBA{155, 89, 182, 255}},
			StrokeWidth: "1",
		})
	case "annotation":
		// Small circle with number
		_ = enc.Encode(svg.Circle{
			P:           p,
			R:           6,
			Fill:        svg.White,
			Stroke:      svg.Gray(128),
			StrokeWidth: "1",
		})
	default:
		// Signal subtypes: signal_accelerating, signal_declining, etc.
		if signalType, ok := strings.CutPrefix(item.Type, "signal_"); ok {
			marshalSignalIcon(enc, signalType, p)
		}
	}
}

// marshalSignalIcon renders the same signal icon used on components, centered at p.
func marshalSignalIcon(enc *xml.Encoder, signalType string, p image.Point) {
	ox, oy := p.X-7, p.Y-7
	switch signalType {
	case "accelerating":
		_ = enc.Encode(svg.Path{
			D:           fmt.Sprintf("M %d,%d l 7,-7 l -7,-7 M %d,%d l 7,-7 l -7,-7", ox, oy+7, ox+8, oy+7),
			Stroke:      svg.Color{Color: color.RGBA{0x27, 0xAE, 0x60, 0xFF}},
			StrokeWidth: "2",
		})
	case "declining":
		_ = enc.Encode(svg.Path{
			D:           fmt.Sprintf("M %d,%d l 0,14 l -5,-5 M %d,%d l 0,14 l 5,-5", ox, oy, ox, oy),
			Stroke:      svg.Color{Color: color.RGBA{0xE7, 0x4C, 0x3C, 0xFF}},
			StrokeWidth: "2",
		})
	case "stagnating":
		_ = enc.Encode(svg.Path{
			D:           fmt.Sprintf("M %d,%d h 14 M %d,%d h 14", ox, oy-3, ox, oy+3),
			Stroke:      svg.Color{Color: color.RGBA{0x95, 0xA5, 0xA6, 0xFF}},
			StrokeWidth: "2",
		})
	case "co-evolution":
		_ = enc.Encode(svg.Path{
			D:           fmt.Sprintf("M %d,%d q 4,-7 8,0 q 4,7 8,0", ox, oy),
			Stroke:      svg.Color{Color: color.RGBA{0x8E, 0x44, 0xAD, 0xFF}},
			StrokeWidth: "2",
		})
	case "red-queen":
		_ = enc.Encode(svg.Path{
			D:           fmt.Sprintf("M %d,%d l 11,0 l -4,-4 M %d,%d l 11,0 l -4,4 M %d,%d h 5 M %d,%d h 4", ox, oy, ox, oy, ox-4, oy-4, ox-4, oy+4),
			Stroke:      svg.Color{Color: color.RGBA{0xE7, 0x4C, 0x3C, 0xFF}},
			StrokeWidth: "1.5",
		})
	case "commoditization":
		_ = enc.Encode(svg.Path{
			D:           fmt.Sprintf("M %d,%d l 11,8 l -5,0 M %d,%d l 11,8 l 0,-5", ox, oy, ox, oy),
			Stroke:      svg.Color{Color: color.RGBA{0x34, 0x49, 0x5E, 0xFF}},
			StrokeWidth: "2",
		})
	case "network-effects":
		cx, cy := ox+7, oy
		_ = enc.Encode(svg.Circle{
			P:           image.Pt(cx, cy),
			R:           3,
			Fill:        svg.Color{Color: color.RGBA{0x29, 0x80, 0xB9, 0xFF}},
			Stroke:      svg.Color{Color: color.RGBA{0x29, 0x80, 0xB9, 0xFF}},
			StrokeWidth: "1.5",
		})
		_ = enc.Encode(svg.Path{
			D:           fmt.Sprintf("M %d,%d l 7,-4 M %d,%d l 7,4 M %d,%d l -7,-4 M %d,%d l -7,4", cx, cy, cx, cy, cx, cy, cx, cy),
			Stroke:      svg.Color{Color: color.RGBA{0x29, 0x80, 0xB9, 0xFF}},
			StrokeWidth: "1.5",
		})
	case "economies-of-scale":
		_ = enc.Encode(svg.Path{
			D:           fmt.Sprintf("M %d,%d a 4,4 0 0,1 8,0 M %d,%d a 7,7 0 0,1 14,0", ox, oy, ox-3, oy),
			Stroke:      svg.Color{Color: color.RGBA{0x16, 0xA0, 0x85, 0xFF}},
			StrokeWidth: "2",
		})
	default:
		marshalDiamond(enc, p, svg.Color{Color: color.RGBA{0xF3, 0x9C, 0x12, 0xFF}})
	}
}

func marshalDiamond(enc *xml.Encoder, p image.Point, fill svg.Color) {
	_ = enc.Encode(svg.Transform{
		Translate: p,
		Components: []any{
			svg.Line{
				F:           image.Point{0, -7},
				T:           image.Point{7, 0},
				Stroke:      fill,
				StrokeWidth: "2",
			},
			svg.Line{
				F:           image.Point{7, 0},
				T:           image.Point{0, 7},
				Stroke:      fill,
				StrokeWidth: "2",
			},
			svg.Line{
				F:           image.Point{0, 7},
				T:           image.Point{-7, 0},
				Stroke:      fill,
				StrokeWidth: "2",
			},
			svg.Line{
				F:           image.Point{-7, 0},
				T:           image.Point{0, -7},
				Stroke:      fill,
				StrokeWidth: "2",
			},
		},
	})
}
