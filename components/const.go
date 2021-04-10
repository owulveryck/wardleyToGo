package components

const (
	Undefined ComponentType = iota
	Wardley                 = 32 << iota
	TeamTopologies
	UnedfinedEdge EdgeType = iota

	maxUint        = ^uint(0)
	maxInt         = int(maxUint >> 1)
	UndefinedCoord = -maxInt - 1
)
