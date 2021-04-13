package components

import "image"

// CalclCoords calculates the coordinates wrt to the bounds.
// it scales accordingly
func CalcCoords(p image.Point, bounds image.Rectangle) image.Point {
	scaleX := float64(bounds.Dx()) / 100
	scaleY := float64(bounds.Dy()) / 100
	dest := image.Pt(int(float64(p.X)*scaleX), int(float64(p.Y)*scaleY))
	return dest.Add(bounds.Min)
}
