package owm

import (
	"fmt"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	"gonum.org/v1/gonum/graph"
)

type edge struct {
	ToLabel   string
	FromLabel string
	T         graph.Node
	F         graph.Node
	EdgeLabel string
	EdgeType  wardleyToGo.EdgeType
}

func (e edge) From() graph.Node {
	return e.F
}

func (e edge) ReversedEdge() graph.Edge {
	return edge{
		F:         e.T,
		T:         e.F,
		ToLabel:   e.FromLabel,
		FromLabel: e.ToLabel,
		EdgeLabel: e.EdgeLabel,
	}
}

func (e edge) To() graph.Node {
	return e.T
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
	e := p.g.Edges()
	for e.Next() {
		if e, ok := e.Edge().(*wardley.Collaboration); ok {
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
		p.g.SetEdge(&wardley.Collaboration{
			F:     edge.F.(wardleyToGo.Component),
			T:     edge.T.(wardleyToGo.Component),
			Type:  wardley.RegularEdge,
			Label: edge.EdgeLabel,
		})

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
		p.g.SetEdge(&wardley.Collaboration{
			F:    node.(wardleyToGo.Component),
			T:    nodeEvolved.(wardleyToGo.Component),
			Type: wardley.EvolvedComponentEdge,
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
			p.g.RemoveEdge(nodeEvolved.ID(), fromIT.Node().ID())
			p.g.SetEdge(&wardley.Collaboration{
				F:    nodeEvolved.(wardleyToGo.Component),
				T:    fromIT.Node().(wardleyToGo.Component),
				Type: wardley.EvolvedEdge,
			})
		}
		toIT := p.g.To(node.ID())
		for toIT.Next() {
			p.g.RemoveEdge(toIT.Node().ID(), nodeEvolved.ID())
			p.g.SetEdge(&wardley.Collaboration{
				F:    toIT.Node().(wardleyToGo.Component),
				T:    nodeEvolved.(wardleyToGo.Component),
				Type: wardley.EvolvedEdge,
			})
		}
	}
	return nil
}
