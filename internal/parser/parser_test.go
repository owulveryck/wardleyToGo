package parser

import (
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/wardley"
)

func TestParse(t *testing.T) {
	const src = `
title Tea Shop
wardley.Anchor Business [0.95, 0.63]
wardley.Anchor Public [0.95, 0.78]
wardley.Component Cup of Tea [0.79, 0.61] label [19, -4]
wardley.Component Cup [0.73, 0.78]
wardley.Component Tea [0.63, 0.81]
wardley.Component Hot Water [0.52, 0.80]
wardley.Component Water [0.38, 0.82]
wardley.Component Kettle [0.43, 0.35] label [-57, 4]
evolve Kettle 0.62 label [16, 7]
wardley.Component Power [0.1, 0.7] label [-27, 20]
evolve Power 0.89 label [-12, 21]
Business->Cup of Tea
Public->Cup of Tea
Cup of Tea->Cup
Cup of Tea->Tea
Cup of Tea->Hot Water
Hot Water->Water
Hot Water->Kettle 
Kettle->Power

annotation 1 [[0.43,0.49],[0.08,0.79]] Standardising power allows Kettles to evolve faster
annotation 2 [0.48, 0.85] Hot water is obvious and well known
annotations [0.60, 0.02]

note +a generic note appeared [0.16, 0.36]

style wardley
`
	p := NewParser(strings.NewReader(src))
	p.Parse()
}

func Test_parser_parseComponent(t *testing.T) {

	newScanner := func(content string) *scanner.Scanner {
		var s scanner.Scanner
		s.Whitespace ^= 1 << '\n' // don't skip tabs and new lines
		s.Init(strings.NewReader(content))
		return &s
	}

	type fields struct {
		s *scanner.Scanner
	}
	tests := []struct {
		name   string
		fields fields
		want   *wardley.Component
	}{
		{
			"simple without coordinates",
			fields{
				s: newScanner(`bla`),
			},
			&wardley.Component{
				Coords:      [2]int{-1, -1},
				Label:       `bla`,
				LabelCoords: [2]int{-1, -1},
			},
		},
		{
			"two words without coordinates",
			fields{
				s: newScanner(`bla   bla`),
			},
			&wardley.Component{
				Coords:      [2]int{-1, -1},
				Label:       `bla bla`,
				LabelCoords: [2]int{-1, -1},
			},
		},
		{
			"two words with coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&wardley.Component{
				Coords:      [2]int{40, 30},
				Label:       `bla bla`,
				LabelCoords: [2]int{-1, -1},
			},
		},
		{
			"two words with coordinates and label coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3] label [12,12]`),
			},
			&wardley.Component{
				Coords:      [2]int{40, 30},
				Label:       `bla bla`,
				LabelCoords: [2]int{12, 12},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				s: tt.fields.s,
			}
			if got := p.parseComponent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.parseComponent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_parseAnchor(t *testing.T) {

	newScanner := func(content string) *scanner.Scanner {
		var s scanner.Scanner
		s.Whitespace ^= 1 << '\n' // don't skip tabs and new lines
		s.Init(strings.NewReader(content))
		return &s
	}

	type fields struct {
		s *scanner.Scanner
	}
	tests := []struct {
		name   string
		fields fields
		want   *wardley.Anchor
	}{
		{
			"simple without coordinates",
			fields{
				s: newScanner(`bla`),
			},
			&wardley.Anchor{
				Coords: [2]int{-1, -1},
				Label:  `bla`,
			},
		},
		{
			"two words without coordinates",
			fields{
				s: newScanner(`bla   bla`),
			},
			&wardley.Anchor{
				Coords: [2]int{-1, -1},
				Label:  `bla bla`,
			},
		},
		{
			"two words with coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&wardley.Anchor{
				Coords: [2]int{40, 30},
				Label:  `bla bla`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				s: tt.fields.s,
			}
			if got := p.parseAnchor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.parseAnchor() = %v, want %v", got, tt.want)
			}
		})
	}
}