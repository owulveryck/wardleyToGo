package labels

import (
	"image"
	"testing"
)

func TestPlaceLabels_SingleComponent(t *testing.T) {
	comps := []Component{
		{Name: "A", Position: image.Pt(50, 50), Label: "Component A"},
	}
	results := PlaceLabels(comps, DefaultOptions())
	r, ok := results["A"]
	if !ok {
		t.Fatal("expected result for component A")
	}
	// Single component should get default right placement.
	if r.Offset.X != 10 || r.Offset.Y != 0 {
		t.Errorf("single component: got offset %v, want (10,0)", r.Offset)
	}
}

func TestPlaceLabels_TwoOverlapping(t *testing.T) {
	comps := []Component{
		{Name: "A", Position: image.Pt(50, 50), Label: "Alpha"},
		{Name: "B", Position: image.Pt(50, 52), Label: "Beta"},
	}
	results := PlaceLabels(comps, DefaultOptions())
	rA := results["A"]
	rB := results["B"]

	// They should not both pick the same position.
	if rA.Offset == rB.Offset && rA.Anchor == rB.Anchor {
		t.Errorf("two overlapping components got identical placement: %+v", rA)
	}
}

func TestPlaceLabels_BoundaryRight(t *testing.T) {
	comps := []Component{
		{Name: "A", Position: image.Pt(98, 50), Label: "Right Edge"},
		{Name: "B", Position: image.Pt(50, 50), Label: "Center"},
	}
	results := PlaceLabels(comps, DefaultOptions())
	rA := results["A"]

	// Component at x=98 should not place label to the right (would go off-map).
	// It should choose left or another position.
	if rA.Anchor == AnchorStart && rA.Offset.X > 0 {
		// Verify the label doesn't extend past map boundary
		bbox := candidateBBox(candidate{dx: rA.Offset.X, dy: rA.Offset.Y, anchor: rA.Anchor},
			comps[0], DefaultOptions())
		if bbox.MaxX > 100 {
			t.Errorf("boundary component label extends past map: bbox.MaxX=%.1f", bbox.MaxX)
		}
	}
}

func TestPlaceLabels_BoundaryTop(t *testing.T) {
	comps := []Component{
		{Name: "A", Position: image.Pt(50, 3), Label: "Top Edge"},
		{Name: "B", Position: image.Pt(50, 50), Label: "Center"},
	}
	results := PlaceLabels(comps, DefaultOptions())
	rA := results["A"]

	// Component near top should not place label above.
	if rA.Offset.Y < 0 {
		bbox := candidateBBox(candidate{dx: rA.Offset.X, dy: rA.Offset.Y, anchor: rA.Anchor},
			comps[0], DefaultOptions())
		if bbox.MinY < 0 {
			t.Errorf("top boundary label goes above map: bbox.MinY=%.1f", bbox.MinY)
		}
	}
}

func TestPlaceLabels_PipelineMembers(t *testing.T) {
	// Three pipeline members at the same Y with close X values.
	// X values are close enough that "right" labels would overlap.
	comps := []Component{
		{Name: "A", Position: image.Pt(40, 50), Label: "Member Alpha"},
		{Name: "B", Position: image.Pt(43, 50), Label: "Member Beta"},
		{Name: "C", Position: image.Pt(46, 50), Label: "Member Gamma"},
	}
	results := PlaceLabels(comps, DefaultOptions())

	// Check that not all three pick the same placement.
	offsets := make(map[image.Point]int)
	for _, name := range []string{"A", "B", "C"} {
		r := results[name]
		offsets[r.Offset]++
	}
	allSame := len(offsets) == 1
	if allSame {
		t.Error("all three pipeline members got the same label placement")
	}
}

func TestPlaceLabels_DenseCluster(t *testing.T) {
	comps := []Component{
		{Name: "A", Position: image.Pt(50, 50), Label: "Alpha"},
		{Name: "B", Position: image.Pt(52, 50), Label: "Beta"},
		{Name: "C", Position: image.Pt(50, 52), Label: "Gamma"},
		{Name: "D", Position: image.Pt(52, 52), Label: "Delta"},
		{Name: "E", Position: image.Pt(51, 51), Label: "Epsilon"},
	}
	results := PlaceLabels(comps, DefaultOptions())

	// All 5 should get results.
	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}

	// Count distinct placements — with 5 components in a tight cluster,
	// we should get at least 3 different positions.
	offsets := make(map[image.Point]int)
	for _, r := range results {
		offsets[r.Offset]++
	}
	if len(offsets) < 3 {
		t.Errorf("dense cluster: only %d distinct placements for 5 components", len(offsets))
	}
}

func TestSplitString(t *testing.T) {
	tests := []struct {
		input string
		max   int
		want  int // expected number of lines
	}{
		{"Hello", 8, 1},
		{"Hello World", 8, 1},                  // "Hello" (5) + " World" fits since 5 < 8
		{"Alertes Trafic en Temps Reel", 8, 3}, // "Alertes Trafic" | "en Temps" | "Reel"
		{"A", 8, 1},
		{"", 8, 1},
		{"Application Mobile", 8, 2}, // "Application" (11 >= 8) | "Mobile"
	}
	for _, tt := range tests {
		lines := splitString(tt.input, tt.max)
		if len(lines) != tt.want {
			t.Errorf("splitString(%q, %d) = %d lines %v, want %d",
				tt.input, tt.max, len(lines), lines, tt.want)
		}
	}
}

func TestEstimateBBox(t *testing.T) {
	opts := DefaultOptions()

	// A short label anchored at start: MinX should be at anchor X minus padding.
	bbox := estimateBBox("Test", 8, 50, 50, AnchorStart, opts)
	expectedMinX := 50.0 - opts.Padding
	if abs64(bbox.MinX-expectedMinX) > 0.01 {
		t.Errorf("AnchorStart: MinX should be %.2f, got %.2f", expectedMinX, bbox.MinX)
	}
	if bbox.MaxX <= bbox.MinX {
		t.Error("bbox width should be positive")
	}

	// AnchorEnd: label extends to the left, MaxX at anchor X plus padding.
	bboxEnd := estimateBBox("Test", 8, 50, 50, AnchorEnd, opts)
	expectedMaxX := 50.0 + opts.Padding
	if abs64(bboxEnd.MaxX-expectedMaxX) > 0.01 {
		t.Errorf("AnchorEnd: MaxX should be %.2f, got %.2f", expectedMaxX, bboxEnd.MaxX)
	}

	// AnchorMiddle: label is centered.
	bboxMid := estimateBBox("Test", 8, 50, 50, AnchorMiddle, opts)
	width := bboxMid.MaxX - bboxMid.MinX
	if abs64(bboxMid.MinX-(50-width/2)) > 0.01 {
		t.Errorf("AnchorMiddle: not centered, MinX=%.2f", bboxMid.MinX)
	}
}

func TestPlaceLabels_Deterministic(t *testing.T) {
	comps := []Component{
		{Name: "Anchor", Position: image.Pt(50, 3), Label: "Automobiliste"},
		{Name: "App", Position: image.Pt(62, 21), Label: "Application Mobile"},
		{Name: "Itin", Position: image.Pt(55, 43), Label: "Itineraire Affiche"},
		{Name: "Alertes", Position: image.Pt(37, 37), Label: "Alertes Trafic"},
		{Name: "CDN", Position: image.Pt(87, 40), Label: "CDN"},
		{Name: "Moteur", Position: image.Pt(42, 61), Label: "Moteur de Calcul"},
		{Name: "Flux", Position: image.Pt(27, 54), Label: "Flux Temps Reel"},
		{Name: "Donnees", Position: image.Pt(52, 80), Label: "Modele de Donnees"},
		{Name: "Cloud", Position: image.Pt(82, 77), Label: "Infrastructure Cloud"},
		{Name: "OSM", Position: image.Pt(70, 95), Label: "Donnees OSM"},
		{Name: "Reseau", Position: image.Pt(92, 97), Label: "Reseau Mobile"},
		// Pipeline members at same Y — triggers the tiebreaker
		{Name: "AlgoD", Position: image.Pt(82, 61), Label: "Algo Dijkstra"},
		{Name: "AlgoP", Position: image.Pt(37, 61), Label: "Algo Predictif IA"},
		{Name: "AlgoQ", Position: image.Pt(5, 61), Label: "Algo Quantique"},
	}

	opts := DefaultOptions()
	reference := PlaceLabels(comps, opts)

	for run := 0; run < 50; run++ {
		result := PlaceLabels(comps, opts)
		for name, r := range reference {
			if result[name] != r {
				t.Fatalf("run %d: %s placement %+v != reference %+v",
					run, name, result[name], r)
			}
		}
	}
}

func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
