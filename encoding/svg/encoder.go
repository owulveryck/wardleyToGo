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
)

type Encoder struct {
	start  xml.StartElement
	box    image.Rectangle
	canvas image.Rectangle
	e      *xml.Encoder
	// Themes are applied in order during encoding to add CSS, JS, or other
	// content to the SVG output. Set to nil for pure SVG without any
	// styling or interactivity.
	Themes []Theme
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
	if err := e.EncodeToken(start); err != nil {
		return nil, err
	}
	return &Encoder{
		start:  start,
		box:    box,
		canvas: canvas,
		e:      e,
		Themes: []Theme{&CSSTheme{}, &JSTheme{}},
	}, nil
}

func (e *Encoder) Close() {
	_ = e.e.EncodeToken(e.start.End())
	_ = e.e.Flush()
}

func (e *Encoder) Init(s SVGStyleMarshaler) {
	s.MarshalStyleSVG(e.e, e.box, e.canvas)
}

func (e *Encoder) Encode(m *wardleyToGo.Map) error {
	for _, t := range e.Themes {
		if err := t.Embed(e.e, m); err != nil {
			return err
		}
	}

	if err := e.e.Encode(svg.Text{
		P:          image.Pt(e.box.Dx()/2, 20),
		Text:       []byte(m.Title),
		TextAnchor: svg.TextAnchorMiddle,
	}); err != nil {
		return err
	}
	elems := make([]SVGMarshaler, 0)
	for _, c := range m.Collaborations() {
		if e, ok := c.(SVGMarshaler); ok {
			elems = append(elems, e)
		}
	}
	for _, n := range m.Components() {
		if n, ok := n.(SVGMarshaler); ok {
			elems = append(elems, n)
		}
	}
	sort.Sort(svgMarshalers(elems))
	currentLayer := makeGroup("layer", encoding.Background)
	if err := e.e.EncodeToken(currentLayer.StartElement); err != nil {
		return err
	}
	for _, element := range elems {
		if elem, ok := element.(encoding.Layer); ok {
			layer := elem.GetLayer()
			if layer != currentLayer.id {
				currentLayer = makeGroup("layer", layer)
				if err := e.e.EncodeToken(currentLayer.End()); err != nil {
					return err
				}
				if err := e.e.EncodeToken(currentLayer.StartElement); err != nil {
					return err
				}
			}
		}
		var g *group
		if elem, ok := element.(wardleyToGo.Component); ok {
			g = makeGroup("element", int(elem.ID()))
			for _, t := range e.Themes {
				if dec, ok := t.(ComponentDecorator); ok {
					g.Attr = append(g.Attr, dec.DecorateComponent(elem)...)
				}
			}
			if chainer, ok := elem.(wardleyToGo.Chainer); ok {
				found := false
				for i := range g.Attr {
					if g.Attr[i].Name.Local == "class" {
						found = true
						g.Attr[i] = xml.Attr{
							Name:  xml.Name{Local: "class"},
							Value: fmt.Sprintf("%v %v", g.Attr[i].Value, chainer.GetAbsoluteVisibility()),
						}
					}
				}
				if !found {
					g.Attr = append(g.Attr, xml.Attr{
						Name:  xml.Name{Local: "class"},
						Value: fmt.Sprintf("visibility%v", chainer.GetAbsoluteVisibility()),
					})
				}
			}
			if err := e.e.EncodeToken(g.StartElement); err != nil {
				return err
			}
		}
		if elem, ok := element.(wardleyToGo.Collaboration); ok {
			g = makeGroup(fmt.Sprintf("edge_%v", int(elem.From().ID())), int(elem.To().ID()))
			if chainer, ok := elem.(wardleyToGo.Chainer); ok {
				found := false
				for i := range g.Attr {
					if g.Attr[i].Name.Local == "class" {
						found = true
						g.Attr[i] = xml.Attr{
							Name:  xml.Name{Local: "class"},
							Value: fmt.Sprintf("%v %v", g.Attr[i].Value, chainer.GetAbsoluteVisibility()),
						}
					}
				}
				if !found {
					g.Attr = append(g.Attr, xml.Attr{
						Name:  xml.Name{Local: "class"},
						Value: fmt.Sprintf("visibility%v", chainer.GetAbsoluteVisibility()),
					})
				}
			}
			if err := e.e.EncodeToken(g.StartElement); err != nil {
				return err
			}
		}
		err := element.MarshalSVG(e.e, e.canvas)
		if err != nil {
			return err
		}
		if g != nil {
			if err := e.e.EncodeToken(g.End()); err != nil {
				return err
			}
		}
	}
	return e.e.EncodeToken(currentLayer.End())
}

type group struct {
	xml.StartElement
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
