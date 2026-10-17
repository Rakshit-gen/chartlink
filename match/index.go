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

	// The index grows, so a random record in it is ever less likely to
	// be this record's person.
	prior := ix.sc.m.PriorAt(id)
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
			pr := ix.sc.pair(o, id, compare(&ix.sc.ps[o], p))
			pr.Prob = prob(prior, pr.Weight)
			ix.pairs = append(ix.pairs, pr)
		}
	}
	for i, k := range keys {
		if k != "" {
			ix.blocks[i][k] = append(ix.blocks[i][k], uint32(id))
		}
	}

	slices.SortFunc(ix.pairs, func(x, y Pair) int { return cmp.Compare(y.Weight, x.Weight) })
	for _, pr := range ix.pairs {
		review = place(ix.ch, ix.sc, pr, &ix.cfg, review)
	}
	return id, review
}

// Charts returns every record's chart, by record ID.
func (ix *Index) Charts() []int { return ix.ch.labels() }

// Len is the number of records added.
func (ix *Index) Len() int { return len(ix.sc.ps) }
