package wtg2bin

import (
	"encoding/base64"
	"os"
	"reflect"
	"testing"

	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

func TestURLRoundTrip(t *testing.T) {
	doc := &wtg2.Document{
		Title:    "Test Map",
		Date:     "2026-01-15",
		Author:   "Author",
		Stages:   [4]string{"A", "B", "C", "D"},
		Doctrine: "context",
		Legend:   true,
		Nodes: []*wtg2.NodeDecl{
			{Name: "Comp A", Kind: wtg2.KindComponent, Evolution: "III.5", Visibility: -1},
			{Name: "Comp B", Kind: wtg2.KindAnchor, Evolution: "II.3", Visibility: 0.8},
		},
		Edges: []*wtg2.EdgeDecl{
			{From: "Comp A", To: "Comp B"},
		},
	}

	s, err := EncodeURL(doc)
	if err != nil {
		t.Fatalf("EncodeURL: %v", err)
	}

	t.Logf("URL string length: %d chars", len(s))

	got, err := DecodeURL(s)
	if err != nil {
		t.Fatalf("DecodeURL: %v", err)
	}

	if !reflect.DeepEqual(doc, got) {
		t.Error("URL round-trip mismatch")
	}
}

func TestURLNoPadding(t *testing.T) {
	doc := &wtg2.Document{Title: "X", Stages: [4]string{"", "", "", ""}}
	s, err := EncodeURL(doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) > 0 && s[len(s)-1] == '=' {
		t.Error("URL string should not have base64 padding")
	}
}

func TestURLWithPaddingInput(t *testing.T) {
	doc := &wtg2.Document{Title: "Padded", Stages: [4]string{"", "", "", ""}}
	data, _ := Encode(doc)
	padded := base64.URLEncoding.EncodeToString(data)

	got, err := DecodeURL(padded)
	if err != nil {
		t.Fatalf("DecodeURL with padding: %v", err)
	}
	if !reflect.DeepEqual(doc, got) {
		t.Error("padded input round-trip mismatch")
	}
}

func TestURLExampleFile(t *testing.T) {
	f, err := os.Open("../../wtg2/example.wtg2")
	if err != nil {
		t.Skipf("cannot open example: %v", err)
	}
	defer f.Close()

	p, err := wtg2.NewParser(f)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	s, err := EncodeURL(doc)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("example.wtg2 URL: %d chars", len(s))

	got, err := DecodeURL(s)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(doc, got) {
		t.Error("example file URL round-trip mismatch")
	}
}

func TestDecodeURLInvalid(t *testing.T) {
	_, err := DecodeURL("!!!not-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64url input")
	}
}
