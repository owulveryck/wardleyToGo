package wtg

import (
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func findLeafs(m *wardleyToGo.Map) []*wardley.Component {
	ret := make([]*wardley.Component, 0)
	nodes := m.Nodes()
	for nodes.Next() {
		n := nodes.Node()
		if m.From(n.ID()).Len() == 0 {
			switch n := n.(type) {
			case *wardley.Component:
				ret = append(ret, n)
			case *wardley.EvolvedComponent:
				ret = append(ret, n.Component)
			}
		}
	}
	return ret
}

func findRoot(m *wardleyToGo.Map) []*wardley.Component {
	ret := make([]*wardley.Component, 0)
	nodes := m.Nodes()
	for nodes.Next() {
		n := nodes.Node()
		if m.To(n.ID()).Len() == 0 {
			switch n := n.(type) {
			case *wardley.Component:
				ret = append(ret, n)
			case *wardley.EvolvedComponent:
				ret = append(ret, n.Component)
			}
		}
	}
	return ret
}
