// Command chartlink runs the experiments behind results/, and links a CSV
// of patient records into charts.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/Rakshit-gen/chartlink/lab"
	"github.com/Rakshit-gen/chartlink/match"
	"github.com/Rakshit-gen/chartlink/score"
	"github.com/Rakshit-gen/chartlink/synth"
)

// populations is how many synthetic populations every table averages over.
const populations = 5

func usage() {
	fmt.Fprint(os.Stderr, `usage: chartlink <command> [-out dir]

commands:
  quality   exact matching, Fellegi-Sunter and chartlink on the same records
  ablate    chartlink with one idea switched off at a time
  all       everything above, written to -out
`)
	os.Exit(2)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	cmd := os.Args[1]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	out := fs.String("out", "", "directory to write markdown and charts to (default: print markdown)")
	fs.Parse(os.Args[2:])
	cfg := match.DefaultConfig()
	switch cmd {
	case "quality":
		emit(*out, "quality.md", func(w io.Writer) { quality(w, cfg) })
	case "ablate":
		emit(*out, "ablation.md", func(w io.Writer) { ablate(w, cfg) })
	case "all":
		if *out == "" {
			*out = "results"
		}
		emit(*out, "quality.md", func(w io.Writer) { quality(w, cfg) })
		emit(*out, "ablation.md", func(w io.Writer) { ablate(w, cfg) })
	default:
		usage()
	}
}

// emit writes a report to dir/name, or to stdout when dir is empty.
func emit(dir, name string, f func(io.Writer)) {
	if dir == "" {
		f(os.Stdout)
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fatal(err)
	}
	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		fatal(err)
	}
	f(file)
	if err := file.Close(); err != nil {
		fatal(err)
	}
	fmt.Fprintln(os.Stderr, "wrote", filepath.Join(dir, name))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "chartlink:", err)
	os.Exit(1)
}

// runAll runs every contender on every population, populations in parallel.
func runAll(cs []lab.Contender) [][]lab.Outcome {
	out := make([][]lab.Outcome, len(cs))
	for i := range out {
		out[i] = make([]lab.Outcome, populations)
	}
	var wg sync.WaitGroup
	for s := range populations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d := lab.Data(uint64(s + 1))
			for i, c := range cs {
				out[i][s] = lab.Run(c, d)
			}
		}()
	}
	wg.Wait()
	return out
}

func mean(xs []lab.Outcome, f func(lab.Outcome) float64) float64 {
	t := 0.0
	for _, x := range xs {
		t += f(x)
	}
	return t / float64(len(xs))
}

func describe(w io.Writer) {
	d := lab.Data(1)
	cfg := synth.DefaultConfig()
	fmt.Fprintf(w, "Each population has %d people in households, about %d records from five sources (hospital, clinic, urgent care, lab, pharmacy), each with its own rates of typos, nicknames, wrong birth dates and missing fields. Numbers are means over %d populations.\n\n",
		cfg.People, len(d.Records)/100*100, populations)
}

func quality(w io.Writer, cfg match.Config) {
	cs := lab.Contenders(cfg)
	res := runAll(cs)
	fmt.Fprintln(w, "## Linking quality")
	fmt.Fprintln(w)
	describe(w)
	fmt.Fprintln(w, "An overlay is a person whose records were merged into someone else's chart: their allergies and results show up on the wrong patient. A duplicate is an extra chart for one person: their history is split. Both are counted per 1000 people. Pairwise precision and recall count record pairs.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| matcher | overlays per 1000 | duplicates per 1000 | people with a split chart | pairwise precision | pairwise recall | review queue |")
	fmt.Fprintln(w, "|---|---:|---:|---:|---:|---:|---:|")
	for i, c := range cs {
		o := res[i]
		fmt.Fprintf(w, "| %s | %.2f | %.1f | %.2f%% | %.4f | %.4f | %.0f |\n", c.Name,
			mean(o, func(x lab.Outcome) float64 { return per1000(x.Before.Overlays, x.Before.People) }),
			mean(o, func(x lab.Outcome) float64 { return per1000(x.Before.Duplicates, x.Before.People) }),
			mean(o, func(x lab.Outcome) float64 { return 100 * x.Before.DuplicatePeople }),
			mean(o, func(x lab.Outcome) float64 { return x.Before.Precision }),
			mean(o, func(x lab.Outcome) float64 { return x.Before.Recall }),
			mean(o, func(x lab.Outcome) float64 { return float64(x.Reviews) }))
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "### Hard cases")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Twins and fathers and sons with the same name must stay apart; people who married, moved or go by their middle name must come together.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| matcher | twin pairs merged | Jr/Sr pairs merged | married, split | moved, split | goes by middle name, split |")
	fmt.Fprintln(w, "|---|---:|---:|---:|---:|---:|")
	for i, c := range cs {
		o := res[i]
		share := func(f func(score.Report) (int, int)) string {
			return pooled(o, func(x lab.Outcome) (int, int) { return f(x.Before) })
		}
		fmt.Fprintf(w, "| %s | %s | %s | %s | %s | %s |\n", c.Name,
			share(func(r score.Report) (int, int) { return r.TwinsMerged, r.TwinPairs }),
			share(func(r score.Report) (int, int) { return r.JrSrMerged, r.JrSrPairs }),
			share(func(r score.Report) (int, int) { return r.MarriedSplit, r.Married }),
			share(func(r score.Report) (int, int) { return r.MovedSplit, r.Moved }),
			share(func(r score.Report) (int, int) { return r.ByMiddleSplit, r.ByMiddle }))
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "### After review")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "A reviewer who is always right works through each matcher's review queue and merges the pairs of charts that are one person. That is the best case for the queue; a real reviewer is slower and sometimes wrong.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| matcher | charts reviewed per 1000 records | of those, one person | overlays per 1000 after | duplicates per 1000 after | goes by middle name, split after |")
	fmt.Fprintln(w, "|---|---:|---:|---:|---:|---:|")
	for i, c := range cs {
		o := res[i]
		same := 0.0
		if r := mean(o, func(x lab.Outcome) float64 { return float64(x.Reviews) }); r > 0 {
			same = mean(o, func(x lab.Outcome) float64 { return float64(x.ReviewsSame) }) / r
		}
		fmt.Fprintf(w, "| %s | %.1f | %.0f%% | %.2f | %.1f | %s |\n", c.Name,
			mean(o, func(x lab.Outcome) float64 { return per1000(x.Reviews, x.Before.Records) }),
			100*same,
			mean(o, func(x lab.Outcome) float64 { return per1000(x.After.Overlays, x.After.People) }),
			mean(o, func(x lab.Outcome) float64 { return per1000(x.After.Duplicates, x.After.People) }),
			pooled(o, func(x lab.Outcome) (int, int) { return x.After.ByMiddleSplit, x.After.ByMiddle }))
	}
}

// pooled adds up a count and its total over every population and gives
// the share, so a rare case isn't averaged over populations that have
// three of them.
func pooled(xs []lab.Outcome, f func(lab.Outcome) (int, int)) string {
	a, b := 0, 0
	for _, x := range xs {
		n, d := f(x)
		a, b = a+n, b+d
	}
	return fmt.Sprintf("%.1f%%", 100*float64(a)/float64(b))
}

func per1000(n, of int) float64 { return 1000 * float64(n) / float64(of) }

func ablate(w io.Writer, cfg match.Config) {
	cs := lab.Ablations(cfg)
	res := runAll(cs)
	fmt.Fprintln(w, "## Ablation")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "chartlink with one idea switched off at a time, on the same %d populations as the quality table.\n\n", populations)
	fmt.Fprintln(w, "| variant | overlays per 1000 | duplicates per 1000 | twin pairs merged | Jr/Sr pairs merged | review queue | duplicates per 1000 after review |")
	fmt.Fprintln(w, "|---|---:|---:|---:|---:|---:|---:|")
	for i, c := range cs {
		o := res[i]
		fmt.Fprintf(w, "| %s | %.2f | %.1f | %s | %s | %.0f | %.1f |\n", c.Name,
			mean(o, func(x lab.Outcome) float64 { return per1000(x.Before.Overlays, x.Before.People) }),
			mean(o, func(x lab.Outcome) float64 { return per1000(x.Before.Duplicates, x.Before.People) }),
			pooled(o, func(x lab.Outcome) (int, int) { return x.Before.TwinsMerged, x.Before.TwinPairs }),
			pooled(o, func(x lab.Outcome) (int, int) { return x.Before.JrSrMerged, x.Before.JrSrPairs }),
			mean(o, func(x lab.Outcome) float64 { return float64(x.Reviews) }),
			mean(o, func(x lab.Outcome) float64 { return per1000(x.After.Duplicates, x.After.People) }))
	}
}
