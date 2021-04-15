package svg

import (
	"bytes"
	"encoding/xml"
	"image/color"
	"testing"
)

func Test_linearGradient(t *testing.T) {
	expected := `<linearGradient x1="0%" y1="0%" x2="100%" y2="0%"><stop offset="0%" stop-color="#c4c4c4"></stop><stop offset="30%" stop-color="#ffffff"></stop><stop offset="70%" stop-color="#ffffff"></stop><stop offset="100%" stop-color="#c4c4c4"></stop></linearGradient>`
	l := LinearGradient{
		X1: "0%",
		Y1: "0%",
		X2: "100%",
		Y2: "0%",
		Stops: []Stop{
			{
				Offset:    "0%",
				StopColor: Color{color.RGBA{196, 196, 196, 255}},
			},
			{
				Offset:    "30%",
				StopColor: Color{color.RGBA{255, 255, 255, 255}},
			},
			{
				Offset:    "70%",
				StopColor: Color{color.RGBA{255, 255, 255, 255}},
			},
			{
				Offset:    "100%",
				StopColor: Color{color.RGBA{196, 196, 196, 255}},
			},
		},
	}
	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	enc.Encode(l)
	if b.String() != expected {
		t.Fatal(b.String())
	}
}
