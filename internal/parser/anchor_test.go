package parser

import (
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/wardley"
)

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
				Coords: [2]int{wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:  `bla`,
			},
		},
		{
			"two words without coordinates",
			fields{
				s: newScanner(`bla   bla`),
			},
			&wardley.Anchor{
				Coords: [2]int{wardley.UndefinedCoord, wardley.UndefinedCoord},
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
			if got, _ := p.parseAnchor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.parseAnchor() = %v, want %v", got, tt.want)
			}
		})
	}
}
