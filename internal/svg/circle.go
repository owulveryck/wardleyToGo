package svg

import (
	"encoding/xml"
	"image"
	"strconv"
)

type Circle struct {
	P           image.Point
	R           int
	Fill        Color
	Stroke      Color
	StrokeWidth string
}

func must(a xml.Attr, _ error) xml.Attr {
	return a
}

func (c Circle) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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
		must(c.Stroke.MarshalXMLAttr(xml.Name{Local: "stroke"})),
		must(c.Stroke.MarshalXMLAttr(xml.Name{Local: "stroke-opacity"})),
	}
	if c.StrokeWidth != "" {
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "stroke-width"},
			Value: c.StrokeWidth,
		})
	}
	e.EncodeToken(element)
	e.EncodeToken(element.End())
	return nil
}
