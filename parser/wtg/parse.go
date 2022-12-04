package main

import (
	"bufio"
	"io"
	"regexp"

	"github.com/owulveryck/wardleyToGo"
	"gonum.org/v1/gonum/graph"
)

func parse(r io.Reader) (graph.Directed, error) {
	inventory := make(map[string]*node, 0)
	edgeInventory := make([]*edge, 0)
	var link = regexp.MustCompile(`^\s*(.*\S)\s+(-+)\s+(.*)$`)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		elements := link.FindStringSubmatch(scanner.Text())
		if len(elements) != 4 {
			// log.Fatal("bad entry", scanner.Text())
			continue
		}
		if _, ok := inventory[elements[1]]; !ok {
			inventory[elements[1]] = &node{
				id:    int64(len(inventory)),
				label: elements[1],
			}
		}
		if _, ok := inventory[elements[3]]; !ok {
			inventory[elements[3]] = &node{
				id:    int64(len(inventory)),
				label: elements[3],
			}
		}
		edgeInventory = append(edgeInventory, &edge{
			from:       inventory[elements[1]],
			to:         inventory[elements[3]],
			visibility: len(elements[2]),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	g := wardleyToGo.NewMap(0)
	for _, n := range inventory {
		g.AddNode(n)
	}
	for _, e := range edgeInventory {
		g.SetEdge(e)
	}
	return g, nil
}
