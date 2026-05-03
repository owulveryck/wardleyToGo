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
	legendFontFamily = "'Outfit', sans-serif"
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
		buildCol := svg.Color{Color: color.RGBA{0x29, 0x80, 0xB9, 0xFF}}
		_ = enc.Encode(svg.Circle{P: p, R: 7, Fill: buildCol, Stroke: buildCol, StrokeWidth: "1"})
		_ = enc.Encode(svg.Path{
			D: fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d Z",
				p.X-1, p.Y+5, p.X-1, p.Y+1, p.X-3, p.Y+1, p.X-4, p.Y-1, p.X-3, p.Y-3, p.X-1, p.Y-3,
				p.X-1, p.Y-1, p.X+1, p.Y-1, p.X+1, p.Y-3, p.X+3, p.Y-3, p.X+4, p.Y-1, p.X+3, p.Y+1,
				p.X+1, p.Y+1, p.X+1, p.Y+5),
			Fill: svg.White,
		})
	case "buy":
		buyCol := svg.Color{Color: color.RGBA{0x27, 0xAE, 0x60, 0xFF}}
		_ = enc.Encode(svg.Circle{P: p, R: 7, Fill: buyCol, Stroke: buyCol, StrokeWidth: "1"})
		_ = enc.Encode(svg.Path{
			D: fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d Z",
				p.X-5, p.Y-4, p.X-3, p.Y-4, p.X-2, p.Y-2, p.X+4, p.Y-2, p.X+3, p.Y+1, p.X-1, p.Y+1),
			Fill: svg.White,
		})
		_ = enc.Encode(svg.Circle{P: image.Pt(p.X-1, p.Y+3), R: 1, Fill: svg.White, Stroke: svg.White, StrokeWidth: "1"})
		_ = enc.Encode(svg.Circle{P: image.Pt(p.X+2, p.Y+3), R: 1, Fill: svg.White, Stroke: svg.White, StrokeWidth: "1"})
	case "outsource":
		outCol := svg.Color{Color: color.RGBA{0x8E, 0x44, 0xAD, 0xFF}}
		_ = enc.Encode(svg.Circle{P: p, R: 7, Fill: outCol, Stroke: outCol, StrokeWidth: "1"})
		_ = enc.Encode(svg.Circle{P: image.Pt(p.X-2, p.Y-3), R: 2, Fill: svg.White, Stroke: svg.White, StrokeWidth: "1"})
		_ = enc.Encode(svg.Path{
			D: fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Q %d,%d %d,%d Z",
				p.X-5, p.Y, p.X-5, p.Y+5, p.X, p.Y+5, p.X, p.Y, p.X-2, p.Y-2, p.X-5, p.Y),
			Fill: svg.White,
		})
		_ = enc.Encode(svg.Circle{P: image.Pt(p.X+2, p.Y-3), R: 2, Fill: svg.White, Stroke: svg.White, StrokeWidth: "1"})
		_ = enc.Encode(svg.Path{
			D: fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Q %d,%d %d,%d Z",
				p.X, p.Y, p.X, p.Y+5, p.X+5, p.Y+5, p.X+5, p.Y, p.X+2, p.Y-2, p.X, p.Y),
			Fill: svg.White,
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
		lx, cy, rx, ry := p.X-10, p.Y, 4, 5
		_ = enc.Encode(svg.Path{
			D:    fmt.Sprintf("M %d,%d A %d,%d 0 1,0 %d,%d A %d,%d 0 1,0 %d,%d Z", lx-rx, cy, rx, ry, lx+rx, cy, rx, ry, lx-rx, cy),
			Fill: svg.Color{Color: color.RGBA{130, 145, 180, 255}},
		})
		bodyEl := xml.StartElement{
			Name: xml.Name{Local: "path"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "d"}, Value: fmt.Sprintf("M %d,%d L %d,%d A %d,%d 0 0 1 %d,%d L %d,%d A %d,%d 0 0 1 %d,%d Z", lx, cy-ry, p.X+10, cy-ry, rx, ry, p.X+10, cy+ry, lx, cy+ry, rx, ry, lx, cy-ry)},
				{Name: xml.Name{Local: "style"}, Value: "fill:url(#pipelineTubeGradient)"},
				{Name: xml.Name{Local: "stroke"}, Value: "rgb(90,110,160)"},
				{Name: xml.Name{Local: "stroke-width"}, Value: "1"},
			},
		}
		_ = enc.EncodeToken(bodyEl)
		_ = enc.EncodeToken(bodyEl.End())
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
		marshalLegendInertiaWall(enc, p, "rgba(60,60,80,0.5)")
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
		} else if _, ok := strings.CutPrefix(item.Type, "inertia_"); ok {
			fill := "rgba(60,60,80,0.5)"
			if item.Color != nil {
				r, g, b, _ := item.Color.RGBA()
				fill = fmt.Sprintf("rgba(%d,%d,%d,0.5)", r>>8, g>>8, b>>8)
			}
			marshalLegendInertiaWall(enc, p, fill)
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
		col := svg.Color{Color: color.RGBA{0x29, 0x80, 0xB9, 0xFF}}
		cx, cy := ox+7, oy
		// Edges
		_ = enc.Encode(svg.Path{
			D:           fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d M %d,%d L %d,%d", cx, cy-6, cx-5, cy+3, cx+5, cy+3, cx, cy-6, cx, cy, cx-5, cy+3, cx, cy, cx+5, cy+3),
			Stroke:      col,
			StrokeWidth: "1",
		})
		// Center node
		_ = enc.Encode(svg.Circle{P: image.Pt(cx, cy), R: 2, Fill: col, Stroke: col, StrokeWidth: "1"})
		// Outer nodes
		_ = enc.Encode(svg.Circle{P: image.Pt(cx, cy-6), R: 1, Fill: col, Stroke: col, StrokeWidth: "1"})
		_ = enc.Encode(svg.Circle{P: image.Pt(cx-5, cy+3), R: 1, Fill: col, Stroke: col, StrokeWidth: "1"})
		_ = enc.Encode(svg.Circle{P: image.Pt(cx+5, cy+3), R: 1, Fill: col, Stroke: col, StrokeWidth: "1"})
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

func marshalLegendInertiaWall(enc *xml.Encoder, p image.Point, frontFill string) {
	cx, cy := p.X, p.Y
	halfThick := 4
	halfH := 5
	shearX := 2
	shearY := -2

	ax, ay := cx-halfThick, cy-halfH
	bx, by := cx+halfThick, cy-halfH
	cx2, cy2 := cx+halfThick, cy+halfH
	dx, dy := cx-halfThick, cy+halfH
	fx, fy := bx+shearX, by+shearY
	ex, ey := ax+shearX, ay+shearY
	gx, gy := cx2+shearX, cy2+shearY

	// Front face
	frontD := fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Z", ax, ay, bx, by, cx2, cy2, dx, dy)
	frontEl := xml.StartElement{
		Name: xml.Name{Local: "path"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "d"}, Value: frontD},
			{Name: xml.Name{Local: "style"}, Value: "fill:" + frontFill},
			{Name: xml.Name{Local: "stroke"}, Value: "rgba(40,40,60,0.6)"},
			{Name: xml.Name{Local: "stroke-width"}, Value: "0.5"},
		},
	}
	_ = enc.EncodeToken(frontEl)
	_ = enc.EncodeToken(frontEl.End())

	// Right side face
	sideD := fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Z", bx, by, fx, fy, gx, gy, cx2, cy2)
	sideEl := xml.StartElement{
		Name: xml.Name{Local: "path"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "d"}, Value: sideD},
			{Name: xml.Name{Local: "style"}, Value: "fill:rgba(40,40,60,0.5)"},
			{Name: xml.Name{Local: "stroke"}, Value: "rgba(30,30,50,0.6)"},
			{Name: xml.Name{Local: "stroke-width"}, Value: "0.5"},
		},
	}
	_ = enc.EncodeToken(sideEl)
	_ = enc.EncodeToken(sideEl.End())

	// Top face
	topD := fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d Z", ax, ay, bx, by, fx, fy, ex, ey)
	topEl := xml.StartElement{
		Name: xml.Name{Local: "path"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "d"}, Value: topD},
			{Name: xml.Name{Local: "style"}, Value: "fill:rgba(100,100,130,0.4)"},
			{Name: xml.Name{Local: "stroke"}, Value: "rgba(80,80,110,0.5)"},
			{Name: xml.Name{Local: "stroke-width"}, Value: "0.5"},
		},
	}
	_ = enc.EncodeToken(topEl)
	_ = enc.EncodeToken(topEl.End())
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
