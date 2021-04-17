package svgmap

import (
	"encoding/xml"
	"fmt"
	"image"
	"io"
	"sort"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/encoding"
	"github.com/owulveryck/wardleyToGo/internal/svg"
	"gonum.org/v1/gonum/graph"
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
	e.e.Encode(svg.Text{
		P:          image.Pt(e.box.Dx()/2, 20),
		Text:       []byte(m.Title),
		TextAnchor: svg.TextAnchorMiddle,
	})
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
	currentLayer := makeGroup("layer", encoding.Background)
	e.e.EncodeToken(currentLayer.StartElement)
	for _, element := range elems {
		if elem, ok := element.(encoding.Layer); ok {
			layer := elem.GetLayer()
			if layer != currentLayer.id {
				currentLayer = makeGroup("layer", layer)
				e.e.EncodeToken(currentLayer.End())
				e.e.EncodeToken(currentLayer.StartElement)
			}
		}
		var g *group
		if elem, ok := element.(graph.Node); ok {
			g = makeGroup("element", int(elem.ID()))
			e.e.EncodeToken(g.StartElement)
		}
		if elem, ok := element.(graph.Edge); ok {
			g = makeGroup(fmt.Sprintf("edge_%v", int(elem.From().ID())), int(elem.To().ID()))
			e.e.EncodeToken(g.StartElement)
		}
		err := element.MarshalSVG(e.e, e.canvas)
		if err != nil {
			return err
		}
		if g != nil {
			e.e.EncodeToken(g.End())
		}
	}
	e.e.EncodeToken(currentLayer.End())
	return nil
}

type group struct {
	xml.StartElement
	s  string
	id int
}

func makeGroup(s string, id int) *group {
	return &group{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: `g`},
			Attr: []xml.Attr{
				{
					Name:  xml.Name{Local: "id"},
					Value: fmt.Sprintf("%v_%v", s, id),
				},
			},
		},
		id: id,
	}
}
