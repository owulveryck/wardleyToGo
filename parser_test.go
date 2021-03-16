package main

import (
	"reflect"
	"strings"
	"testing"
	"text/scanner"
)

func TestParse(t *testing.T) {
	const src = `
title Tea Shop
anchor Business [0.95, 0.63]
anchor Public [0.95, 0.78]
component Cup of Tea [0.79, 0.61] label [19, -4]
component Cup [0.73, 0.78]
component Tea [0.63, 0.81]
component Hot Water [0.52, 0.80]
component Water [0.38, 0.82]
component Kettle [0.43, 0.35] label [-57, 4]
evolve Kettle 0.62 label [16, 7]
component Power [0.1, 0.7] label [-27, 20]
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
	p := newParser(strings.NewReader(src))
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
		want   *component
	}{
		{
			"simple without coordinates",
			fields{
				s: newScanner(`bla`),
			},
			&component{
				coords:     [2]int{-1, -1},
				label:      `bla`,
				labelCoord: [2]int{-1, -1},
			},
		},
		{
			"two words without coordinates",
			fields{
				s: newScanner(`bla   bla`),
			},
			&component{
				coords:     [2]int{-1, -1},
				label:      `bla bla`,
				labelCoord: [2]int{-1, -1},
			},
		},
		{
			"two words with coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&component{
				coords:     [2]int{40, 30},
				label:      `bla bla`,
				labelCoord: [2]int{-1, -1},
			},
		},
		{
			"two words with coordinates and label coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3] label [12,12]`),
			},
			&component{
				coords:     [2]int{40, 30},
				label:      `bla bla`,
				labelCoord: [2]int{12, 12},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
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
		want   *anchor
	}{
		{
			"simple without coordinates",
			fields{
				s: newScanner(`bla`),
			},
			&anchor{
				coords: [2]int{-1, -1},
				label:  `bla`,
			},
		},
		{
			"two words without coordinates",
			fields{
				s: newScanner(`bla   bla`),
			},
			&anchor{
				coords: [2]int{-1, -1},
				label:  `bla bla`,
			},
		},
		{
			"two words with coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&anchor{
				coords: [2]int{40, 30},
				label:  `bla bla`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				s: tt.fields.s,
			}
			if got := p.parseAnchor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.parseAnchor() = %v, want %v", got, tt.want)
			}
		})
	}
}
