package main

import (
	"image"
	"log"
	"os"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/wtg"
)

func main() {

	m, err := wtg.ParseValueChain(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	e, err := svgmap.NewEncoder(os.Stdout, image.Rect(0, 0, 1100, 900), image.Rect(30, 50, 1070, 850))
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	style := svgmap.NewWardleyStyle(svgmap.DefaultEvolution)
	e.Init(style)
	err = e.Encode(m)
	if err != nil {
		log.Fatal(err)
	}
}
