package labels

import "strings"

// Rect is a floating-point rectangle used for overlap calculations.
type Rect struct {
	MinX, MinY, MaxX, MaxY float64
}

// intersectionArea returns the area of overlap between two rectangles.
// Returns 0 if they do not overlap.
func intersectionArea(a, b Rect) float64 {
	overlapX := min64(a.MaxX, b.MaxX) - max64(a.MinX, b.MinX)
	overlapY := min64(a.MaxY, b.MaxY) - max64(a.MinY, b.MinY)
	if overlapX <= 0 || overlapY <= 0 {
		return 0
	}
	return overlapX * overlapY
}

// estimateBBox computes the bounding box of a label in 100-unit map space.
// cx, cy is the label anchor point (component position + offset in unit space).
// anchor is the text-anchor value (AnchorStart, AnchorMiddle, AnchorEnd).
func estimateBBox(label string, maxChars int, cx, cy float64, anchor int, opts Options) Rect {
	if maxChars == 0 {
		maxChars = 8
	}
	lines := splitString(label, maxChars)

	maxLineLen := 0
	for _, line := range lines {
		if len(line) > maxLineLen {
			maxLineLen = len(line)
		}
	}

	width := float64(maxLineLen)*opts.CharWidth + opts.Padding*2
	height := float64(len(lines))*opts.LineHeight + opts.Padding*2

	var left float64
	switch anchor {
	case AnchorStart:
		left = cx - opts.Padding
	case AnchorMiddle:
		left = cx - width/2
	case AnchorEnd:
		left = cx - width + opts.Padding
	default:
		left = cx - opts.Padding
	}

	// Vertical centering: when dy=0 in the renderer, multi-line text is
	// centered around the anchor point. For dy<0 (above) or dy>0 (below),
	// the text starts from the anchor.
	top := cy - height/2

	return Rect{
		MinX: left,
		MinY: top,
		MaxX: left + width,
		MaxY: top + height,
	}
}

// splitString replicates internal/svg.splitString: it splits s into lines
// where each line contains at most max characters (breaking on word
// boundaries).
func splitString(s string, max int) []string {
	output := make([]string, 1)
	words := strings.Fields(s)
	for _, w := range words {
		switch {
		case len(output[len(output)-1]) >= max:
			output = append(output, w)
		case len(output[len(output)-1]) == 0:
			output[len(output)-1] = w
		default:
			output[len(output)-1] = output[len(output)-1] + " " + w
		}
	}
	return output
}

func min64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
