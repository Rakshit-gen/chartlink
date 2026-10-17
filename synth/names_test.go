package synth

import (
	"math/rand/v2"
	"testing"
)

// The head has to be common enough to collide and the tail rare enough to
// be informative; that difference is what term-frequency weighting uses.
func TestSurnameFrequencies(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	n := 100000
	count := map[string]int{}
	for range n {
		count[surname(r)]++
	}
	smith := float64(count["SMITH"]) / float64(n)
	if smith < 0.01 || smith > 0.04 {
		t.Errorf("SMITH is %.1f%% of people", 100*smith)
	}
	once := 0
	for _, c := range count {
		if c == 1 {
			once++
		}
	}
	if share := float64(once) / float64(n); share < 0.1 {
		t.Errorf("only %.1f%% of people have a surname nobody else has", 100*share)
	}
}
