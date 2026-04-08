package labels

// candidate describes one possible label position relative to its component.
type candidate struct {
	name   string // descriptive name for debugging
	dx, dy int    // offset in SVG pixels from component center
	anchor int    // text-anchor: AnchorStart, AnchorMiddle, AnchorEnd
}

// defaultCandidates lists the positions tried for each label.
// Order matters only for tie-breaking (first candidate wins on equal score).
var defaultCandidates = []candidate{
	{name: "right", dx: 12, dy: 0, anchor: AnchorStart},
	{name: "above-right", dx: 12, dy: -8, anchor: AnchorStart},
	{name: "below-right", dx: 12, dy: 18, anchor: AnchorStart},
	{name: "left", dx: -12, dy: 0, anchor: AnchorEnd},
	{name: "above-left", dx: -12, dy: -8, anchor: AnchorEnd},
	{name: "above-center", dx: 0, dy: -8, anchor: AnchorMiddle},
	{name: "below-left", dx: -12, dy: 18, anchor: AnchorEnd},
	{name: "far-right", dx: 25, dy: 0, anchor: AnchorStart},
}

// scoreCandidate evaluates how good a candidate bounding box is.
// Higher is better. Penalties for overlapping placed labels, component
// circles, and going out of map bounds.
func scoreCandidate(bbox Rect, placed []Rect, circles []Rect, preferRight bool) float64 {
	score := 0.0

	for _, p := range placed {
		score -= intersectionArea(bbox, p) * 100.0
	}

	for _, c := range circles {
		score -= intersectionArea(bbox, c) * 50.0
	}

	// Boundary penalties
	if bbox.MinX < 0 {
		score -= -bbox.MinX * 200.0
	}
	if bbox.MinY < 0 {
		score -= -bbox.MinY * 200.0
	}
	if bbox.MaxX > 100 {
		score -= (bbox.MaxX - 100) * 200.0
	}
	if bbox.MaxY > 100 {
		score -= (bbox.MaxY - 100) * 200.0
	}

	if preferRight {
		score += 0.5
	}

	return score
}

// computeDensity counts how many other components are within radius units.
func computeDensity(idx int, positions []componentPos, radius float64) int {
	r2 := radius * radius
	count := 0
	cx := float64(positions[idx].x)
	cy := float64(positions[idx].y)
	for i, p := range positions {
		if i == idx {
			continue
		}
		dx := float64(p.x) - cx
		dy := float64(p.y) - cy
		if dx*dx+dy*dy < r2 {
			count++
		}
	}
	return count
}

type componentPos struct {
	x, y int
}

// candidateBBox converts a candidate's SVG pixel offset to 100-unit space
// and computes the label bounding box.
func candidateBBox(c candidate, comp Component, opts Options) Rect {
	pxPerUnitX := float64(opts.CanvasWidth) / float64(opts.MapSize)
	pxPerUnitY := float64(opts.CanvasHeight) / float64(opts.MapSize)

	unitDX := float64(c.dx) / pxPerUnitX
	unitDY := float64(c.dy) / pxPerUnitY

	cx := float64(comp.Position.X) + unitDX
	cy := float64(comp.Position.Y) + unitDY

	maxChars := comp.MaxChars
	if maxChars == 0 {
		maxChars = opts.MaxCharsPerLine
	}

	return estimateBBox(comp.Label, maxChars, cx, cy, c.anchor, opts)
}

// componentCircle returns a bounding Rect for a component's circle in
// 100-unit map space.
func componentCircle(x, y int, radius float64) Rect {
	cx := float64(x)
	cy := float64(y)
	return Rect{
		MinX: cx - radius,
		MinY: cy - radius,
		MaxX: cx + radius,
		MaxY: cy + radius,
	}
}
