package match

import (
	"testing"

	"github.com/Rakshit-gen/chartlink/record"
)

func pair(a, b record.Record) (*prepped, *prepped, Gamma) {
	pa, pb := prep(a), prep(b)
	return &pa, &pb, compare(&pa, &pb)
}

func TestTwins(t *testing.T) {
	d := record.Date{Y: 2001, M: 6, D: 9}
	for _, c := range []struct {
		a, b string
		dob  record.Date
		want bool
	}{
		{"MARK", "LUKE", d, true},
		{"MARK", "LUKE", record.Date{Y: 2001, M: 6, D: 10}, true}, // DOB typo
		{"MARK", "LUKE", record.Date{Y: 1975, M: 2, D: 1}, false}, // not twins: different birthdays
		{"RYAN", "RYSN", d, false},                                // a typo
		{"LIAS", "LSIA", d, false},                                // two typos of LISA
		{"HUCK", "CHARLIE", d, false},                             // a typo of CHUCK
		{"ROBERT", "BOB", d, false},
	} {
		pa, pb, g := pair(record.Record{First: c.a, Last: "KIM", DOB: d}, record.Record{First: c.b, Last: "KIM", DOB: c.dob})
		if got := twins(pa, pb, g); got != c.want {
			t.Errorf("%s/%s: twins = %v", c.a, c.b, got)
		}
	}
}

func TestGenerations(t *testing.T) {
	dad := record.Record{First: "JOHN", Last: "KIM", DOB: record.Date{Y: 1960, M: 3, D: 1}}
	for _, c := range []struct {
		name string
		son  record.Record
		want bool
	}{
		{"born 30 years later", record.Record{First: "JOHN", Last: "KIM", DOB: record.Date{Y: 1990, M: 8, D: 20}}, true},
		{"only the son has a suffix", record.Record{First: "JOHN", Last: "KIM", Suffix: "JR"}, false}, // dad has none
		{"only the year off", record.Record{First: "JOHN", Last: "KIM", DOB: record.Date{Y: 2026, M: 3, D: 1}}, false},
		{"a decade typo", record.Record{First: "JOHN", Last: "KIM", DOB: record.Date{Y: 1970, M: 3, D: 2}}, false},
	} {
		pa, pb, g := pair(dad, c.son)
		if got := generations(pa, pb, g); got != c.want {
			t.Errorf("%s: generations = %v", c.name, got)
		}
	}
	sr := dad
	sr.Suffix = "SR"
	pa, pb, g := pair(sr, record.Record{First: "JOHN", Last: "KIM", Suffix: "JR"})
	if !generations(pa, pb, g) {
		t.Error("SR and JR should be two people")
	}
}
