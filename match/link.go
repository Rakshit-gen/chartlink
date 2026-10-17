package match

import (
	"cmp"
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
	// Held means the pair would have merged but for the family rules: it
	// looks like twins, or a father and son. Someone who goes by a middle
	// name looks like a twin too, so these always go to review.
	Held bool
}

func (sc *scorer) pair(a, b int, g Gamma, cfg *Config) Pair {
	w, raw := sc.weight(a, b, g)
	p := Pair{A: a, B: b, Weight: w, Prob: sc.m.Prob(w)}
	p.Held = p.Prob < cfg.Merge && sc.m.Prob(raw) >= cfg.Merge
	return p
}

// place merges or queues one pair, strongest pairs first.
func place(ch *charts, sc *scorer, p Pair, cfg *Config, review []Pair) []Pair {
	ra, rb := ch.find(p.A), ch.find(p.B)
	if ra == rb {
		return review
	}
	switch {
	case p.Prob >= cfg.Merge && !(sc.on(Guard) && ch.clash(ra, rb, sc.conflict)):
		ch.union(ra, rb)
	case p.Prob >= cfg.Review || p.Held:
		// Includes strong pairs the guard stopped: a person should see
		// why two charts that look alike can't be merged.
		review = append(review, p)
	}
	return review
}

// Result is the linker's output.
type Result struct {
	Chart      []int  // Chart[i] is record i's chart
	Review     []Pair // pairs a person should look at
	Model      Model
	Candidates int
}

// Fit fits a model to a set of records without linking them, for an Index
// to use on records that arrive later.
func Fit(recs []record.Record, cfg Config) Model {
	m, _, _ := fit(prepAll(recs), cfg)
	return m
}

func prepAll(recs []record.Record) []prepped {
	ps := make([]prepped, len(recs))
	for i, r := range recs {
		ps[i] = prep(r)
	}
	return ps
}

func fit(ps []prepped, cfg Config) (Model, []uint64, []Gamma) {
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
	m.countNames(ps)
	return m, cands, gs
}

// Link groups records into charts, one per person.
func Link(recs []record.Record, cfg Config) Result {
	ps := prepAll(recs)
	m, cands, gs := fit(ps, cfg)
	sc := &scorer{ps: ps, m: &m, off: cfg.Without}
	pairs := make([]Pair, len(cands))
	for i, c := range cands {
		a, b := unpack(c)
		pairs[i] = sc.pair(a, b, gs[i], &cfg)
	}
	// Strongest evidence first, so a chart forms around its surest pairs.
	slices.SortFunc(pairs, func(x, y Pair) int { return cmp.Compare(y.Weight, x.Weight) })

	ch := newCharts(len(ps))
	res := Result{Model: m, Candidates: len(cands)}
	for _, p := range pairs {
		res.Review = place(ch, sc, p, &cfg, res.Review)
	}
	res.Review = settle(ch, res.Review)
	res.Chart = ch.labels()
	return res
}

// settle drops review pairs whose records ended up in one chart anyway, and
// keeps one pair per two charts, the strongest: a reviewer decides about
// charts, not about every pair of records in them.
func settle(ch *charts, review []Pair) []Pair {
	seen := map[[2]int]bool{}
	out := review[:0]
	for _, p := range review { // already strongest first
		a, b := ch.find(p.A), ch.find(p.B)
		if a == b {
			continue
		}
		k := [2]int{min(a, b), max(a, b)}
		if !seen[k] {
			seen[k] = true
			out = append(out, p)
		}
	}
	return out
}
