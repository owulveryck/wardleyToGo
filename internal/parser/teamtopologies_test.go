package parser

import (
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/internal/wardley"
)

func Test_parser_parseStreamAligned(t *testing.T) {

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
		want   *wardley.StreamAlignedTeam
	}{
		{
			"simple without coordinates",
			fields{
				s: newScanner(`bla`),
			},
			&wardley.StreamAlignedTeam{
				Coords: [4]int{wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:  `bla`,
			},
		},
		{
			"two words without coordinates",
			fields{
				s: newScanner(`bla   bla`),
			},
			&wardley.StreamAlignedTeam{
				Coords: [4]int{wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:  `bla bla`,
			},
		},
		{
			"two words with coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&wardley.StreamAlignedTeam{
				Coords: [4]int{40, 30, wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:  `bla bla`,
			},
		},
		{
			"two words with coordinates and label coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3, 0.2, 0.1]`),
			},
			&wardley.StreamAlignedTeam{
				Coords: [4]int{40, 30, 20, 10},
				Label:  `bla bla`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				s: tt.fields.s,
			}
			if got, _ := p.parseStreamAligned(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.parseStreamAligned() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_parsePlatform(t *testing.T) {

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
		want   *wardley.PlatformTeam
	}{
		{
			"simple without coordinates",
			fields{
				s: newScanner(`bla`),
			},
			&wardley.PlatformTeam{
				Coords: [4]int{wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:  `bla`,
			},
		},
		{
			"two words without coordinates",
			fields{
				s: newScanner(`bla   bla`),
			},
			&wardley.PlatformTeam{
				Coords: [4]int{wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:  `bla bla`,
			},
		},
		{
			"two words with coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&wardley.PlatformTeam{
				Coords: [4]int{40, 30, wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:  `bla bla`,
			},
		},
		{
			"two words with coordinates and label coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3, 0.2, 0.1]`),
			},
			&wardley.PlatformTeam{
				Coords: [4]int{40, 30, 20, 10},
				Label:  `bla bla`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				s: tt.fields.s,
			}
			if got, _ := p.parsePlatform(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.parsePlatform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_parseEnabling(t *testing.T) {

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
		want   *wardley.EnablingTeam
	}{
		{
			"simple without coordinates",
			fields{
				s: newScanner(`bla`),
			},
			&wardley.EnablingTeam{
				Coords: [4]int{wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:  `bla`,
			},
		},
		{
			"two words without coordinates",
			fields{
				s: newScanner(`bla   bla`),
			},
			&wardley.EnablingTeam{
				Coords: [4]int{wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:  `bla bla`,
			},
		},
		{
			"two words with coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&wardley.EnablingTeam{
				Coords: [4]int{40, 30, wardley.UndefinedCoord, wardley.UndefinedCoord},
				Label:  `bla bla`,
			},
		},
		{
			"two words with coordinates and label coordinates",
			fields{
				s: newScanner(`bla   bla [0.4, 0.3, 0.2, 0.1]`),
			},
			&wardley.EnablingTeam{
				Coords: [4]int{40, 30, 20, 10},
				Label:  `bla bla`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				s: tt.fields.s,
			}
			if got, _ := p.parseEnabling(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.parseEnabling() = %v, want %v", got, tt.want)
			}
		})
	}
}
