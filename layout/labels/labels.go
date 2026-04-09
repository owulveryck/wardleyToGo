// Package labels provides automatic label placement for Wardley Map
// components.
//
// It uses a greedy candidate-placement algorithm: components are processed
// in density order (most crowded first), and for each component the
// algorithm tries several candidate positions and picks the one with the
// least overlap against already-placed labels and component circles.
package labels

import (
	"image"
	"math"
	"sort"
)

// Anchor constants matching components/wardley text-anchor values.
const (
	AnchorUndefined int = iota
	AnchorStart
	AnchorMiddle
	AnchorEnd
)

// Component describes a positioned component whose label needs placement.
type Component struct {
	Name       string
	Position   image.Point // center on the 100x100 map
	Label      string      // text to display
	MaxChars   int         // wrap threshold (0 means use Options default)
	IsPipeline bool
}

// Result holds the computed label offset and text-anchor for a component.
type Result struct {
	Offset image.Point // relative offset from component center (SVG pixels)
	Anchor int         // AnchorStart, AnchorMiddle, AnchorEnd
}

// Options controls label placement behaviour.
type Options struct {
	// CanvasWidth and CanvasHeight are the SVG drawing area dimensions
	// in pixels, used to convert between 100-unit map space and SVG
	// pixel offsets. Default: 1040x800 (matching the standard wtg2svg canvas).
	CanvasWidth  int
	CanvasHeight int

	// MapSize is the logical coordinate space size. Default: 100.
	MapSize int

	// ComponentRadius in 100-unit map space. Default: 0.5.
	ComponentRadius float64

	// CharWidth in 100-unit map space per character. Default: 0.65.
	CharWidth float64

	// LineHeight in 100-unit map space per text line. Default: 2.25.
	LineHeight float64

	// MaxCharsPerLine is the default wrap threshold. Default: 8.
	MaxCharsPerLine int

	// Padding adds extra space around each label bounding box (in 100-unit
	// map space) to prevent labels from being placed too close together.
	// Default: 0.5.
	Padding float64

	// DensityRadius is the distance (in 100-unit map space) within which
	// components are considered neighbours for density-based priority
	// ordering. Default: 15.
	DensityRadius float64
}

// DefaultOptions returns sensible defaults for label placement on a
// standard Wardley Map.
func DefaultOptions() Options {
	return Options{
		CanvasWidth:     1040,
		CanvasHeight:    800,
		MapSize:         100,
		ComponentRadius: 0.5,
		CharWidth:       0.75,
		LineHeight:      2.25,
		MaxCharsPerLine: 8,
		Padding:         0.8,
		DensityRadius:   15,
	}
}

// PlaceLabels assigns label offsets and text-anchors to minimise overlap.
//
// It returns a map from component Name to the computed Result. Components
// not present in the result should use the default placement.
func PlaceLabels(comps []Component, opts Options) map[string]Result {
	results := make(map[string]Result, len(comps))
	if len(comps) <= 1 {
		for _, c := range comps {
			results[c.Name] = Result{Offset: image.Pt(10, 0), Anchor: AnchorUndefined}
		}
		return results
	}

	// Build component circles in 100-unit space for overlap detection.
	circles := make([]Rect, len(comps))
	positions := make([]componentPos, len(comps))
	for i, c := range comps {
		circles[i] = componentCircle(c.Position.X, c.Position.Y, opts.ComponentRadius)
		positions[i] = componentPos{x: c.Position.X, y: c.Position.Y}
	}

	// Compute density and sort: most crowded first, then by Y (top first).
	type indexedComp struct {
		index   int
		density int
	}
	indexed := make([]indexedComp, len(comps))
	for i := range comps {
		indexed[i] = indexedComp{
			index:   i,
			density: computeDensity(i, positions, opts.DensityRadius),
		}
	}
	sort.Slice(indexed, func(a, b int) bool {
		if indexed[a].density != indexed[b].density {
			return indexed[a].density > indexed[b].density
		}
		return comps[indexed[a].index].Position.Y < comps[indexed[b].index].Position.Y
	})

	placed := make([]Rect, 0, len(comps))

	for _, ic := range indexed {
		comp := comps[ic.index]
		bestScore := math.Inf(-1)
		bestIdx := 0

		for ci, cand := range defaultCandidates {
			bbox := candidateBBox(cand, comp, opts)
			preferRight := cand.name == "right"
			score := scoreCandidate(bbox, placed, circles, preferRight)

			if score > bestScore {
				bestScore = score
				bestIdx = ci
			}
		}

		best := defaultCandidates[bestIdx]
		results[comp.Name] = Result{
			Offset: image.Pt(best.dx, best.dy),
			Anchor: best.anchor,
		}

		// Record the chosen bbox so subsequent labels avoid it.
		placed = append(placed, candidateBBox(best, comp, opts))
	}

	return results
}
