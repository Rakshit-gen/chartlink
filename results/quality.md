## Linking quality

Each population has 20000 people in households, about 43000 records from five sources (hospital, clinic, urgent care, lab, pharmacy), each with its own rates of typos, nicknames, wrong birth dates and missing fields. Numbers are means over 5 populations.

An overlay is a person whose records were merged into someone else's chart: their allergies and results show up on the wrong patient. A duplicate is an extra chart for one person: their history is split. Both are counted per 1000 people. Pairwise precision and recall count record pairs.

| matcher | overlays per 1000 | duplicates per 1000 | people with a split chart | pairwise precision | pairwise recall | review queue |
|---|---:|---:|---:|---:|---:|---:|
| exact first, last, DOB | 0.05 | 354.2 | 28.58% | 0.9999 | 0.6040 | 0 |
| Fellegi-Sunter with EM | 29.27 | 28.3 | 2.79% | 0.9119 | 0.9766 | 957 |
| chartlink | 0.22 | 34.7 | 3.37% | 0.9997 | 0.9660 | 1041 |

### Hard cases

Twins and fathers and sons with the same name must stay apart; people who married, moved or go by their middle name must come together.

| matcher | twin pairs merged | Jr/Sr pairs merged | married, split | moved, split | goes by middle name, split |
|---|---:|---:|---:|---:|---:|
| exact first, last, DOB | 0.0% | 0.0% | 43.3% | 28.4% | 45.0% |
| Fellegi-Sunter with EM | 65.6% | 59.7% | 6.3% | 8.5% | 7.6% |
| chartlink | 1.5% | 0.2% | 7.0% | 7.5% | 33.4% |

### After review

A reviewer who is always right works through each matcher's review queue and merges the pairs of charts that are one person. That is the best case for the queue; a real reviewer is slower and sometimes wrong.

| matcher | charts reviewed per 1000 records | of those, one person | overlays per 1000 after | duplicates per 1000 after | goes by middle name, split after |
|---|---:|---:|---:|---:|---:|
| exact first, last, DOB | 0.0 | 0% | 0.05 | 354.2 | 45.0% |
| Fellegi-Sunter with EM | 22.1 | 43% | 29.26 | 7.9 | 2.0% |
| chartlink | 24.1 | 54% | 0.22 | 7.1 | 1.9% |
