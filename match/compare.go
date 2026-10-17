package match

import "github.com/Rakshit-gen/chartlink/record"

// Field is one thing two records are compared on.
type Field int

const (
	First Field = iota
	Last
	DOB
	Sex
	Phone
	Street
	Zip
	SSN
	Middle
	Suffix
	Fields
)

var fieldNames = [Fields]string{"first", "last", "dob", "sex", "phone", "street", "zip", "ssn4", "middle", "suffix"}

func (f Field) String() string { return fieldNames[f] }

// Missing means one side didn't have the field. It carries no evidence
// either way.
const Missing = -1

// Level names for each field, in order from strongest agreement to
// disagreement. A pair gets one level per field.
var levelNames = [Fields][]string{
	First:  {"exact", "nickname", "close", "differ"},
	Last:   {"exact", "close", "differ"},
	DOB:    {"exact", "month and day swapped", "one part off", "differ"},
	Sex:    {"agree", "differ"},
	Phone:  {"agree", "differ"},
	Street: {"exact", "same house number, close street", "differ"},
	Zip:    {"agree", "differ"},
	SSN:    {"agree", "one digit off or swapped", "differ"},
	Middle: {"agree", "differ"},
	Suffix: {"agree", "differ"},
}

// Levels returns how many levels a field has.
func Levels(f Field) int { return len(levelNames[f]) }

// LevelName describes a level.
func LevelName(f Field, l int) string {
	if l == Missing {
		return "missing"
	}
	return levelNames[f][l]
}

// Gamma is the comparison vector of a pair: one level per field.
type Gamma [Fields]int8

func compare(a, b *prepped) Gamma {
	var g Gamma
	g[First] = names(a.first, b.first, true)
	g[Last] = names(a.last, b.last, false)
	g[DOB] = dates(a.dob, b.dob)
	g[Sex] = same(a.sex != 0 && b.sex != 0, a.sex == b.sex)
	g[Phone] = same(a.phone != "" && b.phone != "", a.phone == b.phone)
	g[Street] = streets(a, b)
	g[Zip] = same(a.zip != "" && b.zip != "", a.zip == b.zip)
	g[SSN] = ssns(a.ssn4, b.ssn4)
	g[Middle] = same(a.middle != "" && b.middle != "", a.middle == b.middle)
	g[Suffix] = same(a.suffix != "" && b.suffix != "", a.suffix == b.suffix)
	return g
}

func same(both, equal bool) int8 {
	switch {
	case !both:
		return Missing
	case equal:
		return 0
	}
	return 1
}

// closeName is the Jaro-Winkler score above which two names are taken as
// one name with a typo. At 0.88, GARCIA/GRACIA (0.95) and SMITH/SMYTH
// (0.89) count and JOHN/JOAN (0.87) doesn't. MARIA/MARIE (0.92) also
// counts, which is the price: two different names can look like a typo.
const closeName = 0.88

func names(a, b string, given bool) int8 {
	if a == "" || b == "" {
		return Missing
	}
	if a == b {
		return 0
	}
	off := int8(0)
	if given {
		if record.SameGiven(a, b) {
			return 1
		}
		off = 1
	}
	if record.JaroWinkler(a, b) >= closeName {
		return 1 + off
	}
	return 2 + off
}

func dates(a, b record.Date) int8 {
	if !a.Known() || !b.Known() {
		return Missing
	}
	switch {
	case a == b:
		return 0
	case a.Y == b.Y && a.M == b.D && a.D == b.M:
		return 1
	}
	parts := 0
	if a.Y != b.Y {
		parts++
	}
	if a.M != b.M {
		parts++
	}
	if a.D != b.D {
		parts++
	}
	if parts == 1 {
		return 2
	}
	return 3
}

func streets(a, b *prepped) int8 {
	if a.street == "" || b.street == "" {
		return Missing
	}
	if a.street == b.street {
		return 0
	}
	if a.house != "" && a.house == b.house && record.JaroWinkler(a.street, b.street) >= 0.9 {
		return 1
	}
	return 2
}

func ssns(a, b string) int8 {
	if a == "" || b == "" {
		return Missing
	}
	if a == b {
		return 0
	}
	diff, at := 0, -1
	for i := range 4 {
		if a[i] != b[i] {
			diff++
			if at < 0 {
				at = i
			}
		}
	}
	if diff == 1 || (diff == 2 && at < 3 && a[at] == b[at+1] && a[at+1] == b[at]) {
		return 1
	}
	return 2
}
