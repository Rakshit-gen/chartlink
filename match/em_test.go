package match

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Pairs drawn from known M and U, mixed 30/70: EM should get M back.
func TestEMRecoversM(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 1))
	var truthM, truthU [Fields][]float64
	for f := range Fields {
		n := Levels(f)
		truthM[f], truthU[f] = make([]float64, n), make([]float64, n)
		truthM[f][0], truthU[f][n-1] = 0.85, 0.7
		for l := 1; l < n; l++ {
			truthM[f][l] = 0.15 / float64(n-1)
		}
		for l := range n - 1 {
			truthU[f][l] = 0.3 / float64(n-1)
		}
	}
	draw := func(p []float64) int8 {
		u := r.Float64()
		for l, x := range p {
			if u < x {
				return int8(l)
			}
			u -= x
		}
		return int8(len(p) - 1)
	}
	gs := make([]Gamma, 20000)
	for i := range gs {
		src := truthU
		if r.Float64() < 0.3 {
			src = truthM
		}
		for f := range Fields {
			gs[i][f] = draw(src[f])
		}
	}
	m := Model{U: truthU}
	share := m.estimateM(gs, 50)
	if math.Abs(share-0.3) > 0.02 {
		t.Errorf("match share %.3f, want 0.3", share)
	}
	for f := range Fields {
		if math.Abs(m.M[f][0]-0.85) > 0.03 {
			t.Errorf("%s: M exact %.3f, want 0.85", f, m.M[f][0])
		}
	}
}
