package compress

import (
	"bufio"
	"bytes"
	"testing"
)

func TestRoundTripUniform(t *testing.T) {
	dist := Uniform(4)
	symbols := []int{0, 1, 2, 3, 0, 3, 2, 1, 0, 0, 1, 3}

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	for _, s := range symbols {
		if err := enc.Encode(s, dist); err != nil {
			t.Fatalf("Encode(%d): %v", s, err)
		}
	}
	if err := enc.Finish(); err != nil {
		t.Fatalf("Finish: %v", err)
	}

	dec := NewDecoder(bufio.NewReader(&buf))
	for i, want := range symbols {
		got, err := dec.Decode(dist)
		if err != nil {
			t.Fatalf("Decode[%d]: %v", i, err)
		}
		if got != want {
			t.Errorf("Decode[%d] = %d, want %d", i, got, want)
		}
	}
}

func TestRoundTripSkewed(t *testing.T) {
	// Highly skewed: symbol 0 has 90% probability
	dist := NewDistribution([]uint32{900, 50, 30, 20})
	symbols := []int{0, 0, 0, 1, 0, 0, 2, 0, 3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	for _, s := range symbols {
		if err := enc.Encode(s, dist); err != nil {
			t.Fatalf("Encode(%d): %v", s, err)
		}
	}
	if err := enc.Finish(); err != nil {
		t.Fatalf("Finish: %v", err)
	}

	encodedSize := buf.Len()
	dec := NewDecoder(bufio.NewReader(&buf))
	for i, want := range symbols {
		got, err := dec.Decode(dist)
		if err != nil {
			t.Fatalf("Decode[%d]: %v", i, err)
		}
		if got != want {
			t.Errorf("Decode[%d] = %d, want %d", i, got, want)
		}
	}

	t.Logf("skewed: %d symbols → %d bytes", len(symbols), encodedSize)
}

func TestRoundTripBinary(t *testing.T) {
	dist := NewDistribution([]uint32{9, 1})
	var symbols []int
	for range 100 {
		symbols = append(symbols, 0)
	}
	symbols = append(symbols, 1) // one rare symbol

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	for _, s := range symbols {
		if err := enc.Encode(s, dist); err != nil {
			t.Fatalf("Encode(%d): %v", s, err)
		}
	}
	if err := enc.Finish(); err != nil {
		t.Fatalf("Finish: %v", err)
	}

	encodedSize := buf.Len()
	dec := NewDecoder(bufio.NewReader(&buf))
	for i, want := range symbols {
		got, err := dec.Decode(dist)
		if err != nil {
			t.Fatalf("Decode[%d]: %v", i, err)
		}
		if got != want {
			t.Errorf("Decode[%d] = %d, want %d", i, got, want)
		}
	}

	t.Logf("binary: %d symbols → %d bytes (%.1f bits/sym)", len(symbols), encodedSize, float64(encodedSize*8)/float64(len(symbols)))
}

func TestRoundTripMixedDistributions(t *testing.T) {
	dist1 := NewDistribution([]uint32{10, 5, 3, 2})
	dist2 := Uniform(8)

	type step struct {
		sym  int
		dist *Distribution
	}

	steps := []step{
		{0, dist1}, {1, dist1}, {0, dist1}, {3, dist1},
		{7, dist2}, {0, dist2}, {5, dist2}, {3, dist2},
		{2, dist1}, {0, dist1},
		{1, dist2}, {6, dist2},
	}

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	for _, s := range steps {
		if err := enc.Encode(s.sym, s.dist); err != nil {
			t.Fatalf("Encode(%d): %v", s.sym, err)
		}
	}
	if err := enc.Finish(); err != nil {
		t.Fatalf("Finish: %v", err)
	}

	dec := NewDecoder(bufio.NewReader(&buf))
	for i, s := range steps {
		got, err := dec.Decode(s.dist)
		if err != nil {
			t.Fatalf("Decode[%d]: %v", i, err)
		}
		if got != s.sym {
			t.Errorf("Decode[%d] = %d, want %d", i, got, s.sym)
		}
	}
}

func TestRoundTripSingleSymbol(t *testing.T) {
	dist := Uniform(1)
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	for range 50 {
		if err := enc.Encode(0, dist); err != nil {
			t.Fatalf("Encode: %v", err)
		}
	}
	if err := enc.Finish(); err != nil {
		t.Fatalf("Finish: %v", err)
	}

	encodedSize := buf.Len()
	dec := NewDecoder(bufio.NewReader(&buf))
	for i := range 50 {
		got, err := dec.Decode(dist)
		if err != nil {
			t.Fatalf("Decode[%d]: %v", i, err)
		}
		if got != 0 {
			t.Errorf("Decode[%d] = %d, want 0", i, got)
		}
	}

	t.Logf("single symbol: 50 symbols → %d bytes", encodedSize)
}

func TestRoundTripLargeAlphabet(t *testing.T) {
	n := 256
	dist := Uniform(n)

	var symbols []int
	for i := range n {
		symbols = append(symbols, i)
	}

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	for _, s := range symbols {
		if err := enc.Encode(s, dist); err != nil {
			t.Fatalf("Encode(%d): %v", s, err)
		}
	}
	if err := enc.Finish(); err != nil {
		t.Fatalf("Finish: %v", err)
	}

	encodedSize := buf.Len()
	dec := NewDecoder(bufio.NewReader(&buf))
	for i, want := range symbols {
		got, err := dec.Decode(dist)
		if err != nil {
			t.Fatalf("Decode[%d]: %v", i, err)
		}
		if got != want {
			t.Errorf("Decode[%d] = %d, want %d", i, got, want)
		}
	}

	t.Logf("large alphabet: %d symbols → %d bytes (%.1f bits/sym, theoretical = 8.0)", n, encodedSize, float64(encodedSize*8)/float64(n))
}

func TestEncodeOutOfRange(t *testing.T) {
	dist := Uniform(4)
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.Encode(-1, dist); err == nil {
		t.Error("expected error for negative symbol")
	}
	if err := enc.Encode(4, dist); err == nil {
		t.Error("expected error for symbol >= alphabet size")
	}
}

func TestDistributionCreation(t *testing.T) {
	d := NewDistribution([]uint32{10, 20, 30})
	if d.Total != 60 {
		t.Errorf("Total = %d, want 60", d.Total)
	}
	if len(d.CumFreq) != 4 {
		t.Errorf("len(CumFreq) = %d, want 4", len(d.CumFreq))
	}
	want := []uint32{0, 10, 30, 60}
	for i, w := range want {
		if d.CumFreq[i] != w {
			t.Errorf("CumFreq[%d] = %d, want %d", i, d.CumFreq[i], w)
		}
	}
}
