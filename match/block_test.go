package match

import (
	"testing"

	"github.com/Rakshit-gen/chartlink/record"
)

func TestCandidates(t *testing.T) {
	d := record.Date{Y: 1980, M: 4, D: 2}
	recs := []record.Record{
		{First: "ANNA", Last: "SMITH", DOB: d},
		{First: "ANNA", Last: "SMYTH", DOB: d},                      // typo: same sound, same year
		{First: "ANNA", Last: "JONES", DOB: d},                      // married: same DOB and initial
		{First: "ANA", Last: "JONES", Phone: "5551112222"},          // no DOB: only the phone links it
		{First: "ZED", Last: "ZULU", Phone: "5551112222", Zip: "1"}, // shares the phone
		{First: "OMAR", Last: "ALI", DOB: record.Date{Y: 1950, M: 1, D: 1}},
	}
	ps := make([]prepped, len(recs))
	for i, r := range recs {
		ps[i] = prep(r)
	}
	got := map[[2]int]bool{}
	for _, p := range candidates(ps) {
		a, b := unpack(p)
		if a >= b {
			t.Errorf("pair %d,%d out of order", a, b)
		}
		got[[2]int{a, b}] = true
	}
	for _, want := range [][2]int{{0, 1}, {0, 2}, {1, 2}, {3, 4}} {
		if !got[want] {
			t.Errorf("missing pair %v", want)
		}
	}
	for p := range got {
		if p[0] == 5 || p[1] == 5 {
			t.Errorf("record 5 shares nothing, but got %v", p)
		}
	}
	if len(got) != 4 {
		t.Errorf("got %d pairs: %v", len(got), got)
	}
}
