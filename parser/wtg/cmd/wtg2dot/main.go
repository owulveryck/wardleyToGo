package main

import (
	"fmt"
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo/parser/wtg"

	"gonum.org/v1/gonum/graph/encoding/dot"
)

func main() {
	p := wtg.NewParser()
	err := p.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	b, err := dot.Marshal(p.WMap, "sample", "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
