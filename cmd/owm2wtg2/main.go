// owm2wtg2 reads an OnlineWardleyMaps (OWM) map description from stdin
// and writes a WTG2 format equivalent to stdout.
//
// Usage:
//
//	cat map.owm | owm2wtg2 > map.wtg2
package main

import (
	"log"
	"os"

	"github.com/owulveryck/wardleyToGo/exp/owm2wtg2"
)

func main() {
	doc, err := owm2wtg2.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	if err := owm2wtg2.Emit(doc, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
