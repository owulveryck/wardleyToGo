package owm2wtg2

import (
	"math"
	"testing"

	"github.com/owulveryck/wardleyToGo/parser/wtg2"
)

func TestMaturityToEvolution(t *testing.T) {
	tests := []struct {
		mat  float64
		want string
	}{
		{0.0, "I.0"},
		{0.125, "I.5"},
		{0.25, "II.0"},
		{0.5, "III.0"},
		{0.625, "III.5"},
		{0.75, "IV.0"},
		{1.0, "IV.9"},   // clamped to IV.9 (digit 10 → 9)
		{0.95, "IV.8"},
		{0.05, "I.2"},
		{0.62, "III.5"},  // 62.0 → III phase, offset 12.0, round(12/2.5)=round(4.8)=5
	}

	for _, tt := range tests {
		got := MaturityToEvolution(tt.mat)
		if got != tt.want {
			t.Errorf("MaturityToEvolution(%v) = %q, want %q", tt.mat, got, tt.want)
		}
	}
}

func TestMaturityToEvolutionRoundTrip(t *testing.T) {
	for i := 0; i <= 100; i++ {
		mat := float64(i) / 100.0
		evo := MaturityToEvolution(mat)

		parsed, err := wtg2.ParsePosition(evo)
		if err != nil {
			t.Errorf("ParsePosition(%q) from mat=%v: %v", evo, mat, err)
			continue
		}

		diff := math.Abs(float64(parsed) - float64(i))
		// IV.9 is the max position (97), so mat=1.0 (100) has a 3-unit gap
		if diff > 3 {
			t.Errorf("round-trip mat=%v → evo=%q → parsed=%d, diff=%v (>3)", mat, evo, parsed, diff)
		}
	}
}

func TestMaturityToEvolutionBoundaries(t *testing.T) {
	tests := []struct {
		mat       float64
		wantPhase string
	}{
		{0.0, "I"},
		{0.24, "I"},
		{0.25, "II"},
		{0.49, "II"},
		{0.50, "III"},
		{0.74, "III"},
		{0.75, "IV"},
		{0.99, "IV"},
	}

	for _, tt := range tests {
		got := MaturityToEvolution(tt.mat)
		if len(got) == 0 {
			t.Errorf("MaturityToEvolution(%v) returned empty", tt.mat)
			continue
		}
		dotIdx := 0
		for i, c := range got {
			if c == '.' {
				dotIdx = i
				break
			}
		}
		phase := got[:dotIdx]
		if phase != tt.wantPhase {
			t.Errorf("MaturityToEvolution(%v) = %q, phase %q, want phase %q", tt.mat, got, phase, tt.wantPhase)
		}
	}
}
