package wtg

type nodeSorter []*node

// Len is the number of elements in the collection.
func (n nodeSorter) Len() int {
	return len(n)
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
func (n nodeSorter) Less(i int, j int) bool {
	if n[i].visibility == n[j].visibility {
		return n[i].c.Label < n[j].c.Label
	}
	return n[i].visibility < n[j].visibility
}

// Swap swaps the elements with indexes i and j.
func (n nodeSorter) Swap(i int, j int) {
	n[i], n[j] = n[j], n[i]
}
