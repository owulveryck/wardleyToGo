package main

import (
	"image"
	"log"
	"os"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"

	"github.com/owulveryck/wardleyToGo/parser/wtg"
)

func main() {
	p := wtg.NewParser()

	err := p.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if len(p.InvalidEntries) != 0 {
		for _, err := range p.InvalidEntries {
			log.Println(err)
		}
	}
	setCoords(*p.WMap, true)
	imgArea := (p.ImageSize.Max.X - p.ImageSize.Min.X) * (p.ImageSize.Max.X - p.ImageSize.Min.Y)
	canvasArea := (p.MapSize.Max.X - p.MapSize.Min.X) * (p.MapSize.Max.X - p.MapSize.Min.Y)
	if imgArea == 0 || canvasArea == 0 {
		p.ImageSize = image.Rect(0, 0, 1100, 900)
		p.MapSize = image.Rect(30, 50, 1070, 850)
	}
	e, err := svgmap.NewEncoder(os.Stdout, p.ImageSize, p.MapSize)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	style := svgmap.NewOctoStyle(p.EvolutionStages)
	e.Init(style)
	err = e.Encode(p.WMap)
	if err != nil {
		log.Fatal(err)
	}

}
