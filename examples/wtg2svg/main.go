package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/wtg"
)

type configuration struct {
	Width          int  `default:"1100"`
	Height         int  `default:"900"`
	WithSpace      bool `default:"true"`
	WithControls   bool `default:"true"`
	WithValueChain bool `default:"true"`
	WithIndicators bool `default:"false"`
}

func main() {
	var config configuration
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *helpFlag {
		fmt.Printf("usage: %v [STDIN]\n", os.Args[0])
		envconfig.Usage("WTG2SVG", &config)
		os.Exit(0)
	}

	err := envconfig.Process("WTG2SVG", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	p := wtg.NewParser()

	err = p.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if len(p.InvalidEntries) != 0 {
		for _, err := range p.InvalidEntries {
			log.Println(err)
		}
	}
	imgArea := (p.ImageSize.Max.X - p.ImageSize.Min.X) * (p.ImageSize.Max.X - p.ImageSize.Min.Y)
	canvasArea := (p.MapSize.Max.X - p.MapSize.Min.X) * (p.MapSize.Max.X - p.MapSize.Min.Y)
	if imgArea == 0 || canvasArea == 0 {
		p.ImageSize = image.Rect(0, 0, config.Width, config.Height)
		p.MapSize = image.Rect(30, 50, config.Width-30, config.Height-50)
	}
	e, err := svgmap.NewEncoder(os.Stdout, p.ImageSize, p.MapSize)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	indicators := []svgmap.Annotator{}
	if config.WithIndicators {
		indicators = svgmap.AllEvolutionIndications()
	}
	style := svgmap.NewOctoStyle(p.EvolutionStages, indicators...)
	style.WithSpace = config.WithSpace
	style.WithControls = config.WithControls
	style.WithValueChain = config.WithValueChain
	e.Init(style)
	err = e.Encode(p.WMap)
	if err != nil {
		log.Fatal(err)
	}
}
