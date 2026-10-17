package match

import (
	"math"

	"github.com/Rakshit-gen/chartlink/record"
)

// Idea is one part of chartlink beyond plain Fellegi-Sunter, so each can
// be switched off to see what it is worth.
type Idea uint8

const (
	// TermFreq weighs agreement on a name by how rare that name is.
	TermFreq Idea = 1 << iota
	// Household counts agreement on street, zip and phone once. Everyone
	// in a household shares them, so they are one piece of evidence, not
	// three.
	Household
	// Family drops the evidence relatives share when a pair looks like
	// twins or a father and son.
	Family
)

// scorer weighs pairs with the fitted model and whichever ideas are on.
type scorer struct {
	ps  []prepped
	m   *Model
	off Idea
	tf  [2]map[string]float64 // share of records with each first, last name
}

func newScorer(ps []prepped, m *Model, off Idea) *scorer {
	s := &scorer{ps: ps, m: m, off: off}
	for i := range s.tf {
		s.tf[i] = map[string]float64{}
	}
	for i := range ps {
		s.tf[0][ps[i].first]++
		s.tf[1][ps[i].last]++
	}
	n := float64(len(ps))
	for i := range s.tf {
		for k, c := range s.tf[i] {
			s.tf[i][k] = c / n
		}
	}
	return s
}

func (s *scorer) on(i Idea) bool { return s.off&i == 0 }

var household = [3]Field{Street, Zip, Phone}

// weight is the pair's evidence in bits: positive for the same person.
func (s *scorer) weight(a, b int, g Gamma) float64 {
	var c [Fields]float64
	for f := range Fields {
		if l := g[f]; l != Missing {
			c[f] = math.Log2(s.m.M[f][l] / s.m.U[f][l])
		}
	}
	pa, pb := &s.ps[a], &s.ps[b]
	if s.on(TermFreq) {
		// Exact agreement on a name is worth log2(M/f), where f is the
		// share of records with that name, rather than log2(M/U) with U
		// averaged over all names. SMITH agreeing is weak; KOVALENKO is
		// strong.
		if g[First] == 0 {
			c[First] = math.Log2(s.m.M[First][0] / s.tf[0][pa.first])
		}
		if g[Last] == 0 {
			c[Last] = math.Log2(s.m.M[Last][0] / s.tf[1][pa.last])
		}
	}
	if s.on(Household) {
		// Keep only the largest positive household weight. Disagreements
		// still count in full.
		best := Street
		for _, f := range household {
			if c[f] > c[best] {
				best = f
			}
		}
		for _, f := range household {
			if f != best && c[f] > 0 {
				c[f] = 0
			}
		}
	}
	if s.on(Family) {
		var shared []Field
		switch {
		case twins(pa, pb, g):
			// Twins share a last name, a birthday and a home. With the
			// first names different, those agreeing is what twins look
			// like, not evidence of one person.
			shared = []Field{Last, DOB, Street, Zip, Phone}
		case generations(pa, pb, g):
			// A father and son share every name and a home.
			shared = []Field{First, Middle, Last, Suffix, Street, Zip, Phone}
		}
		for _, f := range shared {
			c[f] = min(c[f], 0)
		}
	}
	w := 0.0
	for _, x := range c {
		w += x
	}
	return w
}

// twins: same date of birth and first names that are not one keystroke
// apart, even through nicknames. RYAN and RYSN is a typo; RICHARD and ROCH
// is a typo of RICH; MARK and LUKE are twins.
func twins(a, b *prepped, g Gamma) bool {
	if g[DOB] != 0 || g[First] != 3 {
		return false
	}
	for _, x := range record.Variants(a.first) {
		for _, y := range record.Variants(b.first) {
			if record.OSA(x, y) <= 1 {
				return false
			}
		}
	}
	return true
}

// generationGap is how many years apart two birth years must be before a
// pair with the same name looks like a parent and child rather than a typo.
// The DOB mistakes registration desks make move a year by 1 or 10.
const generationGap = 12

// generations: JR against SR, or birth years a generation apart on a
// different day of the year. When only the year is off, the likelier story
// is a mistyped year (some desks type the current one by habit); a father
// and son share a birthday about one time in 365.
func generations(a, b *prepped, g Gamma) bool {
	if g[Suffix] == 1 {
		return true
	}
	if a.dob.Y == 0 || b.dob.Y == 0 || (a.dob.M == b.dob.M && a.dob.D == b.dob.D) {
		return false
	}
	gap := a.dob.Y - b.dob.Y
	return max(gap, -gap) >= generationGap
}
