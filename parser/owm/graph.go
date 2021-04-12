package owm

import (
	"fmt"

	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func (p *Parser) completeEvolve() error {
	for name, nodeEvolve := range p.nodeEvolveDict {
		node, ok := p.nodeDict[name]
		if !ok {
			return fmt.Errorf("bad evolution, non existent component %v", name)
		}
		nodeEvolve.(*wardley.EvolvedComponent).Placement = node.(*wardley.Component).Placement
	}
	return nil
}
