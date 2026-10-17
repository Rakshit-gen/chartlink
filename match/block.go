package match

import (
	"fmt"
	"slices"
)

// A blocking pass gives each record a key; only records that share a key
// in some pass are compared. Each pass is there for a way the others miss
// a true pair.
type pass struct {
	name string
	key  func(p *prepped) string
}

var passes = []pass{
	// The workhorse: typos in the last name mostly keep its Soundex code.
	{"last name sound + birth year", func(p *prepped) string {
		if p.lastSx == "" || p.dob.Y == 0 {
			return ""
		}
		return fmt.Sprint(p.lastSx, p.dob.Y)
	}},
	// Catches a changed last name (marriage) when the DOB was typed right.
	{"date of birth + first initial", func(p *prepped) string {
		if !p.dob.Known() || p.first == "" {
			return ""
		}
		return fmt.Sprint(p.dob, p.first[:1])
	}},
	// Catches a mistyped DOB when the phone survived.
	{"phone", func(p *prepped) string { return p.phone }},
	// A changed last name and a mistyped DOB, at the same address.
	{"first name sound + zip", func(p *prepped) string {
		if p.firstSx == "" || p.zip == "" {
			return ""
		}
		return p.firstSx + p.zip
	}},
	{"last 4 of SSN + birth month", func(p *prepped) string {
		if p.ssn4 == "" || p.dob.M == 0 {
			return ""
		}
		return fmt.Sprint(p.ssn4, p.dob.M)
	}},
}

// maxBlock caps a block's size. A key shared by thousands of records (a
// shared clinic phone number, say) says nothing and would cost millions of
// comparisons, so it is skipped.
const maxBlock = 200

// candidates returns every pair of records that share a key in any pass,
// each pair once, smaller index first.
func candidates(ps []prepped) []uint64 {
	var out []uint64
	for _, ps2 := range passes {
		blocks := map[string][]uint32{}
		for i := range ps {
			if k := ps2.key(&ps[i]); k != "" {
				blocks[k] = append(blocks[k], uint32(i))
			}
		}
		for _, b := range blocks {
			if len(b) < 2 || len(b) > maxBlock {
				continue
			}
			for i, x := range b {
				for _, y := range b[i+1:] {
					out = append(out, pack(x, y))
				}
			}
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}

func pack(a, b uint32) uint64 {
	if a > b {
		a, b = b, a
	}
	return uint64(a)<<32 | uint64(b)
}

func unpack(p uint64) (int, int) { return int(p >> 32), int(uint32(p)) }
