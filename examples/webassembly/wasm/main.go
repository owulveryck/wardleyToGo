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
	js.Global().Set("generateSVG", wtgWrapper())
	<-make(chan bool)
}

func wtgWrapper() js.Func {
	wtgFunc := js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in wtgWrapper", r)
			}
		}()
		if len(args) < 1 {
			return "Invalid no of arguments passed"
		}
		input := args[0].String()
		width := 1300
		height := 900
		withAnnotations := false
		if len(args) >= 3 {
			width = args[1].Int()
			height = args[2].Int()
		}
		if len(args) >= 4 {
			withAnnotations = args[3].Bool()
		}
		if width < 200 || height < 200 {
			return fmt.Sprintf("size too small %vx%v (expected at least 200x200)", width, height)
		}
		svg, err := wtg2SVG(input, width, height, withAnnotations)
		if err != nil {
			fmt.Printf("unable to generate svg %s\n", err)
			return err.Error()
		}
		return svg
	})
	return wtgFunc
}

func wtg2SVG(s string, width int, height int, withAnnotations bool) (string, error) {
	p, err := wtg2.NewParser(bytes.NewBufferString(s))
	if err != nil {
		return "", err
	}
	doc, err := p.Parse()
	if err != nil {
		return "", err
	}

	result, err := wtg2.BuildMap(doc)
	if err != nil {
		return "", err
	}

	output := new(bytes.Buffer)
	imgSize := image.Rect(0, 0, width, height)
	mapSize := image.Rect(30, 50, width-30, height-50)

	e, err := svgmap.NewEncoder(output, imgSize, mapSize)
	if err != nil {
		return "", err
	}
	defer e.Close()

	indicators := []svgmap.Annotator{}
	if withAnnotations {
		indicators = svgmap.AllEvolutionIndications()
	}
	style := svgmap.NewOctoStyle(result.Stages, indicators...)
	style.WithControls = true
	e.Init(style)
	err = e.Encode(result.Map)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
