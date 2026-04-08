package svg

import (
	"encoding/xml"
	"fmt"
)

// Path renders an SVG <path> element. When Stroke is set, it renders
// as a stroked path (used for curved edges). When only Fill is set,
// it renders as a filled shape (used for markers).
type Path struct {
	D               string
	Fill            Color
	Stroke          Color
	StrokeWidth     string
	StrokeDashArray []int
	MarkerEnd       string
	Class           []string
}

func (p Path) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	element := xml.StartElement{
		Name: xml.Name{Local: "path"},
	}
	attrs := newAttributes()
	attrs = attrs.append("d", p.D)

	if p.StrokeWidth != "" {
		// Stroked path (edge)
		attrs = attrs.append("fill", "none")
		element.Attr = attrs
		element.Attr = append(element.Attr, must(p.Stroke.MarshalXMLAttr(xml.Name{Local: "stroke"})))
		element.Attr = append(element.Attr, must(p.Stroke.MarshalXMLAttr(xml.Name{Local: "stroke-opacity"})))
		element.Attr = append(element.Attr, xml.Attr{Name: xml.Name{Local: "stroke-width"}, Value: p.StrokeWidth})
		if len(p.StrokeDashArray) == 2 {
			element.Attr = append(element.Attr, xml.Attr{
				Name:  xml.Name{Local: "stroke-dasharray"},
				Value: fmt.Sprintf("%v %v", p.StrokeDashArray[0], p.StrokeDashArray[1]),
			})
		}
		if p.MarkerEnd != "" {
			element.Attr = append(element.Attr, xml.Attr{Name: xml.Name{Local: "marker-end"}, Value: p.MarkerEnd})
		}
		if len(p.Class) > 0 {
			classes := ""
			for _, c := range p.Class {
				classes += " " + c
			}
			element.Attr = append(element.Attr, xml.Attr{Name: xml.Name{Local: "class"}, Value: classes})
		}
	} else {
		// Filled path (marker shape)
		element.Attr = attrs
		element.Attr = append(element.Attr, must(p.Fill.MarshalXMLAttr(xml.Name{Local: "fill"})))
	}

	e.EncodeToken(element)
	e.EncodeToken(element.End())
	return nil
}
