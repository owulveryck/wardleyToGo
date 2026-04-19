package svgmap

import (
	"encoding/xml"
	"image"
	"image/color"
	"io"
	"sort"
	"strconv"

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
		ViewBox:             strconv.Itoa(box.Min.X) + " " + strconv.Itoa(box.Min.Y) + " " + strconv.Itoa(box.Max.X) + " " + strconv.Itoa(box.Max.Y),
	}.StartSVG()
	if err := e.EncodeToken(start); err != nil {
		return nil, err
	}
	return &Encoder{
		start:  start,
		box:    box,
		canvas: canvas,
		e:      e,
		Themes: func() []Theme {
			css := &CSSTheme{}
			return []Theme{css, &JSTheme{CSS: css}}
		}(),
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
		P:          image.Pt(e.canvas.Min.X+(e.canvas.Dx()/2), e.canvas.Min.Y-15),
		Text:       []byte(m.Title),
		TextAnchor: svg.TextAnchorMiddle,
		FontWeight: "bold",
		FontSize:   "16px",
		Fill:       svg.Color{Color: color.RGBA{0x13, 0x24, 0x54, 0xff}},
	}); err != nil {
		return err
	}
	collabs := m.Collaborations()
	comps := m.Components()
	elems := make([]SVGMarshaler, 0, len(collabs)+len(comps))
	for _, c := range collabs {
		if e, ok := c.(SVGMarshaler); ok {
			elems = append(elems, e)
		}
	}
	for _, n := range comps {
		if n, ok := n.(SVGMarshaler); ok {
			elems = append(elems, n)
		}
	}
	sort.Stable(svgMarshalers(elems))
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
							Value: g.Attr[i].Value + " " + strconv.Itoa(chainer.GetAbsoluteVisibility()),
						}
					}
				}
				if !found {
					g.Attr = append(g.Attr, xml.Attr{
						Name:  xml.Name{Local: "class"},
						Value: "visibility" + strconv.Itoa(chainer.GetAbsoluteVisibility()),
					})
				}
			}
			if err := e.e.EncodeToken(g.StartElement); err != nil {
				return err
			}
		}
		if elem, ok := element.(wardleyToGo.Collaboration); ok {
			g = makeGroup("edge_"+strconv.Itoa(int(elem.From().ID())), int(elem.To().ID()))
			for _, t := range e.Themes {
				if dec, ok := t.(CollaborationDecorator); ok {
					g.Attr = append(g.Attr, dec.DecorateCollaboration(elem)...)
				}
			}
			if chainer, ok := elem.(wardleyToGo.Chainer); ok {
				found := false
				for i := range g.Attr {
					if g.Attr[i].Name.Local == "class" {
						found = true
						g.Attr[i] = xml.Attr{
							Name:  xml.Name{Local: "class"},
							Value: g.Attr[i].Value + " " + strconv.Itoa(chainer.GetAbsoluteVisibility()),
						}
					}
				}
				if !found {
					g.Attr = append(g.Attr, xml.Attr{
						Name:  xml.Name{Local: "class"},
						Value: "visibility" + strconv.Itoa(chainer.GetAbsoluteVisibility()),
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
					Value: s + "_" + strconv.Itoa(id),
				},
			},
		},
		id: id,
	}
}
