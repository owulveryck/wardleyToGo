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
// CSS holds a pointer to the CSSTheme so that visibility data is computed only once.
type JSTheme struct {
	CSS *CSSTheme
}

func (t *JSTheme) Embed(enc *xml.Encoder, m *wardleyToGo.Map) error {
	var buf bytes.Buffer
	data := generateJsData(m)
	if t.CSS != nil && t.CSS.CachedData != nil {
		data.Visibility = t.CSS.CachedData
	} else {
		data.Visibility = generateCSSData(m)
	}
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
	AllLinks   []string // in the form edge_F_T
	Visibility []cssVisibility
}

func generateJsData(w *wardleyToGo.Map) jsData {
	allCollabs := w.Collaborations()
	allLinks := make([]string, len(allCollabs))
	for i, c := range allCollabs {
		allLinks[i] = "edge_" + strconv.FormatInt(c.From().ID(), 10) + "_" + strconv.FormatInt(c.To().ID(), 10)
	}

	return jsData{
		AllLinks: allLinks,
	}
}

type script struct {
	XMLName xml.Name `xml:"script"`
	Data    string   `xml:",cdata"`
	ID      string   `xml:"id,attr"`
}
