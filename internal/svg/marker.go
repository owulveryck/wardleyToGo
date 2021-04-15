package svg

import "encoding/xml"

type Marker struct {
	XMLName      xml.Name `xml:"marker"`
	ID           string   `xml:"id,attr"`
	RefX         int      `xml:"refX,attr"`
	RefY         int      `xml:"refY,attr"`
	MarkerWidth  int      `xml:"markerWidth,attr"`
	MarkerHeight int      `xml:"markerHeight,attr"`
	ViewBox      string   `xml:"viewBox,attr,omitempty"`
	Path         *Path    `xml:",omitempty"`
	Polygon      *Polygon `xml:",omitempty"`
}
