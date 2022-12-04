package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

type node struct {
	id    int64
	label string
}

func (n *node) ID() int64 {
	return n.id
}

type edge struct {
	from       *node
	to         *node
	visibility int
}

// From returns the from node of the edge.
func (e *edge) From() graph.Node {
	return e.from
}

// To returns the to node of the edge.
func (e *edge) To() graph.Node {
	return e.to
}

// ReversedEdge returns the edge reversal of the receiver
// if a reversal is valid for the data type.
// When a reversal is valid an edge of the same type as
// the receiver with nodes of the receiver swapped should
// be returned, otherwise the receiver should be returned
// unaltered.
func (e *edge) ReversedEdge() graph.Edge {
	return &edge{
		from:       e.to,
		to:         e.from,
		visibility: e.visibility,
	}
}

func main() {

	inventory := make(map[string]*node, 0)
	edgeInventory := make([]*edge, 0)
	var link = regexp.MustCompile(`^(.*) (-+) (.*)$`)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("graph G {")
	for scanner.Scan() {
		elements := link.FindStringSubmatch(scanner.Text())
		if len(elements) != 4 {
			log.Fatal("bad entry", scanner.Text())
		}
		fmt.Printf(`"%v" -- "%v" [minlen=%v]`+"\n", elements[1], elements[3], len(elements[2]))
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

	fmt.Println("}")
	g := simple.NewDirectedGraph()
	for _, n := range inventory {
		g.AddNode(n)
	}
	for _, e := range edgeInventory {
		g.SetEdge(e)
	}
	b, err := dot.Marshal(g, "sample", " ", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
