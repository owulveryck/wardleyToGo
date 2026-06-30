// wtg2coords reads a WTG2 Wardley Map description from stdin and writes
// a JSON array of components with their coordinates to stdout.
//
// Usage:
//
//	cat map.wtg2 | wtg2coords
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/owulveryck/wardleyToGo/v2"
	"github.com/owulveryck/wardleyToGo/v2/components/wardley"
	"github.com/owulveryck/wardleyToGo/v2/parser/wtg2"
)

type entry struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Type string `json:"type"`
}

func main() {
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

	var entries []entry
	for _, c := range result.Map.Components() {
		if _, ok := c.(*wardley.Group); ok {
			continue
		}
		pos := c.GetPosition()
		entries = append(entries, entry{
			Name: componentName(c),
			X:    pos.X,
			Y:    pos.Y,
			Type: componentType(c),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	out, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}

func componentName(c wardleyToGo.Component) string {
	switch v := c.(type) {
	case *wardley.EvolvedComponent:
		return v.Label
	case *wardley.Component:
		return v.Label
	case *wardley.Anchor:
		return v.Label
	default:
		if s, ok := c.(fmt.Stringer); ok {
			return s.String()
		}
		return fmt.Sprintf("unknown-%d", c.ID())
	}
}

func componentType(c wardleyToGo.Component) string {
	switch v := c.(type) {
	case *wardley.EvolvedComponent:
		return "evolved"
	case *wardley.Anchor:
		return "anchor"
	case *wardley.Component:
		switch v.Type {
		case wardley.BuildComponent:
			return "build"
		case wardley.BuyComponent:
			return "buy"
		case wardley.OutsourceComponent:
			return "outsource"
		case wardley.DataProductComponent:
			return "data-product"
		case wardley.PipelineComponent:
			return "pipeline"
		default:
			return "component"
		}
	default:
		return "unknown"
	}
}
