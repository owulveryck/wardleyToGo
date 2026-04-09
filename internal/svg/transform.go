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
	Classes    []string
}

func (t Transform) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	g := xml.StartElement{
		Name: xml.Name{Local: "g"},
	}
	attrs := newAttributes()
	if t.Classes != nil {
		classes := ""
		for i := range t.Classes {
			classes = classes + " " + t.Classes[i]
		}
		attrs = attrs.append("class", classes)
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
	attrs = attrs.append("transform", transformation)
	g.Attr = attrs
	/*
		g.Attr = []xml.Attr{
			{
				Name:  xml.Name{Local: "transform"},
				Value: transformation,
			},
		}
	*/
	if err := e.EncodeToken(g); err != nil {
		return err
	}
	if err := e.Encode(t.Components); err != nil {
		return err
	}
	return e.EncodeToken(g.End())
}
