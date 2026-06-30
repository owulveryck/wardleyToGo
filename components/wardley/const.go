package wardley

import (
	"github.com/owulveryck/wardleyToGo/v2"
	"github.com/owulveryck/wardleyToGo/v2/components"
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
	PipelineComponent
	lastComponent
	RegularEdge wardleyToGo.EdgeType = iota + wardleyToGo.EdgeType(lastComponent) | wardleyToGo.EdgeType(components.Wardley)
	EvolvedComponentEdge
	EvolvedEdge
)
