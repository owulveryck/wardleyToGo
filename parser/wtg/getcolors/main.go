package main

import (
	"fmt"

	"github.com/owulveryck/wardleyToGo/parser/wtg"
)

func main() {
	for c, v := range wtg.Colors {
		r, g, b, _ := v.RGBA()
		fmt.Printf(`<p style="color:rgb(%v,%v,%v);display:inline">%v</p>`+"\n", r>>8, g>>8, b>>8, c)
	}
}
