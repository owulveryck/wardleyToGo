package svgmap

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"fmt"
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
		allLinks[i] = fmt.Sprintf("edge_%v_%v", c.From().ID(), c.To().ID())
	}

	paths := make(map[string][]string)
	for _, n := range w.Components() {
		successors := w.From(n.ID())
		if len(successors) == 0 {
			continue
		}
		element := fmt.Sprintf("element_%v", n.ID())
		paths[element] = make([]string, 0)
		// DFS to collect all edges reachable from this node
		visited := make(map[int64]bool)
		var dfs func(id int64)
		dfs = func(id int64) {
			if visited[id] {
				return
			}
			visited[id] = true
			for _, succ := range w.From(id) {
				paths[element] = append(paths[element], fmt.Sprintf("edge_%v_%v", id, succ.ID()))
				dfs(succ.ID())
			}
		}
		dfs(n.ID())
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
