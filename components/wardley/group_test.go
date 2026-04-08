package wardley

import (
	"bytes"
	"encoding/xml"
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestConvexHull(t *testing.T) {
	tests := []struct {
		name   string
		points []fpoint
		want   int // expected number of hull vertices
	}{
		{
			name:   "triangle",
			points: []fpoint{{0, 0}, {10, 0}, {5, 10}},
			want:   3,
		},
		{
			name:   "square with interior point",
			points: []fpoint{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {5, 5}},
			want:   4,
		},
		{
			name:   "collinear points",
			points: []fpoint{{0, 0}, {5, 0}, {10, 0}},
			want:   2,
		},
		{
			name:   "two points",
			points: []fpoint{{0, 0}, {10, 10}},
			want:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hull := convexHull(tt.points)
			if len(hull) != tt.want {
				t.Errorf("convexHull() returned %d vertices, want %d", len(hull), tt.want)
			}
		})
	}
}

func TestGroupGetPosition(t *testing.T) {
	g := &Group{
		MemberPoints: []image.Point{{10, 20}, {30, 40}},
	}
	pos := g.GetPosition()
	if pos.X != 20 || pos.Y != 30 {
		t.Errorf("GetPosition() = %v, want (20, 30)", pos)
	}
}

func TestGroupGetArea(t *testing.T) {
	g := &Group{
		MemberPoints: []image.Point{{10, 20}, {30, 40}},
	}
	area := g.GetArea()
	if area.Min.X != 5 || area.Min.Y != 15 || area.Max.X != 35 || area.Max.Y != 45 {
		t.Errorf("GetArea() = %v, want (5,15)-(35,45)", area)
	}
}

func TestGroupGetLayer(t *testing.T) {
	g := &Group{}
	if g.GetLayer() != 1 {
		t.Errorf("GetLayer() = %d, want 1", g.GetLayer())
	}
}

func TestGroupMarshalSVG_Empty(t *testing.T) {
	g := NewGroup(1, "empty", nil, color.RGBA{0, 0, 0, 0x30})
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	err := g.MarshalSVG(enc, image.Rect(30, 50, 1070, 850))
	if err != nil {
		t.Fatal(err)
	}
	enc.Flush()
	if buf.Len() != 0 {
		t.Errorf("empty group should produce no output, got %q", buf.String())
	}
}

func TestGroupMarshalSVG_SinglePoint(t *testing.T) {
	g := NewGroup(1, "single", []image.Point{{50, 50}}, color.RGBA{0x34, 0x98, 0xDB, 0x30})
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	err := g.MarshalSVG(enc, image.Rect(30, 50, 1070, 850))
	if err != nil {
		t.Fatal(err)
	}
	enc.Flush()
	out := buf.String()
	// Should contain arc commands for ellipse
	if !strings.Contains(out, "<path") {
		t.Error("single-point group should render as path with arcs")
	}
	if !strings.Contains(out, "A ") {
		t.Error("single-point group should contain arc commands")
	}
}

func TestGroupMarshalSVG_TwoPoints(t *testing.T) {
	g := NewGroup(1, "pair", []image.Point{{30, 50}, {70, 50}}, color.RGBA{0x2E, 0xCC, 0x71, 0x30})
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	err := g.MarshalSVG(enc, image.Rect(30, 50, 1070, 850))
	if err != nil {
		t.Fatal(err)
	}
	enc.Flush()
	out := buf.String()
	// Should contain arc commands for capsule
	if !strings.Contains(out, "<path") {
		t.Error("two-point group should render as path")
	}
	if !strings.Contains(out, "A ") {
		t.Error("two-point group should contain arc commands for capsule endcaps")
	}
}

func TestGroupMarshalSVG_MultiplePoints(t *testing.T) {
	g := NewGroup(1, "multi", []image.Point{{20, 30}, {60, 50}, {80, 70}}, color.RGBA{0xE7, 0x4C, 0x3C, 0x30})
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	err := g.MarshalSVG(enc, image.Rect(30, 50, 1070, 850))
	if err != nil {
		t.Fatal(err)
	}
	enc.Flush()
	out := buf.String()
	// Should contain quadratic Bezier commands for smooth corners
	if !strings.Contains(out, "Q ") {
		t.Error("multi-point group should contain quadratic Bezier curves")
	}
	// Should have both fill and stroke paths
	if strings.Count(out, "<path") != 2 {
		t.Errorf("multi-point group should render 2 paths (fill + stroke), got %d", strings.Count(out, "<path"))
	}
}
