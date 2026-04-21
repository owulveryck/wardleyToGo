package owm2wtg2

import (
	"os"
	"strings"
	"testing"
)

func TestParseTitle(t *testing.T) {
	doc, err := Parse(strings.NewReader("title My Map"))
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title != "My Map" {
		t.Errorf("title = %q, want %q", doc.Title, "My Map")
	}
}

func TestParseComponent(t *testing.T) {
	input := "component Cup of Tea [0.79, 0.61]"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Components) != 1 {
		t.Fatalf("got %d components, want 1", len(doc.Components))
	}
	c := doc.Components[0]
	if c.Name != "Cup of Tea" {
		t.Errorf("name = %q, want %q", c.Name, "Cup of Tea")
	}
	if c.Visibility != 0.79 {
		t.Errorf("visibility = %v, want 0.79", c.Visibility)
	}
	if c.Maturity != 0.61 {
		t.Errorf("maturity = %v, want 0.61", c.Maturity)
	}
}

func TestParseComponentWithInertia(t *testing.T) {
	input := "component Kettle [0.43, 0.35] inertia"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	c := doc.Components[0]
	if !c.Inertia {
		t.Error("expected inertia to be true")
	}
}

func TestParseComponentWithType(t *testing.T) {
	tests := []struct {
		input    string
		wantType string
	}{
		{"component X [0.5, 0.5] (buy)", "buy"},
		{"component X [0.5, 0.5] (build)", "build"},
		{"component X [0.5, 0.5] (outsource)", "outsource"},
		{"component X [0.5, 0.5] (market)", "market"},
	}
	for _, tt := range tests {
		doc, err := Parse(strings.NewReader(tt.input))
		if err != nil {
			t.Fatal(err)
		}
		if doc.Components[0].Type != tt.wantType {
			t.Errorf("type = %q, want %q for input %q", doc.Components[0].Type, tt.wantType, tt.input)
		}
	}
}

func TestParseComponentWithLabel(t *testing.T) {
	input := "component Cup of Tea [0.79, 0.61] label [19, -4]"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	c := doc.Components[0]
	if !c.HasLabel {
		t.Error("expected HasLabel to be true")
	}
	if c.LabelX != 19 || c.LabelY != -4 {
		t.Errorf("label = [%d, %d], want [19, -4]", c.LabelX, c.LabelY)
	}
}

func TestParseAnchor(t *testing.T) {
	input := "anchor Business [0.95, 0.63]"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Anchors) != 1 {
		t.Fatalf("got %d anchors, want 1", len(doc.Anchors))
	}
	a := doc.Anchors[0]
	if a.Name != "Business" {
		t.Errorf("name = %q, want %q", a.Name, "Business")
	}
}

func TestParseEvolve(t *testing.T) {
	input := "evolve Kettle 0.62"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Evolves) != 1 {
		t.Fatalf("got %d evolves, want 1", len(doc.Evolves))
	}
	ev := doc.Evolves[0]
	if ev.Name != "Kettle" {
		t.Errorf("name = %q, want %q", ev.Name, "Kettle")
	}
	if ev.Maturity != 0.62 {
		t.Errorf("maturity = %v, want 0.62", ev.Maturity)
	}
}

func TestParseEvolveWithRename(t *testing.T) {
	input := "evolve Physical->Virtual 0.8"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	ev := doc.Evolves[0]
	if ev.Name != "Physical" {
		t.Errorf("name = %q, want %q", ev.Name, "Physical")
	}
	if ev.NewName != "Virtual" {
		t.Errorf("newName = %q, want %q", ev.NewName, "Virtual")
	}
	if ev.Maturity != 0.8 {
		t.Errorf("maturity = %v, want 0.8", ev.Maturity)
	}
}

func TestParseEvolveWithType(t *testing.T) {
	input := "evolve Kettle 0.62 (buy)"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if doc.Evolves[0].Type != "buy" {
		t.Errorf("type = %q, want %q", doc.Evolves[0].Type, "buy")
	}
}

func TestParseRegularEdge(t *testing.T) {
	input := "Customer->Cup of Tea"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Edges) != 1 {
		t.Fatalf("got %d edges, want 1", len(doc.Edges))
	}
	e := doc.Edges[0]
	if e.From != "Customer" || e.To != "Cup of Tea" {
		t.Errorf("edge = %q -> %q, want Customer -> Cup of Tea", e.From, e.To)
	}
	if e.FlowType != "regular" {
		t.Errorf("flowType = %q, want regular", e.FlowType)
	}
}

func TestParseChainedEdge(t *testing.T) {
	input := "A->B->C"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Edges) != 2 {
		t.Fatalf("got %d edges, want 2", len(doc.Edges))
	}
	if doc.Edges[0].From != "A" || doc.Edges[0].To != "B" {
		t.Errorf("edge[0] = %q -> %q, want A -> B", doc.Edges[0].From, doc.Edges[0].To)
	}
	if doc.Edges[1].From != "B" || doc.Edges[1].To != "C" {
		t.Errorf("edge[1] = %q -> %q, want B -> C", doc.Edges[1].From, doc.Edges[1].To)
	}
}

func TestParseBidirectionalFlow(t *testing.T) {
	input := "Customer+<>Cup of Tea"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	e := doc.Edges[0]
	if e.FlowType != "bidirectional" {
		t.Errorf("flowType = %q, want bidirectional", e.FlowType)
	}
}

func TestParseLabeledFlow(t *testing.T) {
	input := "Hot Water+'$0.10'>Kettle"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	e := doc.Edges[0]
	if e.FlowType != "labeled" {
		t.Errorf("flowType = %q, want labeled", e.FlowType)
	}
	if e.Label != "$0.10" {
		t.Errorf("label = %q, want %q", e.Label, "$0.10")
	}
}

func TestParsePastFlow(t *testing.T) {
	input := "Hot Water+<Kettle"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if doc.Edges[0].FlowType != "past" {
		t.Errorf("flowType = %q, want past", doc.Edges[0].FlowType)
	}
}

func TestParseFutureFlow(t *testing.T) {
	input := "Hot Water+>Kettle"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if doc.Edges[0].FlowType != "future" {
		t.Errorf("flowType = %q, want future", doc.Edges[0].FlowType)
	}
}

func TestParsePipeline(t *testing.T) {
	input := "pipeline Customer [0.15, 0.9]"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Pipelines) != 1 {
		t.Fatalf("got %d pipelines, want 1", len(doc.Pipelines))
	}
	p := doc.Pipelines[0]
	if p.Name != "Customer" {
		t.Errorf("name = %q, want Customer", p.Name)
	}
	if !p.HasRange {
		t.Error("expected HasRange to be true")
	}
	if p.StartMat != 0.15 || p.EndMat != 0.9 {
		t.Errorf("range = [%v, %v], want [0.15, 0.9]", p.StartMat, p.EndMat)
	}
}

func TestParsePipelineNoRange(t *testing.T) {
	input := "pipeline Customer"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	p := doc.Pipelines[0]
	if p.HasRange {
		t.Error("expected HasRange to be false")
	}
}

func TestParseNote(t *testing.T) {
	input := "note +future development [0.16, 0.36]"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Notes) != 1 {
		t.Fatalf("got %d notes, want 1", len(doc.Notes))
	}
	n := doc.Notes[0]
	if n.Text != "future development" {
		t.Errorf("text = %q, want %q", n.Text, "future development")
	}
	if n.Visibility != 0.16 || n.Maturity != 0.36 {
		t.Errorf("position = [%v, %v], want [0.16, 0.36]", n.Visibility, n.Maturity)
	}
}

func TestParseSubmap(t *testing.T) {
	input := "submap Website [0.83, 0.50] url(submapUrl)\nurl submapUrl [https://example.com]"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Submaps) != 1 {
		t.Fatalf("got %d submaps, want 1", len(doc.Submaps))
	}
	sm := doc.Submaps[0]
	if sm.Name != "Website" {
		t.Errorf("name = %q, want Website", sm.Name)
	}
	if sm.URLName != "submapUrl" {
		t.Errorf("urlName = %q, want submapUrl", sm.URLName)
	}
	if doc.URLs["submapUrl"] != "https://example.com" {
		t.Errorf("url = %q, want https://example.com", doc.URLs["submapUrl"])
	}
}

func TestParseEvolution(t *testing.T) {
	input := "evolution Genesis->Custom->Product->Commodity"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	want := [4]string{"Genesis", "Custom", "Product", "Commodity"}
	if doc.Evolution != want {
		t.Errorf("evolution = %v, want %v", doc.Evolution, want)
	}
}

func TestParseArea(t *testing.T) {
	input := "pioneers [0.90, 0.10, 0.70, 0.40]"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Areas) != 1 {
		t.Fatalf("got %d areas, want 1", len(doc.Areas))
	}
	a := doc.Areas[0]
	if a.Kind != "pioneers" {
		t.Errorf("kind = %q, want pioneers", a.Kind)
	}
	if a.Vis1 != 0.90 || a.Mat1 != 0.10 || a.Vis2 != 0.70 || a.Mat2 != 0.40 {
		t.Errorf("coords = [%v, %v, %v, %v], want [0.90, 0.10, 0.70, 0.40]",
			a.Vis1, a.Mat1, a.Vis2, a.Mat2)
	}
}

func TestParseSignal(t *testing.T) {
	input := "accelerator foobar [0.1, 0.8]"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Signals) != 1 {
		t.Fatalf("got %d signals, want 1", len(doc.Signals))
	}
	s := doc.Signals[0]
	if s.Kind != "accelerator" || s.Name != "foobar" {
		t.Errorf("signal = {%q, %q}, want {accelerator, foobar}", s.Kind, s.Name)
	}
}

func TestParseStyle(t *testing.T) {
	input := "style wardley"
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if doc.Style != "wardley" {
		t.Errorf("style = %q, want wardley", doc.Style)
	}
}

func TestParseTeashopFile(t *testing.T) {
	doc, err := parseFile("testdata/teashop.owm")
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title != "Tea Shop" {
		t.Errorf("title = %q, want %q", doc.Title, "Tea Shop")
	}
	if len(doc.Anchors) != 2 {
		t.Errorf("anchors = %d, want 2", len(doc.Anchors))
	}
	if len(doc.Components) != 7 {
		t.Errorf("components = %d, want 7", len(doc.Components))
	}
	if len(doc.Evolves) != 1 {
		t.Errorf("evolves = %d, want 1", len(doc.Evolves))
	}
	if len(doc.Edges) != 7 {
		t.Errorf("edges = %d, want 7", len(doc.Edges))
	}
}

func TestParseCompleteFile(t *testing.T) {
	doc, err := parseFile("testdata/complete.owm")
	if err != nil {
		t.Fatal(err)
	}
	if doc.Title != "Tea Shop" {
		t.Errorf("title = %q", doc.Title)
	}
	if doc.Style != "wardley" {
		t.Errorf("style = %q", doc.Style)
	}
	if doc.Evolution[0] != "Genesis" {
		t.Errorf("evolution[0] = %q", doc.Evolution[0])
	}
	if len(doc.Anchors) != 2 {
		t.Errorf("anchors = %d, want 2", len(doc.Anchors))
	}
	if len(doc.Components) != 8 {
		t.Errorf("components = %d, want 8", len(doc.Components))
	}
	if len(doc.Evolves) != 2 {
		t.Errorf("evolves = %d, want 2", len(doc.Evolves))
	}
	if len(doc.Pipelines) != 1 {
		t.Errorf("pipelines = %d, want 1", len(doc.Pipelines))
	}
	if len(doc.Notes) != 1 {
		t.Errorf("notes = %d, want 1", len(doc.Notes))
	}
	if len(doc.Submaps) != 1 {
		t.Errorf("submaps = %d, want 1", len(doc.Submaps))
	}
	if len(doc.Areas) != 3 {
		t.Errorf("areas = %d, want 3", len(doc.Areas))
	}
	if len(doc.Signals) != 2 {
		t.Errorf("signals = %d, want 2", len(doc.Signals))
	}
}

func parseFile(path string) (*OWMDocument, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}
