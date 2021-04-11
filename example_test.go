package wardleyToGo_test

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/internal/utils"
)

type dummyComponent struct {
	id       int64
	position image.Point
}

func (d *dummyComponent) GetPosition() image.Point { return d.position }

func (d *dummyComponent) ID() int64 { return d.id }

func (d *dummyComponent) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	log.Println(dst)
	log.Println(r)
	log.Println(src)
	log.Println(sp)
}

func Example() {
	// Create a new map
	m := wardleyToGo.NewMap(0)
	c0 := &dummyComponent{id: 0, position: image.Pt(25, 25)}
	c1 := &dummyComponent{id: 1, position: image.Pt(50, 50)}
	m.AddComponent(c0)
	m.AddComponent(c1)

	// Very trivial example to draw a map on stdout
	drawMap(m)
}

//drawMap on stdout
func drawMap(m *wardleyToGo.Map) {
	// Creates a picture representation of the map
	const width = 100
	const height = 20

	im := image.NewGray(image.Rectangle{Max: image.Point{X: width, Y: height}})
	nodes := m.Nodes()
	for nodes.Next() {
		n := nodes.Node().(wardleyToGo.Component)
		// CalcCoords adapt the coordinates accrording to the map
		coords := utils.CalcCoords(n.GetPosition(), im.Bounds())
		im.SetGray(coords.X, coords.Y, color.Gray{Y: 255})
	}
	pi := image.NewPaletted(im.Bounds(), []color.Color{
		color.Gray{Y: 255},
		color.Gray{Y: 160},
		color.Gray{Y: 70},
		color.Gray{Y: 35},
		color.Gray{Y: 0},
	})

	draw.FloydSteinberg.Draw(pi, im.Bounds(), im, image.Point{})
	shade := []string{" ", "░", "▒", "▓", "█"}
	for i, p := range pi.Pix {
		fmt.Print(shade[p])
		if (i+1)%width == 0 {
			fmt.Print("\n")
		}
	}
}
