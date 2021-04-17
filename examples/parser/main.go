package main

import (
	"bytes"
	"image"
	"log"
	"os"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"

	"github.com/owulveryck/wardleyToGo/parser/owm"
)

// see https://onlinewardleymaps.com/#Jd8sWxt6ch998gQ630
const teashop = `
title Tea Shop
anchor Business [0.95, 0.63]
anchor Public [0.95, 0.78]
component Cup of Tea [0.79, 0.61] label [19, -4]
component Cup [0.73, 0.78]
component Tea [0.63, 0.81]
component Hot Water [0.52, 0.80]
component Water [0.38, 0.82]
component Kettle [0.43, 0.35] label [-57, 4]
evolve Kettle 0.62 label [16, 7]
component Power [0.1, 0.7] label [-27, 20]
evolve Power 0.89 label [-12, 21]
Business->Cup of Tea
Public->Cup of Tea
Cup of Tea->Cup
Cup of Tea->Tea
Cup of Tea->Hot Water
Hot Water->Water
Hot Water->Kettle 
Kettle->Power

annotation 1 [[0.43,0.49],[0.08,0.79]] Standardising power allows Kettles to evolve faster
annotation 2 [0.48, 0.85] Hot water is obvious and well known
annotations [0.60, 0.02]

note +a generic note appeared [0.23, 0.33]

style wardley
`

func main() {
	p := owm.NewParser(bytes.NewBufferString(teashop))
	m, err := p.Parse() // the map
	if err != nil {
		log.Fatal(err)
	}
	e, err := svgmap.NewEncoder(os.Stdout, image.Rect(0, 0, 1100, 900), image.Rect(30, 50, 1070, 850))
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	style := svgmap.NewWardleyStyle(svgmap.DefaultEvolution)
	e.Init(style)
	e.Encode(m)
}
