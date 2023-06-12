//go:build js && wasm

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"strings"
	"syscall/js"

	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/wtg"
	"github.com/owulveryck/wardleyToGo/parser/wtg2owm"
)

func main() {
	js.Global().Set("generateSVG", wtgWrapper())
	js.Global().Set("toOWM", owmWrapper())
	js.Global().Set("getUnconfiguredComponents", componentsWrapper())
	<-make(chan bool)
}

func componentsWrapper() js.Func {
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
		p := wtg.NewParser()

		buf := bytes.NewBufferString(input)
		err := p.Parse(buf)
		if err != nil {
			fmt.Printf("unable to parse %s\n", err)
			return err.Error()
		}
		if len(p.InvalidEntries) != 0 {
			for _, err := range p.InvalidEntries {
				fmt.Println(err)
			}
		}
		output := make([]string, 0)
		c := p.WMap.Nodes()
		for c.Next() {
			if c, ok := c.Node().(*wardley.Component); ok {
				if !c.Configured {
					output = append(output, c.Label)
				}
			}
		}
		b, _ := json.Marshal(output)
		return string(b)
	})
	return wtgFunc
}

func owmWrapper() js.Func {
	wtgFunc := js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in owmWrapper", r)
			}
		}()
		if len(args) < 1 {
			return "Invalid no of arguments passed"
		}
		input := args[0].String()
		var err error
		var b strings.Builder
		err = wtg2owm.Convert(bytes.NewBufferString(input), &b)
		if err != nil {
			fmt.Printf("unable to generate owm %s\n", err)
			return err.Error()
		}
		return b.String()
	})
	return wtgFunc
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
	p := wtg.NewParser()

	buf := bytes.NewBufferString(s)
	err := p.Parse(buf)
	if err != nil && err != wtg.ErrEmptyMap {
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
		p.ImageSize = image.Rect(0, 0, width, height)
		p.MapSize = image.Rect(30, 50, width-30, height-50)
	}
	e, err := svgmap.NewEncoder(output, p.ImageSize, p.MapSize)
	if err != nil {
		return "", err
	}
	defer e.Close()
	indicators := []svgmap.Annotator{}
	if withAnnotations {
		indicators = svgmap.AllEvolutionIndications()
	}
	style := svgmap.NewOctoStyle(p.EvolutionStages, indicators...)
	style.WithControls = true
	e.Init(style)
	err = e.Encode(p.WMap)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
