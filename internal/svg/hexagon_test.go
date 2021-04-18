package svg

import (
	"image"
	"reflect"
	"testing"
)

func Test_hexCorner(t *testing.T) {
	type args struct {
		center image.Point
		size   float64
		i      float64
	}
	tests := []struct {
		name string
		args args
		want image.Point
	}{
		{
			"initial point",
			args{
				image.Pt(0, 0),
				100,
				0,
			},
			image.Pt(100, 0),
		},
		{
			"third point",
			args{
				image.Pt(0, 0),
				100,
				3,
			},
			image.Pt(-100, 0),
		},
		{
			"second point",
			args{
				image.Pt(0, 0),
				100,
				2,
			},
			image.Pt(-49, 86),
		},
		{
			"first point",
			args{
				image.Pt(0, 0),
				100,
				1,
			},
			image.Pt(50, 86),
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hexCorner(tt.args.center, tt.args.size, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hexCorner() = %v, want %v", got, tt.want)
			}
		})
	}
}
