package wardley

import (
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components"
)

const (
	// This is a RegularComponent
	RegularComponent wardleyToGo.ComponentType = iota | components.Wardley
	// BuildComponent ...
	BuildComponent
	// Off the shelf element
	BuyComponent
	// OutsourceComponent ...
	OutsourceComponent
	// DataProductComponent ...
	DataProductComponent
	lastComponent
	RegularEdge wardleyToGo.EdgeType = iota + wardleyToGo.EdgeType(lastComponent) | wardleyToGo.EdgeType(components.Wardley)
	EvolvedComponentEdge
	EvolvedEdge
)
