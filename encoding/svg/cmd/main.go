package main

import (
	"image"
	"os"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
)

func main() {
	svgmap.Encode(nil, os.Stdout, 1025, 825, image.Rect(25, 20, 1000, 800))
}
