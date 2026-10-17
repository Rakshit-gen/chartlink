package match

import "testing"

func TestCharts(t *testing.T) {
	c := newCharts(5)
	c.union(c.find(0), c.find(1))
	c.union(c.find(3), c.find(1))
	l := c.labels()
	if l[0] != l[1] || l[1] != l[3] || l[2] == l[0] || l[4] == l[0] || l[2] == l[4] {
		t.Errorf("labels %v", l)
	}
	if got := len(c.members[c.find(0)]); got != 3 {
		t.Errorf("chart has %d members", got)
	}
}
