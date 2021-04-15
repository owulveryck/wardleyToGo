package svgmap_test

import (
	"image"
	"log"
	"os"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
)

func Example_style() {
	f, _ := os.Create("/tmp/test.svg")
	defer f.Close()
	e, err := svgmap.NewEncoder(f, image.Rect(0, 0, 1100, 900), image.Rect(30, 50, 1070, 850))
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	style := svgmap.NewWardleyStyle(svgmap.DefaultEvolution)
	e.Init(style)
}
