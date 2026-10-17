## Cost of linking

Machine: Apple M4, go1.25.4. One goroutine. The model is fitted once on the default population; each population below is then fed through an empty index one record at a time, and every Add is timed. Quality is scored on the whole stream at the end.

| records | people | Add p50 | Add p99 | records per second | index memory | overlays per 1000 | duplicates per 1000 |
|---:|---:|---:|---:|---:|---:|---:|---:|
| 43176 | 20000 | 6.1 µs | 26.3 µs | 126993 | 25 MB | 0.50 | 31.1 |
| 215898 | 100002 | 9.8 µs | 60.0 µs | 73286 | 113 MB | 0.37 | 36.5 |
| 1079349 | 500000 | 24.5 µs | 177.6 µs | 26332 | 533 MB | 0.49 | 46.3 |
