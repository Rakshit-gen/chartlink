// Package baseline has the matching rule most registration systems start
// with: same first name, last name and date of birth.
package baseline

import (
	"fmt"

	"github.com/Rakshit-gen/chartlink/record"
)

// Exact puts records in one chart when their cleaned first name, last name
// and date of birth are all equal. Records missing any of the three get a
// chart of their own.
func Exact(recs []record.Record) []int {
	label := make([]int, len(recs))
	seen := map[string]int{}
	for i, r := range recs {
		f, l := record.Clean(r.First), record.Clean(r.Last)
		if f == "" || l == "" || !r.DOB.Known() {
			label[i] = i
			continue
		}
		k := fmt.Sprintf("%s|%s|%d-%d-%d", f, l, r.DOB.Y, r.DOB.M, r.DOB.D)
		if c, ok := seen[k]; ok {
			label[i] = c
		} else {
			seen[k] = i
			label[i] = i
		}
	}
	return label
}
