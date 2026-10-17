package record

import (
	"strconv"
	"strings"
	"testing"
)

func TestReadCSV(t *testing.T) {
	in := "id,First,Last,DOB,Sex,extra\n" +
		"a1,Robert,Smith,1980-04-02,male,x\n" +
		"a2,Bob,Smith,,F,y\n"
	recs, ids, err := ReadCSV(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 2 || ids[0] != "a1" || ids[1] != "a2" {
		t.Fatalf("%+v %v", recs, ids)
	}
	if recs[0].DOB != (Date{1980, 4, 2}) || recs[0].Sex != 'M' || recs[1].DOB.Known() || recs[1].Sex != 'F' || recs[1].ID != 1 {
		t.Errorf("%+v", recs)
	}
}

func TestReadCSVRejects(t *testing.T) {
	for name, in := range map[string]string{
		"no last column": "first,dob\nA,1980-01-01\n",
		"bad date":       "first,last,dob\nA,B,1980/01/01\n",
		"month 13":       "first,last,dob\nA,B,1980-13-01\n",
		"bad source":     "first,last,dob,source\nA,B,,clinic\n",
		"empty":          "",
	} {
		if _, _, err := ReadCSV(strings.NewReader(in)); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

func TestCSVRoundTrip(t *testing.T) {
	want := []Record{
		{ID: 0, Source: 2, First: "ANA", Middle: "M", Last: "DE LA CRUZ", Suffix: "JR", DOB: Date{1990, 1, 2}, Sex: 'F', Phone: "5551234567", Street: "1 MAIN ST, APT 2", Zip: "02139", SSN4: "0042"},
		{ID: 1, First: "LEE", Last: "KIM"},
	}
	var b strings.Builder
	if err := WriteCSV(&b, want); err != nil {
		t.Fatal(err)
	}
	got, ids, err := ReadCSV(strings.NewReader(b.String()))
	if err != nil {
		t.Fatal(err)
	}
	for i := range want {
		if got[i] != want[i] || ids[i] != strconv.Itoa(i) {
			t.Errorf("row %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}
