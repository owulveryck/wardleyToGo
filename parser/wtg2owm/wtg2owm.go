package wtg2owm

import (
	"fmt"
	"image"
	"io"
	"sort"

	"github.com/owulveryck/wardleyToGo/components/wardley"
	"github.com/owulveryck/wardleyToGo/parser/wtg"
)

func Convert(r io.Reader, w io.Writer) error {
	p := wtg.NewParser()
	err := p.Parse(r)
	if err != nil {
		return err
	}
	dict := make(map[int64]string, 0)
	fmt.Fprintf(w, "title %v\n\n", p.WMap.Title)
	allComponents := p.WMap.Nodes()
	for allComponents.Next() {
		if n, ok := allComponents.Node().(*wardley.Component); ok {
			switch n.Type {
			case wardley.BuildComponent:
				fmt.Fprintf(w, "component %v [%v, %v] (build)\n", n.Label, float64(100-n.Placement.Y)/100, float64(n.Placement.X)/100)
				//				fmt.Fprintf(w,"build %v\n", n.Label)
			case wardley.BuyComponent:
				fmt.Fprintf(w, "component %v [%v, %v] (buy)\n", n.Label, float64(100-n.Placement.Y)/100, float64(n.Placement.X)/100)
				//				fmt.Fprintf(w,"buy %v\n", n.Label)
			case wardley.OutsourceComponent:
				fmt.Fprintf(w, "component %v [%v, %v] (outsource)\n", n.Label, float64(100-n.Placement.Y)/100, float64(n.Placement.X)/100)
				//				fmt.Fprintf(w,"outsource %v\n", n.Label)
			case wardley.PipelineComponent:
				fmt.Fprintf(w, "component %v [%v, %v]\n", n.Label, ((float64(100-n.Placement.Y) + 1.3) / 100), float64(n.Placement.X)/100)
				if n.PipelinedComponents != nil && len(n.PipelinedComponents) > 1 {
					rect := getBounds(n.PipelinedComponents)
					fmt.Fprintf(w, "pipeline %v [%v, %v]\n", n.Label, (float64(rect.Max.X)-0.5)/100, (float64(rect.Min.X)+0.5)/100)
					_ = rect
				} else {
					fmt.Fprintf(w, "pipeline %v\n", n.Label)
				}
			default:
				fmt.Fprintf(w, "component %v [%v, %v]\n", n.Label, float64(100-n.Placement.Y)/100, float64(n.Placement.X)/100)
			}

			dict[n.ID()] = n.Label

		}
		if n, ok := allComponents.Node().(*wardley.EvolvedComponent); ok {
			fmt.Fprintf(w, "evolve %v %v\n", n.Label, float64(n.Placement.X)/100)
			dict[n.ID()] = n.Label
		}
	}
	allEdges := p.WMap.Edges()
	for allEdges.Next() {
		if e, ok := allEdges.Edge().(*wardley.Collaboration); ok {
			if dict[e.F.ID()] == dict[e.T.ID()] {
				continue
			}
			fmt.Fprintf(w, "%v -> %v\n", dict[e.F.ID()], dict[e.T.ID()])
		}
	}
	fmt.Fprintf(w, "evolution %v -> %v -> %v -> %v", p.EvolutionStages[0].Label, p.EvolutionStages[1].Label, p.EvolutionStages[2].Label, p.EvolutionStages[3].Label)
	fmt.Fprintln(w, "\nstyle wardley")
	return nil
}

func getBounds(cs []*wardley.Component) image.Rectangle {
	csCopy := make([]*wardley.Component, len(cs))
	i := 0
	for _, c := range cs {
		csCopy[i] = c
		i++
	}
	sort.Sort(csSorter(csCopy))
	return image.Rectangle{
		Min: image.Point{
			X: csCopy[0].GetPosition().X,
			Y: csCopy[0].GetPosition().Y,
		},
		Max: image.Point{
			X: csCopy[len(csCopy)-1].GetPosition().X,
			Y: csCopy[len(csCopy)-1].GetPosition().Y,
		},
	}
}

type csSorter []*wardley.Component

// Len is the number of elements in the collection.
func (cs csSorter) Len() int {
	return len(cs)
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//   - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//   - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (cs csSorter) Less(i int, j int) bool {
	return cs[i].GetPosition().X < cs[j].GetPosition().X
}

// Swap swaps the elements with indexes i and j.
func (cs csSorter) Swap(i int, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}