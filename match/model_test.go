package match

import (
	"math"
	"testing"
)

func TestWeightAndProb(t *testing.T) {
	var m Model
	for f := range Fields {
		n := Levels(f)
		m.M[f], m.U[f] = make([]float64, n), make([]float64, n)
		for l := range n {
			m.M[f][l], m.U[f][l] = 1, 1
		}
	}
	m.M[First][0], m.U[First][0] = 0.9, 0.01
	m.Prior = 0.001
	var g Gamma
	for f := range Fields {
		g[f] = Missing
	}
	if w := m.Weight(g); w != 0 {
		t.Errorf("all missing should weigh 0, got %v", w)
	}
	g[First] = 0
	if w := m.Weight(g); math.Abs(w-math.Log2(90)) > 1e-9 {
		t.Errorf("weight %v", w)
	}
	for _, p := range []float64{0.01, 0.5, 0.99} {
		if got := m.Prob(m.WeightFor(p)); math.Abs(got-p) > 1e-9 {
			t.Errorf("Prob(WeightFor(%v)) = %v", p, got)
		}
	}
}
