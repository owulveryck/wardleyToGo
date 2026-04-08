package svgmap

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"fmt"
	"text/template"

	"github.com/owulveryck/wardleyToGo"
)

//go:embed assets/embeded.css
var embededCSS string

var cssTmpl = template.Must(template.New("CSS").Parse(embededCSS))

// CSSTheme embeds CSS styling for evolution animations, visibility opacity,
// and component text styling.
type CSSTheme struct{}

func (t *CSSTheme) Embed(enc *xml.Encoder, m *wardleyToGo.Map) error {
	var buf bytes.Buffer
	cssData := generateCSSData(m)
	if err := cssTmpl.Execute(&buf, cssData); err != nil {
		return err
	}
	return enc.Encode(style{Data: buf.String()})
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
