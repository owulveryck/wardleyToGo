package wardley

import (
	svg "github.com/ajstarks/svgo"
)

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)
const UndefinedCoord = -maxInt - 1

// A Component is an element of the map
type Component struct {
	Id          int64
	Coords      [2]int
	Label       string
	LabelCoords [2]int
}

func (c *Component) ID() int64 {
	return c.Id
}

func (c *Component) SVG(s *svg.SVG, width, height, padLeft, padBottom int) {
	labelCoordX := c.LabelCoords[0]
	labelCoordY := c.LabelCoords[1]
	if labelCoordX == UndefinedCoord {
		labelCoordX = 10
	}
	if labelCoordY == UndefinedCoord {
		labelCoordY = 10
	}
	s.Translate(c.Coords[1]*(width-padLeft)/100+padLeft, (height-padLeft)-c.Coords[0]*(height-padLeft)/100)
	s.Text(labelCoordX, labelCoordY, c.Label)
	s.Circle(0, 0, 5, `stroke-width="1"`, `stroke="black"`, `fill="white"`)
	s.Gend()
}

func (c *Component) GetCoordinates() []int {
	return []int{c.Coords[0], c.Coords[1]}
}
