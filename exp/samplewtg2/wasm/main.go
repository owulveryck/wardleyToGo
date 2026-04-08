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
	// Third argument: resolution scale percentage (100 = default 1100x900)
	scalePct := 100
	if len(args) >= 3 {
		scalePct = args[2].Int()
		if scalePct < 25 {
			scalePct = 25
		}
		if scalePct > 400 {
			scalePct = 400
		}
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

	// Scale viewBox and canvas proportionally
	boxW := 1100 * scalePct / 100
	boxH := 900 * scalePct / 100
	marginL := 30 * scalePct / 100
	marginT := 50 * scalePct / 100
	marginR := 30 * scalePct / 100
	marginB := 50 * scalePct / 100

	output := new(bytes.Buffer)
	e, err := svgmap.NewEncoder(output,
		image.Rect(0, 0, boxW, boxH),
		image.Rect(marginL, marginT, boxW-marginR, boxH-marginB))
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
