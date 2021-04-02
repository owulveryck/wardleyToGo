package parser

import (
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo"
)

func Test_scanComponent(t *testing.T) {
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
		want    *wardleyToGo.Component
		wantErr bool
	}{
		{
			"simple without coordinates",
			args{
				s: newScanner(`bla`),
			},
			&wardleyToGo.Component{
				Coords:      [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Label:       `bla`,
				LabelCoords: [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
			},
			false,
		},
		{
			"two words without coordinates",
			args{
				s: newScanner(`bla   bla`),
			},
			&wardleyToGo.Component{
				Coords:      [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Label:       `bla bla`,
				LabelCoords: [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
			},
			false,
		},
		{
			"two words with coordinates",
			args{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&wardleyToGo.Component{
				Coords:      [2]int{40, 30},
				Label:       `bla bla`,
				LabelCoords: [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
			},
			false,
		},
		{
			"two words with coordinates and label coordinates",
			args{
				s: newScanner(`bla   bla [0.4, 0.3] label [12,12]`),
			},
			&wardleyToGo.Component{
				Coords:      [2]int{40, 30},
				Label:       `bla bla`,
				LabelCoords: [2]int{12, 12},
			},
			false,
		},
		{
			"two words with coordinates and negative label coordinates",
			args{
				s: newScanner(`bla   bla [0.4, 0.3] label [-12,12]`),
			},
			&wardleyToGo.Component{
				Coords:      [2]int{40, 30},
				Label:       `bla bla`,
				LabelCoords: [2]int{-12, 12},
			},
			false,
		},
		// TODO: Add test cases.
		{
			"two words with with build type",
			args{
				s: newScanner(`bla   bla (build)`),
			},
			&wardleyToGo.Component{
				Coords:      [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Label:       `bla bla`,
				LabelCoords: [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Type:        wardleyToGo.BuildComponent,
			},
			false,
		},
		{
			"two words with with build type",
			args{
				s: newScanner(`bla   bla (build)`),
			},
			&wardleyToGo.Component{
				Coords:      [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Label:       `bla bla`,
				LabelCoords: [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Type:        wardleyToGo.BuildComponent,
			},
			false,
		},
		{
			"two words with with buy type",
			args{
				s: newScanner(`bla   bla (buy)`),
			},
			&wardleyToGo.Component{
				Coords:      [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Label:       `bla bla`,
				LabelCoords: [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Type:        wardleyToGo.BuyComponent,
			},
			false,
		},
		{
			"two words with with outsource type",
			args{
				s: newScanner(`bla   bla (outsource)`),
			},
			&wardleyToGo.Component{
				Coords:      [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Label:       `bla bla`,
				LabelCoords: [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Type:        wardleyToGo.OutsourceComponent,
			},
			false,
		},
		{
			"two words with with dataProduct type",
			args{
				s: newScanner(`bla   bla (dataProduct)`),
			},
			&wardleyToGo.Component{
				Coords:      [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Label:       `bla bla`,
				LabelCoords: [2]int{wardleyToGo.UndefinedCoord, wardleyToGo.UndefinedCoord},
				Type:        wardleyToGo.DataProductComponent,
			},
			false,
		},
		{
			"two words with with unhandled type",
			args{
				s: newScanner(`bla   bla (XXXXX)`),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := scanComponent(tt.args.s, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanComponent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scanComponent() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
