// Package lab runs the experiments behind results/: every matcher on the
// same synthetic populations, scored the same way.
package lab

import (
	"github.com/Rakshit-gen/chartlink/baseline"
	"github.com/Rakshit-gen/chartlink/match"
	"github.com/Rakshit-gen/chartlink/score"
	"github.com/Rakshit-gen/chartlink/synth"
)

// Contender is one way of linking records.
type Contender struct {
	Name string
	Run  func(d synth.Data) (chart []int, review []match.Pair)
}

// All switches every idea off: plain Fellegi-Sunter with EM.
const All = match.TermFreq | match.Household | match.Guard

// Linker is chartlink with some ideas switched off.
func Linker(name string, cfg match.Config) Contender {
	return Contender{name, func(d synth.Data) ([]int, []match.Pair) {
		r := match.Link(d.Records, cfg)
		return r.Chart, r.Review
	}}
}

// Contenders are the three matchers the README compares.
func Contenders(cfg match.Config) []Contender {
	plain := cfg
	plain.Without = All
	return []Contender{
		{"exact first, last, DOB", func(d synth.Data) ([]int, []match.Pair) { return baseline.Exact(d.Records), nil }},
		Linker("Fellegi-Sunter with EM", plain),
		Linker("chartlink", cfg),
	}
}

// Ablations are chartlink with one idea off at a time.
func Ablations(cfg match.Config) []Contender {
	out := []Contender{Linker("chartlink", cfg)}
	for _, a := range []struct {
		name string
		off  match.Idea
	}{
		{"names not weighted by rarity", match.TermFreq},
		{"household fields counted separately", match.Household},
		{"no chart guard (twin, Jr/Sr and SSN conflicts)", match.Guard},
	} {
		c := cfg
		c.Without = a.off
		out = append(out, Linker(a.name, c))
	}
	return out
}

// Data is one synthetic population and its records.
func Data(seed uint64) synth.Data {
	return synth.Records(synth.People(synth.DefaultConfig(), seed), seed+1000)
}

// Outcome is one contender on one population.
type Outcome struct {
	Before, After score.Report
	Reviews       int // chart pairs a person had to look at
	ReviewsSame   int // of those, how many were one person
}

// Run links d and has a reviewer who is always right work through the
// review queue: a pair of charts is merged when most of both charts' records
// are one person. After is the clustering that leaves.
func Run(c Contender, d synth.Data) Outcome {
	chart, review := c.Run(d)
	o := Outcome{Before: score.Score(d, chart), Reviews: len(review)}
	who := majority(d, chart)
	after := append([]int(nil), chart...)
	uf := map[int]int{}
	var root func(int) int
	root = func(x int) int {
		if p, ok := uf[x]; ok && p != x {
			r := root(p)
			uf[x] = r
			return r
		}
		return x
	}
	for _, p := range review {
		a, b := chart[p.A], chart[p.B]
		if who[a] == who[b] {
			o.ReviewsSame++
			if ra, rb := root(a), root(b); ra != rb {
				uf[ra] = rb
			}
		}
	}
	for i := range after {
		after[i] = root(chart[i])
	}
	o.After = score.Score(d, after)
	return o
}

// majority maps each chart to the person most of its records belong to.
func majority(d synth.Data, chart []int) map[int]int {
	count := map[[2]int]int{}
	for i, c := range chart {
		count[[2]int{c, d.Truth[i]}]++
	}
	best := map[int]int{}
	n := map[int]int{}
	for k, v := range count {
		if v > n[k[0]] || (v == n[k[0]] && k[1] < best[k[0]]) {
			n[k[0]], best[k[0]] = v, k[1]
		}
	}
	return best
}
