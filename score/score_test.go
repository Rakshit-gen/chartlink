package score

import (
	"math"
	"testing"

	"github.com/Rakshit-gen/chartlink/record"
	"github.com/Rakshit-gen/chartlink/synth"
)

// Three people: 0 and 1 are twins. Records 0-2 are person 0, 3-4 person
// 1, 5 person 2.
func data() synth.Data {
	return synth.Data{
		People: []synth.Person{
			{ID: 0, Kin: 1, Tags: synth.Twin},
			{ID: 1, Kin: 0, Tags: synth.Twin},
			{ID: 2, Kin: -1, Tags: synth.Moved},
		},
		Records: make([]record.Record, 6),
		Truth:   []int{0, 0, 0, 1, 1, 2},
	}
}

func TestPerfect(t *testing.T) {
	r := Score(data(), []int{7, 7, 7, 8, 8, 9})
	if r.Precision != 1 || r.Recall != 1 || r.Overlays != 0 || r.Duplicates != 0 || r.TwinsMerged != 0 {
		t.Errorf("%+v", r)
	}
}

func TestTwinsMerged(t *testing.T) {
	r := Score(data(), []int{1, 1, 1, 1, 1, 2})
	// 10 pairs predicted, 3 + 1 true.
	if math.Abs(r.Precision-0.4) > 1e-9 || r.Recall != 1 {
		t.Errorf("precision %v recall %v", r.Precision, r.Recall)
	}
	if r.Overlays != 1 || r.TwinsMerged != 1 || r.TwinPairs != 1 {
		t.Errorf("%+v", r)
	}
}

func TestSplit(t *testing.T) {
	r := Score(data(), []int{0, 1, 2, 3, 3, 4})
	if r.Precision != 1 || r.Recall != 0.25 || r.Duplicates != 2 || r.Overlays != 0 {
		t.Errorf("%+v", r)
	}
	if math.Abs(r.DuplicatePeople-1.0/3) > 1e-9 {
		t.Errorf("duplicate people %v", r.DuplicatePeople)
	}
}
