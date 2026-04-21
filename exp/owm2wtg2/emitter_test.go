package owm2wtg2

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestEmitTeashop(t *testing.T) {
	f, err := os.Open("testdata/teashop.owm")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	doc, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := Emit(doc, &buf); err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	checks := []string{
		"title: Tea Shop",
		"anchor Business :",
		"anchor Public :",
		"Cup of Tea :",
		"Kettle :",
		">> ",
		"Business -> Cup of Tea",
		"Cup of Tea -> Cup",
		"Kettle -> Power",
	}

	for _, want := range checks {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q\nfull output:\n%s", want, output)
		}
	}
}

func TestEmitComponentWithInertia(t *testing.T) {
	doc := &OWMDocument{
		URLs: make(map[string]string),
		Components: []*OWMComponent{{
			Name:       "Kettle",
			Visibility: 0.43,
			Maturity:   0.35,
			Inertia:    true,
		}},
	}

	var buf bytes.Buffer
	if err := Emit(doc, &buf); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "!!") {
		t.Errorf("expected !! in output, got:\n%s", buf.String())
	}
}

func TestEmitComponentWithType(t *testing.T) {
	doc := &OWMDocument{
		URLs: make(map[string]string),
		Components: []*OWMComponent{{
			Name:       "Cup",
			Visibility: 0.73,
			Maturity:   0.78,
			Type:       "buy",
		}},
	}

	var buf bytes.Buffer
	if err := Emit(doc, &buf); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "(buy)") {
		t.Errorf("expected (buy) in output, got:\n%s", buf.String())
	}
}

func TestEmitComponentWithEvolve(t *testing.T) {
	doc := &OWMDocument{
		URLs: make(map[string]string),
		Components: []*OWMComponent{{
			Name:       "Kettle",
			Visibility: 0.43,
			Maturity:   0.35,
		}},
		Evolves: []*OWMEvolve{{
			Name:     "Kettle",
			Maturity: 0.62,
		}},
	}

	var buf bytes.Buffer
	if err := Emit(doc, &buf); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), ">>") {
		t.Errorf("expected >> in output, got:\n%s", buf.String())
	}
}

func TestEmitEvolveWithRename(t *testing.T) {
	doc := &OWMDocument{
		URLs: make(map[string]string),
		Components: []*OWMComponent{{
			Name:       "Physical",
			Visibility: 0.5,
			Maturity:   0.3,
		}},
		Evolves: []*OWMEvolve{{
			Name:     "Physical",
			NewName:  "Virtual",
			Maturity: 0.8,
		}},
	}

	var buf bytes.Buffer
	if err := Emit(doc, &buf); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, ">>") {
		t.Errorf("expected >> in output, got:\n%s", output)
	}
	if !strings.Contains(output, "// OWM: evolved name") {
		t.Errorf("expected rename comment in output, got:\n%s", output)
	}
}

func TestEmitEdgeTypes(t *testing.T) {
	doc := &OWMDocument{
		URLs: make(map[string]string),
		Edges: []*OWMEdge{
			{From: "A", To: "B", FlowType: "regular"},
			{From: "C", To: "D", FlowType: "bidirectional"},
			{From: "E", To: "F", FlowType: "labeled", Label: "cost"},
			{From: "G", To: "H", FlowType: "past"},
			{From: "I", To: "J", FlowType: "future"},
		},
	}

	var buf bytes.Buffer
	if err := Emit(doc, &buf); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	checks := []string{
		"A -> B",
		"C <-> D",
		"E -[cost]-> F",
		"G -> H",
		"// OWM: past flow",
		"I -> J",
	}
	for _, want := range checks {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q\nfull output:\n%s", want, output)
		}
	}
}

func TestEmitPipeline(t *testing.T) {
	doc := &OWMDocument{
		URLs: make(map[string]string),
		Components: []*OWMComponent{
			{Name: "X", Visibility: 0.5, Maturity: 0.7},
			{Name: "Y", Visibility: 0.5, Maturity: 0.8},
		},
		Edges: []*OWMEdge{
			{From: "Parent", To: "X", FlowType: "regular"},
			{From: "Parent", To: "Y", FlowType: "regular"},
		},
		Pipelines: []*OWMPipeline{{
			Name:     "Parent",
			StartMat: 0.6,
			EndMat:   0.9,
			HasRange: true,
		}},
	}

	var buf bytes.Buffer
	if err := Emit(doc, &buf); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, "pipeline Parent {") {
		t.Errorf("expected pipeline block in output, got:\n%s", output)
	}
	if !strings.Contains(output, "X :") {
		t.Errorf("expected member X in pipeline, got:\n%s", output)
	}
}

func TestEmitSignals(t *testing.T) {
	doc := &OWMDocument{
		URLs: make(map[string]string),
		Signals: []*OWMSignal{
			{Kind: "accelerator", Name: "AI", Visibility: 0.5, Maturity: 0.3},
			{Kind: "deaccelerator", Name: "Legacy", Visibility: 0.5, Maturity: 0.7},
		},
	}

	var buf bytes.Buffer
	if err := Emit(doc, &buf); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, "signal accelerating on AI") {
		t.Errorf("expected accelerating signal, got:\n%s", output)
	}
	if !strings.Contains(output, "signal declining on Legacy") {
		t.Errorf("expected declining signal, got:\n%s", output)
	}
}

func TestEmitSubmap(t *testing.T) {
	doc := &OWMDocument{
		URLs: map[string]string{"mapUrl": "https://example.com/map"},
		Submaps: []*OWMSubmap{{
			Name:       "Website",
			Visibility: 0.83,
			Maturity:   0.50,
			URLName:    "mapUrl",
		}},
	}

	var buf bytes.Buffer
	if err := Emit(doc, &buf); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, "submap Website :") {
		t.Errorf("expected submap in output, got:\n%s", output)
	}
	if !strings.Contains(output, "// OWM URL: https://example.com/map") {
		t.Errorf("expected URL comment in output, got:\n%s", output)
	}
}

func TestEmitAreas(t *testing.T) {
	doc := &OWMDocument{
		URLs: make(map[string]string),
		Areas: []*OWMArea{{
			Kind: "pioneers",
			Vis1: 0.90, Mat1: 0.10, Vis2: 0.70, Mat2: 0.40,
		}},
	}

	var buf bytes.Buffer
	if err := Emit(doc, &buf); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, "group Pioneers {") {
		t.Errorf("expected group Pioneers in output, got:\n%s", output)
	}
	if !strings.Contains(output, "team: pioneer") {
		t.Errorf("expected team: pioneer in output, got:\n%s", output)
	}
}

func TestEmitCompleteRoundTrip(t *testing.T) {
	f, err := os.Open("testdata/complete.owm")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	doc, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := Emit(doc, &buf); err != nil {
		t.Fatal(err)
	}

	output := buf.String()

	required := []string{
		"title: Tea Shop",
		"stages: Genesis, Custom, Product, Commodity",
		"anchor Business",
		"Cup of Tea :",
		"Kettle :",
		"!!",
		">>",
		"(buy)",
		"submap Supplier",
		"pipeline Cup",
		"Business -> Cup of Tea",
		"Cup of Tea <-> Market Research",
		"-[$0.10]->",
		"// OWM: past flow",
		"group Pioneers",
		"group Settlers",
		"group Townplanners",
		"team: pioneer",
		"signal accelerating",
		"signal declining",
		"// OWM note",
		"// OWM style: wardley",
		"// OWM: market component",
		"// OWM: evolved name",
	}

	for _, want := range required {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q\nfull output:\n%s", want, output)
		}
	}
}
