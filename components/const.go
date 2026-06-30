package components

import "github.com/owulveryck/wardleyToGo/v2"

const (
	Undefined wardleyToGo.ComponentType = iota
	Wardley                             = 32 << iota
	TeamTopologies
	UndefinedEdge uint8 = iota

	maxUint        = ^uint(0)
	maxInt         = int(maxUint >> 1)
	UndefinedCoord = -maxInt - 1
)
