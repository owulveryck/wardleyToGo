package plan

import (
	svg "github.com/ajstarks/svgo"
)

type StreamAlignedTeam struct {
	Id     int64
	Coords [4]int
	Label  string
}

func (s *StreamAlignedTeam) String() string {
	return s.Label
}

func NewStreamAlignedTeam(id int64) *StreamAlignedTeam {
	return &StreamAlignedTeam{
		Id:     id,
		Coords: [4]int{UndefinedCoord, UndefinedCoord, UndefinedCoord, UndefinedCoord},
	}
}

func (s *StreamAlignedTeam) SVG(svg *svg.SVG, width, height, padLeft, padBottom int) {
	x1 := s.Coords[1]*(width-padLeft)/100 + padLeft
	y1 := (height - padLeft) - s.Coords[0]*(height-padLeft)/100
	x2 := s.Coords[3]*(width-padLeft)/100 + padLeft
	y2 := (height - padLeft) - s.Coords[2]*(height-padLeft)/100
	svg.Translate(x1, y1)
	svg.Roundrect(0, 0, abs(x2-x1), abs(y2-y1), 15, 15, `fill="rgb(252, 237, 190)"`, `opacity="0.8"`, `stroke="rgb(250,216,120)"`, `stroke-opacity="0.7"`, `stroke-width="5px"`)
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

func (e *EnablingTeam) String() string {
	return e.Label
}
func NewEnablingTeam(id int64) *EnablingTeam {
	return &EnablingTeam{
		Id:     id,
		Coords: [4]int{UndefinedCoord, UndefinedCoord, UndefinedCoord, UndefinedCoord},
	}
}

func (e *EnablingTeam) SVG(svg *svg.SVG, width, height, padLeft, padBottom int) {
	x1 := e.Coords[1]*(width-padLeft)/100 + padLeft
	y1 := (height - padLeft) - e.Coords[0]*(height-padLeft)/100
	x2 := e.Coords[3]*(width-padLeft)/100 + padLeft
	y2 := (height - padLeft) - e.Coords[2]*(height-padLeft)/100
	svg.Translate(x1, y1)
	svg.Roundrect(0, 0, abs(x2-x1), abs(y2-y1), 15, 15, `fill="rgb(217, 190, 206)"`, `opacity="0.8"`, `stroke="rgb(200,159,182)"`, `stroke-opacity="0.7"`, `stroke-width="5px"`)
	svg.Gend()
}

func (e *EnablingTeam) ID() int64 {
	return e.Id
}

type PlatformTeam struct {
	Id     int64
	Coords [4]int
	Label  string
}

func NewPlatformTeam(id int64) *PlatformTeam {
	return &PlatformTeam{
		Id:     id,
		Coords: [4]int{UndefinedCoord, UndefinedCoord, UndefinedCoord, UndefinedCoord},
	}
}

func (p *PlatformTeam) SVG(svg *svg.SVG, width, height, padLeft, padBottom int) {
	x1 := p.Coords[1]*(width-padLeft)/100 + padLeft
	y1 := (height - padLeft) - p.Coords[0]*(height-padLeft)/100
	x2 := p.Coords[3]*(width-padLeft)/100 + padLeft
	y2 := (height - padLeft) - p.Coords[2]*(height-padLeft)/100
	svg.Translate(x1, y1)
	svg.Rect(0, 0, abs(x2-x1), abs(y2-y1), `fill="rgb(170, 185, 215)"`, `opacity="0.8"`, `stroke="rgb(119,159,229)"`, `stroke-opacity="0.7"`, `stroke-width="5px"`)
	svg.Gend()
}

func (p *PlatformTeam) ID() int64 {
	return p.Id
}
func (p *PlatformTeam) String() string {
	return p.Label
}

type ComplicatedSubsystemTeam struct {
	Id     int64
	Coords [4]int
	Label  string
}

func NewComplicatedSubsystemTeam(id int64) *ComplicatedSubsystemTeam {
	return &ComplicatedSubsystemTeam{
		Id:     id,
		Coords: [4]int{UndefinedCoord, UndefinedCoord, UndefinedCoord, UndefinedCoord},
	}
}

func (p *ComplicatedSubsystemTeam) SVG(svg *svg.SVG, width, height, padLeft, padBottom int) {
	x1 := p.Coords[1]*(width-padLeft)/100 + padLeft
	y1 := (height - padLeft) - p.Coords[0]*(height-padLeft)/100
	x2 := p.Coords[3]*(width-padLeft)/100 + padLeft
	y2 := (height - padLeft) - p.Coords[2]*(height-padLeft)/100
	svg.Translate(x1, y1)
	svg.Roundrect(0, 0, abs(x2-x1), abs(y2-y1), 35, 35, `fill="rgb(236, 210, 177)"`, `opacity="0.8"`, `stroke="rgb(210,149,84)"`, `stroke-opacity="0.7"`, `stroke-width="5px"`)
	svg.Gend()
}

func (p *ComplicatedSubsystemTeam) ID() int64 {
	return p.Id
}
func (p *ComplicatedSubsystemTeam) String() string {
	return p.Label
}

// Abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (p *ComplicatedSubsystemTeam) GetCoordinates() []int {
	return []int{p.Coords[0], p.Coords[1], p.Coords[2], p.Coords[3]}
}

func (p *PlatformTeam) GetCoordinates() []int {
	return []int{p.Coords[0], p.Coords[1], p.Coords[2], p.Coords[3]}
}
func (p *EnablingTeam) GetCoordinates() []int {
	return []int{p.Coords[0], p.Coords[1], p.Coords[2], p.Coords[3]}
}
func (p *StreamAlignedTeam) GetCoordinates() []int {
	return []int{p.Coords[0], p.Coords[1], p.Coords[2], p.Coords[3]}
}
