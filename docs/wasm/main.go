// +build wasm

package main

import (
	"bytes"
	"image"
	"log"
	"syscall/js"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/owm"
)

func main() {
	c := make(chan bool)
	js.Global().Set("generateSVG", js.FuncOf(generateSVG))
	<-c

}

func generateSVG(this js.Value, inputs []js.Value) interface{} {
	message := inputs[0].String()

	p := owm.NewParser(bytes.NewBufferString(message))
	m, err := p.Parse() // the map
	if err != nil {
		log.Println(err)
	}
	var output bytes.Buffer
	e, err := svgmap.NewEncoder(&output, image.Rect(0, 0, 1300, 800), image.Rect(30, 50, 1070, 850))
	if err != nil {
		log.Fatal(err)
	}
	style := svgmap.NewWardleyStyle(svgmap.DefaultEvolution)
	e.Init(style)
	e.Encode(m)
	e.Close()

	return js.ValueOf(output.String())
}
