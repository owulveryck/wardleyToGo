package svgmap

import (
	"image"

	svg "github.com/ajstarks/svgo"
)

// SVGer is any object that can represent itself on a map
type SVGer interface {
	// SVG is a method that represent the object on the svg mag with coordinates relatives to the bounds
	SVG(s *svg.SVG, bounds image.Rectangle)
}
