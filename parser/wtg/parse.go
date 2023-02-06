package wtg

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"

	"github.com/owulveryck/wardleyToGo"
	"github.com/owulveryck/wardleyToGo/components/wardley"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
)

type Parser struct {
	visibilityOnly  bool
	WMap            *wardleyToGo.Map
	EvolutionStages []svgmap.Evolution
	ImageSize       image.Rectangle
	MapSize         image.Rectangle
	// InvalidEntries reports any invalid of unkonwn token
	InvalidEntries []error
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

	inv := &inventorier{}
	err := inv.init(s)
	if err != nil {
		return fmt.Errorf("error in parsing: %w", err)
	}
	err = inv.start()
	if err != nil {
		return fmt.Errorf("error in parsing: %w", err)
	}
	if len(inv.nodeInventory) == 0 {
		return fmt.Errorf("no map")
	}
	m, err := consolidateMap(inv.nodeInventory, inv.edgeInventory)
	if err != nil {
		return fmt.Errorf("cannot consolidate map: %w", err)
	}
	m.Title = inv.title
	log.Println(inv.evolutionStages)
	copy(p.EvolutionStages, inv.evolutionStages)
	p.WMap = m
	SetCoords(*p.WMap, true)
	SetLabelAnchor(*p.WMap)
	return nil
}
