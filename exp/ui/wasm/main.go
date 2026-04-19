//go:build js && wasm

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"syscall/js"

	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

func main() {
	js.Global().Set("generateSVG", js.FuncOf(generate))
	js.Global().Set("parseWTG2ToState", js.FuncOf(parseToState))
	<-make(chan bool)
}

func generate(_ js.Value, args []js.Value) any {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	if len(args) < 1 {
		return "error: no input provided"
	}

	input := args[0].String()
	static := false
	if len(args) >= 2 {
		static = args[1].Bool()
	}
	// Third argument: resolution scale percentage (100 = default)
	scalePct := 100
	if len(args) >= 3 {
		scalePct = args[2].Int()
		if scalePct < 25 {
			scalePct = 25
		}
		if scalePct > 400 {
			scalePct = 400
		}
	}
	// Fourth and fifth arguments: base width and height (before scaling)
	baseW := 1200
	baseH := 900
	if len(args) >= 5 {
		baseW = args[3].Int()
		baseH = args[4].Int()
	}

	p, err := wtg2.NewParser(bytes.NewBufferString(input))
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	doc, err := p.Parse()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}

	result, err := wtg2.BuildMap(doc)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}

	// Scale viewBox and canvas proportionally
	boxW := baseW * scalePct / 100
	boxH := baseH * scalePct / 100
	marginL := 30 * scalePct / 100
	marginT := 50 * scalePct / 100
	marginR := 30 * scalePct / 100
	marginB := 65 * scalePct / 100

	legendWidth := 0
	if result.Legend && len(result.LegendItems) > 0 {
		legendWidth = svgmap.LegendWidth * scalePct / 100
	}

	svgBuf.Reset()
	e, err := svgmap.NewEncoder(&svgBuf,
		image.Rect(0, 0, boxW+legendWidth, boxH),
		image.Rect(marginL, marginT, boxW-marginR, boxH-marginB))
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	// Close() is called explicitly before reading the buffer below,
	// not via defer, because defer runs after the return value is captured.

	if static {
		e.Themes = nil
	}

	if result.Focus != nil {
		focusTheme := &svgmap.FocusTheme{
			ComponentIDs: result.Focus.ComponentIDs,
			EdgeKeys:     result.Focus.EdgeKeys,
			GroupIDs:     result.Focus.GroupIDs,
		}
		if e.Themes == nil {
			e.Themes = []svgmap.Theme{focusTheme}
		} else {
			e.Themes = append(e.Themes, focusTheme)
		}
	}

	var indicators []svgmap.Annotator
	if result.Legend && len(result.LegendItems) > 0 {
		indicators = append(indicators, &svgmap.Legend{Items: result.LegendItems})
	}
	style := svgmap.NewOctoStyle(result.Stages, indicators...)
	e.Init(style)

	if err := e.Encode(result.Map); err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	e.Close()
	return svgBuf.String()
}

// svgBuf is reused across calls to avoid repeated allocation and GC pressure.
// WASM is single-threaded, so no synchronization is needed.
var svgBuf bytes.Buffer

// parseToState parses WTG2 text and returns a JSON object matching the
// JavaScript wizardState structure, enabling editor -> guided mode sync.
func parseToState(_ js.Value, args []js.Value) any {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	if len(args) < 1 {
		return "error: no input provided"
	}

	input := args[0].String()

	p, err := wtg2.NewParser(bytes.NewBufferString(input))
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	doc, err := p.Parse()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}

	state := docToState(doc)

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(data)
}

type jsState struct {
	Meta        jsMeta          `json:"meta"`
	Stages      [4]string       `json:"stages"`
	Components  []jsComponent   `json:"components"`
	Edges       []jsEdge        `json:"edges"`
	Groups      []jsGroup       `json:"groups"`
	Annotations []jsAnnotation  `json:"annotations"`
	Signals     []jsSignal      `json:"signals"`
	Legend      bool            `json:"legend"`
	Focus       string          `json:"focus"`
}

type jsMeta struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Question string `json:"question"`
}

type jsComponent struct {
	Name            string             `json:"name"`
	Kind            string             `json:"kind"`
	Evolution       int                `json:"evolution"`
	Type            string             `json:"type"`
	Evolving        bool               `json:"evolving"`
	EvolvedTo       int                `json:"evolvedTo"`
	Inertia         int                `json:"inertia"`
	InertiaKinds    []string           `json:"inertiaKinds"`
	IsPipeline      bool               `json:"isPipeline"`
	PipelineMembers []jsPipelineMember `json:"pipelineMembers"`
}

type jsPipelineMember struct {
	Name      string `json:"name"`
	Evolution int    `json:"evolution"`
}

type jsEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type jsGroup struct {
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Members []string `json:"members"`
}

type jsAnnotation struct {
	Kind   string `json:"kind"`
	Text   string `json:"text"`
	Target string `json:"target"`
}

type jsSignal struct {
	Type   string `json:"type"`
	Target string `json:"target"`
}

func docToState(doc *wtg2.Document) jsState {
	state := jsState{
		Meta: jsMeta{
			Title:    doc.Title,
			Author:   doc.Author,
			Question: doc.Question,
		},
		Stages: doc.Stages,
		Legend: doc.Legend,
	}

	// Build a set of pipeline names and their members for later lookup.
	pipelineMap := make(map[string][]*wtg2.PipelineMemberDecl)
	for _, pl := range doc.Pipelines {
		pipelineMap[pl.Name] = pl.Members
	}

	// Components
	for _, n := range doc.Nodes {
		kind := "component"
		if n.Kind == wtg2.KindAnchor {
			kind = "anchor"
		}

		evo := 50 // default for anchors without explicit position
		if n.Evolution != "" {
			if v, err := wtg2.ParsePosition(n.Evolution); err == nil {
				evo = v
			}
		}

		evolvedTo := 0
		evolving := n.EvolvedTo != ""
		if evolving {
			if v, err := wtg2.ParsePosition(n.EvolvedTo); err == nil {
				evolvedTo = v
			}
		}
		if !evolving {
			evolvedTo = min(99, evo+15)
		}

		inertiaKinds := n.InertiaKinds
		if inertiaKinds == nil {
			inertiaKinds = []string{}
		}
		comp := jsComponent{
			Name:         n.Name,
			Kind:         kind,
			Evolution:    evo,
			Type:         n.Type,
			Evolving:     evolving,
			EvolvedTo:    evolvedTo,
			Inertia:      n.Inertia,
			InertiaKinds: inertiaKinds,
		}

		// Check if this component is a pipeline parent.
		if members, ok := pipelineMap[n.Name]; ok {
			comp.IsPipeline = true
			for _, m := range members {
				mEvo := evo
				if m.Position != "" {
					if v, err := wtg2.ParsePosition(m.Position); err == nil {
						mEvo = v
					}
				}
				comp.PipelineMembers = append(comp.PipelineMembers, jsPipelineMember{
					Name:      m.Name,
					Evolution: mEvo,
				})
			}
		}

		state.Components = append(state.Components, comp)
	}

	// Edges
	for _, e := range doc.Edges {
		state.Edges = append(state.Edges, jsEdge{From: e.From, To: e.To})
	}

	// Groups
	for _, g := range doc.Groups {
		state.Groups = append(state.Groups, jsGroup{
			Name:    g.Name,
			Color:   g.Color,
			Members: g.Members,
		})
	}

	// Annotations
	for _, a := range doc.Annotations {
		state.Annotations = append(state.Annotations, jsAnnotation{
			Kind:   a.Kind,
			Text:   a.Text,
			Target: a.Target,
		})
	}

	// Signals
	for _, s := range doc.Signals {
		state.Signals = append(state.Signals, jsSignal{
			Type:   s.Type,
			Target: s.Target,
		})
	}

	// Focus (wizard supports a single focus target)
	if len(doc.Focuses) > 0 {
		state.Focus = doc.Focuses[0].Target
	}

	// Ensure nil slices become empty arrays in JSON.
	if state.Components == nil {
		state.Components = []jsComponent{}
	}
	if state.Edges == nil {
		state.Edges = []jsEdge{}
	}
	if state.Groups == nil {
		state.Groups = []jsGroup{}
	}
	if state.Annotations == nil {
		state.Annotations = []jsAnnotation{}
	}
	if state.Signals == nil {
		state.Signals = []jsSignal{}
	}

	return state
}
