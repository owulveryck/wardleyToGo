package wardley

import (
	"bytes"
	"encoding/xml"
	"image"
	"image/color"
	"testing"

	"github.com/owulveryck/wardleyToGo/components"
)

func ExampleComponent_Draw() {
	// draw a circle at the center of the picture
	im := image.NewRGBA(image.Rect(0, 0, 1200, 1000))
	c := &Component{
		Label:          "test",
		LabelPlacement: image.Pt(components.UndefinedCoord, components.UndefinedCoord),
		Placement:      image.Pt(100, 50),
	}
	c.Draw(im, im.Bounds(), im, image.Point{})
}

func TestComponent_MarshalSVG(t *testing.T) {
	c := &Component{
		id:                 0,
		Placement:          image.Point{50, 50},
		Label:              "label",
		LabelPlacement:     image.Point{components.UndefinedCoord, components.UndefinedCoord},
		Type:               PipelineComponent,
		RenderingLayer:     0,
		Configured:         true,
		EvolutionPos:       0,
		Color:              color.Gray{},
		AbsoluteVisibility: 0,
		PipelinedComponents: []*Component{
			{
				id:                  0,
				Placement:           image.Point{20, 50},
				Label:               "p1",
				LabelPlacement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Type:                0,
				RenderingLayer:      0,
				Configured:          true,
				EvolutionPos:        0,
				Color:               color.Gray{},
				AbsoluteVisibility:  0,
				PipelinedComponents: nil,
			},
			{
				id:                  0,
				Placement:           image.Point{70, 50},
				Label:               "p2",
				LabelPlacement:      image.Point{components.UndefinedCoord, components.UndefinedCoord},
				Type:                0,
				RenderingLayer:      0,
				Configured:          true,
				EvolutionPos:        0,
				Color:               color.Gray{},
				AbsoluteVisibility:  0,
				PipelinedComponents: nil,
			},
		},
	}
	var output bytes.Buffer
	e := xml.NewEncoder(&output)
	canvas := image.Rect(0, 0, 100, 100)
	err := c.MarshalSVG(e, canvas)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output.String())
}
