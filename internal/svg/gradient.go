package svg

import "encoding/xml"

type LinearGradient struct {
	XMLName xml.Name `xml:"linearGradient" `
	ID      string   `xml:"id,attr,omitempty"`
	X1      string   `xml:"x1,attr"`
	Y1      string   `xml:"y1,attr"`
	X2      string   `xml:"x2,attr"`
	Y2      string   `xml:"y2,attr"`
	Stops   []Stop
}

type Stop struct {
	XMLName   xml.Name `xml:"stop"`
	Offset    string   `xml:"offset,attr"`
	StopColor Color    `xml:"stop-color,attr"`
}
