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
	wtgFunc := js.FuncOf(func(_ js.Value, args []js.Value) any {
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
	if len(p.InvalidEntries) != 0 {
		for _, err := range p.InvalidEntries {
			fmt.Println(err)
		}
	}

	output := new(bytes.Buffer)
	imgArea := (p.ImageSize.Max.X - p.ImageSize.Min.X) * (p.ImageSize.Max.X - p.ImageSize.Min.Y)
	canvasArea := (p.MapSize.Max.X - p.MapSize.Min.X) * (p.MapSize.Max.X - p.MapSize.Min.Y)
	if imgArea == 0 || canvasArea == 0 {
		p.ImageSize = image.Rect(0, 0, 1100, 900)
		p.MapSize = image.Rect(30, 50, 1070, 850)
	}
	e, err := svgmap.NewEncoder(output, p.ImageSize, p.MapSize)
	if err != nil {
		return "", err
	}
	defer e.Close()
	style := svgmap.NewOctoStyle(p.EvolutionStages)
	e.Init(style)
	err = e.Encode(p.WMap)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
