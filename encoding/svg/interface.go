package svgmap

import (
	"encoding/xml"
	"image"

	"github.com/owulveryck/wardleyToGo"
)

type SVGMarshaler interface {
	MarshalSVG(e *xml.Encoder, bounds image.Rectangle) error
}

type SVGStyleMarshaler interface {
	MarshalStyleSVG(e *xml.Encoder, box, canvas image.Rectangle)
}

// Theme adds styling or interactivity to the SVG output.
// Themes are applied in order during encoding. A nil Themes slice
// on the Encoder produces pure SVG without any CSS or JS.
type Theme interface {
	Embed(enc *xml.Encoder, m *wardleyToGo.Map) error
}

// ComponentDecorator is optionally implemented by themes that need to modify
// component group attributes (e.g., adding onclick handlers).
type ComponentDecorator interface {
	DecorateComponent(c wardleyToGo.Component) []xml.Attr
}
