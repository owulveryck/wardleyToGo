package wardley

import svg "github.com/ajstarks/svgo"

// An Anchor of the map
type Anchor struct {
	Id     int64
	Coords [2]int
	Label  string
}

func (a *Anchor) GetCoordinates() [2]int {
	return a.Coords
}

func (c *Anchor) ID() int64 {
	return c.Id
}

func (c *Anchor) SVG(s *svg.SVG, width, height, padLeft, padBottom int) {
	s.Translate(c.Coords[1]*(width-padLeft)/100+padLeft, (height-padLeft)-c.Coords[0]*(height-padLeft)/100)
	s.Text(0, 0, c.Label, `font-weight="14px"`, `font-size="14px"`, `text-anchor="middle"`)
	s.Gend()
}
