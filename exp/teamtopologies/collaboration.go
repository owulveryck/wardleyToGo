package tt

import (
	"github.com/owulveryck/wardleyToGo/v2"
	"github.com/owulveryck/wardleyToGo/v2/components"
)

const (
	CollaborationEdge wardleyToGo.EdgeType = iota | wardleyToGo.EdgeType(components.TeamTopologies)
	FacilitatingEdge
	XAsAServiceEdge
)
