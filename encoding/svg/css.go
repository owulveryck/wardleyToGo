package svgmap

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"strconv"
	"text/template"

	"github.com/owulveryck/wardleyToGo"
)

//go:embed assets/embeded.css
var embededCSS string

var cssTmpl = template.Must(template.New("CSS").Parse(embededCSS))

// CSSTheme embeds CSS styling for evolution animations, visibility opacity,
// and component text styling. After Embed is called, CachedData holds the
// computed visibility data so that JSTheme can reuse it without recomputation.
type CSSTheme struct {
	CachedData []cssVisibility
}

func (t *CSSTheme) Embed(enc *xml.Encoder, m *wardleyToGo.Map) error {
	var buf bytes.Buffer
	t.CachedData = generateCSSData(m)
	if err := cssTmpl.Execute(&buf, t.CachedData); err != nil {
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
			Visibility: "visibility" + strconv.Itoa(i),
			Opacity:    strconv.FormatFloat(1-float64(i)*step, 'f', 2, 64),
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
