package main

import (
	"fmt"
	"io"
	"runtime"
	"slices"
	"time"

	"github.com/Rakshit-gen/chartlink/match"
	"github.com/Rakshit-gen/chartlink/score"
	"github.com/Rakshit-gen/chartlink/synth"
)

// bench streams populations of three sizes through an index fitted on the
// default population, the way a master patient index takes a registration
// feed, and times each Add.
func bench(w io.Writer, cfg match.Config) {
	model := match.Fit(synth.Records(synth.People(synth.DefaultConfig(), 1), 1001).Records, cfg)
	fmt.Fprintln(w, "## Cost of linking")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Machine: %s, %s. One goroutine. The model is fitted once on the default population; each population below is then fed through an empty index one record at a time, and every Add is timed. Quality is scored on the whole stream at the end.\n\n", cpuName(), runtime.Version())
	fmt.Fprintln(w, "| records | people | Add p50 | Add p99 | records per second | index memory | overlays per 1000 | duplicates per 1000 |")
	fmt.Fprintln(w, "|---:|---:|---:|---:|---:|---:|---:|---:|")
	for _, people := range []int{20000, 100000, 500000} {
		pc := synth.DefaultConfig()
		pc.People = people
		pc.Zips = people / 100 // keep about as many people per zip
		d := synth.Records(synth.People(pc, 7), 8)

		runtime.GC()
		var before, after runtime.MemStats
		runtime.ReadMemStats(&before)
		ix := match.NewIndex(model, cfg)
		lat := make([]time.Duration, len(d.Records))
		start := time.Now()
		for i, r := range d.Records {
			t := time.Now()
			ix.Add(r)
			lat[i] = time.Since(t)
		}
		total := time.Since(start)
		runtime.GC()
		runtime.ReadMemStats(&after)
		mem := float64(after.HeapAlloc-before.HeapAlloc) / 1e6
		rep := score.Score(d, ix.Charts())
		runtime.KeepAlive(ix)

		slices.Sort(lat)
		fmt.Fprintf(w, "| %d | %d | %s | %s | %.0f | %.0f MB | %.2f | %.1f |\n",
			len(d.Records), len(d.People), us(lat[len(lat)/2]), us(lat[len(lat)*99/100]),
			float64(len(d.Records))/total.Seconds(), mem,
			per1000(rep.Overlays, rep.People), per1000(rep.Duplicates, rep.People))
	}
}

func us(d time.Duration) string { return fmt.Sprintf("%.1f µs", float64(d.Nanoseconds())/1e3) }
