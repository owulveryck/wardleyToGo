package svgmap_test

import (
	"encoding/xml"
	"image"
	"image/color"
	"os"

	"github.com/owulveryck/wardleyToGo/internal/svg"
)

type component struct{}

func (c *component) MarshalSVG(e *xml.Encoder, bounds image.Rectangle) error {
	e.Encode(&svg.Transform{
		Translate: image.Point{50, 100},
		Components: []interface{}{
			&svg.Circle{
				Fill: svg.Color{color.Black},
			},
		},
	})
	return nil
}
func Example_customMarshaler() {
	c := &component{}
	enc := xml.NewEncoder(os.Stdout)
	c.MarshalSVG(enc, image.Rectangle{})
}
