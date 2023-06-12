package wtg

import (
	"errors"
	"fmt"
	"image"
	"io"
	"io/ioutil"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
)

var ErrEmptyMap = errors.New("no map")

type Parser struct {
	visibilityOnly  bool
	WMap            *wardleyToGo.Map
	EvolutionStages []svgmap.Evolution
	ImageSize       image.Rectangle
	MapSize         image.Rectangle
	// InvalidEntries reports any invalid of unkonwn token
	InvalidEntries []error
	Docs           []string
}

func NewParser() *Parser {
	return &Parser{
		visibilityOnly:  true,
		WMap:            wardleyToGo.NewMap(0),
		EvolutionStages: svgmap.DefaultEvolution,
		InvalidEntries:  make([]error, 0),
	}
}

func (p *Parser) DumpComponents(w io.Writer) {
	ns := p.WMap.Nodes()
	for ns.Next() {
		if n, ok := ns.Node().(*wardley.Component); ok {
			fmt.Fprintf(w, "%v\n", n)
		}
	}
}

func (p *Parser) Parse(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return p.parse(string(b))
}

func (p *Parser) parse(s string) error {

	inv := NewInventory()
	err := inv.init(s)
	if err != nil {
		return fmt.Errorf("error in parsing: %w", err)
	}
	err = inv.Run()
	if err != nil {
		return fmt.Errorf("error in parsing: %w", err)
	}
	if len(inv.NodeInventory) == 0 {
		return ErrEmptyMap
	}
	m, err := consolidateMap(inv.NodeInventory, inv.EdgeInventory)
	if err != nil {
		return fmt.Errorf("cannot consolidate map: %w", err)
	}
	m.Title = inv.Title
	copy(p.EvolutionStages, inv.EvolutionStages)
	p.WMap = m
	p.Docs = inv.Documentation
	SetCoords(*p.WMap, true)
	SetLabelAnchor(*p.WMap)
	return nil
}
