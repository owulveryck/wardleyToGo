package svg

import (
	"encoding/xml"
	"fmt"
	"image/color"
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
