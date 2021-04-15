package svgmap

import (
	"encoding/xml"
	"image"
)

type SVGMarshaler interface {
	MarshalSVG(e *xml.Encoder, bounds image.Rectangle) error
}

type SVGStyleMarshaler interface {
	MarshalStyleSVG(e *xml.Encoder, box, canvas image.Rectangle)
}
