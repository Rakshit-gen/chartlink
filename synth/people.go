// Package synth makes a population of people and the messy records that
// hospitals, clinics, labs and pharmacies keep about them, with the true
// person behind every record.
package synth

import (
	"fmt"
	"math/rand/v2"

	"github.com/Rakshit-gen/chartlink/record"
)

// Tag marks the people who make matching hard.
type Tag uint8

const (
	Twin    Tag = 1 << iota // has a twin: same last name, birthday and address
	JrSr                    // father and son with the same name and address
	Married                 // changed last name; older records have the old one
	Moved                   // changed address; older records have the old one
)

// Address is where a household lives and its landline, if it has one.
type Address struct {
	Street, Zip, Phone string
}

// Person is the truth behind a set of records.
type Person struct {
	ID        int
	Household int
	First     string
	Middle    string
	Last      string
	Maiden    string // last name before marriage; empty if never changed
	Suffix    string
	DOB       record.Date
	Sex       byte
	SSN4      string
	Mobile    string
	Addr      []Address // oldest first
	Kin       int       // the twin, or the other half of a Jr/Sr pair; -1 if none
	Tags      Tag
}

// Has reports whether every tag in t is set.
func (p *Person) Has(t Tag) bool { return p.Tags&t == t }

// Config sizes the population.
type Config struct {
	People int
	Zips   int // fewer zips means more people share one
	// Today is the year records are made in; ages are counted from it.
	Today int
}

// DefaultConfig is the population the README's numbers come from.
func DefaultConfig() Config {
	return Config{People: 20000, Zips: 200, Today: 2026}
}

var streets = []string{"MAIN", "OAK", "PINE", "MAPLE", "CEDAR", "ELM", "WASHINGTON", "LAKE",
	"HILL", "PARK", "SUNSET", "RIVER", "CHURCH", "HIGHLAND", "MILL", "FOREST", "MEADOW",
	"SPRING", "RIDGE", "LINCOLN", "JEFFERSON", "FRANKLIN", "WALNUT", "CHESTNUT", "WILLOW"}
var streetKinds = []string{"ST", "AVE", "RD", "DR", "LN", "CT", "BLVD", "WAY"}

type gen struct {
	r     *rand.Rand
	cfg   Config
	zips  []string
	out   []Person
	house int
}

// People makes at least cfg.People people in households: people living
// alone, couples, families with children, twins, fathers and sons with the
// same name, and some who married or moved. The last household is kept whole, so there
// can be a few more people than asked for.
func People(cfg Config, seed uint64) []Person {
	g := &gen{r: rand.New(rand.NewPCG(seed, 0x9e3779b97f4a7c15)), cfg: cfg}
	for range max(cfg.Zips, 1) {
		g.zips = append(g.zips, fmt.Sprintf("%05d", 10000+g.r.IntN(89999)))
	}
	for len(g.out) < cfg.People {
		g.household()
	}
	return g.out
}

func (g *gen) address() Address {
	a := Address{
		Street: fmt.Sprintf("%d %s %s", 1+g.r.IntN(9999), streets[g.r.IntN(len(streets))], streetKinds[g.r.IntN(len(streetKinds))]),
		Zip:    g.zips[g.r.IntN(len(g.zips))],
	}
	if g.r.Float64() < 0.4 {
		a.Phone = g.phone()
	}
	return a
}

func (g *gen) phone() string { return fmt.Sprintf("%03d%07d", 200+g.r.IntN(800), g.r.IntN(10000000)) }

func (g *gen) birth(year int) record.Date {
	m := 1 + g.r.IntN(12)
	days := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}[m-1]
	return record.Date{Y: year, M: m, D: 1 + g.r.IntN(days)}
}

func (g *gen) person(sex byte, last string, year int, home Address) *Person {
	p := Person{
		ID: len(g.out), Household: g.house, Sex: sex, Last: last, Kin: -1,
		First: given(g.r, sex), DOB: g.birth(year),
		SSN4: fmt.Sprintf("%04d", g.r.IntN(10000)), Addr: []Address{home},
	}
	if g.r.Float64() < 0.7 {
		p.Middle = given(g.r, sex)
	}
	if g.cfg.Today-year >= 14 {
		p.Mobile = g.phone()
	}
	g.out = append(g.out, p)
	return &g.out[len(g.out)-1]
}

func randSex(r *rand.Rand) byte {
	if r.IntN(2) == 0 {
		return 'M'
	}
	return 'F'
}

func (g *gen) household() {
	g.house++
	home := g.address()
	start := len(g.out)
	adult := func() int { return g.cfg.Today - 22 - g.r.IntN(60) }
	u := g.r.Float64()
	switch {
	case u < 0.35: // living alone
		g.person(randSex(g.r), surname(g.r), adult(), home)
	default: // couple, maybe with children
		year := adult()
		last := surname(g.r)
		dad := g.person('M', last, year, home)
		dadID := dad.ID
		mom := g.person('F', surname(g.r), year+g.r.IntN(7)-3, home)
		if g.r.Float64() < 0.75 {
			mom.Maiden, mom.Last = mom.Last, last
			mom.Tags |= Married
		}
		if u < 0.65 {
			g.children(dadID, last, home)
		}
	}
	if g.r.Float64() < 0.15 {
		next := g.address()
		for i := start; i < len(g.out); i++ {
			g.out[i].Addr = append(g.out[i].Addr, next)
			g.out[i].Tags |= Moved
		}
	}
}

func (g *gen) children(dadID int, last string, home Address) {
	parentYear := g.out[dadID].DOB.Y
	year := parentYear + 20 + g.r.IntN(15)
	named := false
	for range 1 + g.r.IntN(3) {
		if year > g.cfg.Today {
			return
		}
		kid := g.person(randSex(g.r), last, year, home)
		kidID := kid.ID
		// A first son sometimes gets his father's full name.
		if kid.Sex == 'M' && !named && g.r.Float64() < 0.1 {
			named = true
			dad := &g.out[dadID]
			kid.First, kid.Middle, kid.Suffix = dad.First, dad.Middle, "JR"
			if g.r.Float64() < 0.5 {
				dad.Suffix = "SR"
			}
			kid.Kin, dad.Kin = dad.ID, kidID
			kid.Tags |= JrSr
			dad.Tags |= JrSr
		}
		if g.r.Float64() < 0.04 && g.out[kidID].Kin < 0 {
			twin := g.person(randSex(g.r), last, year, home)
			first := &g.out[kidID]
			twin.DOB = first.DOB
			for twin.First == first.First {
				twin.First = given(g.r, twin.Sex)
			}
			twin.Kin, first.Kin = first.ID, twin.ID
			twin.Tags |= Twin
			first.Tags |= Twin
		}
		year += 1 + g.r.IntN(5)
	}
}
