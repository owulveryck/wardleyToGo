package encoding

type Layer interface {
	// GetLayer gives a position of the element on a picture regarding its depth
	GetLayer() int
}

// Background is the lower level of an element
const Background int = 0
