package match

// estimateM fits M by expectation maximization over the candidate pairs.
// Nobody labels which pairs match, so EM alternates: given M and U, how
// likely is each pair a match; given those likelihoods, what are M and U.
//
// The U it fits along the way is for non-matching pairs inside blocks,
// which agree far more often than random pairs (they share a key). It is
// thrown away; weights use the U from random pairs, so a weight measures
// evidence against the whole population, not against the block.
//
// Returns the share of candidate pairs that match.
func (m *Model) estimateM(gs []Gamma, iters int) float64 {
	var mm, uu [Fields][]float64
	for f := range Fields {
		n := Levels(f)
		mm[f], uu[f] = make([]float64, n), make([]float64, n)
		// Start M leaning to agreement and U at the population's rates.
		for l := range n {
			mm[f][l] = 0.1 / float64(n-1)
			uu[f][l] = m.U[f][l]
		}
		mm[f][0] = 0.9
	}
	share := 0.1
	var cm, cu [Fields][]float64
	for f := range Fields {
		cm[f], cu[f] = make([]float64, Levels(f)), make([]float64, Levels(f))
	}
	for range iters {
		for f := range Fields {
			clear(cm[f])
			clear(cu[f])
		}
		total := 0.0
		for _, g := range gs {
			pm, pu := share, 1-share
			for f := range Fields {
				if l := g[f]; l != Missing {
					pm *= mm[f][l]
					pu *= uu[f][l]
				}
			}
			w := pm / (pm + pu)
			total += w
			for f := range Fields {
				if l := g[f]; l != Missing {
					cm[f][l] += w
					cu[f][l] += 1 - w
				}
			}
		}
		share = total / float64(len(gs))
		for f := range Fields {
			// A small floor keeps a level nobody hit from getting M zero,
			// which would make any pair landing there impossible.
			for l := range cm[f] {
				cm[f][l] += 0.5
				cu[f][l] += 0.5
			}
			mm[f], uu[f] = normalize(cm[f]), normalize(cu[f])
		}
	}
	m.M = mm
	return share
}
