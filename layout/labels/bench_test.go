package labels

import (
	"image"
	"testing"
)

// buildBenchComps creates n components spread across the 100x100 map space.
func buildBenchComps(n int) []Component {
	comps := make([]Component, n)
	for i := range comps {
		comps[i] = Component{
			Name:     "Component " + string(rune('A'+i%26)),
			Position: image.Pt(10+80*i/(n+1), 5+90*i/(n+1)),
			Label:    "Component " + string(rune('A'+i%26)),
		}
	}
	return comps
}

func BenchmarkPlaceLabels(b *testing.B) {
	for _, n := range []int{15, 30, 50} {
		comps := buildBenchComps(n)
		opts := DefaultOptions()
		b.Run("N"+string(rune('0'+n/10))+string(rune('0'+n%10)), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = PlaceLabels(comps, opts)
			}
		})
	}
}

func BenchmarkSplitString(b *testing.B) {
	labels := []string{
		"Moteur de Calcul d'Itinéraire",
		"DB",
		"Infrastructure Cloud Provider",
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, l := range labels {
			_ = splitString(l, 8)
		}
	}
}
