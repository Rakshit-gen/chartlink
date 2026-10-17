// Package record holds the patient record every other package works on,
// and the string comparisons used to match records.
package record

import "strings"

// Date is a date of birth. A zero field means the source didn't send it.
type Date struct {
	Y, M, D int
}

func (d Date) Known() bool { return d.Y != 0 && d.M != 0 && d.D != 0 }

// Record is one patient registration as a source system sent it.
type Record struct {
	ID     int
	Source int
	First  string
	Middle string
	Last   string
	Suffix string // JR, SR, II, III; empty if none
	DOB    Date
	Sex    byte // 'M', 'F', or 0 if unknown
	Phone  string
	Street string
	Zip    string
	SSN4   string // last four digits of the SSN, often missing
}

// Clean upper-cases a name, folds accented Latin letters to plain ones and
// drops everything else, so "O'Brien", "OBRIEN" and "o brien" compare
// equal, and so do "José" and "JOSE" (registration systems differ on
// whether they keep accents).
func Clean(s string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(s) {
		if f, ok := fold[r]; ok {
			r = f
		}
		if r >= 'A' && r <= 'Z' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

var fold = map[rune]rune{}

func init() {
	for base, accented := range map[rune]string{
		'A': "ÀÁÂÃÄÅ", 'C': "Ç", 'E': "ÈÉÊË", 'I': "ÌÍÎÏ", 'N': "Ñ",
		'O': "ÒÓÔÕÖØ", 'U': "ÙÚÛÜ", 'Y': "Ý",
	} {
		for _, r := range accented {
			fold[r] = base
		}
	}
}

// Digits keeps only digits, for phones and SSNs.
func Digits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
