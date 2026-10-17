package record

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// CSV columns, by header name. Only first, last and dob are required; dob
// is YYYY-MM-DD and may be empty.
var columns = []string{"id", "source", "first", "middle", "last", "suffix", "dob", "sex", "phone", "street", "zip", "ssn4"}

// ReadCSV reads records with a header row. Unknown columns are ignored. A
// record's ID is its row number, starting at 0, whatever the id column
// says; the id column is returned alongside so output can refer to it.
func ReadCSV(r io.Reader) (recs []Record, ids []string, err error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	head, err := cr.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("csv header: %w", err)
	}
	col := map[string]int{}
	for i, h := range head {
		col[strings.ToLower(strings.TrimSpace(h))] = i
	}
	for _, need := range []string{"first", "last", "dob"} {
		if _, ok := col[need]; !ok {
			return nil, nil, fmt.Errorf("csv header: no %q column (want some of %s)", need, strings.Join(columns, ", "))
		}
	}
	for line := 2; ; line++ {
		row, err := cr.Read()
		if err == io.EOF {
			return recs, ids, nil
		}
		if err != nil {
			return nil, nil, err
		}
		get := func(name string) string {
			if i, ok := col[name]; ok && i < len(row) {
				return strings.TrimSpace(row[i])
			}
			return ""
		}
		rec := Record{
			ID: len(recs), First: get("first"), Middle: get("middle"), Last: get("last"),
			Suffix: get("suffix"), Phone: get("phone"), Street: get("street"), Zip: get("zip"), SSN4: get("ssn4"),
		}
		if s := get("source"); s != "" {
			if rec.Source, err = strconv.Atoi(s); err != nil {
				return nil, nil, fmt.Errorf("line %d: source %q is not a number", line, s)
			}
		}
		if rec.DOB, err = parseDate(get("dob")); err != nil {
			return nil, nil, fmt.Errorf("line %d: %w", line, err)
		}
		switch strings.ToUpper(get("sex")) {
		case "M", "MALE":
			rec.Sex = 'M'
		case "F", "FEMALE":
			rec.Sex = 'F'
		}
		recs = append(recs, rec)
		ids = append(ids, get("id"))
	}
}

func parseDate(s string) (Date, error) {
	if s == "" {
		return Date{}, nil
	}
	p := strings.Split(s, "-")
	if len(p) != 3 {
		return Date{}, fmt.Errorf("dob %q: want YYYY-MM-DD", s)
	}
	y, e1 := strconv.Atoi(p[0])
	m, e2 := strconv.Atoi(p[1])
	d, e3 := strconv.Atoi(p[2])
	if e1 != nil || e2 != nil || e3 != nil || y < 1850 || y > 2200 || m < 1 || m > 12 || d < 1 || d > 31 {
		return Date{}, fmt.Errorf("dob %q: want YYYY-MM-DD", s)
	}
	return Date{Y: y, M: m, D: d}, nil
}
