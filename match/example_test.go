package match_test

import (
	"fmt"

	"github.com/Rakshit-gen/chartlink/match"
	"github.com/Rakshit-gen/chartlink/record"
	"github.com/Rakshit-gen/chartlink/synth"
)

// Fit a model on existing records, then link new registrations as they
// arrive: a nickname with a mistyped last name joins the chart, and a twin
// brother at the same address doesn't.
func Example() {
	d := synth.Records(synth.People(synth.DefaultConfig(), 1), 2)
	cfg := match.DefaultConfig()
	ix := match.NewIndex(match.Fit(d.Records, cfg), cfg)

	dob := record.Date{Y: 1988, M: 7, D: 14}
	home := record.Record{Last: "OKONKWO", DOB: dob, Sex: 'M', Phone: "6175550142", Street: "48 ELM ST", Zip: "02139"}
	robert, bob, twin := home, home, home
	robert.First, robert.SSN4 = "ROBERT", "4417"
	bob.First, bob.Last = "BOB", "OKONKOW"
	twin.First, twin.SSN4 = "DANIEL", "4418"

	a, _ := ix.Add(robert)
	b, _ := ix.Add(bob)
	c, review := ix.Add(twin)
	ch := ix.Charts()
	fmt.Println("BOB OKONKOW joins ROBERT OKONKWO:", ch[a] == ch[b])
	fmt.Println("DANIEL joins them:", ch[a] == ch[c])
	fmt.Println("sent for review:", len(review) > 0)
	// Output:
	// BOB OKONKOW joins ROBERT OKONKWO: true
	// DANIEL joins them: false
	// sent for review: true
}
