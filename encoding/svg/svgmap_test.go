package svgmap

import (
	"image"
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	Encode(nil, os.Stdout, 500, 300, image.Rect(10, 0, 500, 290))
	t.Fatal()
}
