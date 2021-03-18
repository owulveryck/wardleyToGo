package parser

import (
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/wardley"
)

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
				Coords:      [2]int{wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:       `bla`,
				LabelCoords: [2]int{wardley.UndefinedCoord, wardley.UndefinedCoord},
			},
		},
		{
			"two words without coordinates",
			fields{
				s: newScanner(`bla   bla`),
			},
			&wardley.Component{
				Coords:      [2]int{wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:       `bla bla`,
				LabelCoords: [2]int{wardley.UndefinedCoord, wardley.UndefinedCoord},
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
				LabelCoords: [2]int{wardley.UndefinedCoord, wardley.UndefinedCoord},
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
			if got, _ := p.parseComponent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.parseComponent() = %v, want %v", got, tt.want)
			}
		})
	}
}
