package main

import (
	"fmt"
	"image"
	"math"
)

func main() {
	for i := 0.0; i < 2*math.Pi; i += 0.01 {
		fmt.Println(image.Pt(int(math.Sin(i)*100), int(math.Cos(i)*100)))
	}
}
