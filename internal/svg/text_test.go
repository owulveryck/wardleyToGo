package svg

import (
	"bytes"
	"encoding/xml"
	"image"
	"regexp"
	"strconv"
	"testing"
)

func TestTextMarshalXML_LineSpacing(t *testing.T) {
	dyRe := regexp.MustCompile(`dy="(-?\d+)"`)

	tests := []struct {
		name   string
		py     int
		text   string
		minGap int
	}{
		{
			name:   "above candidate dy=-8",
			py:     -8,
			text:   "Collaboration Server",
			minGap: 18,
		},
		{
			name:   "above candidate dy=-20",
			py:     -20,
			text:   "Collaboration Server",
			minGap: 18,
		},
		{
			name:   "right candidate dy=0",
			py:     0,
			text:   "Collaboration Server",
			minGap: 18,
		},
		{
			name:   "below candidate dy=18",
			py:     18,
			text:   "Collaboration Server",
			minGap: 18,
		},
		{
			name:   "three-word label dy=-8",
			py:     -8,
			text:   "My Long Component Name",
			minGap: 18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txt := Text{
				P:          image.Pt(10, tt.py),
				Text:       []byte(tt.text),
				TextAdjust: true,
			}
			var buf bytes.Buffer
			enc := xml.NewEncoder(&buf)
			if err := enc.Encode(txt); err != nil {
				t.Fatal(err)
			}
			enc.Flush()

			matches := dyRe.FindAllStringSubmatch(buf.String(), -1)
			if len(matches) < 2 {
				t.Fatalf("expected at least 2 tspan elements, got %d", len(matches))
			}

			for i := 1; i < len(matches); i++ {
				gap, err := strconv.Atoi(matches[i][1])
				if err != nil {
					t.Fatal(err)
				}
				if gap < tt.minGap {
					t.Errorf("tspan[%d] dy=%d, want >= %d (text=%q, py=%d)", i, gap, tt.minGap, tt.text, tt.py)
				}
			}
		})
	}
}
