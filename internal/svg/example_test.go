package svg_test

import (
	"encoding/xml"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo/internal/svg"
)

func Example() {
	s := svg.SVG{
		Width:               "100%",
		Height:              "100%",
		PreserveAspectRatio: "xMidYMid meet",
		ViewBox:             "0 0 1050 1050",
	}.StartSVG()
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "    ")
	enc.EncodeToken(s)
	enc.Encode(svg.Rectangle{
		R:    image.Rect(0, 0, 1050, 1050),
		Fill: svg.Color{color.Gray{Y: 128}},
	})
	enc.Encode(svg.Defs{
		Gradient: svg.LinearGradient{
			ID: "wardleyGradient",
			X1: "0%", Y1: "0%", X2: "100%", Y2: "0%",
			Stops: []svg.Stop{
				{
					Offset:    "0%",
					StopColor: svg.Color{color.RGBA{196, 196, 196, 255}},
				},
				{
					Offset:    "30%",
					StopColor: svg.White,
				},
				{
					Offset:    "70%",
					StopColor: svg.White,
				},
				{
					Offset:    "100%",
					StopColor: svg.Color{color.RGBA{196, 196, 196, 255}},
				},
			},
		},
		Markers: []svg.Marker{
			{
				ID:           "arrow",
				RefX:         15,
				RefY:         0,
				MarkerWidth:  12,
				MarkerHeight: 12,
				ViewBox:      "0 -5 10 10",
				Path: &svg.Path{
					D:    "M0,-5L10,0L0,5",
					Fill: svg.Red,
				},
			},
			{
				ID:           "graphArrow",
				RefX:         9,
				RefY:         0,
				MarkerWidth:  12,
				MarkerHeight: 12,
				ViewBox:      "0 -5 10 10",
				Path: &svg.Path{
					D:    "M0,-5L10,0L0,5",
					Fill: svg.Color{color.Black}},
			},
		},
	})
	enc.Encode(svg.Rectangle{
		R:     image.Rect(25, 25, 1000, 1000),
		Style: "fill:url(#wardleyGradient)",
	})

	enc.Encode(svg.Transform{
		Rotate:    270,
		Translate: image.Point{25, 1025},
		Components: []interface{}{
			svg.Line{
				F:           image.Point{0, 0},
				T:           image.Point{1000, 0},
				Stroke:      svg.Color{color.Black},
				StrokeWidth: "1",
				MarkerEnd:   "url(#graphArrow)",
			},
			svg.Line{
				F:               image.Point{0, 173},
				T:               image.Point{1000, 173},
				Stroke:          svg.Gray(0xb8),
				StrokeWidth:     "1",
				StrokeDashArray: []int{2, 2},
			},
			svg.Line{
				F:               image.Point{0, 400},
				T:               image.Point{1000, 400},
				Stroke:          svg.Gray(0xb8),
				StrokeWidth:     "1",
				StrokeDashArray: []int{2, 2},
			},
			svg.Line{
				F:               image.Point{0, 700},
				T:               image.Point{1000, 700},
				Stroke:          svg.Gray(0xb8),
				StrokeWidth:     "1",
				StrokeDashArray: []int{2, 2},
			},
			svg.Text{
				P:          image.Point{5, -10},
				Text:       []byte(`Invisible`),
				TextAnchor: svg.TextAnchorStart,
			},
			svg.Text{
				P:          image.Point{995, -10},
				Text:       []byte(`Visible`),
				TextAnchor: svg.TextAnchorEnd,
			},
			svg.Text{
				P:          image.Point{500, -10},
				Text:       []byte(`Value Chain`),
				TextAnchor: svg.TextAnchorMiddle,
				FontWeight: "bold",
			},
		},
	})
	enc.Encode(svg.Line{
		F:         image.Point{25, 1025},
		T:         image.Point{1025, 1025},
		Stroke:    svg.Color{color.Black},
		MarkerEnd: "url(#graphArrrow)",
	})
	enc.Encode(svg.Text{
		P:          image.Point{32, 40},
		FontWeight: "bold",
		FontSize:   "11px",
		Text:       []byte(`Uncharted`),
		TextAnchor: svg.TextAnchorStart,
	})
	enc.Encode(svg.Text{
		P:          image.Point{1020, 40},
		FontWeight: "bold",
		FontSize:   "11px",
		Text:       []byte(`Industrialised`),
		TextAnchor: svg.TextAnchorEnd,
	})
	enc.Encode(svg.Text{
		P:    image.Point{25, 1040},
		Text: []byte(`Genesis`),
	})
	enc.Encode(svg.Text{
		P:    image.Point{198, 1040},
		Text: []byte(`Custom-Built`),
	})
	enc.Encode(svg.Text{
		P:    image.Point{425, 1040},
		Text: []byte(`Product (+rental)`),
	})
	enc.Encode(svg.Text{
		P:    image.Point{725, 1040},
		Text: []byte(`Commodity (+utility)`),
	})
	enc.Encode(svg.Text{
		P:          image.Point{1025, 1040},
		TextAnchor: svg.TextAnchorEnd,
		FontWeight: "bold",
		Text:       []byte(`Evolution`),
	})
	err := enc.EncodeToken(s.End())
	if err != nil {
		log.Fatal(err)
	}
	enc.Flush()
	//	<svg width="100%" height="100%" viewBox="0 0 1050 1050" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="xMidYMid meet">
	//    <rect x="0" y="0" width="1050" height="1050" fill="rgb(128,128,128)" fill-opacity="1.0"></rect>
	//    <defs>
	//        <linearGradient id="wardleyGradient" x1="0%" y1="0%" x2="100%" y2="0%">
	//            <stop offset="0%" stop-color="rgb(196,196,196)"></stop>
	//            <stop offset="30%" stop-color="rgb(255,255,255)"></stop>
	//            <stop offset="70%" stop-color="rgb(255,255,255)"></stop>
	//            <stop offset="100%" stop-color="rgb(196,196,196)"></stop>
	//        </linearGradient>
	//        <marker id="arrow" refX="15" refY="0" markerWidth="12" markerHeight="12" viewBox="0 -5 10 10">
	//            <path d="M0,-5L10,0L0,5" fill="rgb(255,0,0)"></path>
	//        </marker>
	//        <marker id="graphArrow" refX="9" refY="0" markerWidth="12" markerHeight="12" viewBox="0 -5 10 10">
	//            <path d="M0,-5L10,0L0,5" fill="rgb(0,0,0)"></path>
	//        </marker>
	//    </defs>
	//    <rect x="25" y="25" width="975" height="975" style="fill:url(#wardleyGradient)"></rect>
	//    <g transform=" translate(25,1025) rotate(270)">
	//        <line x1="0" y1="0" x2="1000" y2="0" stroke-width="1" marker-end="url(#graphArrow)" stroke="rgb(0,0,0)" stroke-opacity="1.0"></line>
	//        <line x1="0" y1="173" x2="1000" y2="173" stroke-width="1" stroke-dasharray="2 2" stroke="rgb(184,184,184)" stroke-opacity="1.0"></line>
	//        <line x1="0" y1="400" x2="1000" y2="400" stroke-width="1" stroke-dasharray="2 2" stroke="rgb(184,184,184)" stroke-opacity="1.0"></line>
	//        <line x1="0" y1="700" x2="1000" y2="700" stroke-width="1" stroke-dasharray="2 2" stroke="rgb(184,184,184)" stroke-opacity="1.0"></line>
	//        <text text-anchor="start">
	//            <tspan x="5" dy="-10">Invisible</tspan>
	//        </text>
	//        <text text-anchor="end">
	//            <tspan x="995" dy="-10">Visible</tspan>
	//        </text>
	//        <text font-weight="bold" text-anchor="middle">
	//            <tspan x="500" dy="-10">Value Chain</tspan>
	//        </text>
	//    </g>
	//    <line x1="25" y1="1025" x2="1025" y2="1025" marker-end="url(#graphArrrow)" stroke="rgb(0,0,0)" stroke-opacity="1.0"></line>
	//    <text font-weight="bold" font-size="11px" text-anchor="start">
	//        <tspan x="32" dy="40">Uncharted</tspan>
	//    </text>
	//    <text font-weight="bold" font-size="11px" text-anchor="end">
	//        <tspan x="1020" dy="40">Industrialised</tspan>
	//    </text>
	//    <text>
	//        <tspan x="25" dy="1040">Genesis</tspan>
	//    </text>
	//    <text>
	//        <tspan x="198" dy="1040">Custom-Built</tspan>
	//    </text>
	//    <text>
	//        <tspan x="425" dy="1040">Product (+rental)</tspan>
	//    </text>
	//    <text>
	//        <tspan x="725" dy="1040">Commodity (+utility)</tspan>
	//    </text>
	//    <text font-weight="bold" text-anchor="end">
	//        <tspan x="1025" dy="1040">Evolution</tspan>
	//    </text>
	// </svg>
}
