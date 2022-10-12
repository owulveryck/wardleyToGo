package main

import (
	"image"
	"log"
	"os"

	dotMap "github.com/owulveryck/wardleyToGo/encoding/dot"
	"github.com/owulveryck/wardleyToGo/parser/owm"
)

func main() {
	p := owm.NewParser(os.Stdin)
	m, err := p.Parse() // the map
	if err != nil {
		log.Fatal(err)
	}
	e := dotMap.NewEncoder(os.Stdout, image.Rect(0, 0, 1100, 900), image.Rect(30, 50, 1070, 850))
	err = e.Encode(m)
	if err != nil {
		log.Fatal(err)
	}
}
