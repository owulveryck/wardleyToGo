package wardleyToGo

import (
	"image"
	"reflect"
	"testing"
)

func Test_calcCoords(t *testing.T) {
	type args struct {
		p      image.Point
		bounds image.Rectangle
	}
	tests := []struct {
		name string
		args args
		want image.Point
	}{
		{
			"simple",
			args{
				p:      image.Pt(50, 50),
				bounds: image.Rect(0, 0, 100, 100),
			},
			image.Pt(50, 50),
		},
		{
			"simple upper left",
			args{
				p:      image.Pt(0, 0),
				bounds: image.Rect(0, 0, 100, 100),
			},
			image.Pt(0, 0),
		},
		{
			"simple bottom right",
			args{
				p:      image.Pt(100, 100),
				bounds: image.Rect(0, 0, 100, 100),
			},
			image.Pt(100, 100),
		},
		{
			"overflow bottom right",
			args{
				p:      image.Pt(200, 200),
				bounds: image.Rect(0, 0, 100, 100),
			},
			image.Pt(200, 200),
		},
		{
			"translation no scale",
			args{
				p:      image.Pt(50, 50),
				bounds: image.Rect(12, 22, 112, 122),
			},
			image.Pt(62, 72),
		},
		{
			"scale X no translation",
			args{
				p:      image.Pt(50, 50),
				bounds: image.Rect(0, 0, 150, 100),
			},
			image.Pt(75, 50),
		},
		{
			"scale X,Y no translation",
			args{
				p:      image.Pt(50, 50),
				bounds: image.Rect(0, 0, 150, 200),
			},
			image.Pt(75, 100),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcCoords(tt.args.p, tt.args.bounds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calcCoords() = %v, want %v", got, tt.want)
			}
		})
	}
}
