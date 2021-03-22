package plan

import (
	svg "github.com/ajstarks/svgo"
)

type ComponentType int

const (
	RegularComponent ComponentType = iota
	BuildComponent
	BuyComponent
	OutsourceComponent
)

const maxUint = ^uint(0)
const maxInt = int(maxUint >> 1)
const UndefinedCoord = -maxInt - 1

// A Component is an element of the map
type Component struct {
	id          int64
	Coords      [2]int
	Label       string
	LabelCoords [2]int
	Type        ComponentType
}

func NewComponent(id int64) *Component {
	return &Component{
		id:          id,
		Coords:      [2]int{UndefinedCoord, UndefinedCoord},
		LabelCoords: [2]int{UndefinedCoord, UndefinedCoord},
	}
}

func (c *Component) ID() int64 {
	return c.id
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
	switch c.Type {
	case BuildComponent:
		s.Circle(0, 0, 20, `fill="#D6D6D6"`, `stroke="#000000"`)
	case BuyComponent:
		s.Circle(0, 0, 20, `fill="#AAA5A9"`, `stroke="#D6D6D6"`)
	case OutsourceComponent:
		s.Circle(0, 0, 20, `fill="#444444"`, `stroke="#444444"`)
	}
	s.Circle(0, 0, 5, `stroke-width="1"`, `stroke="black"`, `fill="white"`)
	s.Gend()
}

func (c *Component) GetCoordinates() []int {
	return []int{c.Coords[0], c.Coords[1]}
}

func (c *Component) String() string {
	return c.Label
}

type EvolvedComponent struct {
	id          int64
	Coords      [2]int
	Label       string
	LabelCoords [2]int
	Type        ComponentType
}

func (e *EvolvedComponent) ID() int64 {
	return e.id
}

func NewEvolvedComponent(id int64) *EvolvedComponent {
	return &EvolvedComponent{
		id:          id,
		Coords:      [2]int{UndefinedCoord, UndefinedCoord},
		LabelCoords: [2]int{UndefinedCoord, UndefinedCoord},
	}
}

func (e *EvolvedComponent) SVG(s *svg.SVG, width, height, padLeft, padBottom int) {
	labelCoordX := e.LabelCoords[0]
	labelCoordY := e.LabelCoords[1]
	if labelCoordX == UndefinedCoord {
		labelCoordX = 10
	}
	if labelCoordY == UndefinedCoord {
		labelCoordY = 10
	}
	s.Translate(e.Coords[1]*(width-padLeft)/100+padLeft, (height-padLeft)-e.Coords[0]*(height-padLeft)/100)
	s.Text(labelCoordX, labelCoordY, e.Label, `fill="red"`)
	switch e.Type {
	case BuildComponent:
		s.Circle(0, 0, 20, `fill="#D6D6D6"`, `stroke="#000000"`)
	case BuyComponent:
		s.Circle(0, 0, 20, `fill="#AAA5A9"`, `stroke="#D6D6D6"`)
	case OutsourceComponent:
		s.Circle(0, 0, 20, `fill="#444444"`, `stroke="#444444"`)
	}
	s.Circle(0, 0, 5, `stroke-width="1"`, `stroke="red"`, `fill="white"`)
	s.Gend()
}

func (e *EvolvedComponent) GetCoordinates() []int {
	return []int{e.Coords[0], e.Coords[1]}
}

func (c *EvolvedComponent) String() string {
	return "[evolved]" + c.Label
}
