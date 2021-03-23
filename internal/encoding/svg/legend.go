package svgmap

import "github.com/owulveryck/wardleyToGo/internal/plan"

func (w *svgMap) writeLegend() {
	w.Group(`font-family="&quot;Helvetica Neue&quot;,Helvetica,Arial,sans-serif"`, `font-size="13px"`)

	p := &plan.StreamAlignedTeam{
		Coords: [4]int{92, 90, 98, 99},
	}
	p.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	x1 := p.Coords[1]*(w.width+150-w.padLeft)/100 + w.padLeft
	y1 := (w.height - w.padLeft) - p.Coords[0]*(w.height-w.padLeft)/100

	w.SVG.Text(x1+5, y1+15, "Stream Aligned")
	w.SVG.Text(x1+5, y1+35, "Team")

	s := &plan.PlatformTeam{
		Coords: [4]int{82, 90, 88, 99},
	}
	s.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	x1 = s.Coords[1]*(w.width+150-w.padLeft)/100 + w.padLeft
	y1 = (w.height - w.padLeft) - s.Coords[0]*(w.height-w.padLeft)/100

	w.SVG.Text(x1+5, y1+15, "Platform")
	w.SVG.Text(x1+5, y1+35, "Team")

	c := &plan.ComplicatedSubsystemTeam{
		Coords: [4]int{72, 90, 78, 99},
	}
	c.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	x1 = c.Coords[1]*(w.width+150-w.padLeft)/100 + w.padLeft
	y1 = (w.height - w.padLeft) - c.Coords[0]*(w.height-w.padLeft)/100

	w.SVG.Text(x1+11, y1+15, "Complicated")
	w.SVG.Text(x1+11, y1+35, "Subsystem")

	e := &plan.EnablingTeam{
		Coords: [4]int{62, 90, 68, 99},
	}
	e.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	x1 = e.Coords[1]*(w.width+150-w.padLeft)/100 + w.padLeft
	y1 = (w.height - w.padLeft) - e.Coords[0]*(w.height-w.padLeft)/100

	w.SVG.Text(x1+11, y1+15, "Enabling")
	w.SVG.Text(x1+11, y1+35, "Team")

	collaborationEdge := plan.Edge{
		F:        &dummyElement{[]int{52, 90}},
		T:        &dummyElement{[]int{52, 99}},
		EdgeType: plan.CollaborationEdge,
	}
	collaborationEdge.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+5, (w.height-w.padLeft)-52*(w.height-w.padLeft)/100+20, "collaboration")

	facilitatingEdge := plan.Edge{
		F:        &dummyElement{[]int{47, 90}},
		T:        &dummyElement{[]int{47, 99}},
		EdgeType: plan.FacilitatingEdge,
	}
	facilitatingEdge.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+5, (w.height-w.padLeft)-47*(w.height-w.padLeft)/100+20, "facilitating")

	xAsAServiceEdge := plan.Edge{
		F:        &dummyElement{[]int{42, 90}},
		T:        &dummyElement{[]int{42, 99}},
		EdgeType: plan.XAsAServiceEdge,
	}
	xAsAServiceEdge.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+5, (w.height-w.padLeft)-42*(w.height-w.padLeft)/100+20, "xAsAService")
	buildComponent := plan.Component{
		Type:   plan.BuildComponent,
		Coords: [2]int{30, 92},
	}
	buildComponent.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+55, (w.height-w.padLeft)-30*(w.height-w.padLeft)/100+7, "build")
	outsourceComponent := plan.Component{
		Type:   plan.OutsourceComponent,
		Coords: [2]int{23, 92},
	}
	outsourceComponent.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+55, (w.height-w.padLeft)-23*(w.height-w.padLeft)/100+7, "outsource")
	buyComponent := plan.Component{
		Type:   plan.BuyComponent,
		Coords: [2]int{16, 92},
	}
	buyComponent.SVG(w.SVG, w.width+150, w.height, w.padLeft, w.padBottom)
	w.SVG.Text(90*(w.width+150-w.padLeft)/100+w.padLeft+55, (w.height-w.padLeft)-16*(w.height-w.padLeft)/100+7, "buy")

	dataProductComponent := plan.Component{
		Type:   plan.DataProductComponent,
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
