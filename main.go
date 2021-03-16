package main

import (
	"os"
)

func main() {
	width := 1000
	height := 800
	padLeft := 25
	padBottom := 30

	w := newWardley(os.Stdout)
	w.Init(width, height, padLeft, padBottom)
	p := newParser(os.Stdin)
	p.parse()
	it := p.g.Nodes()
	for it.Next() {
		w.writeElement(it.Node().(SVGer))
	}
	edgesIt := p.g.Edges()
	for edgesIt.Next() {
		w.writeElement(edgesIt.Edge().(SVGer))
	}
	w.Close()

}
