package owm

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("bad graph 1", func(t *testing.T) {
		const src = `
component Cup of Tea [0.79, 0.61] label [19, -4]
Cup of Tea->Cup
`
		p := NewParser(strings.NewReader(src))
		_, err := p.Parse()
		if err == nil {
			t.Fatal(err)
		}
	})
	t.Run("bad graph 2", func(t *testing.T) {
		const src = `
component A [0.79, 0.61] label [19, -4]
A->B
`
		p := NewParser(strings.NewReader(src))
		_, err := p.Parse()
		if err == nil {
			t.Fatal(err)
		}
	})
	t.Run("Tea Shop", func(t *testing.T) {
		const src = `
title Tea Shop
anchor Business [0.95, 0.63]
anchor Public [0.95, 0.78]
component Cup of Tea [0.79, 0.61] label [19, -4]
component Cup [0.73, 0.78] label [19,-4] (dataProduct)
component Tea [0.63, 0.81]
component Hot Water [0.52, 0.80]
component Water [0.38, 0.82]
component Kettle [0.43, 0.35] label [-73, 4] (build)
evolve Kettle 0.62 label [22, 9] (buy)
component Power [0.1, 0.7] label [-29, 30] (outsource)
evolve Power 0.89 label [-12, 21]
Business->Cup of Tea
Public->Cup of Tea
Cup of Tea-collaboration>Cup
Cup of Tea-collaboration>Tea
Cup of Tea-collaboration>Hot Water
Hot Water->Water
Hot Water-facilitating>Kettle 
Kettle-xAsAService>Power
build Kettle


annotation 1 [[0.43,0.49],[0.08,0.79]] Standardising power allows Kettles to evolve faster
annotation 2 [0.48, 0.85] Hot water is obvious and well known
annotations [0.60, 0.02]

note +a generic note appeared [0.16, 0.36]

style wardley
streamAlignedTeam team A [0.84, 0.58, 0.74, 0.68]
enablingTeam team B [0.52, 0.23, 0.32, 0.43]
platformTeam team C [0.18, 0.61, 0.02, 0.94]
complicatedSubsystemTeam team D [0.83, 0.73, 0.45, 0.90]
`
		p := NewParser(strings.NewReader(src))
		_, err := p.Parse()
		if err != nil {
			t.Fatal(err)
		}
	})
}
