package compress

// symTable is a growing symbol table for identifier encoding.
// During encoding, new identifiers are added and subsequent references
// use the index. The decoder maintains the same table in lockstep.
type symTable struct {
	entries []string
	index   map[string]int
}

func newSymTable() *symTable {
	return &symTable{index: make(map[string]int)}
}

func (s *symTable) lookup(name string) (int, bool) {
	idx, ok := s.index[name]
	return idx, ok
}

func (s *symTable) add(name string) int {
	idx := len(s.entries)
	s.entries = append(s.entries, name)
	s.index[name] = idx
	return idx
}

func (s *symTable) get(idx int) string {
	return s.entries[idx]
}

func (s *symTable) len() int {
	return len(s.entries)
}
