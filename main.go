package main

import (
	"os"

	svg "github.com/ajstarks/svgo"
)

func main() {
	width := 1000
	height := 730
	padLeft := 10
	padBottom := 10

	lg := []svg.Offcolor{
		{Offset: 0, Color: "rgb(196,196,196)", Opacity: 1.0},
		{Offset: 30, Color: "rgb(255,255,255)", Opacity: 1.0},
		{Offset: 70, Color: "rgb(255,255,255)", Opacity: 1.0},
		{Offset: 100, Color: "rgb(196,196,196)", Opacity: 1.0}}

	g := svg.New(os.Stdout)
	g.Start(width, height)
	g.Title("Wardley")
	g.Rect(0, 0, width, height, "fill:white")
	g.Def()
	g.LinearGradient("wardleyGradient", 0, 0, 100, 0, lg)
	g.Marker("arrow", 15, 0, 12, 12)
	g.Path("M2,2 L2,11 L10,6 L2,2", "fill:red")
	g.MarkerEnd()

	g.Marker("graphArrow", 9, 0, 12, 12)
	g.Path("M2,2 L2,11 L10,6 L2,2", "fill:black")
	g.MarkerEnd()
	g.DefEnd()

	g.Rect(padLeft, 0, width-padLeft, height-padBottom, "fill:url(#wardleyGradient)")
	g.Gtransform(`translate(height,0) rotate(270) `)
	g.Line(padLeft, height-padBottom, padLeft, 5,
		`stroke="black"`,
		`stroke-width="1"`,
		`marker-end="url(#graphArrow)"`,
	)
	g.Line(width/4+padLeft, height-padBottom, width/4+padLeft, 5,
		`stroke="#b8b8b8"`,
		`stroke-dasharray="2,2"`)
	g.Line(width/2+padLeft, height-padBottom, width/2+padLeft, 5,
		`stroke="#b8b8b8"`,
		`stroke-dasharray="2,2"`)
	g.Line(3*width/4+padLeft, height-padBottom, 3*width/4+padLeft, 5,
		`stroke="#b8b8b8"`,
		`stroke-dasharray="2,2"`)
	g.Text(50, 50, "bla")
	g.Gend()

	g.Line(padLeft, height-padBottom, width-10, height-padBottom,
		`stroke="black"`,
		`stroke-width="1"`,
		`marker-end="url(#graphArrow)"`,
	)
	g.End()
}
