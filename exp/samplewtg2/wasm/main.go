//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"image"
	"syscall/js"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

func main() {
	js.Global().Set("generateSVG", js.FuncOf(generate))
	<-make(chan bool)
}

func generate(_ js.Value, args []js.Value) any {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	if len(args) < 1 {
		return "error: no input provided"
	}

	input := args[0].String()
	static := false
	if len(args) >= 2 {
		static = args[1].Bool()
	}

	p, err := wtg2.NewParser(bytes.NewBufferString(input))
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	doc, err := p.Parse()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}

	result, err := wtg2.BuildMap(doc)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}

	output := new(bytes.Buffer)
	e, err := svgmap.NewEncoder(output,
		image.Rect(0, 0, 1100, 900),
		image.Rect(30, 50, 1070, 850))
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	defer e.Close()

	if static {
		e.Themes = nil
	}

	style := svgmap.NewOctoStyle(result.Stages)
	style.WithControls = !static
	e.Init(style)

	if err := e.Encode(result.Map); err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return output.String()
}
