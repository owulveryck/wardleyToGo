package owm

import (
	"image"
	"image/color"
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	"github.com/owulveryck/wardleyToGo/components"
	"github.com/owulveryck/wardleyToGo/components/wardley"
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
		want    *wardley.Component
		wantErr bool
	}{
		{
			"simple without coordinates",
			args{
				s: newScanner(`bla`),
			},
			&wardley.Component{
				Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Label:          `bla`,
				LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
				RenderingLayer: wardley.DefaultComponentRenderingLayer,
				Anchor:         1,
				Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
			},
			false,
		},
		{
			"two words without coordinates",
			args{
				s: newScanner(`bla   bla`),
			},
			&wardley.Component{
				Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Label:          `bla bla`,
				LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
				RenderingLayer: wardley.DefaultComponentRenderingLayer,
				Anchor:         1,
				Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
			},
			false,
		},
		{
			"two words with coordinates",
			args{
				s: newScanner(`bla   bla [0.4, 0.3]`),
			},
			&wardley.Component{
				Placement:      image.Point{30, 60},
				Label:          `bla bla`,
				LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
				RenderingLayer: wardley.DefaultComponentRenderingLayer,
				Anchor:         1,
				Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
			},
			false,
		},
		{
			"two words with coordinates and label coordinates",
			args{
				s: newScanner(`bla   bla [0.4, 0.3] label [12,12]`),
			},
			&wardley.Component{
				Placement:      image.Point{30, 60},
				Label:          `bla bla`,
				LabelPlacement: image.Point{12, 12},
				RenderingLayer: wardley.DefaultComponentRenderingLayer,
				Anchor:         1,
				Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
			},
			false,
		},
		{
			"two words with coordinates and negative label coordinates",
			args{
				s: newScanner(`bla   bla [0.4, 0.3] label [-12,12]`),
			},
			&wardley.Component{
				Placement:      image.Point{30, 60},
				Label:          `bla bla`,
				LabelPlacement: image.Point{-12, 12},
				RenderingLayer: wardley.DefaultComponentRenderingLayer,
				Anchor:         1,
				Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
			},
			false,
		},
		// TODO: Add test cases.
		{
			"two words with with build type",
			args{
				s: newScanner(`bla   bla (build)`),
			},
			&wardley.Component{
				Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Label:          `bla bla`,
				LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Type:           wardley.BuildComponent,
				RenderingLayer: wardley.DefaultComponentRenderingLayer,
				Anchor:         1,
				Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
			},
			false,
		},
		{
			"two words with with build type",
			args{
				s: newScanner(`bla   bla (build)`),
			},
			&wardley.Component{
				Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Label:          `bla bla`,
				LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Type:           wardley.BuildComponent,
				RenderingLayer: wardley.DefaultComponentRenderingLayer,
				Anchor:         1,
				Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
			},
			false,
		},
		{
			"two words with with buy type",
			args{
				s: newScanner(`bla   bla (buy)`),
			},
			&wardley.Component{
				Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Label:          `bla bla`,
				LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Type:           wardley.BuyComponent,
				Anchor:         1,
				RenderingLayer: wardley.DefaultComponentRenderingLayer,
				Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
			},
			false,
		},
		{
			"two words with with outsource type",
			args{
				s: newScanner(`bla   bla (outsource)`),
			},
			&wardley.Component{
				Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Label:          `bla bla`,
				LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Type:           wardley.OutsourceComponent,
				Anchor:         1,
				RenderingLayer: wardley.DefaultComponentRenderingLayer,
				Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
			},
			false,
		},
		{
			"two words with with dataProduct type",
			args{
				s: newScanner(`bla   bla (dataProduct)`),
			},
			&wardley.Component{
				Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Label:          `bla bla`,
				LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Type:           wardley.DataProductComponent,
				RenderingLayer: wardley.DefaultComponentRenderingLayer,
				Anchor:         1,
				Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
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
			//got, err := scanComponent(tt.args.s, tt.args.id)
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
