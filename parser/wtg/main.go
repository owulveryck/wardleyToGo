package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

func main() {

	inventory := make(map[string]*node, 0)
	edgeInventory := make([]*edge, 0)
	var link = regexp.MustCompile(`^\s*(.*\S)\s+(-+)\s+(.*)$`)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		elements := link.FindStringSubmatch(scanner.Text())
		if len(elements) != 4 {
			log.Fatal("bad entry", scanner.Text())
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

	g := simple.NewDirectedGraph()
	for _, n := range inventory {
		g.AddNode(n)
	}
	for _, e := range edgeInventory {
		g.SetEdge(e)
	}
	b, err := dot.Marshal(g, "sample", "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
