package match

import (
	"math/rand/v2"
	"slices"

	"github.com/Rakshit-gen/chartlink/record"
)

// Config sets the linker's thresholds.
type Config struct {
	// Merge is the match probability above which two charts are joined.
	Merge float64
	// Review is the probability above which a pair that wasn't merged is
	// sent to a person to look at.
	Review float64
	// Samples is how many random pairs U is estimated from; Iters is how
	// many rounds of EM fit M.
	Samples, Iters int
	Seed           uint64
	// Without switches ideas off, for the ablation and the baseline.
	Without Idea
}

// DefaultConfig is what the README's numbers come from.
func DefaultConfig() Config {
	return Config{Merge: 0.99, Review: 0.5, Samples: 500000, Iters: 30, Seed: 1}
}

// Pair is two records and the linker's opinion of them.
type Pair struct {
	A, B   int
	Weight float64
	Prob   float64
}

// Result is the linker's output.
type Result struct {
	Chart      []int  // Chart[i] is record i's chart
	Review     []Pair // pairs a person should look at
	Model      Model
	Candidates int
}

// Link groups records into charts, one per person.
func Link(recs []record.Record, cfg Config) Result {
	ps := make([]prepped, len(recs))
	for i, r := range recs {
		ps[i] = prep(r)
	}
	cands := candidates(ps)
	gs := make([]Gamma, len(cands))
	for i, c := range cands {
		a, b := unpack(c)
		gs[i] = compare(&ps[a], &ps[b])
	}

	var m Model
	m.estimateU(ps, cfg.Samples, rand.New(rand.NewPCG(cfg.Seed, 1)))
	share := m.estimateM(gs, cfg.Iters)
	n := float64(len(ps))
	m.Prior = share * float64(len(cands)) / (n * (n - 1) / 2)

	sc := newScorer(ps, &m, cfg.Without)
	pairs := make([]Pair, len(cands))
	for i, c := range cands {
		a, b := unpack(c)
		w := sc.weight(a, b, gs[i])
		pairs[i] = Pair{A: a, B: b, Weight: w, Prob: m.Prob(w)}
	}
	// Strongest evidence first, so a chart forms around its surest pairs.
	slices.SortFunc(pairs, func(x, y Pair) int {
		switch {
		case x.Weight > y.Weight:
			return -1
		case x.Weight < y.Weight:
			return 1
		}
		return 0
	})

	ch := newCharts(len(ps))
	res := Result{Model: m, Candidates: len(cands)}
	for _, p := range pairs {
		ra, rb := ch.find(p.A), ch.find(p.B)
		if ra == rb {
			continue
		}
		switch {
		case p.Prob >= cfg.Merge && !(sc.on(Guard) && ch.clash(ra, rb, sc.conflict)):
			ch.union(ra, rb)
		case p.Prob >= cfg.Review:
			res.Review = append(res.Review, p)
		}
	}
	// A pair queued for review whose records were merged later through
	// other pairs is settled; drop it.
	res.Review = slices.DeleteFunc(res.Review, func(p Pair) bool { return ch.find(p.A) == ch.find(p.B) })
	res.Chart = ch.labels()
	return res
}
