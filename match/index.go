package match

import (
	"cmp"
	"slices"

	"github.com/Rakshit-gen/chartlink/record"
)

// Index links records one at a time as they arrive, the way a master
// patient index sits behind a registration feed. It uses a model fitted
// beforehand (see Fit) and the same thresholds and rules as Link.
type Index struct {
	cfg    Config
	sc     *scorer
	ch     *charts
	blocks []map[string][]uint32
	seen   map[int]bool // scratch, reused across Adds
	pairs  []Pair       // scratch
}

// NewIndex starts an empty index.
func NewIndex(m Model, cfg Config) *Index {
	ix := &Index{
		cfg: cfg, sc: &scorer{m: &m, off: cfg.Without}, ch: newCharts(0),
		blocks: make([]map[string][]uint32, len(passes)), seen: map[int]bool{},
	}
	for i := range ix.blocks {
		ix.blocks[i] = map[string][]uint32{}
	}
	return ix
}

// Add links one record and returns its ID (records are numbered in the
// order they were added) and the pairs it raised for review.
func (ix *Index) Add(r record.Record) (id int, review []Pair) {
	id = len(ix.sc.ps)
	ix.sc.ps = append(ix.sc.ps, prep(r))
	ix.ch.add()
	p := &ix.sc.ps[id]

	clear(ix.seen)
	ix.pairs = ix.pairs[:0]
	keys := make([]string, len(passes))
	for i, ps := range passes {
		k := ps.key(p)
		keys[i] = k
		if k == "" {
			continue
		}
		block := ix.blocks[i][k]
		if len(block) >= maxBlock {
			continue
		}
		for _, other := range block {
			o := int(other)
			if ix.seen[o] {
				continue
			}
			ix.seen[o] = true
			w := ix.sc.weight(o, id, compare(&ix.sc.ps[o], p))
			ix.pairs = append(ix.pairs, Pair{A: o, B: id, Weight: w, Prob: ix.sc.m.Prob(w)})
		}
	}
	for i, k := range keys {
		if k != "" {
			ix.blocks[i][k] = append(ix.blocks[i][k], uint32(id))
		}
	}

	slices.SortFunc(ix.pairs, func(x, y Pair) int { return cmp.Compare(y.Weight, x.Weight) })
	for _, pr := range ix.pairs {
		ra, rb := ix.ch.find(pr.A), ix.ch.find(pr.B)
		if ra == rb {
			continue
		}
		switch {
		case pr.Prob >= ix.cfg.Merge && !(ix.sc.on(Guard) && ix.ch.clash(ra, rb, ix.sc.conflict)):
			ix.ch.union(ra, rb)
		case pr.Prob >= ix.cfg.Review:
			review = append(review, pr)
		}
	}
	return id, review
}

// Charts returns every record's chart, by record ID.
func (ix *Index) Charts() []int { return ix.ch.labels() }

// Len is the number of records added.
func (ix *Index) Len() int { return len(ix.sc.ps) }
