package svg

import (
	"encoding/xml"
	"fmt"
	"image"
)

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
