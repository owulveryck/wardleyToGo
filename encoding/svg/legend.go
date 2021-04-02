package svgmap

import "github.com/owulveryck/wardleyToGo"

func (w *svgMap) writeLegend() {
	w.Group(`font-family="Helvetica,Arial,sans-serif"`, `font-size="13px"`)

	p := &wardleyToGo.StreamAlignedTeam{
		Coords: [4]int{92, 90, 98, 99},
	}
	p.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	x1 := p.Coords[1]*(w.width+150-w.padLeft)/100 + w.padLeft
	y1 := (w.height - w.padLeft) - p.Coords[0]*(w.height-w.padLeft)/100

	w.SVG.Text(x1+5, y1+15, "Stream Aligned")
	w.SVG.Text(x1+5, y1+35, "Team")

	s := &wardleyToGo.PlatformTeam{
		Coords: [4]int{82, 90, 88, 99},
	}
	s.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	x1 = s.Coords[1]*(w.width+150-w.padLeft)/100 + w.padLeft
	y1 = (w.height - w.padLeft) - s.Coords[0]*(w.height-w.padLeft)/100

	w.SVG.Text(x1+5, y1+15, "Platform")
	w.SVG.Text(x1+5, y1+35, "Team")

	c := &wardleyToGo.ComplicatedSubsystemTeam{
		Coords: [4]int{72, 90, 78, 99},
	}
	c.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	x1 = c.Coords[1]*(w.width+150-w.padLeft)/100 + w.padLeft
	y1 = (w.height - w.padLeft) - c.Coords[0]*(w.height-w.padLeft)/100

	w.SVG.Text(x1+11, y1+15, "Complicated")
	w.SVG.Text(x1+11, y1+35, "Subsystem")

	e := &wardleyToGo.EnablingTeam{
		Coords: [4]int{62, 90, 68, 99},
	}
	e.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	x1 = e.Coords[1]*(w.width+150-w.padLeft)/100 + w.padLeft
	y1 = (w.height - w.padLeft) - e.Coords[0]*(w.height-w.padLeft)/100

	w.SVG.Text(x1+11, y1+15, "Enabling")
	w.SVG.Text(x1+11, y1+35, "Team")

	collaborationEdge := wardleyToGo.Edge{
		F:        &dummyElement{[]int{52, 90}},
		T:        &dummyElement{[]int{52, 99}},
		EdgeType: wardleyToGo.CollaborationEdge,
	}
	collaborationEdge.SVGTT(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+5, (w.height-w.padLeft)-52*(w.height-w.padLeft)/100+20, "collaboration")

	facilitatingEdge := wardleyToGo.Edge{
		F:        &dummyElement{[]int{47, 90}},
		T:        &dummyElement{[]int{47, 99}},
		EdgeType: wardleyToGo.FacilitatingEdge,
	}
	facilitatingEdge.SVGTT(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+5, (w.height-w.padLeft)-47*(w.height-w.padLeft)/100+20, "facilitating")

	xAsAServiceEdge := wardleyToGo.Edge{
		F:        &dummyElement{[]int{42, 90}},
		T:        &dummyElement{[]int{42, 99}},
		EdgeType: wardleyToGo.XAsAServiceEdge,
	}
	xAsAServiceEdge.SVGTT(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+5, (w.height-w.padLeft)-42*(w.height-w.padLeft)/100+20, "xAsAService")
	buildComponent := wardleyToGo.Component{
		Type:   wardleyToGo.BuildComponent,
		Coords: [2]int{30, 92},
	}
	buildComponent.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+55, (w.height-w.padLeft)-30*(w.height-w.padLeft)/100+7, "build")
	outsourceComponent := wardleyToGo.Component{
		Type:   wardleyToGo.OutsourceComponent,
		Coords: [2]int{23, 92},
	}
	outsourceComponent.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+55, (w.height-w.padLeft)-23*(w.height-w.padLeft)/100+7, "outsource")
	buyComponent := wardleyToGo.Component{
		Type:   wardleyToGo.BuyComponent,
		Coords: [2]int{16, 92},
	}
	buyComponent.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+55, (w.height-w.padLeft)-16*(w.height-w.padLeft)/100+7, "buy")

	dataProductComponent := wardleyToGo.Component{
		Type:   wardleyToGo.DataProductComponent,
		Coords: [2]int{9, 92},
	}
	dataProductComponent.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+55, (w.height-w.padLeft)-9*(w.height-w.padLeft)/100+7, "dataProduct")
	w.Gend()
}

type dummyElement struct {
	coords []int
}

func (d *dummyElement) ID() int64 {
	return 0
}
func (d *dummyElement) GetCoordinates() []int {
	return d.coords
}
