package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo/internal/drawing"
	"github.com/owulveryck/wardleyToGo/parser/owm"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func createBackground(im draw.Image, canvas image.Rectangle) {

	markers := []struct {
		position float64
		label    string
	}{
		{
			position: 0,
			label:    "Genesis",
		},
		{
			position: (float64(100) / 575),
			label:    "Custom-Built",
		},
		{
			position: (float64(100) / 250),
			label:    "Product\n(+rental)",
		},
		{
			position: (float64(574) / 820),
			label:    "Commodity\n(+utility)",
		},
	}

	for i := 0; i < len(markers); i++ {
		axis := markers[i]
		position := image.Point{
			X: int(float64(canvas.Dx()) * axis.position),
		}
		position = position.Add(canvas.Min)
		drawing.Line(im, position.X, position.Y, position.X, position.Y+canvas.Dy(), color.Gray{Y: 128}, [2]int{2, 5})
		dot := fixed.P(position.X, position.Y+canvas.Dy()+15)
		d := font.Drawer{
			Dst:  im,
			Src:  image.NewUniform(color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}),
			Face: basicfont.Face7x13,
			Dot:  dot,
		}
		d.DrawString(axis.label)
		//		w.Text(int(float64(w.canvas.Dx())*axis.position)+w.canvas.Min.X, w.canvas.Max.Y+15, axis.label)
	}

	drawing.Line(im, canvas.Min.X, canvas.Min.Y, canvas.Min.X, canvas.Max.Y, color.Black, [2]int{0, 0})
	drawing.Line(im, canvas.Min.X, canvas.Max.Y, canvas.Max.X, canvas.Max.Y, color.Black, [2]int{0, 0})
}

func main() {
	p := owm.NewParser(os.Stdin)
	m, err := p.Parse() // the map
	if err != nil {
		log.Fatal(err)
	}
	im := image.NewRGBA(image.Rect(0, 0, 1400, 1100))
	canvas := image.Rect(100, 100, 1300, 1000)
	createBackground(im, canvas)

	m.Draw(im, canvas, im, image.Point{})
	png.Encode(os.Stdout, im)
}
