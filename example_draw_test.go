package wardleyToGo_test

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/internal/utils"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

type dummyComponent struct {
	id       int64
	position image.Point
}

func (d *dummyComponent) GetPosition() image.Point { return d.position }

func (d *dummyComponent) ID() int64 { return d.id }

func (d *dummyComponent) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	coords := utils.CalcCoords(d.position, r)
	dst.Set(coords.X, coords.Y, color.Gray{Y: 255})
}

type dummyCollaboration struct{ simple.Edge }

func (d *dummyCollaboration) GetType() wardleyToGo.EdgeType { return 0 }

func newCollaboration(a, b wardleyToGo.Component) wardleyToGo.Collaboration {
	return &dummyCollaboration{Edge: simple.Edge{F: a, T: b}}
}

func Example_canvas() {
	// Create a new map
	m := wardleyToGo.NewMap(0)
	c0 := &dummyComponent{id: 0, position: image.Pt(25, 25)}
	c1 := &dummyComponent{id: 1, position: image.Pt(50, 50)}
	c2 := &dummyComponent{id: 2, position: image.Pt(50, 75)}
	c3 := &dummyComponent{id: 3, position: image.Pt(75, 75)}
	m.AddComponent(c0)
	m.AddComponent(c1)
	m.AddComponent(c2)
	m.AddComponent(c3)
	// c0 -> c1
	// c1 -> c2
	// c2 -> c3
	// c1 -> c3
	m.SetCollaboration(newCollaboration(c0, c1))
	m.SetCollaboration(newCollaboration(c1, c2))
	m.SetCollaboration(newCollaboration(c2, c3))
	m.SetCollaboration(newCollaboration(c1, c3))

	// Creates a picture representation of the map
	const width = 80
	const height = 40

	im := image.NewGray(image.Rectangle{Max: image.Point{X: width, Y: height}})
	m.Canvas = &simpleCanvas{}

	m.Draw(im, image.Rect(5, 2, 75, 38), im, image.Point{X: 0, Y: 0})
	//m.Draw(im, im.Bounds(), im, image.Point{X: 0, Y: 0})
	// Very trivial example to draw a map on stdout
	render(im)

	//	drawMap(m)

	// Find the shortest path betwen c0 and c3
	p, _ := path.AStar(c0, c3, m, euclideanDistance)
	c0Toc3, _ := p.To(c3.ID())
	fmt.Printf("Shortest path from c0 to c3: ")
	for _, c := range c0Toc3 {
		fmt.Printf("-%v", c.ID())
	}
	//Output:
	//████████████████████████████████████████████████████████████████████████████████
	//████████████████████████████████████████████████████████████████████████████████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█████
	//█████▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█████
	//█████▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█████
	//█████▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█████
	//█████▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█████
	//█████▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓ ▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█████
	//█████▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█████
	//█████▒▓▒▒▓▒▒▓▒▒▓▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒█████
	//█████▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▓▒▒▓▒▒▓▒▓▒▓▒▓▒▒▓▒▓▒▓▒▓▒▒▓▒▓▒▓▒▓▒▒▓▒▓▒▓▒▒▓▒▒▒▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒█████
	//█████▒▓▒▒▓▒▒▓▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▓▒▒▓▒▒▒▒▒▓▒▒▒▓▒▒▒▒▓▒▒▒▓▒▒▒▒▓▒▒▒▓▒▒▓▒▒▒▒▒▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ▓▒▒▒▓▒▒▒▒▓▒▒▒▓▒▒▒▒▓▒▒▒▓▒▒▒▒▒▒▓▒▒▓▒█████
	//█████▒▒▓▒▒▓▒▒▒▓▒▒▓▒▒▓▒▒▓▒▓▒▒▓▒▒▓▒▒▒▓▒▒▓▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒█████
	//█████▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▓▒▒▓▒▒▓▒▒▒▒▒▓▒▒▓▒▒▒▒▒▓▒▒▓▒▒▒▒▒▓▒▒▓▒█████
	//█████▒▓▒▒▓▒▒▒▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒█████
	//█████▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▒▒▓▒▒▓▒▒▓▒▒▒▓▒▒▒▒▓▒▒▒▓▒▒▒▒▓▒▒▓▒▒▒▒▓▒▒▓▒█████
	//█████▒▒▓▒▒▓▒▒▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒█████
	//█████▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▓▒▒▓▒▒▒▒▒▒▓▒▒▓▒▒▒▒▒▓▒▒▓▒▒▓▒▒▒▒▒▓▒▒▒█████
	//█████▒▒▓▒▒▓▒▒▒▒▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▒▒▓▒▒▓▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▒▓▒▒▒▒▓█████
	//█████▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▓▒▒▓▒▒▒▒▓▒▒▓▒▒▓▒▒▒▒▒▓▒▒▓▒▒▓▒▒▒▒▒▓▒▒█████
	//█████▒▓▒▒▓▒▒▓▒▒▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▒▒▒▒▓▒▒▒ ▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▓ ▒▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▓▒▒▒▓▒▒▓▒▒▒▒▒▒▓▒▒▓▒▒▓▒▒▒▒▓▒▒▓▒▒▓▒▒▒▒▓▒▒▒▒█████
	//█████▒▒▓▒▒▓▒▒▓▒▒▒▒▓▒▒▓▒▒▓▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▒▓▒▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▒▒▒▒▒▓▒▒▓▒▒▓▒▓▒▒▓▒▒▓▒▒▓▒▒▒▒▒▒▒▒█████
	//█████▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▓▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▒▒▒▒▒▒▓▒▒▓▒▒▒▓▒▒▓▒▒▓▒▒▓▒▒▒▒▒▒▒▒█████
	//█████▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▒▓▒▒▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▓▒▒▒█████
	//█████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▓▒▒▒▓▒▒▒▒▒▓▒▒▒▒▒▓▒▒▓▒▒▓▒▒▒▒▒▒▒▒▓▒█████
	//█████▒▓▒▒▓▒▒▓▒▒▓▒▓▒▒▓▒▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒▒▒▒▓▒▒▒▒▒▒▒▒▒▒▓▒▒▓▒▒▒▒▒█████
	//████████████████████████████████████████████████████████████████████████████████
	//████████████████████████████████████████████████████████████████████████████████
	//Shortest path from c0 to c3: -0-1-3
}

type simpleCanvas struct{}

func (s *simpleCanvas) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	draw.Draw(dst, r, image.NewUniform(color.Gray{Y: 64}), sp, draw.Src)
}

func render(im image.Image) {
	width := im.Bounds().Dx()
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

var euclideanDistance path.Heuristic = func(x, y graph.Node) float64 {
	xC := x.(wardleyToGo.Component).GetPosition()
	yC := y.(wardleyToGo.Component).GetPosition()
	a := xC.X - yC.X
	b := xC.Y - yC.Y
	return math.Sqrt(float64(a*a) + float64(b*b))
}
