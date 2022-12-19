package main

import (
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo/parser/wtg"
)

func main() {
	p := wtg.NewParser()
	err := p.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	p.DumpComponents(os.Stdout)
}
