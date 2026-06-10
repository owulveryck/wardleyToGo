package compress

import "fmt"

const z85Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.-:+=^!/*?&<>()[]{}@%$#"

var z85Decode [256]byte

func init() {
	for i := range z85Decode {
		z85Decode[i] = 0xFF
	}
	for i, c := range z85Alphabet {
		z85Decode[c] = byte(i)
	}
}

func Base85Encode(src []byte) string {
	if len(src) == 0 {
		return ""
	}

	// Number of full 4-byte groups
	nFull := len(src) / 4
	remainder := len(src) % 4

	// Output length: 5 chars per 4 bytes, plus ceil(remainder*5/4) for partial
	outLen := nFull * 5
	if remainder > 0 {
		outLen += remainder + 1
	}

	out := make([]byte, outLen)
	pos := 0

	for i := 0; i < nFull; i++ {
		off := i * 4
		val := uint32(src[off])<<24 | uint32(src[off+1])<<16 | uint32(src[off+2])<<8 | uint32(src[off+3])
		for j := 4; j >= 0; j-- {
			out[pos+j] = z85Alphabet[val%85]
			val /= 85
		}
		pos += 5
	}

	if remainder > 0 {
		// Pad with zeros to form a full 4-byte group
		var val uint32
		off := nFull * 4
		for i := 0; i < 4; i++ {
			val <<= 8
			if i < remainder {
				val |= uint32(src[off+i])
			}
		}
		// Encode full 5-char group, then take only remainder+1 chars
		var tmp [5]byte
		for j := 4; j >= 0; j-- {
			tmp[j] = z85Alphabet[val%85]
			val /= 85
		}
		copy(out[pos:], tmp[:remainder+1])
	}

	return string(out)
}

func Base85Decode(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, nil
	}

	// Number of full 5-char groups
	nFull := len(s) / 5
	remainder := len(s) % 5

	if remainder == 1 {
		return nil, fmt.Errorf("base85: invalid input length %d", len(s))
	}

	// Output length: 4 bytes per 5 chars, plus remainder-1 for partial
	outLen := nFull * 4
	if remainder > 0 {
		outLen += remainder - 1
	}

	out := make([]byte, outLen)
	pos := 0

	for i := 0; i < nFull; i++ {
		off := i * 5
		var val uint32
		for j := 0; j < 5; j++ {
			d := z85Decode[s[off+j]]
			if d == 0xFF {
				return nil, fmt.Errorf("base85: invalid character %q at position %d", s[off+j], off+j)
			}
			val = val*85 + uint32(d)
		}
		out[pos] = byte(val >> 24)
		out[pos+1] = byte(val >> 16)
		out[pos+2] = byte(val >> 8)
		out[pos+3] = byte(val)
		pos += 4
	}

	if remainder > 0 {
		off := nFull * 5
		var val uint32
		for j := 0; j < remainder; j++ {
			d := z85Decode[s[off+j]]
			if d == 0xFF {
				return nil, fmt.Errorf("base85: invalid character %q at position %d", s[off+j], off+j)
			}
			val = val*85 + uint32(d)
		}
		// Pad with 84 (highest digit) to reconstruct the high bytes
		for j := remainder; j < 5; j++ {
			val = val*85 + 84
		}
		nBytes := remainder - 1
		for j := 0; j < nBytes; j++ {
			out[pos+j] = byte(val >> (24 - j*8))
		}
	}

	return out, nil
}
