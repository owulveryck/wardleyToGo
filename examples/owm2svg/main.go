package main

import (
	"image"
	"log"
	"os"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/owm"
)

func main() {
	p := owm.NewParser(os.Stdin)
	m, err := p.Parse() // the map
	if err != nil {
		log.Fatal(err)
	}
	e, err := svgmap.NewEncoder(os.Stdout, image.Rect(0, 0, 1955, 1100), image.Rect(30, 50, 1925, 1050))
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	style := svgmap.NewOctoStyle(svgmap.DefaultEvolution)
	e.Init(style)
	err = e.Encode(m)
	if err != nil {
		log.Fatal(err)
	}
}
