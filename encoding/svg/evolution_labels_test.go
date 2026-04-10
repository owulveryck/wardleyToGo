package svgmap_test

import (
	"bytes"
	"image"
	"strings"
	"testing"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
)

func TestEvolutionAxisLabels_WithZoneLabels(t *testing.T) {
	var buf bytes.Buffer
	e, err := svgmap.NewEncoder(&buf, image.Rect(0, 0, 1100, 900), image.Rect(30, 50, 1070, 835))
	if err != nil {
		t.Fatal(err)
	}
	stages := []svgmap.Evolution{
		{Position: 0.0, Label: "I", ZoneLabel: "Genesis"},
		{Position: 0.25, Label: "II", ZoneLabel: "Custom-Built"},
		{Position: 0.50, Label: "III", ZoneLabel: "Product"},
		{Position: 0.75, Label: "IV", ZoneLabel: "Commodity"},
	}
	style := svgmap.NewOctoStyle(stages)
	e.Init(style)
	e.Close()
	output := buf.String()

	// Roman numerals must be present
	for _, numeral := range []string{">I<", ">II<", ">III<", ">IV<"} {
		if !strings.Contains(output, numeral) {
			t.Errorf("missing Roman numeral %q in SVG output", numeral)
		}
	}
	// Zone labels must be present
	for _, label := range []string{"Genesis", "Custom-Built", "Product", "Commodity"} {
		if !strings.Contains(output, label) {
			t.Errorf("missing zone label %q in SVG output", label)
		}
	}
}

func TestEvolutionAxisLabels_WithoutZoneLabels(t *testing.T) {
	var buf bytes.Buffer
	e, err := svgmap.NewEncoder(&buf, image.Rect(0, 0, 1100, 900), image.Rect(30, 50, 1070, 835))
	if err != nil {
		t.Fatal(err)
	}
	stages := []svgmap.Evolution{
		{Position: 0.0, Label: "I", ZoneLabel: ""},
		{Position: 0.25, Label: "II", ZoneLabel: ""},
		{Position: 0.50, Label: "III", ZoneLabel: ""},
		{Position: 0.75, Label: "IV", ZoneLabel: ""},
	}
	style := svgmap.NewOctoStyle(stages)
	e.Init(style)
	e.Close()
	output := buf.String()

	// Roman numerals must be present
	for _, numeral := range []string{">I<", ">II<", ">III<", ">IV<"} {
		if !strings.Contains(output, numeral) {
			t.Errorf("missing Roman numeral %q in SVG output", numeral)
		}
	}
	// "Genesis" etc. should NOT appear
	for _, label := range []string{"Genesis", "Custom-Built", "Product", "Commodity"} {
		if strings.Contains(output, label) {
			t.Errorf("unexpected zone label %q in SVG output when ZoneLabel is empty", label)
		}
	}
}
