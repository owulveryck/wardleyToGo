package tt

import (
	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components"
)

const (
	CollaborationEdge wardleyToGo.EdgeType = iota | wardleyToGo.EdgeType(components.TeamTopologies)
	FacilitatingEdge
	XAsAServiceEdge
)
