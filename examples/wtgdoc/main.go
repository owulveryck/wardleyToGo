package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/wtg"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	var buf bytes.Buffer
	fmt.Fprint(&buf, `
	<!DOCTYPE html>
<html>
<head>
  <title>Page Title</title>
  <style>
    body {
      margin: 0;
      padding: 0;
      display: flex;
      flex-direction: row;
    }

    nav {
      position: fixed;
      top: 50px;
      left: 0;
      width: 200px;
      height: calc(100% - 50px);
      overflow-y: auto;
      z-index: 1;
    }

    main {
      margin: 50px 0 0 200px;
      padding: 20px;
      flex: 1;
      background-color: #f9f9f9;
    }

    nav ul {
      list-style-type: none;
      padding: 0;
    }

    nav ul lvi {
      margin: 0;
      padding: 0;
    }

    nav ul li a {
      display: block;
      padding: 10px;
      color: #333;
      text-decoration: none;
    }

    nav ul li a:hover {
      background-color: #f5f5f5;
    }
.componentText {
  background: transparent;
  color: #fff;
  resize: none;
  border: 0 none;
  width: 100%;
  font-size: 0.8em;
  outline: none;
  height: 100%;
  position: absolute;
}


  </style>
</head>

<body>
`)
	p := wtg.NewParser()

	err := p.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if len(p.InvalidEntries) != 0 {
		for _, err := range p.InvalidEntries {
			log.Println(err)
		}
	}
	fmt.Fprintln(&buf, `<nav>
    <ul>
	<li><a href="#Map">Map</a></li>
	`)
	nodes := p.WMap.Nodes()
	for nodes.Next() {
		if n, ok := nodes.Node().(*wardley.Component); ok {
			fmt.Fprintf(&buf, `<li><a href="#%v">%v</a></li>`+"\n", n.Label, n.Label)
		}
	}
	fmt.Fprintln(&buf, `<nav>
    </ul>
  </nav>
  <main>
	`)

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	if err = md.Convert([]byte(strings.Join(p.Docs, "\n")), &buf); err != nil {
		log.Fatal(err)
	}
	var tmp bytes.Buffer
	fmt.Fprint(&tmp, "## Components\n")
	if err = md.Convert(tmp.Bytes(), &buf); err != nil {
		log.Fatal(err)
	}
	nodes.Reset()
	for nodes.Next() {
		tmp.Reset()
		if n, ok := nodes.Node().(*wardley.Component); ok {
			fmt.Fprintf(&tmp, "### %v\n", n.Label)
			fmt.Fprintln(&tmp, n.Description)
			if err = md.Convert(tmp.Bytes(), &buf); err != nil {
				log.Fatal(err)
			}
			fmt.Fprintln(&buf, `<div style="height: 170px; overflow: hidden;">`)

			m := wardleyToGo.NewMap(0)
			m.AddComponent(&wardley.Component{
				Placement: image.Point{
					X: n.Placement.X,
					Y: 50,
				},
				Label:               n.Label,
				LabelPlacement:      n.LabelPlacement,
				Type:                n.Type,
				RenderingLayer:      n.RenderingLayer,
				Configured:          n.Configured,
				EvolutionPos:        n.EvolutionPos,
				Inertia:             n.Inertia,
				Color:               n.Color,
				AbsoluteVisibility:  n.AbsoluteVisibility,
				Anchor:              n.Anchor,
				PipelinedComponents: n.PipelinedComponents,
				PipelineReference:   n.PipelineReference,
				Description:         n.Description,
			})
			imageSize := image.Rect(0, 0, 1100, 150)
			mapSize := image.Rect(30, 50, 1070, 100)
			e, err := svgmap.NewEncoder(&buf, imageSize, mapSize)
			if err != nil {
				log.Fatal(err)
			}
			style := svgmap.NewOctoStyle(p.EvolutionStages)
			style.WithValueChain = false
			e.Init(style)
			err = e.Encode(m)
			if err != nil {
				log.Fatal(err)
			}
			e.Close()
			fmt.Fprintln(&buf, `</div>`)
		}
	}
	fmt.Fprintln(&buf, "<h2 id=\"map\">The Map<//h2>")

	imgArea := (p.ImageSize.Max.X - p.ImageSize.Min.X) * (p.ImageSize.Max.X - p.ImageSize.Min.Y)
	canvasArea := (p.MapSize.Max.X - p.MapSize.Min.X) * (p.MapSize.Max.X - p.MapSize.Min.Y)
	if imgArea == 0 || canvasArea == 0 {
		p.ImageSize = image.Rect(0, 0, 1500, 1000)
		p.MapSize = image.Rect(30, 50, 1470, 950)
	}
	e, err := svgmap.NewEncoder(&buf, p.ImageSize, p.MapSize)
	if err != nil {
		log.Fatal(err)
	}
	style := svgmap.NewOctoStyle(p.EvolutionStages)
	style.WithControls = true
	e.Init(style)
	err = e.Encode(p.WMap)
	if err != nil {
		log.Fatal(err)
	}
	e.Close()
	fmt.Fprint(&buf, `
  </main>
</body>
</html>
	`)
	fmt.Println(buf.String())
}
