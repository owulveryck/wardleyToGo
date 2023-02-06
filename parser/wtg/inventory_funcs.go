package wtg

import (
	"fmt"
	"strings"

	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func (inv *inventorier) start() error {
	for inv.offset = 0; inv.offset < len(inv.tokens); inv.offset++ {
		//log.Printf("%v: %v", inv.offset, inv.peek(0).Value)
		var err error
		switch inv.currentToken().Type {
		case identifierToken:
			err = inv.sourceNodeState()
		case titleItem:
			inv.title = strings.TrimSpace(inv.peek(0).Value)
		case stage1Item:
			inv.evolutionStages[0].Label = inv.peek(0).Value
		case stage2Item:
			inv.evolutionStages[1].Label = inv.peek(0).Value
		case stage3Item:
			inv.evolutionStages[2].Label = inv.peek(0).Value
		case stage4Item:
			inv.evolutionStages[3].Label = inv.peek(0).Value
		case startBlockToken:
			err = inv.inComment()
		case eofToken:
			return nil
		case evolutionItem:
			return fmt.Errorf("unhandled element in first col %v", inv.peek(0))
		case visibilityToken:
			return fmt.Errorf("unhandled element in first col %v", inv.peek(0))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (inv *inventorier) inComment() error {
	openComment := 1
	for inv.offset++; openComment > 0; inv.offset++ {
		switch inv.peek(0).Type {
		case startBlockToken:
			openComment++
		case endBlockToken:
			openComment--
		case eofToken:
			return fmt.Errorf("unbalanced comment")
		}
	}
	return nil

}

// offset is a node and offset+1 is a colon
func (inv *inventorier) nodeConfiguration() error {
	n := inv.upsertNode(inv.tokens[inv.offset].Value)
	//log.Printf(".... %v", inv.peek(2))
	switch inv.peek(2).Type {
	case evolutionItem:
		pos, evolutionPos, err := computeEvolutionPosition(inv.peek(2).Value)
		if err != nil {
			return err
		}
		n.Placement.X = pos
		n.Configured = true
		n.EvolutionPos = evolutionPos
		inv.offset += 2
		return nil
	case startBlockToken:
		return inv.nodeBlock()
	case identifierToken:
		a := inv.upsertNode(inv.peek(2).Value)
		n.Type = wardley.PipelineComponent
		upsertAnchor(n, a)
	default:
		//return fmt.Errorf("expected evolution or configuration, got %v %v %v", inv.peek(0).Value, inv.peek(1).Value, inv.peek(2).Value)
	}
	return nil
}

// offset is node, offset:1 is colon offset+2 is open bracket
func (inv *inventorier) nodeBlock() error {
	n := inv.upsertNode(inv.tokens[inv.offset].Value)
	openBrackets := 1
	for inv.offset += 3; openBrackets > 0; inv.offset++ {
		//log.Printf("nodeBlock: %v: %v", inv.offset, inv.peek(0).Value)
		switch inv.peek(0).Type {
		case startBlockToken:
			openBrackets++
		case endBlockToken:
			openBrackets--
		case evolutionItem:
			pos, evolutionPos, err := computeEvolutionPosition(inv.peek(0).Value)
			if err != nil {
				return err
			}
			n.Placement.X = pos
			n.Configured = true
			n.EvolutionPos = evolutionPos
		case typeItem:
			switch inv.peek(0).Value {
			case "pipeline":
				n.Type = wardley.PipelineComponent
				n.LabelPlacement.Y = -15
				n.LabelPlacement.X = 0
				n.Anchor = wardley.AdjustMiddle
			case "build":
				n.Type = wardley.BuildComponent
			case "buy":
				n.Type = wardley.BuyComponent
			case "outsource":
				n.Type = wardley.OutsourceComponent
			default:
				return fmt.Errorf("unknown type %v", inv.peek(0).Value)
			}
		case labelItem:
			err := setLabelPlacement(n, inv.peek(0).Value)
			if err != nil {
				return err
			}
		case colorItem:
			if col, ok := Colors[inv.peek(0).Value]; ok {
				n.Color = col
				continue
			}
			return fmt.Errorf("unknown color %v", inv.peek(0).Value)
		case eofToken:
			return fmt.Errorf("unbalanced block in %v %v", n.Label, inv.peek(0).Value)
		}
	}
	inv.offset--
	return nil
}

// sourceNodeState is a state where we have an initial node that can act as a source of a link
// or the node to be configured
// startOffset is the offset from where it has been called (and i.tokens[startOFfset] must be an identifierToken)
func (inv *inventorier) sourceNodeState() error {
	if inv.visibilitySeek() {
		return nil
	}
	if inv.peek(1).Type == colonToken {
		return inv.nodeConfiguration()
	}
	inv.upsertNode(inv.tokens[inv.offset].Value)
	return nil
}

func (inv *inventorier) titleState() error {
	return fmt.Errorf("not implemented")
}
