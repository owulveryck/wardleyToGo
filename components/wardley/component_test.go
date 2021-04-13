package wardley

import (
	"image"

	"github.com/owulveryck/wardleyToGo/components"
)

func ExampleComponent_Draw() {
	// draw a circle at the center of the picture
	im := image.NewRGBA(image.Rect(0, 0, 1200, 1000))
	c := &Component{
		Label:          "test",
		LabelPlacement: image.Pt(components.UndefinedCoord, components.UndefinedCoord),
		Placement:      image.Pt(100, 50),
	}
	c.Draw(im, im.Bounds(), im, image.Point{})
}
