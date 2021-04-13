package drawing

import (
	"image/color"
	"image/draw"
)

type pen struct {
	dashArray [2]int
	cur       int
	img       draw.Image
}

func (p *pen) Set(x, y int, col color.Color) {
	if p.dashArray == [2]int{0, 0} {
		p.img.Set(x, y, col)
		return
	}
	if p.cur%(p.dashArray[0]+p.dashArray[1]) == 0 {
		p.cur = 0
	}
	if p.cur < p.dashArray[0] {
		p.img.Set(x, y, col)
	}
	p.cur++
}

func Arrow(img draw.Image, x1, y1, x2, y2 int, col color.Color, dashArray [2]int) {
	Line(img, x1, y1, x2, y2, col, dashArray)
	Line(img, x2-10, y2, x2-20, y2-4, col, [2]int{0, 0})
	Line(img, x2-10, y2, x2-20, y2+4, col, [2]int{0, 0})

}

// The Line algorithm is taken from github.com/StephaneBunel/bresenham
func Line(img draw.Image, x1, y1, x2, y2 int, col color.Color, dashArray [2]int) {
	p := &pen{
		img:       img,
		dashArray: dashArray,
	}

	var dx, dy, e, slope int

	// Because drawing p1 -> p2 is equivalent to draw p2 -> p1,
	// I sort points in x-axis order to handle only half of possible cases.
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	// Because point is x-axis ordered, dx cannot be negative
	if dy < 0 {
		dy = -dy
	}

	switch {

	// Is line a point ?
	case x1 == x2 && y1 == y2:
		p.Set(x1, y1, col)

	// Is line an horizontal ?
	case y1 == y2:
		for ; dx != 0; dx-- {
			p.Set(x1, y1, col)
			x1++
		}
		p.Set(x1, y1, col)

	// Is line a vertical ?
	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for ; dy != 0; dy-- {
			p.Set(x1, y1, col)
			y1++
		}
		p.Set(x1, y1, col)

	// Is line a diagonal ?
	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				p.Set(x1, y1, col)
				x1++
				y1++
			}
		} else {
			for ; dx != 0; dx-- {
				p.Set(x1, y1, col)
				x1++
				y1--
			}
		}
		p.Set(x1, y1, col)

	// wider than high ?
	case dx > dy:
		if y1 < y2 {
			// BresenhamDxXRYD(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				p.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
			}
		} else {
			// BresenhamDxXRYU(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				p.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
			}
		}
		img.Set(x2, y2, col)

	// higher than wide.
	default:
		if y1 < y2 {
			// BresenhamDyXRYD(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				p.Set(x1, y1, col)
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		} else {
			// BresenhamDyXRYU(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				p.Set(x1, y1, col)
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		}
		p.Set(x2, y2, col)
	}
}
