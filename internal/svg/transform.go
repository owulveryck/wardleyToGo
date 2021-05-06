package svg

import (
	"encoding/xml"
	"fmt"
	"image"
)

type Transform struct {
	Translate  image.Point
	Rotate     int
	Scale      float32
	Components []interface{}
}

func (t Transform) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	g := xml.StartElement{
		Name: xml.Name{Local: "g"},
	}
	var transformation string
	if !t.Translate.Eq(image.Point{}) {
		transformation = transformation + fmt.Sprintf(` translate(%v,%v)`, t.Translate.X, t.Translate.Y)
	}
	if t.Rotate != 0 {
		transformation = transformation + fmt.Sprintf(` rotate(%v)`, t.Rotate)
	}
	if t.Scale != 0 {
		transformation = transformation + fmt.Sprintf(` scale(%.2f)`, t.Scale)
	}

	g.Attr = []xml.Attr{
		{
			Name:  xml.Name{Local: "transform"},
			Value: transformation,
		},
	}
	e.EncodeToken(g)
	e.Encode(t.Components)
	e.EncodeToken(g.End())
	return nil
}
