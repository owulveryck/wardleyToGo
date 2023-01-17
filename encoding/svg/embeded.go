package svgmap

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"text/template"

	"github.com/owulveryck/wardleyToGo"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

//go:embed assets/embeded.js
var embededJS string

//go:embed assets/embeded.css
var embededCSS string

var jsTmpl *template.Template
var cssTmpl *template.Template

func init() {
	var err error
	jsTmpl, err = template.New("JS").Parse(embededJS)
	if err != nil {
		panic(err)
	}
	cssTmpl, err = template.New("CSS").Parse(embededCSS)
	if err != nil {
		panic(err)
	}
}

type jsData struct {
	AllLinks   []string            // in the form edge_F_T
	G          map[string][]string // 'F': {'edge_F_T', 'edge_T_T2'}
	Visibility []cssVisibility
}

func generateJsData(w *wardleyToGo.Map) jsData {
	allEdges := w.Edges()
	allLinks := make([]string, allEdges.Len())
	for i := 0; allEdges.Next(); i++ {
		allLinks[i] = fmt.Sprintf("edge_%v_%v", allEdges.Edge().From().ID(), allEdges.Edge().To().ID())
	}
	allNodes := w.Nodes()
	paths := make(map[string][]string, allNodes.Len())
	for allNodes.Next() {
		currentNode := allNodes.Node()
		if w.From(currentNode.ID()).Len() == 0 {
			continue
		}
		element := fmt.Sprintf("element_%v", currentNode.ID())
		paths[element] = make([]string, 0)
		df := &traverse.DepthFirst{
			Visit: func(n graph.Node) {
				t := w.From(n.ID())
				for t.Next() {
					paths[element] = append(paths[element], fmt.Sprintf("edge_%v_%v", n.ID(), t.Node().ID()))
				}
			},
		}
		df.Walk(w, currentNode, nil)

	}
	return jsData{
		AllLinks: allLinks,
		G:        paths,
	}
}

func generateCSSData(w *wardleyToGo.Map) []cssVisibility {
	maxVisibility := 0
	nodes := w.Nodes()
	for nodes.Next() {
		if n, ok := nodes.Node().(wardleyToGo.Chainer); ok {
			if n.GetAbsoluteVisibility() > maxVisibility {
				maxVisibility = n.GetAbsoluteVisibility()
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
