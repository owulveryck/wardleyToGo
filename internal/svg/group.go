package svg

import "encoding/xml"

// Group wraps child elements in a <g> element with arbitrary attributes.
type Group struct {
	Attrs    []xml.Attr
	Children []interface{}
}

func (g Group) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s := xml.StartElement{
		Name: xml.Name{Local: "g"},
		Attr: g.Attrs,
	}
	if err := e.EncodeToken(s); err != nil {
		return err
	}
	for _, child := range g.Children {
		if err := e.Encode(child); err != nil {
			return err
		}
	}
	return e.EncodeToken(s.End())
}
