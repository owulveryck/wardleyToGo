// Package compress implements grammar-aware arithmetic compression for WTG2 files.
//
// Instead of compressing characters, it compresses grammatical decisions:
// at each point in the BNF grammar, it encodes which production rule was chosen
// using arithmetic coding with static probability distributions.
package compress

const (
	magic0  = 0x57 // 'W'
	magic1  = 0x41 // 'A'
	version = 0x01
)
