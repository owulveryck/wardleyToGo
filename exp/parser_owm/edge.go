package owm

import (
	"fmt"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

type edge struct {
	ToLabel   string
	FromLabel string
	T         wardleyToGo.Component
	F         wardleyToGo.Component
	EdgeLabel string
	EdgeType  wardleyToGo.EdgeType
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
	for _, c := range p.m.Collaborations() {
		if e, ok := c.(*wardley.Collaboration); ok {
			e.RenderingLayer = 5
		}
	}
	return nil
}

func (p *Parser) createRegularEdges() error {
	var ok bool
	for _, edge := range p.edges {
		edge.F, ok = p.nodeDict[edge.FromLabel]
		if !ok {
			return fmt.Errorf("graph is inconsistent, %v is referencing a non-defined node", edge)
		}
		edge.T, ok = p.nodeDict[edge.ToLabel]
		if !ok {
			return fmt.Errorf("graph is inconsistent, %v is referencing a non-defined node", edge)
		}
		p.m.SetCollaboration(&wardley.Collaboration{
			F:     edge.F,
			T:     edge.T,
			Type:  wardley.RegularEdge,
			Label: edge.EdgeLabel,
		})

	}
	return nil
}

func (p *Parser) createEvolvingComponentEdges() error {
	for name, nodeEvolved := range p.nodeEvolveDict {
		node, ok := p.nodeDict[name]
		if !ok {
			return fmt.Errorf("bad evolution, non existent component %v", name)
		}
		p.m.SetCollaboration(&wardley.Collaboration{
			F:    node,
			T:    nodeEvolved,
			Type: wardley.EvolvedComponentEdge,
		})
	}
	return nil
}

func (p *Parser) createEvolvingEdges() error {
	for name, nodeEvolved := range p.nodeEvolveDict {
		node, ok := p.nodeDict[name]
		if !ok {
			return fmt.Errorf("bad evolution, non existent component %v", name)
		}
		for _, succ := range p.m.From(node.ID()) {
			p.m.RemoveEdge(nodeEvolved.ID(), succ.ID())
			p.m.SetCollaboration(&wardley.Collaboration{
				F:    nodeEvolved,
				T:    succ,
				Type: wardley.EvolvedEdge,
			})
		}
		for _, pred := range p.m.To(node.ID()) {
			p.m.RemoveEdge(pred.ID(), nodeEvolved.ID())
			p.m.SetCollaboration(&wardley.Collaboration{
				F:    pred,
				T:    nodeEvolved,
				Type: wardley.EvolvedEdge,
			})
		}
	}
	return nil
}
