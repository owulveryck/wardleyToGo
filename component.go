package main

import svg "github.com/ajstarks/svgo"

type component struct {
	id         int64
	coords     [2]int
	label      string
	labelCoord [2]int
}

func (c *component) ID() int64 {
	return c.id
}

func (c *component) SVG(s *svg.SVG, width, height, padLeft, padBottom int) {
	labelCoordX := c.labelCoord[0]
	labelCoordY := c.labelCoord[1]
	if labelCoordX <= 0 {
		labelCoordX = 10
	}
	if labelCoordY <= 0 {
		labelCoordY = 10
	}
	s.Translate(c.coords[1]*(width-padLeft)/100+padLeft, (height-padLeft)-c.coords[0]*(height-padLeft)/100)
	s.Text(labelCoordX, labelCoordY, c.label)
	s.Circle(0, 0, 5, `stroke-width="1"`, `stroke="black"`, `fill="white"`)
	s.Gend()
}

type anchor struct {
	id     int64
	coords [2]int
	label  string
}

func (c *anchor) ID() int64 {
	return c.id
}

func (c *anchor) SVG(s *svg.SVG, width, height, padLeft, padBottom int) {
	s.Translate(c.coords[1]*(width-padLeft)/100+padLeft, (height-padLeft)-c.coords[0]*(height-padLeft)/100)
	s.Text(0, 0, c.label)
	s.Gend()
}
