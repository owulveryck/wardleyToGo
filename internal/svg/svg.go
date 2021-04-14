package svg

import "encoding/xml"

type SVG struct {
	Width               string
	Height              string
	Class               string
	PreserveAspectRatio string
	ViewBox             string

	/*
		width="100%"
		height="100%"
		class="wardley-map"
		preserveAspectRatio="xMidYMid meet"
		viewBox="0 0 1050 1050"
		xmlns="http://www.w3.org/2000/svg"
		xmlns:xlink="http://www.w3.org/1999/xlink">
	*/
}

func (s *SVG) StartSVG() xml.StartElement {
	element := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "svg",
		},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "width"},
				Value: s.Width,
			},
			{
				Name:  xml.Name{Local: "height"},
				Value: s.Height,
			},
			{
				Name:  xml.Name{Local: "viewBox"},
				Value: s.ViewBox,
			},
			{
				Name:  xml.Name{Local: "xmlns"},
				Value: "http://www.w3.org/2000/svg",
			},
			{
				Name:  xml.Name{Space: "xmlns", Local: "xlink"},
				Value: "http://www.w3.org/1999/xlink",
			},
			{
				Name:  xml.Name{Local: "preserveAspectRatio"},
				Value: s.PreserveAspectRatio,
			},
		},
	}
	return element
}
