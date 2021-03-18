package wardley

import (
	svg "github.com/ajstarks/svgo"
)

type StreamAlignedTeam struct {
	Id     int64
	Coords [4]int
	Label  string
}

func (s *StreamAlignedTeam) SVG(svg *svg.SVG, width, height, padLeft, padBottom int) {
	x1 := s.Coords[1]*(width-padLeft)/100 + padLeft
	y1 := (height - padLeft) - s.Coords[0]*(height-padLeft)/100
	x2 := s.Coords[3]*(width-padLeft)/100 + padLeft
	y2 := (height - padLeft) - s.Coords[2]*(height-padLeft)/100
	svg.Translate(x1, y1)
	svg.Roundrect(0, 0, abs(x2-x1), abs(y2-y1), 15, 15, `fill="rgb(240, 215, 155)"`, `opacity="0.7"`, `stroke="rgb(245,198,93)"`, `stroke-opacity="0.7"`, `stroke-width="5px"`)
	svg.Gend()
}

func (s *StreamAlignedTeam) ID() int64 {
	return s.Id
}

type EnablingTeam struct {
	Id     int64
	Coords [4]int
	Label  string
}

func (e *EnablingTeam) SVG(svg *svg.SVG, width, height, padLeft, padBottom int) {
	x1 := e.Coords[1]*(width-padLeft)/100 + padLeft
	y1 := (height - padLeft) - e.Coords[0]*(height-padLeft)/100
	x2 := e.Coords[3]*(width-padLeft)/100 + padLeft
	y2 := (height - padLeft) - e.Coords[2]*(height-padLeft)/100
	svg.Translate(x1, y1)
	svg.Roundrect(0, 0, abs(x2-x1), abs(y2-y1), 15, 15, `fill="rgb(184, 71, 214)"`, `opacity="0.7"`, `stroke="rgb(111,85,155)"`, `stroke-opacity="0.7"`, `stroke-width="5px"`)
	svg.Gend()
}

func (e *EnablingTeam) ID() int64 {
	return e.Id
}

type ComplicatedSubsystemTeam struct{}

type PlatformTeam struct{}

// Abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
