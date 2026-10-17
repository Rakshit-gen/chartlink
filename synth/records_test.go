package synth

import (
	"math/rand/v2"
	"testing"

	"github.com/Rakshit-gen/chartlink/record"
)

func TestRecords(t *testing.T) {
	ps := People(Config{People: 3000, Zips: 50, Today: 2026}, 3)
	d := Records(ps, 4)
	if len(d.Records) != len(d.Truth) {
		t.Fatal("records and truth differ in length")
	}
	perPerson := map[int]int{}
	exactSame := 0
	for i, r := range d.Records {
		if r.ID != i {
			t.Fatalf("record %d has ID %d", i, r.ID)
		}
		p := ps[d.Truth[i]]
		perPerson[p.ID]++
		if r.First == p.First && r.Last == p.Last && r.DOB == p.DOB {
			exactSame++
		}
	}
	if len(perPerson) != len(ps) {
		t.Errorf("%d of %d people have records", len(perPerson), len(ps))
	}
	mean := float64(len(d.Records)) / float64(len(ps))
	share := float64(exactSame) / float64(len(d.Records))
	t.Logf("%d records, %.2f per person, %.0f%% with first, last and DOB exactly right", len(d.Records), mean, 100*share)
	if mean < 1.8 || mean > 2.8 {
		t.Errorf("%.2f records per person", mean)
	}
	// Most records are clean; enough are not to make exact matching miss.
	if share < 0.6 || share > 0.9 {
		t.Errorf("%.0f%% of records exactly right", 100*share)
	}
}

func TestTypo(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 1))
	for range 1000 {
		got := typo(r, "GARCIA")
		if d := len(got) - 6; d < -1 || d > 1 {
			t.Fatalf("typo changed length by %d: %q", d, got)
		}
		if record.JaroWinkler(got, "GARCIA") < 0.8 {
			t.Errorf("typo %q is too far from GARCIA", got)
		}
	}
}

func TestWrongDateStaysADate(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 1000 {
		d := wrongDate(r, record.Date{Y: 1980, M: 3, D: 31})
		if d.M < 1 || d.M > 12 || d.D < 1 || d.D > 31 {
			t.Fatalf("made %+v", d)
		}
	}
}

func TestDigitTypo(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 3))
	for range 1000 {
		got := digitTypo(r, "1234")
		diff := 0
		for i := range 4 {
			if got[i] != "1234"[i] {
				diff++
			}
			if got[i] < '0' || got[i] > '9' {
				t.Fatalf("made %q", got)
			}
		}
		if diff < 1 || diff > 2 {
			t.Fatalf("made %q", got)
		}
	}
}
