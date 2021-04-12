package svgmap

import (
	"fmt"
	"image"
	"io"

	svg "github.com/ajstarks/svgo"
	"github.com/owulveryck/wardleyToGo"
)

// svgMap is an object representing the map in SVG
type svgMap struct {
	*svg.SVG
	// canvas is the actual map are
	width, height int
	canvas        image.Rectangle
}

// newSvgMap creates an empty map. The caller must call Init to fill it with the axis and Close.
func newSvgMap(w io.Writer, canvas image.Rectangle) *svgMap {
	return &svgMap{
		SVG:    svg.New(w),
		canvas: canvas,
	}
}

// Encode the map
func Encode(m *wardleyToGo.Map, w io.Writer, width, height int, canvas image.Rectangle) {
	out := newSvgMap(w, canvas)
	out.width = width
	out.height = height
	out.canvas = canvas
	out.init()
	e := m.Edges()
	for e.Next() {
		if e, ok := e.Edge().(SVGDrawer); ok {
			e.SVGDraw(out.SVG, canvas)
		}
	}
	n := m.Nodes()
	for n.Next() {
		if n, ok := n.Node().(SVGDrawer); ok {
			n.SVGDraw(out.SVG, canvas)
		}
	}
	out.Gend()
	out.SVG.End()
}

func (w *svgMap) init() {

	markers := []struct {
		position float64
		label    string
	}{
		{
			position: 0,
			label:    "Genesis",
		},
		{
			position: (float64(100) / 575),
			label:    "Custom-Built",
		},
		{
			position: (float64(100) / 250),
			label:    "Product\n(+rental)",
		},
		{
			position: (float64(574) / 820),
			label:    "Commodity\n(+utility)",
		},
	}
	lg := []svg.Offcolor{
		{Offset: 0, Color: "rgb(196,196,196)", Opacity: 1.0},
		{Offset: 30, Color: "rgb(255,255,255)", Opacity: 1.0},
		{Offset: 70, Color: "rgb(255,255,255)", Opacity: 1.0},
		{Offset: 100, Color: "rgb(196,196,196)", Opacity: 1.0}}

	w.Startraw(`width="100%"`, `height="100%"`, `class="wardley-map"`, `preserveAspectRatio="xMidYMid meet"`, fmt.Sprintf(`viewBox="0 0 %v %v"`, w.width, w.height))

	w.Rect(0, 0, w.width, w.height, "fill:grey")
	w.Def()
	w.LinearGradient("wardleyGradient", 0, 0, 100, 0, lg)
	w.Marker("arrow", 15, 0, 12, 12, `viewBox="0 -5 10 10"`)
	w.Path("M0,-5L10,0L0,5", "fill:red")
	w.MarkerEnd()

	w.Marker("graphArrow", 9, 0, 12, 12, `viewBox="0 -5 10 10"`)
	w.Path("M0,-5L10,0L0,5", "fill:black")
	w.MarkerEnd()
	w.Marker("hexagon", 0, 0, 15, 15)
	w.Polygon([]int{723, 543, 183, 3, 183, 543, 723}, []int{314, 625, 625, 314, 2, 2, 314}, `fill="white"`, `stroke="black"`, `stroke-width="2"`)
	w.MarkerEnd()

	w.DefEnd()

	// The canvas background
	w.Rect(w.canvas.Min.X, w.canvas.Min.Y, w.canvas.Dx(), w.canvas.Dy(), "fill:url(#wardleyGradient)")
	w.TranslateRotate(w.canvas.Min.X, w.canvas.Max.Y, 270)
	w.Group(`font-family="Helvetica,Arial,sans-serif" font-size="13px">`)
	// The value chain
	w.Line(0, 0, w.canvas.Dy(), 0, `stroke="black"`, `stroke-width="1"`, `marker-end="url(#graphArrow)"`)
	for i := 1; i < len(markers); i++ {
		axis := markers[i]
		w.Line(0, int(float64(w.canvas.Dx())*axis.position), w.canvas.Dy(), int(float64(w.canvas.Dx())*axis.position), `stroke="#b8b8b8"`, `stroke-dasharray="2,2"`)

	}
	w.Text(5, -10, "Invisible", `text-anchor="start"`)
	w.Text(w.canvas.Dy()-5, -10, "Visible", `text-anchor="end"`)
	w.Text(w.canvas.Dy()/2, -10, "Value Chain", `text-anchor="middle" font-weight="bold"`)
	w.Gend()
	w.Gend()

	// horizontal axis
	w.Line(w.canvas.Min.X, w.canvas.Max.Y, w.canvas.Max.X, w.canvas.Max.Y, `stroke="black"`, `stroke-width="1"`, `marker-end="url(#graphArrow)"`)
	w.Group(`font-family="Helvetica,Arial,sans-serif"`, `font-size="13px"`, `font-style="italic"`)
	w.Text(w.canvas.Min.X+7, w.canvas.Min.Y+15, "Uncharted", `font-style="normal"`, `font-size="11px"`, `font-weight="bold"`)
	w.Text(w.canvas.Max.X-5, w.canvas.Min.Y+15, "Industrialised", `font-style="normal"`, `font-size="11px"`, `font-weight="bold"`, `text-anchor="end"`)

	for i := 0; i < len(markers); i++ {
		axis := markers[i]
		w.Text(int(float64(w.canvas.Dx())*axis.position)+w.canvas.Min.X, w.canvas.Max.Y+15, axis.label)
	}

	w.Text(w.canvas.Max.X, w.canvas.Max.Y+15, "Evolution", `text-anchor="end"`, `font-weight="bold"`, `font-style="normal"`)
	w.Gend()
	w.Group(`font-family="Consolas, Lucida Console, monospace"`, `font-weight="14px"`, `font-size="13px"`)
}
