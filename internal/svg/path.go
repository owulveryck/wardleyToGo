package svg

import "encoding/xml"

type Path struct {
	XMLName xml.Name `xml:"path"`
	D       string   `xml:"d,attr"`
	Fill    Color    `xml:"fill,attr"`
}
