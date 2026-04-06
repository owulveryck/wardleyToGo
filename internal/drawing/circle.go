package drawing

import (
	"image"
	"image/color"
	"image/draw"
)

// DrawCircle draws a filled circle with a stroke border directly at the
// target size. Pixels within radius (r-1) are filled with fill; pixels
// between (r-1) and r are drawn with stroke.
func DrawCircle(dst draw.Image, r int, p image.Point, stroke, fill color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			dist := x*x + y*y
			if dist <= (r-1)*(r-1) {
				dst.Set(p.X+x, p.Y+y, fill)
			} else if dist <= r*r {
				dst.Set(p.X+x, p.Y+y, stroke)
			}
		}
	}
}
