package svgmap

import "github.com/owulveryck/wardleyToGo/encoding"

type elements []SVGDrawer

// Len is the number of elements in the collection.
func (e elements) Len() int {
	return len(e)
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
//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (e elements) Less(i int, j int) bool {
	iLayout := encoding.Background
	jLayout := encoding.Background
	if e, ok := e[i].(encoding.Layer); ok {
		iLayout = e.GetLayer()
	}
	if e, ok := e[j].(encoding.Layer); ok {
		jLayout = e.GetLayer()
	}
	return iLayout < jLayout
}

// Swap swaps the elements with indexes i and j.
func (e elements) Swap(i int, j int) {
	e[i], e[j] = e[j], e[i]
}
