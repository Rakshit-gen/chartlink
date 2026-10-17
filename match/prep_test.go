package match

import (
	"testing"

	"github.com/Rakshit-gen/chartlink/record"
)

func TestPrep(t *testing.T) {
	p := prep(record.Record{
		First: "Mary-Ann", Middle: "Louise", Last: "O'Neil", Suffix: "Jr.",
		DOB: record.Date{Y: 1900, M: 1, D: 1}, Phone: "+1 (555) 201-3344",
		Street: "12 Oak Street, Apt. 4", Zip: "02139-1234", SSN4: "12",
	})
	want := prepped{
		first: "MARYANN", middle: "L", last: "ONEIL", suffix: "JR",
		phone: "5552013344", street: "12 OAK ST APT 4", house: "12", zip: "02139",
		firstSx: "M650", lastSx: "O540",
	}
	if p != want {
		t.Errorf("got  %+v\nwant %+v", p, want)
	}
}
