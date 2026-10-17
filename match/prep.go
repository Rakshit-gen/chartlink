// Package match links patient records that belong to the same person.
package match

import (
	"strings"

	"github.com/Rakshit-gen/chartlink/record"
)

// prepped is a record with every field cleaned once, so comparing a pair
// does no string work beyond the comparisons themselves.
type prepped struct {
	first, middle, last, suffix string
	firstSx, lastSx             string // Soundex codes
	dob                         record.Date
	sex                         byte
	phone, street, house, zip   string
	ssn4                        string
}

// Placeholder dates of birth that registration systems fill in when the
// real one isn't known. They say nothing about who the patient is.
var placeholderDOB = map[record.Date]bool{
	{Y: 1900, M: 1, D: 1}: true,
	{Y: 1901, M: 1, D: 1}: true,
}

func prep(r record.Record) prepped {
	p := prepped{
		first: record.Clean(r.First), middle: record.Clean(r.Middle), last: record.Clean(r.Last),
		suffix: normSuffix(r.Suffix), dob: r.DOB, sex: r.Sex,
		phone: record.Digits(r.Phone), zip: record.Digits(r.Zip), ssn4: record.Digits(r.SSN4),
	}
	if len(p.middle) > 1 {
		p.middle = p.middle[:1] // sources disagree on full middle names vs initials
	}
	if !p.dob.Known() || placeholderDOB[p.dob] {
		p.dob = record.Date{}
	}
	if len(p.phone) == 11 && p.phone[0] == '1' {
		p.phone = p.phone[1:]
	}
	if len(p.phone) != 10 {
		p.phone = ""
	}
	if len(p.ssn4) != 4 {
		p.ssn4 = ""
	}
	if len(p.zip) > 5 {
		p.zip = p.zip[:5]
	}
	p.street, p.house = normStreet(r.Street)
	p.firstSx, p.lastSx = record.Soundex(p.first), record.Soundex(p.last)
	return p
}

func normSuffix(s string) string {
	switch record.Clean(s) {
	case "JR", "JUNIOR", "II":
		return "JR"
	case "SR", "SENIOR":
		return "SR"
	case "III":
		return "III"
	case "IV":
		return "IV"
	}
	return ""
}

var streetWords = map[string]string{
	"STREET": "ST", "AVENUE": "AVE", "AV": "AVE", "ROAD": "RD", "DRIVE": "DR",
	"LANE": "LN", "COURT": "CT", "BOULEVARD": "BLVD", "PLACE": "PL",
	"NORTH": "N", "SOUTH": "S", "EAST": "E", "WEST": "W",
	"APARTMENT": "APT", "UNIT": "APT", "#": "APT",
}

// normStreet upper-cases a street address, shortens the common words the
// way the USPS does, and returns it with its house number split out.
func normStreet(s string) (street, house string) {
	words := strings.Fields(strings.ToUpper(strings.NewReplacer(".", "", ",", " ").Replace(s)))
	for i, w := range words {
		if short, ok := streetWords[w]; ok {
			words[i] = short
		}
	}
	if len(words) > 0 && strings.Trim(words[0], "0123456789") == "" {
		house = words[0]
	}
	return strings.Join(words, " "), house
}
