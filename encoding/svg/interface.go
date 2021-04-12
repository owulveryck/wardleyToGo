package svgmap

import (
	"image"

	svg "github.com/ajstarks/svgo"
)

// SVGDrawer is any object that can represent itself on a map
type SVGDrawer interface {
	// SVG is a method that represent the object on the svg mag with coordinates relatives to the bounds
	SVGDraw(s *svg.SVG, bounds image.Rectangle)
}
