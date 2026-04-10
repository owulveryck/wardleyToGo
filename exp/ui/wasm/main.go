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
	// Third argument: resolution scale percentage (100 = default)
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
	// Fourth and fifth arguments: base width and height (before scaling)
	baseW := 1200
	baseH := 900
	if len(args) >= 5 {
		baseW = args[3].Int()
		baseH = args[4].Int()
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
	boxW := baseW * scalePct / 100
	boxH := baseH * scalePct / 100
	marginL := 30 * scalePct / 100
	marginT := 50 * scalePct / 100
	marginR := 30 * scalePct / 100
	marginB := 65 * scalePct / 100

	legendWidth := 0
	if result.Legend && len(result.LegendItems) > 0 {
		legendWidth = svgmap.LegendWidth * scalePct / 100
	}

	svgBuf.Reset()
	e, err := svgmap.NewEncoder(&svgBuf,
		image.Rect(0, 0, boxW+legendWidth, boxH),
		image.Rect(marginL, marginT, boxW-marginR, boxH-marginB))
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	// Close() is called explicitly before reading the buffer below,
	// not via defer, because defer runs after the return value is captured.

	if static {
		e.Themes = nil
	}

	if result.Focus != nil {
		focusTheme := &svgmap.FocusTheme{
			ComponentIDs: result.Focus.ComponentIDs,
			EdgeKeys:     result.Focus.EdgeKeys,
			GroupIDs:     result.Focus.GroupIDs,
		}
		if e.Themes == nil {
			e.Themes = []svgmap.Theme{focusTheme}
		} else {
			e.Themes = append(e.Themes, focusTheme)
		}
	}

	var indicators []svgmap.Annotator
	if result.Legend && len(result.LegendItems) > 0 {
		indicators = append(indicators, &svgmap.Legend{Items: result.LegendItems})
	}
	style := svgmap.NewOctoStyle(result.Stages, indicators...)
	e.Init(style)

	if err := e.Encode(result.Map); err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	e.Close()
	return svgBuf.String()
}

// svgBuf is reused across calls to avoid repeated allocation and GC pressure.
// WASM is single-threaded, so no synchronization is needed.
var svgBuf bytes.Buffer
