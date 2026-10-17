package match

import "github.com/Rakshit-gen/chartlink/record"

// PassStat is what one blocking pass contributes.
type PassStat struct {
	Name  string
	Pairs int // candidate pairs this pass proposes
	Found int // true pairs it proposes
	Only  int // true pairs no other pass proposes
}

// BlockingStats measures blocking against the truth (truth[i] is the
// person behind recs[i]): what each pass finds, and the share of true pairs
// any pass finds. A true pair no pass proposes can never be linked.
func BlockingStats(recs []record.Record, truth []int) (stats []PassStat, recall float64) {
	ps := prepAll(recs)
	by := map[uint64]int{} // true pair -> bitmask of passes that found it
	for i, p := range passes {
		st := PassStat{Name: p.name}
		blocks := map[string][]uint32{}
		for j := range ps {
			if k := p.key(&ps[j]); k != "" {
				blocks[k] = append(blocks[k], uint32(j))
			}
		}
		seen := map[uint64]bool{}
		for _, b := range blocks {
			if len(b) < 2 || len(b) > maxBlock {
				continue
			}
			for x, a := range b {
				for _, c := range b[x+1:] {
					k := pack(a, c)
					if seen[k] {
						continue
					}
					seen[k] = true
					st.Pairs++
					if truth[a] == truth[c] {
						st.Found++
						by[k] |= 1 << i
					}
				}
			}
		}
		stats = append(stats, st)
	}
	for _, mask := range by {
		for i := range stats {
			if mask == 1<<i {
				stats[i].Only++
			}
		}
	}
	people := map[int]int{}
	for _, t := range truth {
		people[t]++
	}
	total := 0
	for _, n := range people {
		total += n * (n - 1) / 2
	}
	return stats, float64(len(by)) / float64(total)
}
