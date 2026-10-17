package synth

import (
	"math/rand/v2"
	"slices"

	"github.com/Rakshit-gen/chartlink/record"
)

// Source is one system that registers patients, with how often its front
// desk gets each thing wrong or leaves it out.
type Source struct {
	Name     string
	Typo     float64 // per name field
	Nickname float64 // registers BOB instead of ROBERT
	DOB      float64 // date of birth entered wrong
	NoSSN    float64
	NoPhone  float64
	NoAddr   float64
	NoMiddle float64
	Stale    float64 // has the old name or address after a marriage or move
}

// Sources are the five systems in the benchmark. A hospital registration
// desk is careful and asks for an SSN; a lab or pharmacy takes what is on
// the order.
var Sources = []Source{
	{Name: "hospital", Typo: 0.02, Nickname: 0.03, DOB: 0.01, NoSSN: 0.3, NoPhone: 0.1, NoAddr: 0.05, NoMiddle: 0.3, Stale: 0.3},
	{Name: "clinic", Typo: 0.04, Nickname: 0.15, DOB: 0.02, NoSSN: 0.6, NoPhone: 0.1, NoAddr: 0.1, NoMiddle: 0.5, Stale: 0.4},
	{Name: "urgent care", Typo: 0.08, Nickname: 0.15, DOB: 0.04, NoSSN: 0.7, NoPhone: 0.2, NoAddr: 0.2, NoMiddle: 0.7, Stale: 0.5},
	{Name: "lab", Typo: 0.06, Nickname: 0.1, DOB: 0.03, NoSSN: 0.9, NoPhone: 0.5, NoAddr: 0.5, NoMiddle: 0.9, Stale: 0.5},
	{Name: "pharmacy", Typo: 0.04, Nickname: 0.2, DOB: 0.02, NoSSN: 0.95, NoPhone: 0.05, NoAddr: 0.3, NoMiddle: 0.8, Stale: 0.6},
}

// Data is a set of records and the truth behind them.
type Data struct {
	People  []Person
	Records []record.Record
	Truth   []int // Truth[i] is the person behind Records[i]
}

// Records makes between one and six records per person, each from a random
// source with that source's mistakes, in a random order. Record IDs are
// their index.
func Records(people []Person, seed uint64) Data {
	r := rand.New(rand.NewPCG(seed, 0xc2b2ae3d27d4eb4f))
	d := Data{People: people}
	for i := range people {
		p := &people[i]
		n := 1
		for n < 6 && r.Float64() < 0.55 {
			n++
		}
		for range n {
			d.Records = append(d.Records, render(r, p, &Sources[r.IntN(len(Sources))]))
			d.Truth = append(d.Truth, p.ID)
		}
	}
	r.Shuffle(len(d.Records), func(i, j int) {
		d.Records[i], d.Records[j] = d.Records[j], d.Records[i]
		d.Truth[i], d.Truth[j] = d.Truth[j], d.Truth[i]
	})
	for i := range d.Records {
		d.Records[i].ID = i
	}
	return d
}

// render is what one source wrote down about a person.
func render(r *rand.Rand, p *Person, s *Source) record.Record {
	rec := record.Record{
		Source: sourceIndex(s), First: p.First, Middle: p.Middle, Last: p.Last,
		Suffix: p.Suffix, DOB: p.DOB, Sex: p.Sex, SSN4: p.SSN4,
	}
	addr := p.Addr[len(p.Addr)-1]
	if len(p.Addr) > 1 && r.Float64() < s.Stale {
		addr = p.Addr[r.IntN(len(p.Addr)-1)]
	}
	if p.Maiden != "" && r.Float64() < s.Stale {
		rec.Last = p.Maiden
	}
	rec.Street, rec.Zip = addr.Street, addr.Zip
	rec.Phone = p.Mobile
	if rec.Phone == "" || (addr.Phone != "" && r.Float64() < 0.3) {
		rec.Phone = addr.Phone
	}

	if p.Has(ByMiddle) && r.Float64() < 0.5 {
		// Registered under the name the patient gave: ROBERT for JAMES
		// ROBERT SMITH.
		rec.First, rec.Middle = p.Middle, ""
	}
	if r.Float64() < s.Nickname {
		rec.First = nickname(r, rec.First)
	}
	if r.Float64() < s.Typo {
		rec.First = typo(r, rec.First)
	}
	if r.Float64() < s.Typo {
		rec.Last = typo(r, rec.Last)
	}
	if r.Float64() < s.DOB {
		rec.DOB = wrongDate(r, rec.DOB)
	}
	if r.Float64() < s.NoSSN {
		rec.SSN4 = ""
	} else if r.Float64() < s.Typo {
		rec.SSN4 = digitTypo(r, rec.SSN4)
	}
	if r.Float64() < s.NoPhone {
		rec.Phone = ""
	}
	if r.Float64() < s.NoAddr {
		rec.Street, rec.Zip = "", ""
	}
	if r.Float64() < s.NoMiddle {
		rec.Middle = ""
	} else if rec.Middle != "" && r.Float64() < 0.5 {
		rec.Middle = rec.Middle[:1] // just the initial
	}
	if rec.Suffix != "" && r.Float64() < 0.4 {
		rec.Suffix = "" // suffixes get dropped all the time
	}
	if r.Float64() < 0.02 {
		rec.Sex = 0
	}
	return rec
}

func sourceIndex(s *Source) int {
	for i := range Sources {
		if &Sources[i] == s {
			return i
		}
	}
	return -1
}

// nickname swaps a formal name for one of its nicknames, if it has one.
func nickname(r *rand.Rand, first string) string {
	var opts []string
	for nick, formals := range record.Nicknames() {
		for _, f := range formals {
			if f == first {
				opts = append(opts, nick)
			}
		}
	}
	if len(opts) == 0 {
		return first
	}
	// Map order is random; sort so the same seed gives the same records.
	slices.Sort(opts)
	return opts[r.IntN(len(opts))]
}

// keyboard neighbours on a QWERTY layout, for substitution typos.
var near = map[byte]string{
	'Q': "WA", 'W': "QES", 'E': "WRD", 'R': "ETF", 'T': "RYG", 'Y': "TUH", 'U': "YIJ",
	'I': "UOK", 'O': "IPL", 'P': "OL", 'A': "QSZ", 'S': "ADWX", 'D': "SFEC", 'F': "DGRV",
	'G': "FHTB", 'H': "GJYN", 'J': "HKUM", 'K': "JLI", 'L': "KO", 'Z': "AX", 'X': "ZCS",
	'C': "XVD", 'V': "CBF", 'B': "VNG", 'N': "BMH", 'M': "NJ",
}

// typo makes one keying mistake: a neighbouring key, a dropped letter, a
// doubled letter or two letters swapped.
func typo(r *rand.Rand, s string) string {
	if len(s) < 2 {
		return s
	}
	b := []byte(s)
	i := r.IntN(len(b))
	switch r.IntN(4) {
	case 0:
		if n := near[b[i]]; n != "" {
			b[i] = n[r.IntN(len(n))]
		}
	case 1:
		b = append(b[:i], b[i+1:]...)
	case 2:
		b = append(b[:i+1], b[i:]...)
	default:
		if i == len(b)-1 {
			i--
		}
		b[i], b[i+1] = b[i+1], b[i]
	}
	return string(b)
}

// wrongDate makes the date-of-birth mistakes registration desks make: day
// and month swapped, a digit off in the day or year, or the year of the
// visit's own date typed by habit.
func wrongDate(r *rand.Rand, d record.Date) record.Date {
	switch r.IntN(4) {
	case 0:
		if d.D <= 12 && d.D != d.M {
			d.M, d.D = d.D, d.M
			return d
		}
		fallthrough
	case 1:
		if d.D < 28 {
			d.D++
		} else {
			d.D--
		}
	case 2:
		d.Y += []int{-10, -1, 1, 10}[r.IntN(4)]
	default:
		d.Y = 2026
	}
	return d
}

// digitTypo changes one digit or swaps two neighbours.
func digitTypo(r *rand.Rand, s string) string {
	b := []byte(s)
	i := r.IntN(len(b) - 1)
	if r.IntN(2) == 0 && b[i] != b[i+1] {
		b[i], b[i+1] = b[i+1], b[i]
	} else {
		b[i] = '0' + (b[i]-'0'+1+byte(r.IntN(9)))%10
	}
	return string(b)
}
