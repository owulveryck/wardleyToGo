package drawing

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	drawX "golang.org/x/image/draw"
)

func DrawCircle(dst draw.Image, r int, p image.Point, stroke, fill color.Color) {
	rayon := 80
	sr := image.Rect(0, 0, rayon*2, rayon*2)
	circle := image.NewRGBA(sr)
	for j := 0.0; j <= float64(75); j++ {
		for i := 0.0; i < 2*math.Pi; i += 0.01 {
			pt := image.Pt(int(math.Sin(i)*j), int(math.Cos(i)*j)).Add(image.Pt(rayon, rayon))
			//circle.Set(pt.X, pt.Y, color.RGBA{0xff, 0x00, 0x00, 0xff})
			circle.Set(pt.X, pt.Y, fill)
		}
	}
	for j := 75.0; j <= float64(rayon); j++ {
		for i := 0.0; i < 2*math.Pi; i += 0.01 {
			pt := image.Pt(int(math.Sin(i)*j), int(math.Cos(i)*j)).Add(image.Pt(rayon, rayon))
			//circle.Set(pt.X, pt.Y, color.RGBA{0xff, 0x00, 0x00, 0xff})
			circle.Set(pt.X, pt.Y, stroke)
		}
	}
	drawX.CatmullRom.Scale(dst, image.Rect(-r, -r, r, r).Add(p), circle, sr, draw.Over, nil)
}
