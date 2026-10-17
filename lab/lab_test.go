package lab

import (
	"testing"

	"github.com/Rakshit-gen/chartlink/match"
	"github.com/Rakshit-gen/chartlink/record"
	"github.com/Rakshit-gen/chartlink/synth"
)

// Person 0 has records 0-2, person 1 has 3. The contender split person 0
// and flagged both of its charts, and person 0 against person 1, for
// review. The reviewer should merge only the first.
func TestRunReviewer(t *testing.T) {
	d := synth.Data{
		People:  []synth.Person{{ID: 0, Kin: -1}, {ID: 1, Kin: -1}},
		Records: make([]record.Record, 4),
		Truth:   []int{0, 0, 0, 1},
	}
	c := Contender{"split", func(synth.Data) ([]int, []match.Pair) {
		return []int{10, 10, 11, 12}, []match.Pair{{A: 0, B: 2}, {A: 2, B: 3}}
	}}
	o := Run(c, d)
	if o.Before.Duplicates != 1 || o.After.Duplicates != 0 || o.After.Overlays != 0 {
		t.Errorf("before %+v\nafter %+v", o.Before, o.After)
	}
	if o.Reviews != 2 || o.ReviewsSame != 1 {
		t.Errorf("reviews %d, same %d", o.Reviews, o.ReviewsSame)
	}
}
