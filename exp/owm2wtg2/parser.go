package owm2wtg2

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Parse reads OWM format from r and returns an OWMDocument.
func Parse(r io.Reader) (*OWMDocument, error) {
	doc := &OWMDocument{
		URLs: make(map[string]string),
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		if err := parseLine(doc, line); err != nil {
			return nil, fmt.Errorf("parsing %q: %w", line, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return doc, nil
}

func parseLine(doc *OWMDocument, line string) error {
	switch {
	case strings.HasPrefix(line, "title "):
		doc.Title = strings.TrimPrefix(line, "title ")
	case strings.HasPrefix(line, "style "):
		doc.Style = strings.TrimPrefix(line, "style ")
	case strings.HasPrefix(line, "evolution "):
		return parseEvolution(doc, line)
	case strings.HasPrefix(line, "component "):
		return parseComponent(doc, line)
	case strings.HasPrefix(line, "anchor "):
		return parseAnchor(doc, line)
	case strings.HasPrefix(line, "evolve "):
		return parseEvolve(doc, line)
	case strings.HasPrefix(line, "pipeline "):
		return parsePipeline(doc, line)
	case strings.HasPrefix(line, "note "):
		return parseNote(doc, line)
	case strings.HasPrefix(line, "submap "):
		return parseSubmap(doc, line)
	case strings.HasPrefix(line, "url "):
		return parseURL(doc, line)
	case strings.HasPrefix(line, "pioneers "):
		return parseArea(doc, "pioneers", line)
	case strings.HasPrefix(line, "settlers "):
		return parseArea(doc, "settlers", line)
	case strings.HasPrefix(line, "townplanners "):
		return parseArea(doc, "townplanners", line)
	case strings.HasPrefix(line, "accelerator "):
		return parseSignal(doc, "accelerator", line)
	case strings.HasPrefix(line, "deaccelerator "):
		return parseSignal(doc, "deaccelerator", line)
	case strings.HasPrefix(line, "annotations "):
		// annotations placement — ignored (no WTG2 equivalent)
		return nil
	case strings.HasPrefix(line, "annotation "):
		// positioned annotation — ignored (no WTG2 equivalent)
		return nil
	default:
		return parseEdge(doc, line)
	}
	return nil
}

var bracketRe = regexp.MustCompile(`\[([^\]]+)\]`)
var labelRe = regexp.MustCompile(`label\s*\[([^\]]+)\]`)
var typeRe = regexp.MustCompile(`\((build|buy|outsource|market|dataProduct)\)`)

func parseComponent(doc *OWMDocument, line string) error {
	rest := strings.TrimPrefix(line, "component ")

	loc := bracketRe.FindStringIndex(rest)
	if loc == nil {
		return fmt.Errorf("component missing coordinates: %s", line)
	}

	name := strings.TrimSpace(rest[:loc[0]])
	coordStr := rest[loc[0]+1 : loc[1]-1]
	after := rest[loc[1]:]

	vis, mat, err := parseCoords(coordStr)
	if err != nil {
		return fmt.Errorf("component %q coords: %w", name, err)
	}

	comp := &OWMComponent{
		Name:       name,
		Visibility: vis,
		Maturity:   mat,
	}

	if strings.Contains(after, "inertia") {
		comp.Inertia = true
	}

	if m := typeRe.FindStringSubmatch(after); m != nil {
		comp.Type = m[1]
	}

	if m := labelRe.FindStringSubmatch(after); m != nil {
		lx, ly, err := parseLabelCoords(m[1])
		if err == nil {
			comp.LabelX = lx
			comp.LabelY = ly
			comp.HasLabel = true
		}
	}

	doc.Components = append(doc.Components, comp)
	return nil
}

func parseAnchor(doc *OWMDocument, line string) error {
	rest := strings.TrimPrefix(line, "anchor ")

	loc := bracketRe.FindStringIndex(rest)
	if loc == nil {
		return fmt.Errorf("anchor missing coordinates: %s", line)
	}

	name := strings.TrimSpace(rest[:loc[0]])
	coordStr := rest[loc[0]+1 : loc[1]-1]

	vis, mat, err := parseCoords(coordStr)
	if err != nil {
		return fmt.Errorf("anchor %q coords: %w", name, err)
	}

	doc.Anchors = append(doc.Anchors, &OWMAnchor{
		Name:       name,
		Visibility: vis,
		Maturity:   mat,
	})
	return nil
}

func parseEvolve(doc *OWMDocument, line string) error {
	rest := strings.TrimPrefix(line, "evolve ")

	// Remove label and type suffixes to isolate name and maturity
	clean := labelRe.ReplaceAllString(rest, "")
	clean = typeRe.ReplaceAllString(clean, "")
	clean = strings.TrimSpace(clean)

	// Extract the maturity value (last float in the clean string)
	parts := strings.Fields(clean)
	if len(parts) < 2 {
		return fmt.Errorf("evolve missing maturity: %s", line)
	}

	matStr := parts[len(parts)-1]
	mat, err := strconv.ParseFloat(matStr, 64)
	if err != nil {
		return fmt.Errorf("evolve maturity %q: %w", matStr, err)
	}

	nameStr := strings.TrimSpace(strings.Join(parts[:len(parts)-1], " "))

	ev := &OWMEvolve{
		Maturity: mat,
	}

	// Check for rename: "Name->NewName"
	if idx := strings.Index(nameStr, "->"); idx >= 0 {
		ev.Name = strings.TrimSpace(nameStr[:idx])
		ev.NewName = strings.TrimSpace(nameStr[idx+2:])
	} else {
		ev.Name = nameStr
	}

	// Remove "evolve" keyword if accidentally included in name
	ev.Name = strings.TrimPrefix(ev.Name, "evolve ")

	if m := typeRe.FindStringSubmatch(rest); m != nil {
		ev.Type = m[1]
	}

	if m := labelRe.FindStringSubmatch(rest); m != nil {
		lx, ly, err := parseLabelCoords(m[1])
		if err == nil {
			ev.LabelX = lx
			ev.LabelY = ly
			ev.HasLabel = true
		}
	}

	doc.Evolves = append(doc.Evolves, ev)
	return nil
}

var (
	flowBiRe      = regexp.MustCompile(`^(.+?)\+<>(.+)$`)
	flowPastRe    = regexp.MustCompile(`^(.+?)\+<(.+)$`)
	flowFutureRe  = regexp.MustCompile(`^(.+?)\+>(.+)$`)
	flowLabeledRe = regexp.MustCompile(`^(.+?)\+'([^']+)'>(.+)$`)
	regularEdgeRe = regexp.MustCompile(`^(.+?)->(.+)$`)
)

func parseEdge(doc *OWMDocument, line string) error {
	if m := flowLabeledRe.FindStringSubmatch(line); m != nil {
		doc.Edges = append(doc.Edges, &OWMEdge{
			From:     strings.TrimSpace(m[1]),
			To:       strings.TrimSpace(m[3]),
			Label:    m[2],
			FlowType: "labeled",
		})
		return nil
	}
	if m := flowBiRe.FindStringSubmatch(line); m != nil {
		doc.Edges = append(doc.Edges, &OWMEdge{
			From:     strings.TrimSpace(m[1]),
			To:       strings.TrimSpace(m[2]),
			FlowType: "bidirectional",
		})
		return nil
	}
	if m := flowPastRe.FindStringSubmatch(line); m != nil {
		doc.Edges = append(doc.Edges, &OWMEdge{
			From:     strings.TrimSpace(m[1]),
			To:       strings.TrimSpace(m[2]),
			FlowType: "past",
		})
		return nil
	}
	if m := flowFutureRe.FindStringSubmatch(line); m != nil {
		doc.Edges = append(doc.Edges, &OWMEdge{
			From:     strings.TrimSpace(m[1]),
			To:       strings.TrimSpace(m[2]),
			FlowType: "future",
		})
		return nil
	}
	if m := regularEdgeRe.FindStringSubmatch(line); m != nil {
		// Might be a chain: A->B->C
		parts := strings.Split(line, "->")
		for i := 0; i < len(parts)-1; i++ {
			doc.Edges = append(doc.Edges, &OWMEdge{
				From:     strings.TrimSpace(parts[i]),
				To:       strings.TrimSpace(parts[i+1]),
				FlowType: "regular",
			})
		}
		return nil
	}

	// Unrecognized line — silently ignore
	return nil
}

func parsePipeline(doc *OWMDocument, line string) error {
	rest := strings.TrimPrefix(line, "pipeline ")

	loc := bracketRe.FindStringIndex(rest)
	p := &OWMPipeline{}

	if loc != nil {
		p.Name = strings.TrimSpace(rest[:loc[0]])
		coordStr := rest[loc[0]+1 : loc[1]-1]
		parts := strings.Split(coordStr, ",")
		if len(parts) == 2 {
			start, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			end, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				p.StartMat = start
				p.EndMat = end
				p.HasRange = true
			}
		}
	} else {
		p.Name = strings.TrimSpace(rest)
	}

	doc.Pipelines = append(doc.Pipelines, p)
	return nil
}

func parseNote(doc *OWMDocument, line string) error {
	rest := strings.TrimPrefix(line, "note ")

	loc := bracketRe.FindStringIndex(rest)
	if loc == nil {
		doc.Notes = append(doc.Notes, &OWMNote{Text: rest})
		return nil
	}

	text := strings.TrimSpace(rest[:loc[0]])
	// Remove leading "+" from note text
	text = strings.TrimPrefix(text, "+")
	text = strings.TrimSpace(text)

	coordStr := rest[loc[0]+1 : loc[1]-1]
	vis, mat, err := parseCoords(coordStr)
	if err != nil {
		doc.Notes = append(doc.Notes, &OWMNote{Text: text})
		return nil
	}

	doc.Notes = append(doc.Notes, &OWMNote{
		Text:       text,
		Visibility: vis,
		Maturity:   mat,
	})
	return nil
}

func parseSubmap(doc *OWMDocument, line string) error {
	rest := strings.TrimPrefix(line, "submap ")

	loc := bracketRe.FindStringIndex(rest)
	if loc == nil {
		return fmt.Errorf("submap missing coordinates: %s", line)
	}

	name := strings.TrimSpace(rest[:loc[0]])
	coordStr := rest[loc[0]+1 : loc[1]-1]
	after := rest[loc[1]:]

	vis, mat, err := parseCoords(coordStr)
	if err != nil {
		return fmt.Errorf("submap %q coords: %w", name, err)
	}

	sm := &OWMSubmap{
		Name:       name,
		Visibility: vis,
		Maturity:   mat,
	}

	// Extract url(urlName)
	urlRe := regexp.MustCompile(`url\((\w+)\)`)
	if m := urlRe.FindStringSubmatch(after); m != nil {
		sm.URLName = m[1]
	}

	doc.Submaps = append(doc.Submaps, sm)
	return nil
}

func parseURL(doc *OWMDocument, line string) error {
	rest := strings.TrimPrefix(line, "url ")

	loc := bracketRe.FindStringIndex(rest)
	if loc == nil {
		return fmt.Errorf("url missing value: %s", line)
	}

	name := strings.TrimSpace(rest[:loc[0]])
	value := strings.TrimSpace(rest[loc[0]+1 : loc[1]-1])

	doc.URLs[name] = value
	return nil
}

func parseEvolution(doc *OWMDocument, line string) error {
	rest := strings.TrimPrefix(line, "evolution ")
	parts := strings.Split(rest, "->")
	if len(parts) != 4 {
		return fmt.Errorf("evolution expects 4 stages separated by ->, got %d", len(parts))
	}
	for i, p := range parts {
		doc.Evolution[i] = strings.TrimSpace(p)
	}
	return nil
}

func parseArea(doc *OWMDocument, kind string, line string) error {
	rest := strings.TrimPrefix(line, kind+" ")

	loc := bracketRe.FindStringIndex(rest)
	if loc == nil {
		return fmt.Errorf("%s missing coordinates: %s", kind, line)
	}

	coordStr := rest[loc[0]+1 : loc[1]-1]
	parts := strings.Split(coordStr, ",")
	if len(parts) != 4 {
		return fmt.Errorf("%s expects 4 coordinates, got %d", kind, len(parts))
	}

	vals := make([]float64, 4)
	for i, p := range parts {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return fmt.Errorf("%s coordinate %d: %w", kind, i, err)
		}
		vals[i] = v
	}

	doc.Areas = append(doc.Areas, &OWMArea{
		Kind: kind,
		Vis1: vals[0],
		Mat1: vals[1],
		Vis2: vals[2],
		Mat2: vals[3],
	})
	return nil
}

func parseSignal(doc *OWMDocument, kind string, line string) error {
	rest := strings.TrimPrefix(line, kind+" ")

	loc := bracketRe.FindStringIndex(rest)
	if loc == nil {
		return fmt.Errorf("%s missing coordinates: %s", kind, line)
	}

	name := strings.TrimSpace(rest[:loc[0]])
	coordStr := rest[loc[0]+1 : loc[1]-1]

	vis, mat, err := parseCoords(coordStr)
	if err != nil {
		return fmt.Errorf("%s %q coords: %w", kind, name, err)
	}

	doc.Signals = append(doc.Signals, &OWMSignal{
		Kind:       kind,
		Name:       name,
		Visibility: vis,
		Maturity:   mat,
	})
	return nil
}

func parseCoords(s string) (vis, mat float64, err error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected 2 coordinates, got %d in %q", len(parts), s)
	}
	vis, err = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("visibility: %w", err)
	}
	mat, err = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("maturity: %w", err)
	}
	return vis, mat, nil
}

func parseLabelCoords(s string) (x, y int, err error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected 2 label coordinates, got %d", len(parts))
	}
	x, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	y, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, err
	}
	return x, y, nil
}
