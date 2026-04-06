package svgmap

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"text/template"

	"github.com/owulveryck/wardleyToGo"
)

//go:embed assets/embeded.js
var embededJS string

//go:embed assets/embeded.css
var embededCSS string

var jsTmpl = template.Must(template.New("JS").Parse(embededJS))
var cssTmpl = template.Must(template.New("CSS").Parse(embededCSS))

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

func generateCSSData(w *wardleyToGo.Map) []cssVisibility {
	maxVisibility := 0
	for _, n := range w.Components() {
		if c, ok := n.(wardleyToGo.Chainer); ok {
			if c.GetAbsoluteVisibility() > maxVisibility {
				maxVisibility = c.GetAbsoluteVisibility()
			}
		}
	}
	step := 0.80 / float64(maxVisibility)
	output := make([]cssVisibility, maxVisibility+1)
	for i := 0; i <= maxVisibility; i++ {
		output[i] = cssVisibility{
			Visibility: fmt.Sprintf("visibility%v", i),
			Opacity:    fmt.Sprintf("%0.2f", 1-float64(i)*step),
		}
	}
	return output
}

type cssVisibility struct {
	Visibility string
	Opacity    string
}

type style struct {
	XMLName xml.Name `xml:"style"`
	Data    string   `xml:",cdata"`
}

type script struct {
	XMLName xml.Name `xml:"script"`
	//Data    string   `xml:",cdata"`
	Data string `xml:",innerxml"`
	ID   string `xml:"id,attr"`
}
