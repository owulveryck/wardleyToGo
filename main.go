package main

import (
	"log"
	"os"

	svgmap "github.com/owulveryck/wardleyToGo/internal/encoding/svg"
	"github.com/owulveryck/wardleyToGo/internal/parser"
)

func main() {
	width := 1000
	height := 700
	padLeft := 25
	padBottom := 30

	p := parser.NewParser(os.Stdin)
	m, err := p.Parse() // the map
	if err != nil {
		log.Fatal(err)
	}
	svgmap.Encode(m, os.Stdout, width, height, padLeft, padBottom)
}
