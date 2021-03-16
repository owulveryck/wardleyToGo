package main

import (
	"os"
)

func main() {
	width := 1000
	height := 800
	padLeft := 25
	padBottom := 30

	topLeft := &component{
		id:         0,
		coords:     [2]int{10, 10},
		label:      "TL",
		labelCoord: [2]int{-4, 19},
	}

	bottomRight := &component{
		id:     1,
		coords: [2]int{90, 90},
		label:  "BR",
	}

	topRight := &component{
		id:     2,
		coords: [2]int{10, 90},
		label:  "TR",
	}
	bottomLeft := &component{
		id:     2,
		coords: [2]int{90, 10},
		label:  "TR",
	}

	w := newWardley(os.Stdout)
	w.Init(width, height, padLeft, padBottom)
	w.writeElement(topLeft)
	w.writeElement(bottomLeft)
	w.writeElement(topRight)
	w.writeElement(bottomRight)
	w.Close()

}
