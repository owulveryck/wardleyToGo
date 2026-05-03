package svg

import "encoding/xml"

type Defs struct {
	XMLName  xml.Name `xml:"defs"`
	Gradients []LinearGradient
	Markers  []Marker
}
