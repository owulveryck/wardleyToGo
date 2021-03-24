// +build wasm

package main

import (
	"bytes"
	"log"
	"syscall/js"

	svgmap "github.com/owulveryck/wardleyToGo/internal/encoding/svg"
	"github.com/owulveryck/wardleyToGo/internal/parser"
)

func main() {
	c := make(chan bool)
	js.Global().Set("generateSVG", js.FuncOf(generateSVG))
	<-c

}

func generateSVG(this js.Value, inputs []js.Value) interface{} {
	message := inputs[0].String()
	width := 1300
	height := 800
	padLeft := 25
	padBottom := 30

	p := parser.NewParser(bytes.NewBufferString(message))
	m, err := p.Parse() // the map
	if err != nil {
		log.Println(err)
	}
	var output bytes.Buffer
	svgmap.Encode(m, &output, width, height, padLeft, padBottom, false)
	return js.ValueOf(output.String())
}
