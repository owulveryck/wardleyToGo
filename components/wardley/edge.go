package wardley

import (
	"github.com/owulveryck/wardleyToGo/components"
)

const (
	RegularEdge components.EdgeType = iota | components.EdgeType(components.Wardley)
	EvolvedComponentEdge
	EvolvedEdge
)
