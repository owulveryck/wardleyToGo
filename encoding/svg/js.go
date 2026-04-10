package svgmap

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"strconv"
	"text/template"

	"github.com/owulveryck/wardleyToGo"
)

//go:embed assets/embeded.js
var embededJS string

var jsTmpl = template.Must(template.New("JS").Parse(embededJS))

// JSTheme embeds JavaScript for interactive toggling of links and visibility.
// It also implements ComponentDecorator to add onclick handlers to components.
type JSTheme struct{}

func (t *JSTheme) Embed(enc *xml.Encoder, m *wardleyToGo.Map) error {
	var buf bytes.Buffer
	data := generateJsData(m)
	data.Visibility = generateCSSData(m)
	if err := jsTmpl.Execute(&buf, data); err != nil {
		return err
	}
	return enc.Encode(script{Data: buf.String(), ID: "SVGScript"})
}

// DecorateComponent adds an onclick handler that toggles linked edges.
func (t *JSTheme) DecorateComponent(c wardleyToGo.Component) []xml.Attr {
	return []xml.Attr{
		{
			Name:  xml.Name{Local: "onclick"},
			Value: "toggleLink(this.id)",
		},
	}
}

type jsData struct {
	AllLinks   []string            // in the form edge_F_T
	G          map[string][]string // 'F': {'edge_F_T', 'edge_T_T2'}
	Visibility []cssVisibility
}

func generateJsData(w *wardleyToGo.Map) jsData {
	allCollabs := w.Collaborations()
	allLinks := make([]string, len(allCollabs))
	for i, c := range allCollabs {
		allLinks[i] = "edge_" + strconv.FormatInt(c.From().ID(), 10) + "_" + strconv.FormatInt(c.To().ID(), 10)
	}

	// Build adjacency list once to avoid repeated w.From() calls
	// (each call allocates and sorts).
	components := w.Components()
	adj := make(map[int64][]int64, len(components))
	for _, n := range components {
		succs := w.From(n.ID())
		if len(succs) > 0 {
			ids := make([]int64, len(succs))
			for i, s := range succs {
				ids[i] = s.ID()
			}
			adj[n.ID()] = ids
		}
	}

	// Memoized DFS: compute reachable edges per node in O(V+E) total.
	memo := make(map[int64][]string, len(adj))
	var dfs func(id int64) []string
	dfs = func(id int64) []string {
		if cached, ok := memo[id]; ok {
			return cached
		}
		succs := adj[id]
		edges := make([]string, 0, len(succs))
		for _, succID := range succs {
			edges = append(edges, "edge_"+strconv.FormatInt(id, 10)+"_"+strconv.FormatInt(succID, 10))
			edges = append(edges, dfs(succID)...)
		}
		memo[id] = edges
		return edges
	}

	paths := make(map[string][]string, len(adj))
	for _, n := range components {
		if len(adj[n.ID()]) == 0 {
			continue
		}
		element := "element_" + strconv.FormatInt(n.ID(), 10)
		paths[element] = dfs(n.ID())
	}

	return jsData{
		AllLinks: allLinks,
		G:        paths,
	}
}

type script struct {
	XMLName xml.Name `xml:"script"`
	Data    string   `xml:",innerxml"`
	ID      string   `xml:"id,attr"`
}
