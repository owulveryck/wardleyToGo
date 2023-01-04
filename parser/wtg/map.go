package wtg

import "github.com/owulveryck/wardleyToGo/components/wardley"

func (p *Parser) consolidateMap() error {
	currentID := len(p.nodeInventory)
	for _, n := range p.nodeInventory {
		err := p.WMap.AddComponent(n)
		if err != nil {
			return err
		}
		if n.EvolutionPos != 0 {
			c := wardley.NewEvolvedComponent(int64(currentID))
			c.Placement.X = n.Placement.X
			c.Placement.Y = n.Placement.Y
			c.Label = n.Label
			err := p.WMap.AddComponent(c)
			if err != nil {
				return err
			}
			p.WMap.SetCollaboration(&wardley.Collaboration{
				F:    n,
				T:    c,
				Type: wardley.EvolvedEdge,
			})
		}
		currentID++
	}
	for _, e := range p.edgeInventory {
		err := p.WMap.SetCollaboration(e)
		if err != nil {
			return err
		}
	}
	return nil
}
