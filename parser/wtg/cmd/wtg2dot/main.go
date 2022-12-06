package main

import (
	"fmt"
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo/parser/wtg"

	"gonum.org/v1/gonum/graph/encoding/dot"
)

func main() {

	m, err := wtg.ParseValueChain(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	b, err := dot.Marshal(m, "sample", "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
