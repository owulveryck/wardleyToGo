package wtg

func (inv *Inventory) visibilitySetN1N2() {
	n1 := inv.upsertNode(inv.tokens[inv.offset].Value)
	n2 := inv.upsertNode(inv.tokens[inv.offset+VisibilityParseOffset].Value)
	e := inv.insertEdge(inv.tokens[inv.offset+1].Value)
	e.F = n1
	e.T = n2
	inv.offset += VisibilityParseOffset
}
func (inv *Inventory) visibilitySetP1a1N2() {
	P1 := inv.upsertNode(inv.tokens[inv.offset].Value)
	a1 := inv.upsertNode(inv.tokens[inv.offset+VisibilityParseOffset].Value)
	upsertAnchor(P1, a1)
	P2 := inv.upsertNode(inv.tokens[inv.offset+PipelineParseOffset].Value)
	e := inv.insertEdge(inv.tokens[inv.offset+3].Value)
	e.F = a1
	e.T = P2
	inv.offset += PipelineParseOffset
}
func (inv *Inventory) visibilitySetN1P2a2() {
	N1 := inv.upsertNode(inv.tokens[inv.offset].Value)
	e := inv.insertEdge(inv.tokens[inv.offset+1].Value)
	P2 := inv.upsertNode(inv.tokens[inv.offset+VisibilityParseOffset].Value)
	a2 := inv.upsertNode(inv.tokens[inv.offset+PipelineParseOffset].Value)
	upsertAnchor(P2, a2)
	e.F = N1
	e.T = a2
	inv.offset += PipelineParseOffset
}
func (inv *Inventory) visibilitySetP1a1P2a2() {
	e := inv.insertEdge(inv.tokens[inv.offset+3].Value)
	P1 := inv.upsertNode(inv.tokens[inv.offset].Value)
	a1 := inv.upsertNode(inv.tokens[inv.offset+2].Value)
	P2 := inv.upsertNode(inv.tokens[inv.offset+4].Value)
	a2 := inv.upsertNode(inv.tokens[inv.offset+6].Value)
	upsertAnchor(P1, a1)
	upsertAnchor(P2, a2)
	e.F = a1
	e.T = a2
	inv.offset += 6
}
func (inv *Inventory) visibilitySeek() bool {
	visibilityFunc := []func(){inv.visibilitySetP1a1P2a2, inv.visibilitySetP1a1N2, inv.visibilitySetN1P2a2, inv.visibilitySetN1N2}
	truthTable := []bool{true, true, true, true}
	visibilities := [][]tokenType{
		{
			identifierToken, colonToken, identifierToken,
			visibilityToken,
			identifierToken, colonToken, identifierToken,
		},
		{
			identifierToken, colonToken, identifierToken,
			visibilityToken,
			identifierToken,
		},
		{
			identifierToken,
			visibilityToken,
			identifierToken, colonToken, identifierToken,
		},
		{
			identifierToken,
			visibilityToken,
			identifierToken,
		},
	}
	for v, visibility := range visibilities {
		for j := range visibility {
			if visibility[j] != inv.peek(j).Type {
				truthTable[v] = false
			}
		}
	}
	for j, t := range truthTable {
		// return on the first match
		if t {
			visibilityFunc[j]()
			return true
		}
	}
	return false
}
