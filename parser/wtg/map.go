package wtg

import (
	"errors"
	"fmt"
	"image"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph/topo"
)

func consolidateMap(nodeInventory map[string]*wardley.Component, edgeInventory []*wardley.Collaboration) (*wardleyToGo.Map, error) {
	wmap := wardleyToGo.NewMap(0)
	currentID := len(nodeInventory)
	for _, n := range nodeInventory {
		err := wmap.AddComponent(n)
		if err != nil {
			return nil, err
		}
		if n.EvolutionPos != 0 {
			c := wardley.NewEvolvedComponent(int64(currentID))
			c.Placement.X = n.EvolutionPos
			c.Placement.Y = n.Placement.Y
			c.Label = n.Label
			err := wmap.AddComponent(c)
			if err != nil {
				return nil, err
			}
			inertia := image.Point{}
			if n.Inertia != 0 {
				inertia = image.Point{
					X: n.Inertia,
				}
			}
			wmap.SetCollaboration(&wardley.Collaboration{
				F:       n,
				T:       c,
				Type:    wardley.EvolvedComponentEdge,
				Inertia: inertia,
			})
		}
		currentID++
	}
	for _, e := range edgeInventory {
		if e.F == nil || e.T == nil {
			return nil, fmt.Errorf("bad edge: %v", e)
		}
		if e.F == e.T {
			return nil, fmt.Errorf("self edge: F: %v, T: %v", e.F, e.T)
		}
		err := wmap.SetCollaboration(e)
		if err != nil {
			return nil, err
		}
	}
	cycles := topo.DirectedCyclesIn(wmap)
	if len(cycles) != 0 {
		return nil, errors.New("cycles detected in the map")
	}
	return wmap, nil
}
