package svgmap

import (
	"encoding/xml"

	"github.com/owulveryck/wardleyToGo/v2"
)

// FocusTheme reduces opacity of unfocused elements when focus is active.
// It implements Theme, ComponentDecorator, and CollaborationDecorator.
type FocusTheme struct {
	ComponentIDs map[int64]bool
	EdgeKeys     map[[2]int64]bool
	GroupIDs     map[int64]bool
}

func (t *FocusTheme) Embed(enc *xml.Encoder, _ *wardleyToGo.Map) error {
	return enc.Encode(style{Data: `.unfocused { opacity: 0.15 !important; }`})
}

func (t *FocusTheme) DecorateComponent(c wardleyToGo.Component) []xml.Attr {
	if t.ComponentIDs[c.ID()] || t.GroupIDs[c.ID()] {
		return nil
	}
	return []xml.Attr{{Name: xml.Name{Local: "class"}, Value: "unfocused"}}
}

func (t *FocusTheme) DecorateCollaboration(c wardleyToGo.Collaboration) []xml.Attr {
	key := [2]int64{c.From().ID(), c.To().ID()}
	if t.EdgeKeys[key] {
		return nil
	}
	return []xml.Attr{{Name: xml.Name{Local: "class"}, Value: "unfocused"}}
}
