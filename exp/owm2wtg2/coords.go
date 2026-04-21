package owm2wtg2

import (
	"fmt"
	"math"
)

// MaturityToEvolution converts an OWM maturity float (0.0–1.0) to a WTG2
// evolution position string (e.g. "III.5").
func MaturityToEvolution(mat float64) string {
	val := mat * 100.0
	if val < 0 {
		val = 0
	}
	if val > 100 {
		val = 100
	}

	var roman string
	var phaseBase float64
	switch {
	case val < 25:
		roman, phaseBase = "I", 0
	case val < 50:
		roman, phaseBase = "II", 25
	case val < 75:
		roman, phaseBase = "III", 50
	default:
		roman, phaseBase = "IV", 75
	}

	offset := val - phaseBase
	digit := int(math.Round(offset / 2.5))
	if digit < 0 {
		digit = 0
	}
	if digit > 9 {
		digit = 9
	}

	return fmt.Sprintf("%s.%d", roman, digit)
}
