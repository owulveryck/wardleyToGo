package compress

import (
	"fmt"
	"io"
)

const (
	precBits  = 32
	fullRange = uint64(1) << precBits // 2^32
	halfRange = fullRange >> 1        // 2^31
	qtrRange  = halfRange >> 1        // 2^30
	maxRange  = fullRange - 1         // 0xFFFFFFFF
)

// Distribution represents a cumulative probability distribution.
// CumFreq has length N+1 where N is the number of symbols.
// CumFreq[0] = 0 and CumFreq[N] = Total.
// Symbol i has frequency CumFreq[i+1] - CumFreq[i].
type Distribution struct {
	CumFreq []uint32
	Total   uint32
}

// Uniform returns a uniform distribution over n symbols.
func Uniform(n int) *Distribution {
	cum := make([]uint32, n+1)
	for i := 1; i <= n; i++ {
		cum[i] = uint32(i)
	}
	return &Distribution{CumFreq: cum, Total: uint32(n)}
}

// NewDistribution creates a distribution from a slice of frequencies.
func NewDistribution(freqs []uint32) *Distribution {
	cum := make([]uint32, len(freqs)+1)
	for i, f := range freqs {
		cum[i+1] = cum[i] + f
	}
	return &Distribution{CumFreq: cum, Total: cum[len(freqs)]}
}

// bitWriter accumulates individual bits and flushes whole bytes to w.
type bitWriter struct {
	w     io.Writer
	buf   byte
	nbits int
}

func (bw *bitWriter) writeBit(bit byte) error {
	bw.buf = (bw.buf << 1) | (bit & 1)
	bw.nbits++
	if bw.nbits == 8 {
		_, err := bw.w.Write([]byte{bw.buf})
		bw.buf = 0
		bw.nbits = 0
		return err
	}
	return nil
}

func (bw *bitWriter) flush() error {
	if bw.nbits > 0 {
		bw.buf <<= (8 - bw.nbits)
		_, err := bw.w.Write([]byte{bw.buf})
		bw.buf = 0
		bw.nbits = 0
		return err
	}
	return nil
}

// bitReader reads individual bits from a byte stream.
type bitReader struct {
	r     io.ByteReader
	buf   byte
	nbits int
}

func (br *bitReader) readBit() byte {
	if br.nbits == 0 {
		b, err := br.r.ReadByte()
		if err != nil {
			br.buf = 0
			br.nbits = 8
			return 0
		}
		br.buf = b
		br.nbits = 8
	}
	br.nbits--
	return (br.buf >> br.nbits) & 1
}

// Encoder writes arithmetic-coded symbols to a byte stream.
type Encoder struct {
	bw      bitWriter
	lo      uint64
	hi      uint64
	pending int
}

// NewEncoder creates an arithmetic encoder writing to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		bw: bitWriter{w: w},
		lo: 0,
		hi: maxRange,
	}
}

// Encode encodes symbol sym using the given distribution.
func (e *Encoder) Encode(sym int, dist *Distribution) error {
	if sym < 0 || sym >= len(dist.CumFreq)-1 {
		return fmt.Errorf("symbol %d out of range [0, %d)", sym, len(dist.CumFreq)-1)
	}

	rng := e.hi - e.lo + 1
	total := uint64(dist.Total)
	symLo := uint64(dist.CumFreq[sym])
	symHi := uint64(dist.CumFreq[sym+1])

	e.hi = e.lo + (rng*symHi)/total - 1
	e.lo = e.lo + (rng*symLo)/total

	return e.renormalize()
}

func (e *Encoder) renormalize() error {
	for {
		if e.hi < halfRange {
			if err := e.bw.writeBit(0); err != nil {
				return err
			}
			for e.pending > 0 {
				if err := e.bw.writeBit(1); err != nil {
					return err
				}
				e.pending--
			}
		} else if e.lo >= halfRange {
			if err := e.bw.writeBit(1); err != nil {
				return err
			}
			for e.pending > 0 {
				if err := e.bw.writeBit(0); err != nil {
					return err
				}
				e.pending--
			}
			e.lo -= halfRange
			e.hi -= halfRange
		} else if e.lo >= qtrRange && e.hi < 3*qtrRange {
			e.pending++
			e.lo -= qtrRange
			e.hi -= qtrRange
		} else {
			break
		}
		e.lo <<= 1
		e.hi = (e.hi << 1) | 1
	}
	return nil
}

// Finish flushes remaining state to produce a decodable stream.
func (e *Encoder) Finish() error {
	e.pending++
	if e.lo < qtrRange {
		if err := e.bw.writeBit(0); err != nil {
			return err
		}
		for e.pending > 0 {
			if err := e.bw.writeBit(1); err != nil {
				return err
			}
			e.pending--
		}
	} else {
		if err := e.bw.writeBit(1); err != nil {
			return err
		}
		for e.pending > 0 {
			if err := e.bw.writeBit(0); err != nil {
				return err
			}
			e.pending--
		}
	}
	return e.bw.flush()
}

// Decoder reads arithmetic-coded symbols from a byte stream.
type Decoder struct {
	br    bitReader
	lo    uint64
	hi    uint64
	value uint64
}

// NewDecoder creates an arithmetic decoder reading from r.
func NewDecoder(r io.ByteReader) *Decoder {
	d := &Decoder{
		br: bitReader{r: r},
		lo: 0,
		hi: maxRange,
	}
	for range precBits {
		d.value = (d.value << 1) | uint64(d.br.readBit())
	}
	return d
}

// Decode decodes one symbol using the given distribution.
func (d *Decoder) Decode(dist *Distribution) (int, error) {
	rng := d.hi - d.lo + 1
	total := uint64(dist.Total)

	offset := d.value - d.lo
	scaledValue := ((offset + 1) * total - 1) / rng

	// Linear search for the symbol
	sym := 0
	for sym < len(dist.CumFreq)-2 && uint64(dist.CumFreq[sym+1]) <= scaledValue {
		sym++
	}

	symLo := uint64(dist.CumFreq[sym])
	symHi := uint64(dist.CumFreq[sym+1])

	d.hi = d.lo + (rng*symHi)/total - 1
	d.lo = d.lo + (rng*symLo)/total

	// Renormalize (mirrors encoder)
	for {
		if d.hi < halfRange {
			// top bits both 0: shift out
		} else if d.lo >= halfRange {
			d.lo -= halfRange
			d.hi -= halfRange
			d.value -= halfRange
		} else if d.lo >= qtrRange && d.hi < 3*qtrRange {
			d.lo -= qtrRange
			d.hi -= qtrRange
			d.value -= qtrRange
		} else {
			break
		}
		d.lo <<= 1
		d.hi = (d.hi << 1) | 1
		d.value = (d.value << 1) | uint64(d.br.readBit())
	}

	return sym, nil
}
