package svg

import (
	"encoding/xml"
	"image"
	"strconv"
)

type Rectangle struct {
	R           image.Rectangle
	Rx, Ry      int
	Fill        Color
	Stroke      Color
	StrokeWidth string
	Style       string
}

func (r Rectangle) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	element := xml.StartElement{
		Name: xml.Name{Local: "rect"},
	}
	element.Attr = []xml.Attr{
		{
			Name:  xml.Name{Local: "x"},
			Value: strconv.Itoa(r.R.Min.X),
		},
		{
			Name:  xml.Name{Local: "y"},
			Value: strconv.Itoa(r.R.Min.Y),
		},
		{
			Name:  xml.Name{Local: "width"},
			Value: strconv.Itoa(r.R.Dx()),
		},
		{
			Name:  xml.Name{Local: "height"},
			Value: strconv.Itoa(r.R.Dy()),
		},
		must(r.Fill.MarshalXMLAttr(xml.Name{Local: "fill"})),
		must(r.Fill.MarshalXMLAttr(xml.Name{Local: "fill-opacity"})),
		must(r.Stroke.MarshalXMLAttr(xml.Name{Local: "stroke"})),
		must(r.Stroke.MarshalXMLAttr(xml.Name{Local: "stroke-opacity"})),
	}
	if r.StrokeWidth != "" {
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "stroke-width"},
			Value: r.StrokeWidth,
		})
	}
	if r.Style != "" {
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "style"},
			Value: r.Style,
		})
	}
	if r.Ry != 0 {
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "ry"},
			Value: strconv.Itoa(r.Ry),
		})
	}
	if r.Rx != 0 {
		element.Attr = append(element.Attr, xml.Attr{
			Name:  xml.Name{Local: "rx"},
			Value: strconv.Itoa(r.Rx),
		})
	}
	e.EncodeToken(element)
	e.EncodeToken(element.End())
	return nil
}
