package main

import (
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo/parser/wtg2owm"
)

func main() {
	err := wtg2owm.Convert(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
