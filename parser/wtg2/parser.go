package wtg2

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Parser reads WTG2 source and produces a Document AST.
type Parser struct {
	lexer *Lexer
}

// NewParser creates a parser from an io.Reader.
func NewParser(r io.Reader) (*Parser, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}
	return &Parser{
		lexer: NewLexer(string(data)),
	}, nil
}

// Parse reads the entire WTG2 input and produces a Document.
func (p *Parser) Parse() (*Document, error) {
	doc := &Document{
		Stages: [4]string{"", "", "", ""},
	}

	for {
		tok := p.lexer.Next()
		if tok.Type == TokenEOF {
			break
		}

		var err error
		switch tok.Type {
		case TokenMeta:
			err = p.parseMeta(doc, tok.Text)
		case TokenStages:
			err = p.parseStages(doc, tok.Text)
		case TokenAnchor:
			err = p.parseAnchor(doc, tok.Text, tok.Block)
		case TokenComponent:
			err = p.parseComponent(doc, tok.Text, tok.Block)
		case TokenSubmap:
			err = p.parseSubmap(doc, tok.Text, tok.Block)
		case TokenPipeline:
			err = p.parsePipeline(doc, tok.Text, tok.Block)
		case TokenEdge:
			err = p.parseEdgeChain(doc, tok.Text)
		case TokenGroup:
			err = p.parseGroup(doc, tok.Text, tok.Block)
		case TokenNote:
			err = p.parseAnnotation(doc, "note", tok.Text)
		case TokenWarning:
			err = p.parseAnnotation(doc, "warning", tok.Text)
		case TokenSignal:
			err = p.parseSignal(doc, tok.Text)
		case TokenGameplay:
			err = p.parseGameplay(doc, tok.Text)
		case TokenLegend:
			doc.Legend = true
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", tok.Line, err)
		}
	}

	return doc, nil
}

func (p *Parser) parseMeta(doc *Document, line string) error {
	colonIdx := strings.Index(line, ":")
	if colonIdx < 0 {
		return fmt.Errorf("meta line missing colon: %q", line)
	}
	key := strings.TrimSpace(strings.ToLower(line[:colonIdx]))
	value := strings.TrimSpace(line[colonIdx+1:])
	value = strings.Trim(value, `"`)

	switch key {
	case "title":
		doc.Title = value
	case "date":
		doc.Date = value
	case "author":
		doc.Author = value
	case "scope":
		doc.Scope = value
	case "question":
		doc.Question = value
	case "doctrine":
		doc.Doctrine = value
	}
	return nil
}

func (p *Parser) parseStages(doc *Document, line string) error {
	colonIdx := strings.Index(line, ":")
	if colonIdx < 0 {
		return fmt.Errorf("stages line missing colon: %q", line)
	}
	value := strings.TrimSpace(line[colonIdx+1:])
	parts := strings.Split(value, ",")
	if len(parts) != 4 {
		return fmt.Errorf("stages must have exactly 4 labels, got %d", len(parts))
	}
	for i, part := range parts {
		doc.Stages[i] = strings.TrimSpace(part)
	}
	return nil
}

func (p *Parser) parseAnchor(doc *Document, line string, block []string) error {
	name := strings.TrimSpace(strings.TrimPrefix(line, "anchor"))

	node := &NodeDecl{
		Kind:       KindAnchor,
		Visibility: -1,
	}

	if idx := strings.Index(name, " : "); idx >= 0 {
		node.Name = strings.TrimSpace(name[:idx])
		rest := strings.TrimSpace(name[idx+3:])
		if err := parseShorthand(node, rest); err != nil {
			return err
		}
	} else {
		node.Name = name
	}

	if node.Name == "" {
		return fmt.Errorf("anchor missing name")
	}

	if len(block) > 0 {
		if err := parseBlockConfig(node, block); err != nil {
			return err
		}
	}

	doc.Nodes = append(doc.Nodes, node)
	return nil
}

func (p *Parser) parseComponent(doc *Document, line string, block []string) error {
	// Strip "component " prefix if present
	name := line
	lower := strings.ToLower(line)
	if strings.HasPrefix(lower, "component ") {
		name = strings.TrimSpace(line[len("component "):])
	}

	node := &NodeDecl{
		Kind:       KindComponent,
		Visibility: -1,
	}

	// Check for shorthand: "Name : evolution_expr"
	if idx := strings.Index(name, " : "); idx >= 0 {
		node.Name = strings.TrimSpace(name[:idx])
		rest := strings.TrimSpace(name[idx+3:])
		if err := parseShorthand(node, rest); err != nil {
			return err
		}
	} else {
		node.Name = name
	}

	// Parse block config if present
	if len(block) > 0 {
		if err := parseBlockConfig(node, block); err != nil {
			return err
		}
	}

	doc.Nodes = append(doc.Nodes, node)
	return nil
}

func (p *Parser) parseSubmap(doc *Document, line string, block []string) error {
	name := strings.TrimSpace(strings.TrimPrefix(line, "submap"))

	node := &NodeDecl{
		Kind:       KindSubmap,
		Visibility: -1,
	}

	if idx := strings.Index(name, " : "); idx >= 0 {
		node.Name = strings.TrimSpace(name[:idx])
		rest := strings.TrimSpace(name[idx+3:])
		if err := parseShorthand(node, rest); err != nil {
			return err
		}
	} else {
		node.Name = name
	}

	if len(block) > 0 {
		if err := parseBlockConfig(node, block); err != nil {
			return err
		}
	}

	doc.Nodes = append(doc.Nodes, node)
	return nil
}

// parseShorthand parses the part after " : " in a node declaration.
// Examples: "III.5", "III.5 (buy)", "II.7 !! >> III.5", "II.7 (buy) @0.8"
func parseShorthand(node *NodeDecl, s string) error {
	s = strings.TrimSpace(s)

	// Extract @visibility if present
	if atIdx := strings.LastIndex(s, "@"); atIdx >= 0 {
		visStr := strings.TrimSpace(s[atIdx+1:])
		s = strings.TrimSpace(s[:atIdx])
		if v, err := strconv.ParseFloat(visStr, 64); err == nil {
			node.Visibility = v
		}
	}

	// Extract (type) if present — but not inertia qualifiers like !(tech)
	if openParen := strings.Index(s, "("); openParen >= 0 {
		// Skip if preceded by '!' (inertia qualifier)
		isInertiaQualifier := openParen > 0 && s[openParen-1] == '!'
		if !isInertiaQualifier {
			closeParen := strings.Index(s, ")")
			if closeParen > openParen {
				node.Type = strings.TrimSpace(s[openParen+1 : closeParen])
				s = strings.TrimSpace(s[:openParen]) + " " + strings.TrimSpace(s[closeParen+1:])
				s = strings.TrimSpace(s)
			}
		}
	}

	// Parse evolution expression: "II.7 !! >> III.5"
	return parseEvolutionExpr(node, s)
}

// parseEvolutionExpr parses "II.7", "II.7 >> III.5", "II.7 !! >> III.5",
// or "II.7 !!(tech,human) >> III.5".
func parseEvolutionExpr(node *NodeDecl, s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	// Check for ">>" (evolution movement)
	parts := strings.SplitN(s, ">>", 2)
	if len(parts) == 2 {
		left := strings.TrimSpace(parts[0])
		node.EvolvedTo = strings.TrimSpace(parts[1])

		// Extract inertia (!) and optional kinds from left part
		node.Inertia, node.InertiaKinds = parseInertia(left)
		left = stripInertia(left)
		node.Evolution = strings.TrimSpace(left)
	} else {
		// Extract inertia
		node.Inertia, node.InertiaKinds = parseInertia(s)
		s = stripInertia(s)
		node.Evolution = strings.TrimSpace(s)
	}

	return nil
}

// parseInertia extracts inertia level and optional kinds from a string like
// "II.7 !!(tech,human)" → (2, ["tech","human"]) or "II.7 !!" → (2, nil).
func parseInertia(s string) (int, []string) {
	count := 0
	for _, c := range s {
		if c == '!' {
			count++
		}
	}
	if count == 0 {
		return 0, nil
	}

	// Check for qualified inertia: !!(kind1,kind2)
	// Find the last '!' then look for '(' immediately after
	lastBang := strings.LastIndex(s, "!")
	if lastBang < len(s)-1 {
		after := strings.TrimSpace(s[lastBang+1:])
		if strings.HasPrefix(after, "(") {
			closeIdx := strings.Index(after, ")")
			if closeIdx > 0 {
				inner := after[1:closeIdx]
				var kinds []string
				for _, k := range strings.Split(inner, ",") {
					k = strings.TrimSpace(k)
					if k != "" {
						kinds = append(kinds, k)
					}
				}
				return count, kinds
			}
		}
	}

	return count, nil
}

// stripInertia removes inertia markers (!, !!(kind,...)) from a string,
// returning only the position part.
func stripInertia(s string) string {
	// Find the first '!' to know where inertia starts
	bangIdx := strings.Index(s, "!")
	if bangIdx < 0 {
		return s
	}
	return strings.TrimSpace(s[:bangIdx])
}

// parseBlockConfig parses key-value lines inside a node's {...} block.
func parseBlockConfig(node *NodeDecl, lines []string) error {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		colonIdx := strings.Index(line, ":")
		if colonIdx < 0 {
			continue
		}
		key := strings.TrimSpace(strings.ToLower(line[:colonIdx]))
		value := strings.TrimSpace(line[colonIdx+1:])
		value = strings.Trim(value, `"`)

		switch key {
		case "evolution":
			if err := parseEvolutionExpr(node, value); err != nil {
				return err
			}
		case "type":
			node.Type = value
		case "asset":
			node.Asset = value
		case "color":
			node.Color = value
		case "visibility":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				node.Visibility = v
			}
		case "cost":
			node.Cost = value
		case "note":
			node.Note = value
		}
	}
	return nil
}

func (p *Parser) parsePipeline(doc *Document, line string, block []string) error {
	name := strings.TrimSpace(strings.TrimPrefix(line, "pipeline"))

	pipeline := &PipelineDecl{Name: name}
	for _, bline := range block {
		bline = strings.TrimSpace(bline)
		if bline == "" || strings.HasPrefix(bline, "//") {
			continue
		}
		colonIdx := strings.Index(bline, " : ")
		if colonIdx < 0 {
			// Try without spaces around colon
			colonIdx = strings.Index(bline, ":")
			if colonIdx < 0 {
				continue
			}
			mName := strings.TrimSpace(bline[:colonIdx])
			mPos := strings.TrimSpace(bline[colonIdx+1:])
			pipeline.Members = append(pipeline.Members, &PipelineMemberDecl{
				Name:     mName,
				Position: mPos,
			})
			continue
		}
		mName := strings.TrimSpace(bline[:colonIdx])
		mPos := strings.TrimSpace(bline[colonIdx+3:])
		pipeline.Members = append(pipeline.Members, &PipelineMemberDecl{
			Name:     mName,
			Position: mPos,
		})
	}

	doc.Pipelines = append(doc.Pipelines, pipeline)
	return nil
}

// edgeSplitter matches -> , <-> , and -[text]->
var edgeSplitter = regexp.MustCompile(`\s+(<->|-\[([^\]]*)\]->|->)\s+`)

func (p *Parser) parseEdgeChain(doc *Document, line string) error {
	// Find all link operators and split the line
	matches := edgeSplitter.FindAllStringSubmatchIndex(line, -1)
	if len(matches) == 0 {
		return fmt.Errorf("no edge operator found in %q", line)
	}

	// Extract node references between operators
	var nodes []string
	var operators []struct {
		op    string
		label string
	}

	prevEnd := 0
	for _, match := range matches {
		nodeRef := strings.TrimSpace(line[prevEnd:match[0]])
		nodes = append(nodes, nodeRef)

		op := line[match[2]:match[3]]
		var label string
		if match[4] >= 0 && match[5] >= 0 {
			label = line[match[4]:match[5]]
		}
		operators = append(operators, struct {
			op    string
			label string
		}{op: op, label: label})
		prevEnd = match[1]
	}
	// Last node
	lastNode := strings.TrimSpace(line[prevEnd:])
	nodes = append(nodes, lastNode)

	// Create edges for each adjacent pair
	for i, op := range operators {
		edge := &EdgeDecl{
			From:          nodes[i],
			To:            nodes[i+1],
			Label:         op.label,
			Bidirectional: op.op == "<->",
		}
		doc.Edges = append(doc.Edges, edge)
	}

	return nil
}

func (p *Parser) parseGroup(doc *Document, line string, block []string) error {
	name := strings.TrimSpace(strings.TrimPrefix(line, "group"))

	group := &GroupDecl{Name: name}
	for _, bline := range block {
		bline = strings.TrimSpace(bline)
		if bline == "" || strings.HasPrefix(bline, "//") {
			continue
		}
		// Parse optional color directive
		if strings.HasPrefix(strings.ToLower(bline), "color:") {
			group.Color = strings.TrimSpace(strings.TrimPrefix(bline, bline[:len("color:")]))
			continue
		}
		// Parse optional team directive
		if strings.HasPrefix(strings.ToLower(bline), "team:") {
			group.Team = strings.TrimSpace(bline[len("team:"):])
			continue
		}
		// Members can be node names or edge declarations (ignored for now)
		if strings.Contains(bline, " -> ") || strings.Contains(bline, " <-> ") {
			continue
		}
		group.Members = append(group.Members, bline)
	}

	doc.Groups = append(doc.Groups, group)
	return nil
}

// parseAnnotation handles both "note" and "warning" lines.
// Format: note "text" on NodeName
//
//	warning "text" on NodeName
func (p *Parser) parseAnnotation(doc *Document, kind string, line string) error {
	// Remove the keyword prefix
	rest := strings.TrimSpace(strings.TrimPrefix(line, kind))

	// Extract quoted text
	firstQuote := strings.Index(rest, `"`)
	if firstQuote < 0 {
		return fmt.Errorf("annotation missing quoted text: %q", line)
	}
	secondQuote := strings.Index(rest[firstQuote+1:], `"`)
	if secondQuote < 0 {
		return fmt.Errorf("annotation missing closing quote: %q", line)
	}
	text := rest[firstQuote+1 : firstQuote+1+secondQuote]

	// Find "on" keyword
	afterQuote := rest[firstQuote+1+secondQuote+1:]
	onIdx := strings.Index(strings.ToLower(afterQuote), " on ")
	if onIdx < 0 {
		return fmt.Errorf("annotation missing 'on' keyword: %q", line)
	}
	target := strings.TrimSpace(afterQuote[onIdx+4:])

	doc.Annotations = append(doc.Annotations, &AnnotationDecl{
		Kind:   kind,
		Text:   text,
		Target: target,
	})
	return nil
}

func (p *Parser) parseSignal(doc *Document, line string) error {
	rest := strings.TrimSpace(strings.TrimPrefix(line, "signal"))

	// Format: signal type on NodeName
	onIdx := strings.Index(strings.ToLower(rest), " on ")
	if onIdx < 0 {
		return fmt.Errorf("signal missing 'on' keyword: %q", line)
	}

	signalType := strings.TrimSpace(rest[:onIdx])
	target := strings.TrimSpace(rest[onIdx+4:])

	doc.Signals = append(doc.Signals, &SignalDecl{
		Type:   signalType,
		Target: target,
	})
	return nil
}

func (p *Parser) parseGameplay(doc *Document, line string) error {
	rest := strings.TrimSpace(strings.TrimPrefix(line, "gameplay"))

	gp := &GameplayDecl{}

	// Check for optional quoted description before "on"
	if quoteIdx := strings.Index(rest, `"`); quoteIdx >= 0 {
		// Extract gameplay type before the quote
		gp.Type = strings.TrimSpace(rest[:quoteIdx])

		// Extract quoted text
		secondQuote := strings.Index(rest[quoteIdx+1:], `"`)
		if secondQuote < 0 {
			return fmt.Errorf("gameplay missing closing quote: %q", line)
		}
		gp.Text = rest[quoteIdx+1 : quoteIdx+1+secondQuote]

		// Find "on" keyword after the closing quote
		afterQuote := rest[quoteIdx+1+secondQuote+1:]
		onIdx := strings.Index(strings.ToLower(afterQuote), " on ")
		if onIdx < 0 {
			return fmt.Errorf("gameplay missing 'on' keyword: %q", line)
		}
		gp.Target = strings.TrimSpace(afterQuote[onIdx+4:])
	} else {
		// No quoted text: "gameplay type on NodeName"
		onIdx := strings.Index(strings.ToLower(rest), " on ")
		if onIdx < 0 {
			return fmt.Errorf("gameplay missing 'on' keyword: %q", line)
		}
		gp.Type = strings.TrimSpace(rest[:onIdx])
		gp.Target = strings.TrimSpace(rest[onIdx+4:])
	}

	doc.Gameplays = append(doc.Gameplays, gp)
	return nil
}
