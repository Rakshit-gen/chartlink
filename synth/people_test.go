package synth

import (
	"reflect"
	"testing"
)

func TestPeopleDeterministic(t *testing.T) {
	cfg := Config{People: 500, Zips: 20, Today: 2026}
	if !reflect.DeepEqual(People(cfg, 7), People(cfg, 7)) {
		t.Error("same seed, different people")
	}
}

func TestPeopleHardCases(t *testing.T) {
	ps := People(DefaultConfig(), 1)
	count := map[Tag]int{}
	for i, p := range ps {
		if p.ID != i {
			t.Fatalf("person %d has ID %d", i, p.ID)
		}
		for _, tag := range []Tag{Twin, JrSr, Married, Moved} {
			if p.Has(tag) {
				count[tag]++
			}
		}
		if p.Kin < 0 {
			continue
		}
		k := ps[p.Kin]
		if k.Kin != p.ID {
			t.Errorf("%d's kin is %d, whose kin is %d", p.ID, k.ID, k.Kin)
		}
		if p.Has(Twin) && (k.DOB != p.DOB || k.Last != p.Last || k.First == p.First) {
			t.Errorf("twins %d and %d: %+v / %+v", p.ID, k.ID, p, k)
		}
		if p.Has(JrSr) && (k.First != p.First || k.DOB == p.DOB) {
			t.Errorf("Jr/Sr %d and %d: %+v / %+v", p.ID, k.ID, p, k)
		}
		if p.Has(Married) && p.Maiden == "" {
			t.Errorf("married %d: %+v", p.ID, p)
		}
	}
	t.Logf("%d people: %d twins, %d Jr/Sr, %d married, %d moved", len(ps), count[Twin], count[JrSr], count[Married], count[Moved])
	for _, tag := range []Tag{Twin, JrSr, Married, Moved} {
		if count[tag] < 50 {
			t.Errorf("only %d people with tag %d", count[tag], tag)
		}
	}
}
