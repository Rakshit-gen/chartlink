## Blocking

Comparing every pair of 43043 records would be 926 million comparisons. Blocking compares only records that share a key in at least one pass. A true pair that no pass proposes can never be linked, so this is a ceiling on recall. First population only.

| pass | candidate pairs | true pairs found | found by no other pass |
|---|---:|---:|---:|
| last name sound + birth year | 94640 | 37898 | 1039 |
| date of birth + first initial | 44116 | 40890 | 1368 |
| phone | 27146 | 23421 | 426 |
| first name sound + zip | 76357 | 20736 | 216 |
| last 4 of SSN + birth month | 4868 | 4080 | 21 |

True pairs found by at least one pass: 99.13%. Candidate pairs across passes, before removing repeats: 247127, about 5.7 per record.
