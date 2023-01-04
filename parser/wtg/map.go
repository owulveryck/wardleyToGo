package wtg

import (
	"fmt"

	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func (p *Parser) consolidateMap() error {
	currentID := len(p.nodeInventory)
	for _, n := range p.nodeInventory {
		err := p.WMap.AddComponent(n)
		if err != nil {
			return err
		}
		if n.EvolutionPos != 0 {
			c := wardley.NewEvolvedComponent(int64(currentID))
			c.Placement.X = n.EvolutionPos
			c.Placement.Y = n.Placement.Y
			c.Label = n.Label
			err := p.WMap.AddComponent(c)
			if err != nil {
				return err
			}
			p.WMap.SetCollaboration(&wardley.Collaboration{
				F:    n,
				T:    c,
				Type: wardley.EvolvedComponentEdge,
			})
		}
		currentID++
	}
	for _, e := range p.edgeInventory {
		if e.F == nil || e.T == nil {
			return fmt.Errorf("bad edge: %v", e)
		}
		err := p.WMap.SetCollaboration(e)
		if err != nil {
			return err
		}
	}
	return nil
}
