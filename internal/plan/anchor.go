package plan

import svg "github.com/ajstarks/svgo"

// An Anchor of the map
type Anchor struct {
	id     int64
	Coords [2]int
	Label  string
}

// GetCoordinates fulfils the Element interface
func (a *Anchor) GetCoordinates() []int {
	return []int{a.Coords[0], a.Coords[1]}
}

func NewAnchor(id int64) *Anchor {
	return &Anchor{
		id:     id,
		Coords: [2]int{UndefinedCoord, UndefinedCoord},
	}
}

func (c *Anchor) ID() int64 {
	return c.id
}

func (c *Anchor) SVG(s *svg.SVG, width, height, padLeft, padBottom int) {
	s.Translate(c.Coords[1]*(width-padLeft)/100+padLeft, (height-padLeft)-c.Coords[0]*(height-padLeft)/100)
	s.Text(0, 0, c.Label, `font-weight="14px"`, `font-size="14px"`, `text-anchor="middle"`)
	s.Gend()
}

func (c *Anchor) String() string {
	return c.Label
}
