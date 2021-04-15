package svgmap

import (
	"encoding/xml"
	"fmt"
	"image"
	"io"
	"sort"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/internal/svg"
)

type Encoder struct {
	start  xml.StartElement
	box    image.Rectangle
	canvas image.Rectangle
	e      *xml.Encoder
}

func NewEncoder(w io.Writer, box image.Rectangle, canvas image.Rectangle) (*Encoder, error) {
	// TODO check canvas
	e := xml.NewEncoder(w)
	e.Indent("", "    ")
	start := svg.SVG{
		Width:               "100%",
		Height:              "100%",
		PreserveAspectRatio: "xMidYMid meet",
		ViewBox:             fmt.Sprintf("%v %v %v %v", box.Min.X, box.Min.Y, box.Max.X, box.Max.Y),
	}.StartSVG()
	e.EncodeToken(start)
	return &Encoder{
		start:  start,
		box:    box,
		canvas: canvas,
		e:      e,
	}, nil
}

func (e *Encoder) Close() {
	e.e.EncodeToken(e.start.End())
	e.e.Flush()
}

func (e *Encoder) Init(s SVGStyleMarshaler) {
	s.MarshalStyleSVG(e.e, e.box, e.canvas)
}

func (e *Encoder) Encode(m *wardleyToGo.Map) error {
	elems := make([]SVGMarshaler, 0)
	edges := m.Edges()
	for edges.Next() {
		if e, ok := edges.Edge().(SVGMarshaler); ok {
			elems = append(elems, e)
		}
	}
	n := m.Nodes()
	for n.Next() {
		if n, ok := n.Node().(SVGMarshaler); ok {
			elems = append(elems, n)
		}
	}
	sort.Sort(svgMarshalers(elems))
	for _, element := range elems {
		err := element.MarshalSVG(e.e, e.canvas)
		if err != nil {
			return err
		}
	}
	return nil
}
