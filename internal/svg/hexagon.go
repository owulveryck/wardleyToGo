package svg

import (
	"image"
	"math"
)

func hexCorner(center image.Point, size, i float64) image.Point {
	angleRad := 2 * math.Pi * i / 6
	return image.Point{
		X: center.X + int(size*math.Cos(angleRad)),
		Y: center.Y + int(size*math.Sin(angleRad)),
	}
}
