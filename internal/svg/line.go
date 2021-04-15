package svg

import (
	"encoding/xml"
	"errors"
	"fmt"
	"image"
	"strconv"
)

type Line struct {
	F, T            image.Point
	Stroke          Color
	StrokeWidth     string
	StrokeDashArray []int
	MarkerEnd       string
}

func (l Line) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	element := xml.StartElement{
		Name: xml.Name{Local: "line"},
	}
	attrs := newAttributes()
	attrs = attrs.append("x1", strconv.Itoa(l.F.X))
	attrs = attrs.append("y1", strconv.Itoa(l.F.Y))
	attrs = attrs.append("x2", strconv.Itoa(l.T.X))
	attrs = attrs.append("y2", strconv.Itoa(l.T.Y))
	if l.StrokeWidth != "" {
		attrs = attrs.append("stroke-width", l.StrokeWidth)
	}
	if l.StrokeDashArray != nil {
		if len(l.StrokeDashArray) != 2 {
			return errors.New("bad dash-array option")
		}
		attrs = attrs.append("stroke-dasharray", fmt.Sprintf("%v %v", l.StrokeDashArray[0], l.StrokeDashArray[1]))
	}
	if l.MarkerEnd != "" {
		attrs = attrs.append("marker-end", l.MarkerEnd)
	}
	element.Attr = attrs
	element.Attr = append(element.Attr, must(l.Stroke.MarshalXMLAttr(xml.Name{Local: "stroke"})))
	element.Attr = append(element.Attr, must(l.Stroke.MarshalXMLAttr(xml.Name{Local: "stroke-opacity"})))
	e.EncodeToken(element)
	e.EncodeToken(element.End())
	return nil
}

type attributes []xml.Attr

func newAttributes() attributes {
	attrs := make([]xml.Attr, 0)
	return (attrs)
}

func (a attributes) append(name, value string) attributes {
	return append(a, xml.Attr{
		Name:  xml.Name{Local: name},
		Value: value,
	})
}
