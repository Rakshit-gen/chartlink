// Package score measures a clustering of records against the truth, in the
// terms a health information department uses: duplicates (one patient,
// several charts) and overlays (one chart, several patients).
package score

import "github.com/Rakshit-gen/chartlink/synth"

// Report is how one clustering did.
type Report struct {
	Records, People int

	// Pairwise: of all record pairs put in the same cluster, how many are
	// the same person (precision), and of all same-person pairs, how many
	// ended up together (recall).
	Precision, Recall, F1 float64

	// Overlays is the number of extra people merged into charts: a chart
	// holding three people's records counts two. Every one is a patient
	// whose allergies or results can show up on someone else's chart.
	Overlays int
	// Duplicates is the number of extra charts: a patient whose records
	// sit in three charts counts two.
	Duplicates int
	// DuplicatePeople is the share of people with more than one chart.
	DuplicatePeople float64

	// Hard cases. Twins and Jr/Sr are pairs that must stay apart; married
	// and moved are people whose records must come together.
	TwinsMerged, TwinPairs  int
	JrSrMerged, JrSrPairs   int
	MarriedSplit, Married   int
	MovedSplit, Moved       int
	ByMiddleSplit, ByMiddle int
	Charts                  int
}

func pairs(n int) float64 { return float64(n) * float64(n-1) / 2 }

// Score compares cluster labels (one per record) with the truth.
func Score(d synth.Data, cluster []int) Report {
	rep := Report{Records: len(d.Records), People: len(d.People)}
	type cell struct{ c, p int }
	cells := map[cell]int{}
	byCluster := map[int]int{}
	byPerson := make([]int, len(d.People))
	people := map[int]map[int]bool{}              // cluster -> people in it
	charts := make([]map[int]bool, len(d.People)) // person -> clusters
	for i, c := range cluster {
		p := d.Truth[i]
		cells[cell{c, p}]++
		byCluster[c]++
		byPerson[p]++
		if people[c] == nil {
			people[c] = map[int]bool{}
		}
		people[c][p] = true
		if charts[p] == nil {
			charts[p] = map[int]bool{}
		}
		charts[p][c] = true
	}
	var tp, predicted, actual float64
	for _, n := range cells {
		tp += pairs(n)
	}
	for _, n := range byCluster {
		predicted += pairs(n)
	}
	for _, n := range byPerson {
		actual += pairs(n)
	}
	rep.Precision, rep.Recall = 1, 1
	if predicted > 0 {
		rep.Precision = tp / predicted
	}
	if actual > 0 {
		rep.Recall = tp / actual
	}
	if rep.Precision+rep.Recall > 0 {
		rep.F1 = 2 * rep.Precision * rep.Recall / (rep.Precision + rep.Recall)
	}
	rep.Charts = len(byCluster)
	for _, ps := range people {
		rep.Overlays += len(ps) - 1
	}
	dupPeople := 0
	for _, cs := range charts {
		if len(cs) > 1 {
			rep.Duplicates += len(cs) - 1
			dupPeople++
		}
	}
	rep.DuplicatePeople = float64(dupPeople) / float64(len(d.People))

	together := func(a, b int) bool {
		for c := range charts[a] {
			if charts[b][c] {
				return true
			}
		}
		return false
	}
	for _, p := range d.People {
		if p.Kin > p.ID { // count each pair once
			if p.Has(synth.Twin) {
				rep.TwinPairs++
				if together(p.ID, p.Kin) {
					rep.TwinsMerged++
				}
			}
			if p.Has(synth.JrSr) {
				rep.JrSrPairs++
				if together(p.ID, p.Kin) {
					rep.JrSrMerged++
				}
			}
		}
		split := len(charts[p.ID]) > 1
		if p.Has(synth.Married) {
			rep.Married++
			if split {
				rep.MarriedSplit++
			}
		}
		if p.Has(synth.ByMiddle) {
			rep.ByMiddle++
			if split {
				rep.ByMiddleSplit++
			}
		}
		if p.Has(synth.Moved) {
			rep.Moved++
			if split {
				rep.MovedSplit++
			}
		}
	}
	return rep
}
