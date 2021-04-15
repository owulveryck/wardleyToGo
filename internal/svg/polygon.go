package svg

import (
	"encoding/xml"
	"fmt"
	"image"
	"strconv"
	"strings"
)

type Polygon struct {
	XMLName     xml.Name `xml:"polygon"`
	Points      []image.Point
	Fill        Color
	Stroke      Color
	StrokeWidth int
}

func (p *Polygon) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var b strings.Builder
	for _, p := range p.Points {
		b.WriteString(fmt.Sprintf("%v,%v ", p.X, p.Y))
	}
	attrs := newAttributes()
	attrs = attrs.append("points", b.String())
	if p.StrokeWidth != 0 {
		attrs = attrs.append(" stroke-width", strconv.Itoa(p.StrokeWidth))
	}
	start.Attr = attrs
	start.Attr = append(start.Attr, must(p.Stroke.MarshalXMLAttr(xml.Name{Local: "stroke"})))
	start.Attr = append(start.Attr, must(p.Stroke.MarshalXMLAttr(xml.Name{Local: "stroke-opacity"})))
	start.Attr = append(start.Attr, must(p.Stroke.MarshalXMLAttr(xml.Name{Local: "fill"})))
	start.Attr = append(start.Attr, must(p.Stroke.MarshalXMLAttr(xml.Name{Local: "fill-opacity"})))
	e.EncodeToken(start)
	e.EncodeToken(start.End())
	return nil
}
