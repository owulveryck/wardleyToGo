package parser

import (
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo"
)

/*
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
		want   *plan.Anchor
	}{
		{
			"simple without coordinates",
			fields{
				s: newScanner(`bla`),
			},
			&plan.Anchor{
				Coords: [2]int{plan.UndefinedCoord, plan.UndefinedCoord},
				Label:  `bla`,
			},
		},
		{
			"two words without coordinates",
			fields{
				s: newScanner(`bla   bla`),
			},
			&plan.Anchor{
				Coords: [2]int{plan.UndefinedCoord, plan.UndefinedCoord},
				Label:  `bla bla`,
			},
		},
		{
			"two words with coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&plan.Anchor{
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

*/

func Test_scanAnchor(t *testing.T) {
	newScanner := func(content string) *scanner.Scanner {
		var s scanner.Scanner
		s.Whitespace ^= 1 << '\n' // don't skip tabs and new lines
		s.Init(strings.NewReader(content))
		return &s
	}
	type args struct {
		s  *scanner.Scanner
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *wardleyToGo.Anchor
		wantErr bool
	}{
		{
			"simple without coordinates",
			args{
				s: newScanner(`bla`),
			},
			&wardleyToGo.Anchor{
				Coords: [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Label:  `bla`,
			},
			false,
		},
		{
			"two words without coordinates",
			args{
				s: newScanner(`bla   bla`),
			},
			&wardleyToGo.Anchor{
				Coords: [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Label:  `bla bla`,
			},
			false,
		},
		{
			"two words with coordinates",
			args{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&wardleyToGo.Anchor{
				Coords: [2]int{40, 30},
				Label:  `bla bla`,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := scanAnchor(tt.args.s, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanAnchor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scanAnchor() = %v, want %v", got, tt.want)
			}
		})
	}
}
