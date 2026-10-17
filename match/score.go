package match

import "math"

// Idea is one part of chartlink beyond plain Fellegi-Sunter, so each can
// be switched off to see what it is worth.
type Idea uint8

const (
	// TermFreq weighs agreement on a name by how rare that name is.
	TermFreq Idea = 1 << iota
)

// scorer weighs pairs with the fitted model and whichever ideas are on.
type scorer struct {
	ps  []prepped
	m   *Model
	off Idea
	tf  [2]map[string]float64 // share of records with each first, last name
	n   float64
}

func newScorer(ps []prepped, m *Model, off Idea) *scorer {
	s := &scorer{ps: ps, m: m, off: off, n: float64(len(ps))}
	for i := range s.tf {
		s.tf[i] = map[string]float64{}
	}
	for i := range ps {
		s.tf[0][ps[i].first]++
		s.tf[1][ps[i].last]++
	}
	for i := range s.tf {
		for k, c := range s.tf[i] {
			s.tf[i][k] = c / s.n
		}
	}
	return s
}

func (s *scorer) on(i Idea) bool { return s.off&i == 0 }

// weight is the pair's evidence in bits: positive for the same person.
func (s *scorer) weight(a, b int, g Gamma) float64 {
	w := s.m.Weight(g)
	if s.on(TermFreq) {
		// Exact agreement on a name is worth log2(M/f), where f is the
		// share of records with that name, rather than log2(M/U) with U
		// averaged over all names. SMITH agreeing is weak; KOVALENKO is
		// strong.
		for i, f := range [2]Field{First, Last} {
			if g[f] == 0 {
				name := s.ps[a].first
				if f == Last {
					name = s.ps[a].last
				}
				w += math.Log2(s.m.U[f][0] / s.tf[i][name])
			}
		}
	}
	return w
}
