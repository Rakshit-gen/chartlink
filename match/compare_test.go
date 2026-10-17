package match

import (
	"testing"

	"github.com/Rakshit-gen/chartlink/record"
)

func TestCompare(t *testing.T) {
	a := prep(record.Record{First: "Robert", Last: "Garcia", DOB: record.Date{Y: 1970, M: 3, D: 4}, Sex: 'M',
		Phone: "5551234567", Street: "12 Oak St", Zip: "10001", SSN4: "1234", Middle: "J"})
	b := prep(record.Record{First: "Bob", Last: "Gracia", DOB: record.Date{Y: 1970, M: 4, D: 3}, Sex: 'M',
		Street: "12 Oak Street", Zip: "10002", SSN4: "1243", Suffix: "JR"})
	g := compare(&a, &b)
	want := Gamma{First: 1, Last: 1, DOB: 1, Sex: 0, Phone: Missing, Street: 0, Zip: 1, SSN: 1, Middle: Missing, Suffix: Missing}
	if g != want {
		for f := range Fields {
			if g[f] != want[f] {
				t.Errorf("%s: got %s, want %s", f, LevelName(f, int(g[f])), LevelName(f, int(want[f])))
			}
		}
	}
}

func TestDates(t *testing.T) {
	d := func(y, m, dd int) record.Date { return record.Date{Y: y, M: m, D: dd} }
	for _, c := range []struct {
		a, b record.Date
		want int8
	}{
		{d(1980, 5, 6), d(1980, 5, 6), 0},
		{d(1980, 5, 6), d(1980, 6, 5), 1},
		{d(1980, 5, 6), d(1980, 5, 7), 2},
		{d(1980, 5, 6), d(1990, 5, 6), 2},
		{d(1980, 5, 6), d(2005, 9, 6), 3},
		{d(1980, 5, 6), record.Date{}, Missing},
	} {
		if got := dates(c.a, c.b); got != c.want {
			t.Errorf("%v vs %v: %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestSSN(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want int8
	}{
		{"1234", "1234", 0}, {"1234", "1235", 1}, {"1234", "2134", 1}, {"1234", "1243", 1},
		{"1234", "4321", 2}, {"1234", "1256", 2}, {"1234", "", Missing},
	} {
		if got := ssns(c.a, c.b); got != c.want {
			t.Errorf("%s vs %s: %d, want %d", c.a, c.b, got, c.want)
		}
	}
}
