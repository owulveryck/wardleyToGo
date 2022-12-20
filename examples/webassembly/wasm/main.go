//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"image"
	"syscall/js"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/wtg"
)

func main() {
	js.Global().Set("generateSVG", wtgWrapper())
	<-make(chan bool)
}

func wtgWrapper() js.Func {
	wtgFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		input := args[0].String()
		svg, err := wtg2SVG(input)
		if err != nil {
			fmt.Printf("unable to generate svg %s\n", err)
			return err.Error()
		}
		return svg
	})
	return wtgFunc
}

func wtg2SVG(s string) (string, error) {
	p := wtg.NewParser()

	buf := bytes.NewBufferString(s)
	err := p.Parse(buf)
	if err != nil {
		return "", err
	}
	output := new(bytes.Buffer)
	e, err := svgmap.NewEncoder(output, image.Rect(0, 0, 1100, 900), image.Rect(30, 50, 1070, 850))
	if err != nil {
		return "", err
	}
	defer e.Close()
	style := svgmap.NewOctoStyle(svgmap.DefaultEvolution)
	e.Init(style)
	err = e.Encode(p.WMap)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
