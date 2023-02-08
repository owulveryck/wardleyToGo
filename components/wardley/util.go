package wardley

import (
	"image"
	"sort"
)

func getBounds(cs []*Component) image.Rectangle {
	csCopy := make([]*Component, len(cs))
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

type csSorter []*Component

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
