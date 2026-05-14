// Package wtg2bin encodes and decodes WTG2 Document ASTs in a compact binary format.
//
// The format uses a string table for deduplication, enum codes for fixed vocabularies,
// and compact evolution position encoding. An optional deflate compression layer
// further reduces size. Typical compression is ~20% of the original WTG2 text size.
package wtg2bin

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	magic0  = 0x57 // 'W'
	magic1  = 0x42 // 'B'
	version = 0x01

	flagLegend     = 1 << 0
	flagCompressed = 1 << 1

	maxSectionCount = 10_000

	evolAbsent   = 0x00
	evolFallback = 0xFF

	enumFallback = 0xFF
)

// Doctrine enum encoding.
var (
	doctrineToCode = map[string]byte{
		"":           0,
		"hygiene":    1,
		"context":    2,
		"excellence": 3,
		"evolution":  4,
	}
	codeToDoctrine = reverseByteMap(doctrineToCode)
)

// NodeType enum encoding.
var (
	nodeTypeToCode = map[string]byte{
		"":           0,
		"build":      1,
		"buy":        2,
		"outsource":  3,
	}
	codeToNodeType = reverseByteMap(nodeTypeToCode)
)

// Asset enum encoding.
var (
	assetToCode = map[string]byte{
		"":           0,
		"tech":       1,
		"financial":  2,
		"human":      3,
		"relational": 4,
		"social":     5,
	}
	codeToAsset = reverseByteMap(assetToCode)
)

// Team enum encoding.
var (
	teamToCode = map[string]byte{
		"":             0,
		"explorer":     1,
		"settler":      2,
		"town-planner": 3,
		"pioneer":      4,
		"villager":     5,
	}
	codeToTeam = reverseByteMap(teamToCode)
)

// AnnotationKind enum encoding.
var (
	annotationKindToCode = map[string]byte{
		"note":    0,
		"warning": 1,
	}
	codeToAnnotationKind = reverseByteMap(annotationKindToCode)
)

// SignalType enum encoding.
var (
	signalTypeToCode = map[string]byte{
		"accelerating":       0,
		"stagnating":         1,
		"declining":          2,
		"co-evolution":       3,
		"red-queen":          4,
		"commoditization":    5,
		"network-effects":    6,
		"economies-of-scale": 7,
	}
	codeToSignalType = reverseByteMap(signalTypeToCode)
)

// GameplayType enum encoding.
var (
	gameplayTypeToCode = map[string]byte{
		"ILC":               0,
		"open-source":       1,
		"land-grab":         2,
		"embrace-extend":    3,
		"tower-moat":        4,
		"FUD":               5,
		"strangler-fig":     6,
		"signal-distortion": 7,
	}
	codeToGameplayType = reverseByteMap(gameplayTypeToCode)
)

// InertiaKind bitmask encoding.
var inertiaKindBit = map[string]byte{
	"tech":       1 << 0,
	"financial":  1 << 1,
	"human":      1 << 2,
	"relational": 1 << 3,
	"social":     1 << 4,
}

// inertiaKindOrder defines the canonical order for decoding bitmask back to slice.
var inertiaKindOrder = []struct {
	name string
	bit  byte
}{
	{"tech", 1 << 0},
	{"financial", 1 << 1},
	{"human", 1 << 2},
	{"relational", 1 << 3},
	{"social", 1 << 4},
}

func encodeInertiaKinds(kinds []string) (byte, error) {
	var mask byte
	for _, k := range kinds {
		b, ok := inertiaKindBit[k]
		if !ok {
			return 0, fmt.Errorf("unknown inertia kind %q", k)
		}
		mask |= b
	}
	return mask, nil
}

func decodeInertiaKinds(mask byte) []string {
	if mask == 0 {
		return nil
	}
	var out []string
	for _, entry := range inertiaKindOrder {
		if mask&entry.bit != 0 {
			out = append(out, entry.name)
		}
	}
	return out
}

// Evolution compact encoding.
//
// Roman numeral phases I-IV map to phase index 0-3.
// Each phase has 11 slots: digits 0-9 plus 10 for "bare" (no decimal).
// Compact byte = phase*11 + slot + 1, giving range 0x01..0x2C (1..44).
// 0x00 = absent (empty string), 0xFF = fallback to string table.

var romanToPhase = map[string]int{"I": 0, "II": 1, "III": 2, "IV": 3}
var phaseToRoman = [4]string{"I", "II", "III", "IV"}

func encodeEvolution(s string) (byte, bool) {
	if s == "" {
		return evolAbsent, true
	}
	dotIdx := strings.Index(s, ".")
	if dotIdx < 0 {
		phase, ok := romanToPhase[s]
		if !ok {
			return evolFallback, false
		}
		return byte(phase*11 + 10 + 1), true
	}
	roman := s[:dotIdx]
	decimal := s[dotIdx+1:]
	phase, ok := romanToPhase[roman]
	if !ok {
		return evolFallback, false
	}
	d, err := strconv.Atoi(decimal)
	if err != nil || d < 0 || d > 9 {
		return evolFallback, false
	}
	return byte(phase*11 + d + 1), true
}

func decodeEvolution(b byte) (string, bool) {
	if b == evolAbsent {
		return "", true
	}
	if b == evolFallback {
		return "", false
	}
	v := int(b) - 1
	phase := v / 11
	slot := v % 11
	if phase < 0 || phase > 3 || slot > 10 {
		return "", false
	}
	roman := phaseToRoman[phase]
	if slot == 10 {
		return roman, true
	}
	return roman + "." + strconv.Itoa(slot), true
}

func reverseByteMap(m map[string]byte) map[byte]string {
	r := make(map[byte]string, len(m))
	for k, v := range m {
		r[v] = k
	}
	return r
}
