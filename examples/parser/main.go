package main

import (
	"fmt"
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo/parser/owm"
)

func main() {
	p := owm.NewParser(os.Stdin)
	m, err := p.Parse() // the map
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)
}
