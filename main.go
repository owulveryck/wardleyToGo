package main

import (
	"os"

	svgmap "github.com/owulveryck/wardleyToGo/internal/encoding/svg"
	"github.com/owulveryck/wardleyToGo/internal/parser"
)

func main() {
	width := 1000
	height := 800
	padLeft := 25
	padBottom := 30

	p := parser.NewParser(os.Stdin)
	m := p.Parse() // the map
	svgmap.Encode(m, os.Stdout, width, height, padLeft, padBottom)
}
