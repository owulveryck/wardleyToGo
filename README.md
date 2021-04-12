# WardleyToGo

[![Go Reference](https://pkg.go.dev/badge/github.com/owulveryck/wardleyToGo.svg)](https://pkg.go.dev/github.com/owulveryck/wardleyToGo)
[![Go](https://github.com/owulveryck/wardleyToGo/actions/workflows/go.yml/badge.svg)](https://github.com/owulveryck/wardleyToGo/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/owulveryck/wardleyToGo/branch/main/graph/badge.svg?token=9BQW1KMGJS)](https://codecov.io/gh/owulveryck/wardleyToGo)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fowulveryck%2FwardleyToGo.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fowulveryck%2FwardleyToGo?ref=badge_shield)

A set of primitives to "code a map". In the context of the package "a map" represents a landscape.
The landscape is made of "Components". Each component knows its own location on a map.
Components can collaborate, meaning that they may be linked together. Therefore a map is also a graph.
The entrypoint of this API is the 'Map' structure

## Demo

Check the online demo at [https://owulveryck.github.io/wardleyToGo/](https://owulveryck.github.io/wardleyToGo/)

## Example

First, create a component type
[embedmd]:# (example_draw_test.go /type dummyComponent.*/ /d.id }/)
```go
type dummyComponent struct {
	id       int64
	position image.Point
}

func (d *dummyComponent) GetPosition() image.Point { return d.position }

func (d *dummyComponent) ID() int64 { return d.id }
```

Then a collaboration structure (an edge)

[embedmd]:# (example_draw_test.go /type dummyCollaboration.*/ /^}$/)
```go
type dummyCollaboration struct{ simple.Edge }

func (d *dummyCollaboration) GetType() wardleyToGo.EdgeType { return 0 }

func newCollaboration(a, b wardleyToGo.Component) wardleyToGo.Collaboration {
	return &dummyCollaboration{Edge: simple.Edge{F: a, T: b}}
}
```

And finally create the map

[embedmd]:# (example_draw_test.go /.*m \:= wardleyToGo.NewMap.*/ /.*c1, c3.*/)
```go
	m := wardleyToGo.NewMap(0)
	c0 := &dummyComponent{id: 0, position: image.Pt(25, 25)}
	c1 := &dummyComponent{id: 1, position: image.Pt(50, 50)}
	c2 := &dummyComponent{id: 2, position: image.Pt(50, 75)}
	c3 := &dummyComponent{id: 3, position: image.Pt(75, 75)}
	m.AddComponent(c0)
	m.AddComponent(c1)
	m.AddComponent(c2)
	m.AddComponent(c3)
	// c0 -> c1
	// c1 -> c2
	// c2 -> c3
	// c1 -> c3
	m.SetCollaboration(newCollaboration(c0, c1))
	m.SetCollaboration(newCollaboration(c1, c2))
	m.SetCollaboration(newCollaboration(c2, c3))
	m.SetCollaboration(newCollaboration(c1, c3))
```
