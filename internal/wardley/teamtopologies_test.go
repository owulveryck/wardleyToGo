package wardley

import (
	"testing"

	svg "github.com/ajstarks/svgo"
)

func TestStreamAlignedTeam_SVG(t *testing.T) {
	type fields struct {
		Id     int64
		Coords [4]int
		Label  string
	}
	type args struct {
		svg       *svg.SVG
		width     int
		height    int
		padLeft   int
		padBottom int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamAlignedTeam{
				Id:     tt.fields.Id,
				Coords: tt.fields.Coords,
				Label:  tt.fields.Label,
			}
			s.SVG(tt.args.svg, tt.args.width, tt.args.height, tt.args.padLeft, tt.args.padBottom)
		})
	}
}
