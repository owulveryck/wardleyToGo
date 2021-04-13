package owm

import (
	"image"
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

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
		want    *wardley.Anchor
		wantErr bool
	}{
		{
			"simple without coordinates",
			args{
				s: newScanner(`bla`),
			},
			&wardley.Anchor{
				Placement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Label:     `bla`,
			},
			false,
		},
		{
			"two words without coordinates",
			args{
				s: newScanner(`bla   bla`),
			},
			&wardley.Anchor{
				Placement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Label:     `bla bla`,
			},
			false,
		},
		{
			"two words with coordinates",
			args{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&wardley.Anchor{
				Placement: image.Point{30, 60},
				Label:     `bla bla`,
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
