package record

import "testing"

func TestSameGiven(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want bool
	}{
		{"BOB", "ROBERT", true},
		{"ROBERT", "BOB", true},
		{"ROBERT", "ROBERT", true},
		{"BOB", "BILL", false}, // two nicknames, different names
		{"BOB", "ROB", true},   // both ROBERT
		{"JOHN", "JON", false}, // JON is JONATHAN here, not JOHN
		{"MARIA", "MARIO", false},
	} {
		if got := SameGiven(c.a, c.b); got != c.want {
			t.Errorf("SameGiven(%s, %s) = %v", c.a, c.b, got)
		}
	}
}
