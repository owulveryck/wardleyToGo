package wtg

import (
	"bufio"
	"image"
	"io"
	"regexp"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph/path"
)

func ParseValueChain(r io.Reader) (*wardleyToGo.Map, error) {
	inventory := make(map[string]*wardley.Component, 0)
	edgeInventory := make([]*wardley.Collaboration, 0)
	var link = regexp.MustCompile(`^\s*(.*\S)\s+(-+)\s+(.*)$`)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		elements := link.FindStringSubmatch(scanner.Text())
		if len(elements) != 4 {
			// log.Fatal("bad entry", scanner.Text())
			continue
		}
		if _, ok := inventory[elements[1]]; !ok {
			c := wardley.NewComponent(int64(len(inventory)))
			c.Label = elements[1]
			c.Placement = image.Pt(50, 50)
			inventory[elements[1]] = c
		}
		if _, ok := inventory[elements[3]]; !ok {
			c := wardley.NewComponent(int64(len(inventory)))
			c.Label = elements[3]
			c.Placement = image.Pt(50, 50)
			inventory[elements[3]] = c
		}
		edgeInventory = append(edgeInventory, &wardley.Collaboration{
			F:          inventory[elements[1]],
			T:          inventory[elements[3]],
			Type:       wardley.RegularEdge,
			Visibility: len(elements[2]),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	m := wardleyToGo.NewMap(0)
	for _, c := range inventory {
		err := m.AddComponent(c)
		if err != nil {
			return nil, err
		}
	}
	for _, e := range edgeInventory {
		err := m.SetCollaboration(e)
		if err != nil {
			return nil, err
		}
	}
	err := setCoords(m)
	if err != nil {
		return m, nil
	}
	return m, nil
}

func setCoords(m *wardleyToGo.Map) error {

	allShortestPaths := path.DijkstraAllPaths(m)
	roots := findRoot(m)
	leafs := findLeafs(m)
	var maxDepth int
	for _, r := range roots {
		for _, l := range leafs {
			paths, _ := allShortestPaths.AllBetween(r.ID(), l.ID())
			for _, path := range paths {
				currentVisibility := 0
				for i := 0; i < len(path)-1; i++ {
					e := m.Edge(path[i].ID(), path[i+1].ID())
					currentVisibility += e.(*wardley.Collaboration).Visibility
				}
				if currentVisibility > maxDepth {
					maxDepth = currentVisibility
				}
			}
		}
	}

	step := 100 / maxDepth
	cs := &coordSetter{
		verticalStep: step,
	}
	nroots := len(roots)
	hsteps := 100 / (nroots + 1)
	for i, n := range roots {
		n.Placement.X = hsteps * (i + 1)
		cs.walk(m, n, 0)
	}

	return nil
}

type coordSetter struct {
	verticalStep int
}

func (c *coordSetter) walk(m *wardleyToGo.Map, n *wardley.Component, visibility int) {
	n.Placement.Y = visibility * c.verticalStep
	from := m.From(n.ID())
	hsteps := 100 / (from.Len() + 1)
	i := 1
	for from.Next() {
		from.Node().(*wardley.Component).Placement.X = hsteps * i
		c.walk(m, from.Node().(*wardley.Component), m.Edge(n.ID(), from.Node().ID()).(*wardley.Collaboration).Visibility+visibility)
		i++
	}
}
