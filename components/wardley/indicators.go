package wardley

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/internal/svg"
)

// signalIndicator returns SVG elements for a signal icon at the given offset.
// Signals render as small arrows or symbols to the upper-right of a component.
func signalIndicator(signalType string, offsetX, offsetY int) []interface{} {
	var elements []interface{}

	switch signalType {
	case "accelerating":
		// Double right chevron ">>"
		elements = append(elements, svg.Path{
			D:           fmt.Sprintf("M %d,%d l 10,-10 l -10,-10 M %d,%d l 10,-10 l -10,-10", offsetX, offsetY+10, offsetX+12, offsetY+10),
			Stroke:      svg.Color{Color: color.RGBA{0x27, 0xAE, 0x60, 0xFF}}, // green
			StrokeWidth: "3",
			Class:       []string{"signal-accelerating"},
		})
	case "declining":
		// Downward arrow
		elements = append(elements, svg.Path{
			D:           fmt.Sprintf("M %d,%d l 0,20 l -8,-8 M %d,%d l 0,20 l 8,-8", offsetX, offsetY, offsetX, offsetY),
			Stroke:      svg.Color{Color: color.RGBA{0xE7, 0x4C, 0x3C, 0xFF}}, // red
			StrokeWidth: "3",
			Class:       []string{"signal-declining"},
		})
	case "stagnating":
		// Horizontal double bar "="
		elements = append(elements, svg.Path{
			D:           fmt.Sprintf("M %d,%d h 20 M %d,%d h 20", offsetX, offsetY-4, offsetX, offsetY+4),
			Stroke:      svg.Color{Color: color.RGBA{0x95, 0xA5, 0xA6, 0xFF}}, // gray
			StrokeWidth: "3",
			Class:       []string{"signal-stagnating"},
		})
	case "co-evolution":
		// Intertwined arrows (double helix simplified)
		elements = append(elements, svg.Path{
			D:           fmt.Sprintf("M %d,%d q 6,-10 12,0 q 6,10 12,0", offsetX, offsetY),
			Stroke:      svg.Color{Color: color.RGBA{0x8E, 0x44, 0xAD, 0xFF}}, // purple
			StrokeWidth: "3",
			Class:       []string{"signal-co-evolution"},
		})
	case "red-queen":
		// Running arrow (right arrow with speed lines)
		elements = append(elements, svg.Path{
			D:           fmt.Sprintf("M %d,%d l 16,0 l -6,-6 M %d,%d l 16,0 l -6,6 M %d,%d h 8 M %d,%d h 6", offsetX, offsetY, offsetX, offsetY, offsetX-6, offsetY-6, offsetX-6, offsetY+6),
			Stroke:      svg.Color{Color: color.RGBA{0xE7, 0x4C, 0x3C, 0xFF}}, // red
			StrokeWidth: "2.5",
			Class:       []string{"signal-red-queen"},
		})
	case "commoditization":
		// Downward-right gravity arrow
		elements = append(elements, svg.Path{
			D:           fmt.Sprintf("M %d,%d l 16,12 l -8,0 M %d,%d l 16,12 l 0,-8", offsetX, offsetY, offsetX, offsetY),
			Stroke:      svg.Color{Color: color.RGBA{0x34, 0x49, 0x5E, 0xFF}}, // dark blue
			StrokeWidth: "3",
			Class:       []string{"signal-commoditization"},
		})
	case "network-effects":
		// Interconnected nodes forming a network graph
		col := svg.Color{Color: color.RGBA{0x29, 0x80, 0xB9, 0xFF}}
		cx, cy := offsetX+10, offsetY
		// Edges connecting all nodes
		elements = append(elements, svg.Path{
			D:           fmt.Sprintf("M %d,%d L %d,%d L %d,%d L %d,%d L %d,%d L %d,%d", cx, cy-10, cx-9, cy+5, cx+9, cy+5, cx, cy-10, cx, cy, cx-9, cy+5) + fmt.Sprintf(" M %d,%d L %d,%d", cx, cy, cx+9, cy+5),
			Stroke:      col,
			StrokeWidth: "1.5",
			Class:       []string{"signal-network-effects"},
		})
		// Center node
		elements = append(elements, svg.Circle{
			P: image.Pt(cx, cy), R: 3,
			Fill: col, Stroke: col, StrokeWidth: "1",
			Class: []string{"signal-network-effects"},
		})
		// Outer nodes
		elements = append(elements, svg.Circle{
			P: image.Pt(cx, cy-10), R: 2,
			Fill: col, Stroke: col, StrokeWidth: "1",
			Class: []string{"signal-network-effects"},
		})
		elements = append(elements, svg.Circle{
			P: image.Pt(cx-9, cy+5), R: 2,
			Fill: col, Stroke: col, StrokeWidth: "1",
			Class: []string{"signal-network-effects"},
		})
		elements = append(elements, svg.Circle{
			P: image.Pt(cx+9, cy+5), R: 2,
			Fill: col, Stroke: col, StrokeWidth: "1",
			Class: []string{"signal-network-effects"},
		})
	case "economies-of-scale":
		// Expanding concentric arcs
		elements = append(elements, svg.Path{
			D:           fmt.Sprintf("M %d,%d a 6,6 0 0,1 12,0 M %d,%d a 10,10 0 0,1 20,0", offsetX, offsetY, offsetX-4, offsetY),
			Stroke:      svg.Color{Color: color.RGBA{0x16, 0xA0, 0x85, 0xFF}}, // teal
			StrokeWidth: "2.5",
			Class:       []string{"signal-economies-of-scale"},
		})
	default:
		// Generic signal: small diamond
		elements = append(elements, svg.Path{
			D:           fmt.Sprintf("M %d,%d l 8,-8 l 8,8 l -8,8 Z", offsetX, offsetY),
			Stroke:      svg.Color{Color: color.RGBA{0xF3, 0x9C, 0x12, 0xFF}}, // amber
			StrokeWidth: "2.5",
			Class:       []string{"signal-unknown"},
		})
	}

	return elements
}

// warningIndicator returns SVG elements for a warning triangle at the given offset.
func warningIndicator(offsetX, offsetY int, text string) []interface{} {
	triColor := svg.Color{Color: color.RGBA{0xE6, 0x7E, 0x22, 0xFF}} // orange
	elements := []interface{}{
		// Warning triangle
		svg.Path{
			D:    fmt.Sprintf("M %d,%d l 10,-20 l 10,20 Z", offsetX-10, offsetY+10),
			Fill: triColor,
		},
		// Exclamation mark inside triangle
		svg.Text{
			P:          image.Pt(offsetX, offsetY+6),
			Text:       []byte("!"),
			FontSize:   "14px",
			FontWeight: "bold",
			TextAnchor: svg.TextAnchorMiddle,
			Fill:       svg.White,
		},
	}
	// Tooltip with warning text
	elements = append(elements, struct {
		XMLName xml.Name `xml:"title"`
		Text    string   `xml:",chardata"`
	}{
		Text: "⚠ " + text,
	})
	return elements
}

// gameplayBadge returns SVG elements for a gameplay annotation badge.
func gameplayBadge(gameplayType string, offsetX, offsetY int) []interface{} {
	// Short label for the badge
	label := gameplayType
	badgeWidth := len(label)*6 + 8

	badgeColor := gameplayColor(gameplayType)
	col := svg.Color{Color: badgeColor}

	elements := []interface{}{
		// Rounded background rectangle
		svg.Rectangle{
			R:           image.Rect(offsetX-badgeWidth/2, offsetY-6, offsetX+badgeWidth/2, offsetY+6),
			Rx:          3,
			Ry:          3,
			Fill:        col,
			Stroke:      col,
			StrokeWidth: "1",
		},
		// Label text
		svg.Text{
			P:          image.Pt(offsetX, offsetY+4),
			Text:       []byte(label),
			FontSize:   "8px",
			FontWeight: "bold",
			TextAnchor: svg.TextAnchorMiddle,
			Fill:       svg.White,
		},
	}
	return elements
}

// gameplayColor returns the badge color for a gameplay type.
func gameplayColor(gpType string) color.RGBA {
	switch gpType {
	case "ILC":
		return color.RGBA{0x8E, 0x44, 0xAD, 0xFF} // purple
	case "open-source":
		return color.RGBA{0x27, 0xAE, 0x60, 0xFF} // green
	case "land-grab":
		return color.RGBA{0xE7, 0x4C, 0x3C, 0xFF} // red
	case "embrace-extend":
		return color.RGBA{0xE6, 0x7E, 0x22, 0xFF} // orange
	case "tower-moat":
		return color.RGBA{0x34, 0x49, 0x5E, 0xFF} // dark blue
	case "FUD":
		return color.RGBA{0x7F, 0x8C, 0x8D, 0xFF} // gray
	case "strangler-fig":
		return color.RGBA{0x16, 0xA0, 0x85, 0xFF} // teal
	case "signal-distortion":
		return color.RGBA{0xF3, 0x9C, 0x12, 0xFF} // amber
	default:
		return color.RGBA{0x95, 0xA5, 0xA6, 0xFF} // light gray
	}
}

// typeIndicator returns SVG elements that replace the base circle for typed components.
// Icons are centered on (0,0) and sized visibly (~16px diameter).
func typeIndicator(compType wardleyToGo.ComponentType) []interface{} {
	switch compType {
	case BuildComponent:
		col := svg.Color{Color: color.RGBA{0x29, 0x80, 0xB9, 0xFF}}
		return []interface{}{
			svg.Circle{R: 8, Fill: col, Stroke: col, StrokeWidth: "1"},
			svg.Path{
				D:    "M -1.5,6 L -1.5,1 L -4,1 L -5,-1 L -4,-4 L -1.5,-4 L -1.5,-1.5 L 1.5,-1.5 L 1.5,-4 L 4,-4 L 5,-1 L 4,1 L 1.5,1 L 1.5,6 Z",
				Fill: svg.White,
			},
		}
	case BuyComponent:
		col := svg.Color{Color: color.RGBA{0x27, 0xAE, 0x60, 0xFF}}
		return []interface{}{
			svg.Circle{R: 8, Fill: col, Stroke: col, StrokeWidth: "1"},
			svg.Path{
				D:    "M -6,-5 L -4,-5 L -3,-2 L 5,-2 L 4,2 L -2,2 Z",
				Fill: svg.White,
			},
			svg.Circle{P: image.Pt(-1, 4), R: 1, Fill: svg.White, Stroke: svg.White, StrokeWidth: "1"},
			svg.Circle{P: image.Pt(3, 4), R: 1, Fill: svg.White, Stroke: svg.White, StrokeWidth: "1"},
		}
	case OutsourceComponent:
		col := svg.Color{Color: color.RGBA{0x8E, 0x44, 0xAD, 0xFF}}
		return []interface{}{
			svg.Circle{R: 8, Fill: col, Stroke: col, StrokeWidth: "1"},
			svg.Circle{P: image.Pt(-3, -3), R: 2, Fill: svg.White, Stroke: svg.White, StrokeWidth: "1"},
			svg.Path{
				D:    "M -6,0 L -6,6 L 0,6 L 0,0 Q -3,-2 -6,0 Z",
				Fill: svg.White,
			},
			svg.Circle{P: image.Pt(3, -3), R: 2, Fill: svg.White, Stroke: svg.White, StrokeWidth: "1"},
			svg.Path{
				D:    "M 0,0 L 0,6 L 6,6 L 6,0 Q 3,-2 0,0 Z",
				Fill: svg.White,
			},
		}
	case DataProductComponent:
		return []interface{}{
			svg.Circle{R: 14, StrokeWidth: "1", Fill: svg.Color{Color: color.RGBA{246, 72, 22, 0xff}}},
			svg.Circle{R: 5, StrokeWidth: "1", Stroke: svg.Black, Fill: svg.White},
		}
	default:
		return nil
	}
}

// noteIndicator returns SVG elements for a small note dot at the given offset.
func noteIndicator(offsetX, offsetY int, text string) []interface{} {
	col := svg.Color{Color: color.RGBA{0x34, 0x98, 0xDB, 0xFF}} // light blue
	elements := []interface{}{
		svg.Circle{
			P:           image.Pt(offsetX, offsetY),
			R:           4,
			Fill:        col,
			Stroke:      svg.White,
			StrokeWidth: "1",
		},
		svg.Text{
			P:          image.Pt(offsetX, offsetY+3),
			Text:       []byte("i"),
			FontSize:   "7px",
			FontWeight: "bold",
			TextAnchor: svg.TextAnchorMiddle,
			Fill:       svg.White,
		},
	}
	elements = append(elements, struct {
		XMLName xml.Name `xml:"title"`
		Text    string   `xml:",chardata"`
	}{
		Text: text,
	})
	return elements
}

// inertiaKindColor returns a color for an inertia kind.
func inertiaKindColor(kind string) color.RGBA {
	switch kind {
	case "tech":
		return color.RGBA{0x29, 0x80, 0xB9, 0xFF} // blue
	case "financial":
		return color.RGBA{0x27, 0xAE, 0x60, 0xFF} // green
	case "human":
		return color.RGBA{0xE6, 0x7E, 0x22, 0xFF} // orange
	case "relational":
		return color.RGBA{0x8E, 0x44, 0xAD, 0xFF} // purple
	case "social":
		return color.RGBA{0x16, 0xA0, 0x85, 0xFF} // teal
	default:
		return color.RGBA{0x00, 0x00, 0x00, 0xFF} // black
	}
}
