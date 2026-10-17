package match

import (
	"testing"

	"github.com/Rakshit-gen/chartlink/score"
	"github.com/Rakshit-gen/chartlink/synth"
)

func TestScratch(t *testing.T) {
	d := synth.Records(synth.People(synth.DefaultConfig(), 1), 2)
	r := Link(d.Records, DefaultConfig())
	t.Logf("cands %d review %d prior %g", r.Candidates, len(r.Review), r.Model.Prior)
	for f := range Fields {
		t.Logf("%-7s M %.3f U %.4f", f, r.Model.M[f], r.Model.U[f])
	}
	t.Logf("%+v", score.Score(d, r.Chart))
}
