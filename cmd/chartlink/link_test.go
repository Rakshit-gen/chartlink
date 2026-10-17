package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Rakshit-gen/chartlink/match"
	"github.com/Rakshit-gen/chartlink/record"
	"github.com/Rakshit-gen/chartlink/synth"
)

func TestLinkCSV(t *testing.T) {
	cfg := synth.DefaultConfig()
	cfg.People, cfg.Zips = 1500, 15
	d := synth.Records(synth.People(cfg, 3), 4)
	var in, out bytes.Buffer
	if err := record.WriteCSV(&in, d.Records); err != nil {
		t.Fatal(err)
	}
	rev := filepath.Join(t.TempDir(), "review.csv")
	mc := match.DefaultConfig()
	mc.Samples = 50000
	if err := link(&in, &out, rev, mc); err != nil {
		t.Fatal(err)
	}
	rows := strings.Split(strings.TrimSpace(out.String()), "\n")
	if rows[0] != "id,chart" || len(rows) != len(d.Records)+1 {
		t.Fatalf("%d rows, header %q", len(rows), rows[0])
	}
	// Charts should be close to people: no more than 10% extra.
	got := map[string]bool{}
	for _, r := range rows[1:] {
		got[strings.Split(r, ",")[1]] = true
	}
	if n := len(got); n < len(d.People) || n > len(d.People)*11/10 {
		t.Errorf("%d charts for %d people", n, len(d.People))
	}
	b, err := os.ReadFile(rev)
	if err != nil || !strings.HasPrefix(string(b), "a,b,probability\n") {
		t.Errorf("review file: %v %q", err, b)
	}
}

func TestLinkRejectsBadInput(t *testing.T) {
	var out bytes.Buffer
	if err := link(strings.NewReader("first,last,dob\nA,B,1980-01-01\n"), &out, "", match.DefaultConfig()); err == nil {
		t.Error("linked a single record")
	}
	if err := link(strings.NewReader("name\nx\n"), &out, "", match.DefaultConfig()); err == nil {
		t.Error("accepted a CSV without the needed columns")
	}
}
