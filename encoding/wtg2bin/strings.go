package wtg2bin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// stringTable collects unique strings during encoding and assigns varint IDs.
type stringTable struct {
	strings []string
	index   map[string]uint64
}

func newStringTable() *stringTable {
	return &stringTable{index: make(map[string]uint64)}
}

func (st *stringTable) intern(s string) uint64 {
	if id, ok := st.index[s]; ok {
		return id
	}
	id := uint64(len(st.strings))
	st.strings = append(st.strings, s)
	st.index[s] = id
	return id
}

func (st *stringTable) writeTo(buf *bytes.Buffer) {
	writeUvarint(buf, uint64(len(st.strings)))
	for _, s := range st.strings {
		writeUvarint(buf, uint64(len(s)))
		buf.WriteString(s)
	}
}

// stringIndex reads a string table and provides ID-to-string lookup during decoding.
type stringIndex struct {
	strings []string
}

func readStringIndex(r io.ByteReader, full io.Reader) (*stringIndex, error) {
	count, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, fmt.Errorf("reading string table count: %w", err)
	}
	if count > maxSectionCount {
		return nil, fmt.Errorf("string table count %d exceeds maximum %d", count, maxSectionCount)
	}
	si := &stringIndex{strings: make([]string, 0, count)}
	for i := uint64(0); i < count; i++ {
		length, err := binary.ReadUvarint(r)
		if err != nil {
			return nil, fmt.Errorf("reading string %d length: %w", i, err)
		}
		if length > 1<<20 {
			return nil, fmt.Errorf("string %d length %d exceeds 1MB limit", i, length)
		}
		buf := make([]byte, length)
		if _, err := io.ReadFull(full, buf); err != nil {
			return nil, fmt.Errorf("reading string %d data: %w", i, err)
		}
		si.strings = append(si.strings, string(buf))
	}
	return si, nil
}

func (si *stringIndex) get(id uint64) (string, error) {
	if id >= uint64(len(si.strings)) {
		return "", fmt.Errorf("string ID %d out of range (table has %d entries)", id, len(si.strings))
	}
	return si.strings[id], nil
}

func writeUvarint(buf *bytes.Buffer, v uint64) {
	var tmp [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(tmp[:], v)
	buf.Write(tmp[:n])
}

func writeStringID(buf *bytes.Buffer, st *stringTable, s string) {
	writeUvarint(buf, st.intern(s))
}
