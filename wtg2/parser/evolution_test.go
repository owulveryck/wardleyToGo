package parser

import "testing"

func TestParsePosition(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"I.0", 0},
		{"I.5", 12},
		{"II.0", 25},
		{"II.3", 32},
		{"II.7", 42},
		{"III.0", 50},
		{"III.1", 52},
		{"III.2", 55},
		{"III.5", 62},
		{"III.8", 70},
		{"IV.0", 75},
		{"IV.3", 82},
		{"IV.5", 87},
		{"IV.7", 92},
		{"IV.9", 97},
		// Without decimal: center of phase
		{"I", 12},
		{"II", 37},
		{"III", 62},
		{"IV", 87},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParsePosition(tt.input)
			if err != nil {
				t.Fatalf("ParsePosition(%q) returned error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ParsePosition(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestParsePositionErrors(t *testing.T) {
	errorCases := []string{
		"",
		"V.0",
		"X.5",
		"hello",
	}
	for _, input := range errorCases {
		t.Run(input, func(t *testing.T) {
			_, err := ParsePosition(input)
			if err == nil {
				t.Errorf("ParsePosition(%q) expected error, got nil", input)
			}
		})
	}
}
