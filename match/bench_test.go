package match

import (
	"testing"

	"github.com/Rakshit-gen/chartlink/synth"
)

// BenchmarkIndexAdd times linking one record into an index that already
// holds a 20000-person population.
func BenchmarkIndexAdd(b *testing.B) {
	d := synth.Records(synth.People(synth.DefaultConfig(), 1), 2)
	cfg := DefaultConfig()
	cfg.Samples = 100000
	m := Fit(d.Records, cfg)
	half := len(d.Records) / 2
	var ix *Index
	b.ResetTimer()
	for i := range b.N {
		k := i % half
		if k == 0 {
			b.StopTimer()
			ix = NewIndex(m, cfg)
			for _, r := range d.Records[:half] {
				ix.Add(r)
			}
			b.StartTimer()
		}
		ix.Add(d.Records[half+k])
	}
}

func BenchmarkLink(b *testing.B) {
	d := synth.Records(synth.People(synth.DefaultConfig(), 1), 2)
	cfg := DefaultConfig()
	for range b.N {
		Link(d.Records, cfg)
	}
}
