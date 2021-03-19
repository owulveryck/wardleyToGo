package parser

import (
	"fmt"

	"github.com/owulveryck/wardleyToGo/internal/wardley"
)

func (p *Parser) completeEvolve() error {
	for name, nodeEvolve := range p.nodeEvolveDict {
		node, ok := p.nodeDict[name]
		if !ok {
			return fmt.Errorf("bad evolution, non existent component %v", name)
		}
		nodeEvolve.(*wardley.EvolvedComponent).Coords[0] = node.(*wardley.Component).Coords[0]
	}
	return nil
}
func (p *Parser) createEdges() error {
	err := p.createRegularEdges()
	if err != nil {
		return err
	}
	err = p.createEvolvingEdges()
	if err != nil {
		return err
	}
	err = p.createEvolvingComponentEdges()
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) createRegularEdges() error {
	var ok bool
	for _, edge := range p.edges {
		edge.F, ok = p.nodeDict[edge.fromLabel]
		if !ok {
			return fmt.Errorf("graph is inconsistent, %v is referencing a non-defined node", edge)
		}
		edge.T, ok = p.nodeDict[edge.toLabel]
		if !ok {
			return fmt.Errorf("graph is inconsistent, %v is referencing a non-defined node", edge)
		}
		p.g.SetEdge(edge)
	}
	return nil
}

func (p *Parser) createEvolvingComponentEdges() error {
	// TODO
	for name, nodeEvolved := range p.nodeEvolveDict {
		node, ok := p.nodeDict[name]
		if !ok {
			return fmt.Errorf("bad evolution, non existent component %v", name)
		}
		p.g.SetEdge(edge{
			F:        node,
			T:        nodeEvolved,
			edgeType: wardley.EvolvedComponentEdge,
		})
	}
	return nil
}

func (p *Parser) createEvolvingEdges() error {
	// TODO
	for name, nodeEvolved := range p.nodeEvolveDict {
		node, ok := p.nodeDict[name]
		if !ok {
			return fmt.Errorf("bad evolution, non existent component %v", name)
		}
		fromIT := p.g.From(node.ID())
		for fromIT.Next() {
			p.g.SetEdge(edge{
				F:        nodeEvolved,
				T:        fromIT.Node(),
				edgeType: wardley.EvolvedEdge,
			})
		}
		toIT := p.g.To(node.ID())
		for toIT.Next() {
			p.g.SetEdge(edge{
				F:        toIT.Node(),
				T:        nodeEvolved,
				edgeType: wardley.EvolvedEdge,
			})
		}
	}
	return nil
}