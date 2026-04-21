package svgmap

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/owulveryck/wardleyToGo"
)

// AnimationData holds precomputed metadata for step-by-step SVG animation.
// It is populated by the parser/builder and consumed by AnimationTheme.
type AnimationData struct {
	Depth    map[int64]int     // component ID -> topological rank (0 = anchor)
	Type     map[int64]string  // component ID -> "anchor", "component", "pipeline", "evolved", "group"
	ParentID map[int64]int64   // component ID -> ID of nearest parent (lowest depth), -1 if root
	YRank    map[int64]int     // component ID -> Y-position bucket (alternative grouping mode)
	Members  map[int64][]int64 // group ID -> member component IDs
}

// AnimationTheme enriches SVG output with data-* attributes for step-by-step
// animation. It implements Theme, ComponentDecorator, and CollaborationDecorator.
type AnimationTheme struct {
	Data *AnimationData
}

func (t *AnimationTheme) Embed(_ *xml.Encoder, _ *wardleyToGo.Map) error {
	return nil
}

func (t *AnimationTheme) DecorateComponent(c wardleyToGo.Component) []xml.Attr {
	id := c.ID()
	var attrs []xml.Attr

	if depth, ok := t.Data.Depth[id]; ok {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "data-depth"}, Value: strconv.Itoa(depth)})
	}
	if typ, ok := t.Data.Type[id]; ok {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "data-type"}, Value: typ})
	}
	if parentID, ok := t.Data.ParentID[id]; ok && parentID >= 0 {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "data-parent-id"}, Value: "element_" + strconv.FormatInt(parentID, 10)})
	}
	if yRank, ok := t.Data.YRank[id]; ok {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "data-y-rank"}, Value: strconv.Itoa(yRank)})
	}
	if members, ok := t.Data.Members[id]; ok && len(members) > 0 {
		ids := make([]string, len(members))
		for i, mid := range members {
			ids[i] = "element_" + strconv.FormatInt(mid, 10)
		}
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: "data-members"}, Value: strings.Join(ids, ",")})
	}

	return attrs
}

func (t *AnimationTheme) DecorateCollaboration(c wardleyToGo.Collaboration) []xml.Attr {
	toID := c.To().ID()
	return []xml.Attr{
		{Name: xml.Name{Local: "data-child-id"}, Value: "element_" + strconv.FormatInt(toID, 10)},
	}
}
