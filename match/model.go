package match

import (
	"math"
	"math/rand/v2"
)

// Model is a Fellegi-Sunter model: for every field and level, how likely a
// pair is to land there if the two records are the same person (M) and if
// they are two random records (U). A pair's weight is the sum over fields
// of log2(M/U), the log of how much more likely its comparison vector is
// for a match than for a random pair.
type Model struct {
	M, U [Fields][]float64
	// Prior is the chance that two random records are the same person.
	Prior float64
	// Names holds, for first and last names, the share of records with
	// each name. A name not in the table is taken to be as rare as one
	// record in the set the model was fitted on (NameFloor).
	Names     map[Field]map[string]float64
	NameFloor float64
}

// countNames fills Names from the records the model is fitted on.
func (m *Model) countNames(ps []prepped) {
	m.Names = map[Field]map[string]float64{First: {}, Last: {}}
	for i := range ps {
		m.Names[First][ps[i].first]++
		m.Names[Last][ps[i].last]++
	}
	n := float64(len(ps))
	for _, t := range m.Names {
		for k, c := range t {
			t[k] = c / n
		}
	}
	m.NameFloor = 1 / n
}

// Weight is the pair's log2 Bayes factor. Missing fields add nothing.
func (m *Model) Weight(g Gamma) float64 {
	w := 0.0
	for f := range Fields {
		if l := g[f]; l != Missing {
			w += math.Log2(m.M[f][l] / m.U[f][l])
		}
	}
	return w
}

// Prob turns a weight into the chance the pair is the same person.
func (m *Model) Prob(w float64) float64 {
	odds := m.Prior / (1 - m.Prior) * math.Exp2(w)
	return odds / (1 + odds)
}

// WeightFor is the weight that gives probability p.
func (m *Model) WeightFor(p float64) float64 {
	return math.Log2(p/(1-p)) - math.Log2(m.Prior/(1-m.Prior))
}

// estimateU fills U from random pairs of records. Almost every random pair
// is two different people, so how often random pairs agree is how often
// non-matches agree. Each level starts at a count of one so no level has
// U zero.
func (m *Model) estimateU(ps []prepped, samples int, r *rand.Rand) {
	var counts [Fields][]float64
	for f := range Fields {
		counts[f] = make([]float64, Levels(f))
		for l := range counts[f] {
			counts[f][l] = 1
		}
	}
	n := len(ps)
	for range samples {
		a, b := r.IntN(n), r.IntN(n)
		if a == b {
			continue
		}
		g := compare(&ps[a], &ps[b])
		for f := range Fields {
			if l := g[f]; l != Missing {
				counts[f][l]++
			}
		}
	}
	for f := range Fields {
		m.U[f] = normalize(counts[f])
	}
}

func normalize(c []float64) []float64 {
	t := 0.0
	for _, x := range c {
		t += x
	}
	out := make([]float64, len(c))
	for i, x := range c {
		out[i] = x / t
	}
	return out
}
