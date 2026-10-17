package match

import (
	"testing"

	"github.com/Rakshit-gen/chartlink/score"
	"github.com/Rakshit-gen/chartlink/synth"
)

func TestLinkSmallPopulation(t *testing.T) {
	d := synth.Records(synth.People(synth.Config{People: 4000, Zips: 40, Today: 2026}, 5), 6)
	cfg := DefaultConfig()
	cfg.Samples = 100000
	r := score.Score(d, Link(d.Records, cfg).Chart)
	t.Logf("%+v", r)
	if r.Precision < 0.995 || r.Recall < 0.95 {
		t.Errorf("precision %.4f recall %.4f", r.Precision, r.Recall)
	}
	if r.TwinsMerged > 1 || r.JrSrMerged > 1 {
		t.Errorf("%d of %d twins and %d of %d Jr/Sr merged", r.TwinsMerged, r.TwinPairs, r.JrSrMerged, r.JrSrPairs)
	}
}

func TestLinkDeterministic(t *testing.T) {
	d := synth.Records(synth.People(synth.Config{People: 1000, Zips: 20, Today: 2026}, 1), 1)
	cfg := DefaultConfig()
	cfg.Samples = 20000
	a, b := Link(d.Records, cfg), Link(d.Records, cfg)
	for i := range a.Chart {
		if a.Chart[i] != b.Chart[i] {
			t.Fatal("two runs on the same records disagree")
		}
	}
}
