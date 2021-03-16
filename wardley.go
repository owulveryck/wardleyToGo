package main

import (
	"io"

	svg "github.com/ajstarks/svgo"
)

type wardley struct {
	*svg.SVG
	width     int
	height    int
	padLeft   int
	padBottom int
}

func newWardley(w io.Writer) *wardley {
	return &wardley{
		SVG: svg.New(w),
	}
}

func (w *wardley) Init(width, height, padLeft, padBottom int) {
	w.width = width
	w.height = height
	w.padLeft = padLeft
	w.padBottom = padBottom
	lg := []svg.Offcolor{
		{Offset: 0, Color: "rgb(196,196,196)", Opacity: 1.0},
		{Offset: 30, Color: "rgb(255,255,255)", Opacity: 1.0},
		{Offset: 70, Color: "rgb(255,255,255)", Opacity: 1.0},
		{Offset: 100, Color: "rgb(196,196,196)", Opacity: 1.0}}

	w.Start(width, height)
	w.Title("Wardley")
	w.Rect(0, 0, width, height, "fill:white")
	w.Def()
	w.LinearGradient("wardleyGradient", 0, 0, 100, 0, lg)
	w.Marker("arrow", 15, 0, 12, 12, `viewBox="0 -5 10 10"`)
	w.Path("M0,-5L10,0L0,5", "fill:red")
	w.MarkerEnd()

	w.Marker("graphArrow", 9, 0, 12, 12, `viewBox="0 -5 10 10"`)
	w.Path("M0,-5L10,0L0,5", "fill:black")
	w.MarkerEnd()
	w.DefEnd()

	w.Rect(padLeft, 0, width-padLeft, height-padBottom, "fill:url(#wardleyGradient)")
	w.TranslateRotate(0, height, 270)
	w.Group(`font-family="&quot;Helvetica Neue&quot;,Helvetica,Arial,sans-serif" font-size="13px">`)
	w.Line(padBottom, padLeft, height, padLeft, `stroke="black"`, `stroke-width="1"`, `marker-end="url(#graphArrow)"`)
	w.Line(height, width*100/575+padLeft, padBottom, width*100/575+padLeft, `stroke="#b8b8b8"`, `stroke-dasharray="2,2"`)
	w.Line(height, width*100/250+padLeft, padBottom, width*100/250+padLeft, `stroke="#b8b8b8"`, `stroke-dasharray="2,2"`)
	w.Line(height, 574*width/820+padLeft, padBottom, 574*width/820+padLeft, `stroke="#b8b8b8"`, `stroke-dasharray="2,2"`)
	w.Text(padBottom+10, padLeft-10, "Invisible", `text-anchor="start"`)
	w.Text(height-padBottom-10, padLeft-10, "Visible", `text-anchor="end"`)
	w.Text(height/2+padBottom, padLeft-10, "Value Chain", `text-anchor="middle" font-weight="bold"`)
	w.Gend()
	w.Gend()

	w.Line(padLeft, height-padBottom, width, height-padBottom, `stroke="black"`, `stroke-width="1"`, `marker-end="url(#graphArrow)"`)
	w.Group(`font-family="&quot;Helvetica Neue&quot;,Helvetica,Arial,sans-serif"`, `font-size="13px"`, `font-style="italic"`)
	w.Text(padLeft+10, 15, "Uncharted", `font-style="normal"`, `font-size="11px"`, `font-weight="bold"`)
	w.Text(width-20, 15, "Industrialised", `font-style="normal"`, `font-size="11px"`, `font-weight="bold"`, `text-anchor="end"`)

	w.Text(padLeft+5, height-padBottom/2, "Gensesis")
	w.Text(padLeft+5+width/4, height-padBottom/2, "Custom-Build")
	w.Text(padLeft+5+width/2, height-padBottom/2, "Product")
	w.Text(padLeft+5+width/2, height-3, "(+rental)")
	w.Text(padLeft+5+3*width/4, height-padBottom/2, "Commodity")
	w.Text(padLeft+5+3*width/4, height-3, "(+utility)")
	w.Text(width, height-padBottom/2+5, "Evolution", `text-anchor="end"`, `font-weight="bold"`, `font-style="normal"`)
	w.Gend()
	w.Group(`font-family="Consolas, Lucida Console, monospace"`, `font-weight="14px"`, `font-size="13px"`)
}

func (w *wardley) Close() {
	w.Gend()
	w.End()
}

// SVGer is any object that can represent itself on a map
type SVGer interface {
	SVG(s *svg.SVG, width, height, padLeft, padBottom int)
}

func (w *wardley) writeElement(e SVGer) {
	e.SVG(w.SVG, w.width, w.height, w.padLeft, w.padBottom)
}
