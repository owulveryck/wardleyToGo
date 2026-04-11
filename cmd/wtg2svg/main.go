// wtg2svg reads a WTG2 Wardley Map description from stdin and writes
// an SVG rendering to stdout.
//
// Usage:
//
//	cat map.wtg2 | wtg2svg > map.svg
package main

import (
	"flag"
	"image"
	"log"
	"math"
	"os"

	"github.com/owulveryck/wardleyToGo"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

func main() {
	static := flag.Bool("static", false, "produce static SVG without CSS/JS interactivity")
	width := flag.Int("width", 1100, "SVG viewBox width (impacts rendering definition)")
	height := flag.Int("height", 900, "SVG viewBox height (impacts rendering definition)")
	auto := flag.Bool("auto", false, "automatically compute viewBox dimensions based on component count and density")
	flag.Parse()

	p, err := wtg2.NewParser(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	result, err := wtg2.BuildMap(doc)
	if err != nil {
		log.Fatal(err)
	}

	if *auto {
		widthSet := false
		heightSet := false
		flag.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "width":
				widthSet = true
			case "height":
				heightSet = true
			}
		})
		autoW, autoH := autoSize(result.Map)
		if !widthSet {
			*width = autoW
		}
		if !heightSet {
			*height = autoH
		}
	}

	legendWidth := 0
	if result.Legend && len(result.LegendItems) > 0 {
		legendWidth = int(float64(svgmap.LegendWidth) * float64(*width) / 1100.0)
	}

	marginL := int(float64(*width) * 30.0 / 1100.0)
	marginR := int(float64(*width) * 30.0 / 1100.0)
	marginT := int(float64(*height) * 50.0 / 900.0)
	marginB := int(float64(*height) * 65.0 / 900.0)

	e, err := svgmap.NewEncoder(os.Stdout,
		image.Rect(0, 0, *width+legendWidth, *height),
		image.Rect(marginL, marginT, *width-marginR, *height-marginB))
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()

	if *static {
		e.Themes = nil
	}

	if result.Focus != nil {
		focusTheme := &svgmap.FocusTheme{
			ComponentIDs: result.Focus.ComponentIDs,
			EdgeKeys:     result.Focus.EdgeKeys,
			GroupIDs:     result.Focus.GroupIDs,
		}
		if e.Themes == nil {
			e.Themes = []svgmap.Theme{focusTheme}
		} else {
			e.Themes = append(e.Themes, focusTheme)
		}
	}

	var indicators []svgmap.Annotator
	if result.Legend && len(result.LegendItems) > 0 {
		indicators = append(indicators, &svgmap.Legend{Items: result.LegendItems})
	}
	style := svgmap.NewOctoStyle(result.Stages, indicators...)
	style.WithSpace = true
	e.Init(style)

	if err := e.Encode(result.Map); err != nil {
		log.Fatal(err)
	}
}

// autoSize computes viewBox dimensions based on component count and spatial density.
func autoSize(m *wardleyToGo.Map) (int, int) {
	var points []image.Point
	for _, c := range m.Components() {
		if _, isArea := c.(wardleyToGo.Area); isArea {
			continue
		}
		points = append(points, c.GetPosition())
	}

	n := len(points)
	if n == 0 {
		return 1100, 900
	}

	// Heuristic A: count-based (calibrated so N=20 → width=1100)
	countWidth := 246.0 * math.Sqrt(float64(n))

	// Heuristic B: density-based (ensure ≥60px gap between closest components)
	densityWidth := 0.0
	if n >= 2 {
		minDist := math.Inf(1)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				dx := float64(points[i].X - points[j].X)
				dy := float64(points[i].Y - points[j].Y)
				d := math.Sqrt(dx*dx + dy*dy)
				if d > 0 && d < minDist {
					minDist = d
				}
			}
		}
		if !math.IsInf(minDist, 1) {
			densityWidth = 6000.0 / minDist // 100 * 60 / minDist
		}
	}

	w := math.Max(countWidth, densityWidth)
	if w < 600 {
		w = 600
	}
	if w > 3000 {
		w = 3000
	}

	width := int(math.Round(w))
	height := int(math.Round(w * 900.0 / 1100.0))
	return width, height
}
