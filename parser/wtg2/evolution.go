package wtg2

import (
	"fmt"
	"strconv"
	"strings"
)

// romanBases maps roman numerals to their base fraction on the evolution axis.
var romanBases = map[string]float64{
	"I":   0.00,
	"II":  0.25,
	"III": 0.50,
	"IV":  0.75,
}

// ParsePosition converts a WTG2 position notation (e.g. "III.5") to a 0-100 integer coordinate.
// Each roman numeral phase spans 25% of the axis:
//
//	I = [0, 25), II = [25, 50), III = [50, 75), IV = [75, 100)
//
// The decimal part subdivides within the phase: 0 = start, 9 = end.
// Without a decimal (e.g. "III"), the center of the phase is used.
func ParsePosition(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty position")
	}

	roman, decimal, hasDecimal := parsePositionParts(s)

	base, ok := romanBases[roman]
	if !ok {
		return 0, fmt.Errorf("invalid roman numeral %q in position %q", roman, s)
	}

	var frac float64
	if hasDecimal {
		d, err := strconv.Atoi(decimal)
		if err != nil {
			return 0, fmt.Errorf("invalid decimal %q in position %q: %w", decimal, s, err)
		}
		frac = float64(d) / 10.0
	} else {
		frac = 0.5 // center of phase
	}

	return int((base + frac*0.25) * 100), nil
}

// parsePositionParts splits "III.5" into ("III", "5", true) or "III" into ("III", "", false).
func parsePositionParts(s string) (roman, decimal string, hasDecimal bool) {
	dotIdx := strings.Index(s, ".")
	if dotIdx < 0 {
		return s, "", false
	}
	return s[:dotIdx], s[dotIdx+1:], true
}
