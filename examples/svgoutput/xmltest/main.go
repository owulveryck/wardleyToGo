package main

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"
)

type Color struct {
	color.Color
}

func (c Color) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if c.Color == nil {
		return xml.Attr{}, nil
	}
	r, g, b, a := c.RGBA()
	if strings.Contains(name.Local, "-opacity") {
		opacity := float64(a) / float64(65535)
		return xml.Attr{
			Name:  name,
			Value: fmt.Sprintf(`%.1f`, opacity),
		}, nil
	} else {
		return xml.Attr{
			Name:  name,
			Value: fmt.Sprintf(`#%x%x%x`, uint8(r), uint8(g), uint8(b)),
		}, nil
	}
}

type Circle struct {
	P           image.Point
	R           int    `xml:"r,attr"`
	Fill        Color  `xml:"fill,attr,omitempty"`
	Stroke      Color  `xml:"stroke,attr,omitempty"`
	StrokeWidth string `xml:"stroke-width,attr,omitempty"`
}

func must(a xml.Attr, _ error) xml.Attr {
	return a
}

func (c *Circle) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	element := xml.StartElement{
		Name: xml.Name{Local: "circle"},
	}
	element.Attr = []xml.Attr{
		{
			Name:  xml.Name{Local: "cx"},
			Value: strconv.Itoa(c.P.X),
		},
		{
			Name:  xml.Name{Local: "cy"},
			Value: strconv.Itoa(c.P.Y),
		},
		{
			Name:  xml.Name{Local: "r"},
			Value: strconv.Itoa(c.R),
		},
		must(c.Fill.MarshalXMLAttr(xml.Name{Local: "fill"})),
		must(c.Fill.MarshalXMLAttr(xml.Name{Local: "fill-opacity"})),
		must(c.Stroke.MarshalXMLAttr(xml.Name{Local: "fill"})),
		must(c.Stroke.MarshalXMLAttr(xml.Name{Local: "fill-opacity"})),
	}
	if c.StrokeWidth != "" {
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "stroke-width"},
			Value: c.StrokeWidth,
		})
	}
	e.EncodeToken(element)
	e.EncodeToken(element.End())
	//	e.EncodeToken(start)
	return nil
}

type Translate struct {
	image.Point
	Components []interface{}
}

func (t *Translate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	g := xml.StartElement{
		Name: xml.Name{Local: "g"},
	}
	g.Attr = []xml.Attr{
		{
			Name:  xml.Name{Local: "transform"},
			Value: fmt.Sprintf(`translate(%v,%v)`, t.X, t.Y),
		},
	}
	e.EncodeToken(g)
	e.Encode(t.Components)
	e.EncodeToken(g.End())
	return nil
}

func main() {

	SVG(os.Stdout, &component{})

	/*
			Components: []interface{}{
			},
		})
	*/
}

func StartSVG() xml.StartElement {
	return xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "svg",
		},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "width"},
				Value: "100%",
			},
		},
	}
}

func SVG(w io.Writer, c SVGMarshaler) {
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "    ")
	startElement := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "svg",
		},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "width"},
				Value: "100%",
			},
		},
	}
	enc.EncodeToken(startElement)
	c.MarshalSVG(enc, image.Rectangle{})
	enc.EncodeToken(startElement.End())
	enc.Flush()
}

type SVGMarshaler interface {
	MarshalSVG(e *xml.Encoder, bounds image.Rectangle) error
}

type component struct{}

func (c *component) MarshalSVG(e *xml.Encoder, bounds image.Rectangle) error {
	e.Encode(&Translate{
		image.Point{50, 100},
		[]interface{}{
			&Circle{
				Fill: Color{color.Black},
			},
		},
	})
	return nil
}
