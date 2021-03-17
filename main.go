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

	w := svgmap.NewSVGMap(os.Stdout)
	w.Init(width, height, padLeft, padBottom)
	p := parser.NewParser(os.Stdin)
	p.Parse()
	it := p.G.Nodes()
	for it.Next() {
		w.WriteElement(it.Node().(svgmap.SVGer))
	}
	edgesIt := p.G.Edges()
	for edgesIt.Next() {
		w.WriteElement(edgesIt.Edge().(svgmap.SVGer))
	}
	w.Close()

}
