package parser

import (
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/plan"
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
		want   *plan.Component
	}{
		{
			"simple without coordinates",
			fields{
				s: newScanner(`bla`),
			},
			&plan.Component{
				Coords:      [2]int{plan.UndefinedCoord, plan.UndefinedCoord},
				Label:       `bla`,
				LabelCoords: [2]int{plan.UndefinedCoord, plan.UndefinedCoord},
			},
		},
		{
			"two words without coordinates",
			fields{
				s: newScanner(`bla   bla`),
			},
			&plan.Component{
				Coords:      [2]int{plan.UndefinedCoord, plan.UndefinedCoord},
				Label:       `bla bla`,
				LabelCoords: [2]int{plan.UndefinedCoord, plan.UndefinedCoord},
			},
		},
		{
			"two words with coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&plan.Component{
				Coords:      [2]int{40, 30},
				Label:       `bla bla`,
				LabelCoords: [2]int{plan.UndefinedCoord, plan.UndefinedCoord},
			},
		},
		{
			"two words with coordinates and label coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3] label [12,12]`),
			},
			&plan.Component{
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
