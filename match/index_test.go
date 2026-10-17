package match

import (
	"testing"

	"github.com/Rakshit-gen/chartlink/score"
	"github.com/Rakshit-gen/chartlink/synth"
)

// Fitted on the first half of the records, then fed all of them one by
// one, the index should do about as well as linking everything at once.
func TestIndexMatchesBatch(t *testing.T) {
	d := synth.Records(synth.People(synth.Config{People: 4000, Zips: 40, Today: 2026}, 5), 6)
	cfg := DefaultConfig()
	cfg.Samples = 100000
	batch := score.Score(d, Link(d.Records, cfg).Chart)

	ix := NewIndex(Fit(d.Records[:len(d.Records)/2], cfg), cfg)
	for i, r := range d.Records {
		if id, _ := ix.Add(r); id != i {
			t.Fatalf("record %d got ID %d", i, id)
		}
	}
	inc := score.Score(d, ix.Charts())
	t.Logf("batch:       precision %.4f recall %.4f overlays %d", batch.Precision, batch.Recall, batch.Overlays)
	t.Logf("incremental: precision %.4f recall %.4f overlays %d", inc.Precision, inc.Recall, inc.Overlays)
	if inc.Precision < batch.Precision-0.003 || inc.Recall < batch.Recall-0.01 {
		t.Error("incremental linking is much worse than batch")
	}
}
