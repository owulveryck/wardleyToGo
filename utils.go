package wardleyToGo

import "image"

func calcCoords(p image.Point, bounds image.Rectangle) image.Point {
	scale := bounds.Max.Sub(bounds.Min)
	scaleX := float64(scale.X) / 100
	scaleY := float64(scale.Y) / 100
	dest := image.Pt(int(float64(p.X)*scaleX), int(float64(p.Y)*scaleY))
	return dest.Add(bounds.Min)
}
