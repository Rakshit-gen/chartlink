package baseline

import (
	"testing"

	"github.com/Rakshit-gen/chartlink/record"
)

func TestExact(t *testing.T) {
	d := record.Date{Y: 1980, M: 4, D: 2}
	recs := []record.Record{
		{First: "Robert", Last: "Smith", DOB: d},
		{First: "ROBERT", Last: "smith", DOB: d},
		{First: "BOB", Last: "SMITH", DOB: d}, // nickname: exact can't see it
		{First: "ROBERT", Last: "SMITH"},      // no DOB
		{First: "ROBERT", Last: "SMITH"},      // no DOB either; still not merged
	}
	got := Exact(recs)
	if got[0] != got[1] {
		t.Error("case and spacing should not matter")
	}
	for i := 2; i < len(got); i++ {
		for j := range i {
			if got[i] == got[j] {
				t.Errorf("records %d and %d merged", i, j)
			}
		}
	}
}
