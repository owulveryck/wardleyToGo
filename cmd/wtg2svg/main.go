// wtg2svg reads a WTG2 Wardley Map description from stdin and writes
// an SVG rendering to stdout.
//
// Usage:
//
//	cat map.wtg2 | wtg2svg > map.svg
package main

import (
	"image"
	"log"
	"os"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

func main() {
	p, err := wtg2.NewParser(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	result, err := wtg2.BuildMap(doc)
	if err != nil {
		log.Fatal(err)
	}

	e, err := svgmap.NewEncoder(os.Stdout,
		image.Rect(0, 0, 1100, 900),
		image.Rect(30, 50, 1070, 850))
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()

	style := svgmap.NewOctoStyle(result.Stages)
	style.WithSpace = true
	e.Init(style)

	if err := e.Encode(result.Map); err != nil {
		log.Fatal(err)
	}
}
