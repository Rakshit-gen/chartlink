package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/Rakshit-gen/chartlink/match"
	"github.com/Rakshit-gen/chartlink/record"
	"github.com/Rakshit-gen/chartlink/synth"
)

// exportCmd writes a synthetic population's records as CSV.
func exportCmd(args []string) {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	people := fs.Int("people", 2000, "people in the population")
	seed := fs.Uint64("seed", 1, "random seed")
	fs.Parse(args)
	cfg := synth.DefaultConfig()
	cfg.People, cfg.Zips = *people, max(*people/100, 1)
	d := synth.Records(synth.People(cfg, *seed), *seed+1000)
	if err := record.WriteCSV(os.Stdout, d.Records); err != nil {
		fatal(err)
	}
}

// linkCmd reads records as CSV on stdin, fits a model to them, links them,
// and writes one row per record with its chart. Pairs for review go to the
// -review file.
func linkCmd(args []string) {
	fs := flag.NewFlagSet("link", flag.ExitOnError)
	merge := fs.Float64("merge", match.DefaultConfig().Merge, "match probability at which records are merged")
	reviewFile := fs.String("review", "", "write pairs for review to this CSV file")
	fs.Parse(args)
	cfg := match.DefaultConfig()
	if *merge <= 0 || *merge >= 1 {
		fatal(fmt.Errorf("-merge must be between 0 and 1"))
	}
	cfg.Merge = *merge
	cfg.Review = min(cfg.Review, cfg.Merge)
	if err := link(os.Stdin, os.Stdout, *reviewFile, cfg); err != nil {
		fatal(err)
	}
}

func link(in io.Reader, out io.Writer, reviewFile string, cfg match.Config) error {
	recs, ids, err := record.ReadCSV(in)
	if err != nil {
		return err
	}
	if len(recs) < 2 {
		return fmt.Errorf("need at least two records, got %d", len(recs))
	}
	res := match.Link(recs, cfg)
	label := func(i int) string {
		if ids[i] != "" {
			return ids[i]
		}
		return strconv.Itoa(i)
	}
	w := csv.NewWriter(out)
	w.Write([]string{"id", "chart"})
	for i, c := range res.Chart {
		w.Write([]string{label(i), label(c)})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%d records, %d charts, %d pairs for review\n", len(recs), charts(res.Chart), len(res.Review))
	if reviewFile == "" {
		return nil
	}
	f, err := os.Create(reviewFile)
	if err != nil {
		return err
	}
	rw := csv.NewWriter(f)
	rw.Write([]string{"a", "b", "probability"})
	for _, p := range res.Review {
		rw.Write([]string{label(p.A), label(p.B), strconv.FormatFloat(p.Prob, 'f', 4, 64)})
	}
	rw.Flush()
	if err := rw.Error(); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func charts(chart []int) int {
	seen := map[int]bool{}
	for _, c := range chart {
		seen[c] = true
	}
	return len(seen)
}
