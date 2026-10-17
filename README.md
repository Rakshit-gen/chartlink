# chartlink

Link patient records from different hospitals, clinics, labs and pharmacies into one chart per person, without putting twins, fathers and sons, or two people named Nancy Brown into the same chart.

On synthetic populations of 20,000 people with realistic registration mistakes, plain Fellegi-Sunter (the textbook probabilistic matcher) put 29 people per 1000 on someone else's chart. chartlink put 0.22 per 1000 there, about 130 times fewer, while still finding 96.6% of the record pairs that belong together. After a review queue of about 24 charts per 1000 records, it ends with fewer split charts than Fellegi-Sunter too.

## The problem

The US has no national patient identifier. Each system that registers a patient makes its own record, from whatever the patient says and the desk types: ROBERT J SMITH at the hospital, BOB SMYTH at urgent care, the maiden name at the old clinic. A master patient index has to decide which of these are one person.

It can go wrong two ways, and they are not equally bad:

- **A duplicate** is one patient with two charts. Their history is split: a doctor misses an allergy recorded at the other clinic, a test is done twice.
- **An overlay** is two patients in one chart. One person's results, medications and allergies show up on the other's chart. It's the dangerous one, and it's slow and expensive to undo.

Exact matching on name and birth date makes almost no overlays and leaves 28% of people with a split chart. Probabilistic matching (Fellegi-Sunter) finds nearly everyone, and on real households it overlays twins, siblings, and fathers and sons with the same name. The reason is in its core assumption: fields are taken as independent evidence. Twins agree on last name, birthday, street, zip and home phone, and the model counts all five.

## The idea

Fellegi-Sunter with EM, plus three things. [DESIGN.md](DESIGN.md) has the long version.

1. **A household is one piece of evidence.** Street, zip and phone agreeing are what everyone in a household shares. Only the strongest of them counts. This alone takes overlays from 7.4 to 0.22 per 1000.
2. **Rare names count for more.** Agreeing on SMITH is weak, agreeing on KOVALENKO is strong. Exact name matches are weighed by how many records have that name.
3. **A guard on every merge.** Before two charts join, every pair of records across them is checked for a hard conflict: twins (same birthday, first names more than two edits apart even through nicknames), a father and son (JR against SR, or birth years a generation apart on different days), or two different SSNs. One conflict anywhere stops the merge, which stops chains: A matches B, B matches C, but A and C are twins.

Pairs that are likely but not certain, and pairs the guard stopped, go to a review queue, one entry per pair of charts. Someone who goes by their middle name looks exactly like a twin, so a person decides.

It runs as a batch over a whole set of records or as an index that links each record as it arrives.

## Usage

Link a CSV of records into charts:

```
go run ./cmd/chartlink export -people 2000 > records.csv
go run ./cmd/chartlink link -review review.csv < records.csv > charts.csv
```

```
4291 records, 2056 charts, 136 pairs for review
```

The input needs `first`, `last` and `dob` (YYYY-MM-DD) columns. It can also have `id`, `source`, `middle`, `suffix`, `sex`, `phone`, `street`, `zip` and `ssn4`. The output has one row per record with the chart it belongs to, named after one of its records. `link` fits its model on the records it is given, so give it a few thousand at least.

Run the experiments and write every table and chart:

```
go run ./cmd/chartlink all -out results
```

The commands one at a time are `quality`, `ablate`, `sweep`, `blocking` and `bench`. `all` takes about a minute on an M4.

As a library, linking records as they arrive:

```go
ix := match.NewIndex(match.Fit(existing, match.DefaultConfig()), match.DefaultConfig())
id, review := ix.Add(rec)
chart := ix.Charts()[id]
```

## Running the tests

```
go test -race ./...
go test -run XXX -bench . -benchtime 3x ./match
```

## Layout

```
record/        the record, name cleaning, Jaro-Winkler, Soundex, nicknames, CSV
synth/         a synthetic population and the records five sources keep about it
match/         blocking, Fellegi-Sunter with EM, the guard, batch linking, the index
baseline/      exact first name, last name and birth date
score/         overlays, duplicates, pairwise precision and recall, hard cases
lab/           the experiments, and a reviewer who works the queue
cmd/chartlink  the command line
results/       what `chartlink all` wrote
```

## License

MIT
