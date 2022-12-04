package main

import (
	"bufio"
	"image"
	"io"
	"regexp"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func initialize(r io.Reader) (*wardleyToGo.Map, error) {
	inventory := make(map[string]*wardley.Component, 0)
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
			c := wardley.NewComponent(int64(len(inventory)))
			c.Label = elements[1]
			c.Placement = image.Pt(50, 50)
			inventory[elements[1]] = c
		}
		if _, ok := inventory[elements[3]]; !ok {
			c := wardley.NewComponent(int64(len(inventory)))
			c.Label = elements[3]
			c.Placement = image.Pt(50, 50)
			inventory[elements[3]] = c
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
	m := wardleyToGo.NewMap(0)
	for _, c := range inventory {
		err := m.AddComponent(c)
		if err != nil {
			return nil, err
		}
	}
	for _, e := range edgeInventory {
		err := m.SetCollaboration(e)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
