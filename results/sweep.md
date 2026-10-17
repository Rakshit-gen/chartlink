## Moving the merge threshold

Each matcher merges a pair when the model gives it at least this probability of being one person. Means over 5 populations, before review.

| matcher | merge at | overlays per 1000 | duplicates per 1000 |
|---|---:|---:|---:|
| Fellegi-Sunter | 0.5 | 55.58 | 7.3 |
| Fellegi-Sunter | 0.9 | 43.11 | 20.1 |
| Fellegi-Sunter | 0.99 | 29.27 | 28.3 |
| Fellegi-Sunter | 0.999 | 26.47 | 33.8 |
| Fellegi-Sunter | 0.9999 | 21.90 | 72.4 |
| Fellegi-Sunter | 0.99999 | 10.28 | 180.5 |
| chartlink | 0.5 | 4.57 | 17.0 |
| chartlink | 0.9 | 0.78 | 24.6 |
| chartlink | 0.99 | 0.22 | 34.7 |
| chartlink | 0.999 | 0.15 | 44.3 |
| chartlink | 0.9999 | 0.12 | 89.2 |
| chartlink | 0.99999 | 0.07 | 182.1 |
