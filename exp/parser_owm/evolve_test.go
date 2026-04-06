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

func Test_scanEvolve(t *testing.T) {
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
		want    *wardley.EvolvedComponent
		wantErr bool
	}{
		{
			"simple without coordinates",
			args{
				s: newScanner(`bla`),
			},
			&wardley.EvolvedComponent{
				&wardley.Component{
					Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
					Label:          `bla`,
					Anchor:         1,
					LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
					RenderingLayer: wardley.DefaultComponentRenderingLayer,
					Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
				},
			},
			false,
		},
		{
			"two words without coordinates",
			args{
				s: newScanner(`bla   bla`),
			},
			&wardley.EvolvedComponent{
				&wardley.Component{
					Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
					Label:          `bla bla`,
					Anchor:         1,
					LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
					RenderingLayer: wardley.DefaultComponentRenderingLayer,
					Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
				},
			},
			false,
		},
		{
			"two words with coordinates",
			args{
				s: newScanner(`bla   bla 0.3`),
			},
			&wardley.EvolvedComponent{
				&wardley.Component{
					Placement:      image.Point{30, components.UndefinedCoord},
					Anchor:         1,
					Label:          `bla bla`,
					LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
					RenderingLayer: wardley.DefaultComponentRenderingLayer,
					Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
				},
			},
			false,
		},
		{
			"two words with coordinates and label coordinates",
			args{
				s: newScanner(`bla   bla 0.3 label [12,12]`),
			},
			&wardley.EvolvedComponent{
				&wardley.Component{
					Placement:      image.Point{30, components.UndefinedCoord},
					Label:          `bla bla`,
					Anchor:         1,
					LabelPlacement: image.Point{12, 12},
					RenderingLayer: wardley.DefaultComponentRenderingLayer,
					Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
				},
			},
			false,
		},
		{
			"two words with coordinates and negative label coordinates",
			args{
				s: newScanner(`bla   bla 0.3 label [-12,12]`),
			},
			&wardley.EvolvedComponent{
				&wardley.Component{
					Placement:      image.Point{30, components.UndefinedCoord},
					Label:          `bla bla`,
					Anchor:         1,
					LabelPlacement: image.Point{-12, 12},
					RenderingLayer: wardley.DefaultComponentRenderingLayer,
					Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
				},
			},
			false,
		},
		// TODO: Add test cases.
		{
			"two words with with build type",
			args{
				s: newScanner(`bla   bla (build)`),
			},
			&wardley.EvolvedComponent{
				&wardley.Component{
					Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
					Label:          `bla bla`,
					Anchor:         1,
					LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
					Type:           wardley.BuildComponent,
					RenderingLayer: wardley.DefaultComponentRenderingLayer,
					Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
				},
			},
			false,
		},
		{
			"two words with with buy type",
			args{
				s: newScanner(`bla   bla (buy)`),
			},
			&wardley.EvolvedComponent{
				&wardley.Component{
					Anchor:         1,
					Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
					Label:          `bla bla`,
					LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
					Type:           wardley.BuyComponent,
					RenderingLayer: wardley.DefaultComponentRenderingLayer,
					Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
				},
			},
			false,
		},
		{
			"two words with with outsource type",
			args{
				s: newScanner(`bla   bla (outsource)`),
			},
			&wardley.EvolvedComponent{
				&wardley.Component{
					Anchor:         1,
					Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
					Label:          `bla bla`,
					LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
					Type:           wardley.OutsourceComponent,
					RenderingLayer: wardley.DefaultComponentRenderingLayer,
					Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
				},
			},
			false,
		},
		{
			"two words with with dataProduct type",
			args{
				s: newScanner(`bla   bla (dataProduct)`),
			},
			&wardley.EvolvedComponent{
				&wardley.Component{
					Placement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
					Anchor:         1,
					Label:          `bla bla`,
					LabelPlacement: image.Point{components.UndefinedCoord, components.UndefinedCoord},
					Type:           wardley.DataProductComponent,
					RenderingLayer: wardley.DefaultComponentRenderingLayer,
					Configured:     false, EvolutionPos: 0, Color: color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff},
				},
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := scanEvolve(tt.args.s, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanEvolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scanEvolve() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
