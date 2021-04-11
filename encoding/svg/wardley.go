package svgmap

import (
	"fmt"
	"io"

	svg "github.com/ajstarks/svgo"
	"github.com/owulveryck/wardleyToGo"
)

// svgMap is an object representing the map in SVG
type svgMap struct {
	*svg.SVG
	width     int
	height    int
	padLeft   int
	padBottom int
}

// newSvgMap creates an empty map. The caller must call Init to fill it with the axis and Close.
func newSvgMap(w io.Writer) *svgMap {
	return &svgMap{
		SVG: svg.New(w),
	}
}

//init the map according to this geometry
func (w *svgMap) init(width, height, padLeft, padBottom int) {
	w.width = width
	w.height = height
	w.padLeft = padLeft
	w.padBottom = padBottom
	lg := []svg.Offcolor{
		{Offset: 0, Color: "rgb(196,196,196)", Opacity: 1.0},
		{Offset: 30, Color: "rgb(255,255,255)", Opacity: 1.0},
		{Offset: 70, Color: "rgb(255,255,255)", Opacity: 1.0},
		{Offset: 100, Color: "rgb(196,196,196)", Opacity: 1.0}}

	w.Startraw(`width="100%"`, `height="100%"`, `class="wardley-map"`, `preserveAspectRatio="xMidYMid meet"`, fmt.Sprintf(`viewBox="0 0 %v %v"`, width+150, height+padBottom))
	w.writeLegend()

	w.Rect(0, 0, width, height, "fill:white")
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
	//<polygon points="723,314 543,625.769145 183,625.769145 3,314 183,2.230855 543,2.230855 723,314" fill="white" stroke="black" stroke-width="4"/>
	w.MarkerEnd()

	w.DefEnd()

	w.Rect(padLeft, 0, width-padLeft, height-padBottom, "fill:url(#wardleyGradient)")
	w.TranslateRotate(0, height, 270)
	w.Group(`font-family="Helvetica,Arial,sans-serif" font-size="13px">`)
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
	w.Group(`font-family="Helvetica,Arial,sans-serif"`, `font-size="13px"`, `font-style="italic"`)
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

func (w *svgMap) tooltips(width, height, padLeft, padBottom int) {
	sizeW := float64(width-padLeft) / 100
	sizeH := float64(height-padBottom) / 100
	var placementW, placementH float64
	for i := 0; i < 100; i++ {
		placementH = 0
		for j := 0; j < 100; j++ {
			placementH += sizeH
			w.SVG.Translate(int((placementW))+padLeft, int((placementH)))
			w.SVG.Rect(0, 0, int(sizeW), int(sizeH),
				`fill-opacity="0.0"`,
				fmt.Sprintf(`onmousemove="showTooltip(evt, '[%.2f,%.2f]');"`, 1-placementH/float64(height-padBottom), placementW/float64(width-padLeft)), `onmouseout="hideTooltip();"`)
			w.SVG.Gend()

		}
		placementW += sizeW
	}

}

// close the map (add the closing tags to the SVG)
func (w *svgMap) close() {
	w.Gend()
	w.End()
}

// writeElement on the map
func (w *svgMap) writeElement(e SVGer) {
	e.SVG(w.SVG, w.width, w.height, w.padLeft, w.padBottom)
}

// Encode the map
func Encode(m *wardleyToGo.Map, w io.Writer, width, height, padLeft, padBottom int, withToolTip bool) {
	if m == nil {
		return
	}
	out := newSvgMap(w)
	out.init(width, height, padLeft, padBottom)
	out.Title(m.Title)
	edgesIt := m.Edges()
	out.Gid("teamRelations")
	for edgesIt.Next() {
		edgesIt.Edge().(wardleyToGo.Edge).SVGTT(out.SVG, width, height, padLeft, padBottom)
	}
	out.Gend()

	out.Gid("teamtopologies")
	it := m.Nodes()
	// First place the orphan nodes as they are probably anotations
	for it.Next() {
		n := it.Node()
		if m.To(n.ID()).Len() == 0 && m.From(n.ID()).Len() == 0 {
			out.writeElement(it.Node().(SVGer))
		}
	}
	out.Gend()
	out.Gid("links")
	edgesIt.Reset()
	for edgesIt.Next() {
		out.writeElement(edgesIt.Edge().(SVGer))
	}
	out.Gend()
	it.Reset()
	out.Gid("components")
	for it.Next() {
		n := it.Node()
		if m.To(n.ID()).Len() != 0 || m.From(n.ID()).Len() != 0 {
			out.writeElement(it.Node().(SVGer))
		}
	}
	out.Gend()
	if len(m.Annotations) != 0 {
		out.Gid("annotations")
		for _, annotation := range m.Annotations {
			out.writeElement(annotation)
		}
		// Add the annotation box
		writeAnnotations(out, m, width, height, padLeft, padBottom)
		out.Gend()
	}

	if withToolTip {
		out.Group("tooltips")
		out.tooltips(width, height, padLeft, padBottom)
		out.Gend()
	}
	out.close()
}

func writeAnnotations(out *svgMap, m *wardleyToGo.Map, width, height, padLeft, padBottom int) {
	maxLen := 0
	for _, annotation := range m.Annotations {
		if len(annotation.Label) > maxLen {
			maxLen = len(annotation.Label)
		}
	}
	out.Translate(m.AnnotationsPlacement[1]*(width-padLeft)/100+padLeft, (height-padLeft)-m.AnnotationsPlacement[0]*(height-padLeft)/100)
	out.Rect(-14, -14, 9*maxLen, len(m.Annotations)*19+19, `stroke="#595959"`, `stroke-width="1"`, `fill="#FFFFFF"`)
	out.Text(0, 0, "Annotations:", `text-decoration="underline"`)
	for i, annotation := range m.Annotations {
		out.Text(0, 18*(i+1), fmt.Sprintf(" %v. %v", annotation.Identifier, annotation.Label))
	}
	/*
		<rect x="-14" id="annotationsBoxWrap" y="-14" class="draggable" width="452.875" height="55" stroke="#595959" stroke-width="1" fill="#FFFFFF"></rect>
	*/
	out.Gend()
}