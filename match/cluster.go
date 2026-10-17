package match

// charts is a union-find over records that also keeps each chart's
// members, so a merge can look at every pair it would join.
type charts struct {
	parent  []int32
	members [][]int32 // only valid at a root
}

func newCharts(n int) *charts {
	c := &charts{parent: make([]int32, n), members: make([][]int32, n)}
	for i := range n {
		c.parent[i] = int32(i)
		c.members[i] = []int32{int32(i)}
	}
	return c
}

func (c *charts) find(x int) int {
	for int(c.parent[x]) != x {
		c.parent[x] = c.parent[c.parent[x]]
		x = int(c.parent[x])
	}
	return x
}

// union joins two charts, given their roots; the bigger one absorbs the
// smaller so finds stay short.
func (c *charts) union(a, b int) {
	if len(c.members[a]) < len(c.members[b]) {
		a, b = b, a
	}
	c.parent[b] = int32(a)
	c.members[a] = append(c.members[a], c.members[b]...)
	c.members[b] = nil
}

// labels gives every record its chart's root.
func (c *charts) labels() []int {
	out := make([]int, len(c.parent))
	for i := range out {
		out[i] = c.find(i)
	}
	return out
}

// clash reports whether any record in chart a conflicts with any record in
// chart b. A chart grows one pair at a time, so without this, A links to
// B and B links to C can put A and C in one chart though they are twins.
func (c *charts) clash(a, b int, conflict func(x, y int) bool) bool {
	for _, x := range c.members[a] {
		for _, y := range c.members[b] {
			if conflict(int(x), int(y)) {
				return true
			}
		}
	}
	return false
}
