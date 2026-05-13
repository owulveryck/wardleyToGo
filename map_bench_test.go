package wardleyToGo

import (
	"fmt"
	"image"
	"testing"
)

type benchComponent struct {
	id       int64
	position image.Point
}

func (c *benchComponent) ID() int64              { return c.id }
func (c *benchComponent) GetPosition() image.Point { return c.position }

type benchCollab struct {
	from, to Component
}

func (c *benchCollab) From() Component      { return c.from }
func (c *benchCollab) To() Component        { return c.to }
func (c *benchCollab) GetType() EdgeType    { return 0 }

func buildBenchGraph(n int) *Map {
	m := NewMap(0)
	comps := make([]*benchComponent, n)
	for i := 0; i < n; i++ {
		c := &benchComponent{
			id:       int64(i),
			position: image.Pt(5+90*i/(n+1), 5+90*i/(n+1)),
		}
		if c.position.X > 99 {
			c.position.X = 99
		}
		if c.position.Y > 99 {
			c.position.Y = 99
		}
		comps[i] = c
		_ = m.AddComponent(c)
	}
	for i := 0; i < n-1; i++ {
		_ = m.SetCollaboration(&benchCollab{from: comps[i], to: comps[i+1]})
	}
	for i := 0; i < n-3; i += 3 {
		_ = m.SetCollaboration(&benchCollab{from: comps[i], to: comps[i+2]})
	}
	return m
}

func BenchmarkMapComponents(b *testing.B) {
	for _, size := range []int{15, 50, 100} {
		m := buildBenchGraph(size)
		b.Run(fmt.Sprintf("N%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = m.Components()
			}
		})
	}
}

func BenchmarkMapCollaborations(b *testing.B) {
	for _, size := range []int{15, 50, 100} {
		m := buildBenchGraph(size)
		b.Run(fmt.Sprintf("N%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = m.Collaborations()
			}
		})
	}
}

func BenchmarkMapFrom(b *testing.B) {
	for _, size := range []int{15, 50, 100} {
		m := buildBenchGraph(size)
		b.Run(fmt.Sprintf("N%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = m.From(0)
			}
		})
	}
}

func BenchmarkMapAddComponent(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := NewMap(0)
		for j := 0; j < 50; j++ {
			c := &benchComponent{
				id:       int64(j),
				position: image.Pt(5+j, 5+j),
			}
			_ = m.AddComponent(c)
		}
	}
}
