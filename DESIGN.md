# How chartlink decides

This is the long version of the README's "the idea" section: how a pair of records turns into a merge, a review, or nothing, and what goes wrong without each piece. Numbers come from `results/`.

## The records

Every hospital, clinic, lab and pharmacy registers patients on its own. The same person turns up as ROBERT J SMITH at the hospital and BOB SMYTH at urgent care. The synthetic data in `synth/` has the mistakes registration desks actually make:

- typos: a neighbouring key, a dropped or doubled letter, two letters swapped
- nicknames (BOB for ROBERT) and people registered under their middle name
- wrong birth dates: day and month swapped, a digit off, the current year typed by habit
- a maiden name on older records, an old address after a move
- missing middle names, phones, addresses and SSNs, at different rates per source

And it has the people who make matching dangerous: twins (same last name, birthday and home), fathers and sons with the same name, siblings in one household, and common names. About a third of surnames are rare, made up so almost nobody shares them; the rest come from a list of common US surnames. SMITH is about 2.4% of people, around three times the real US rate, so name collisions are harder here than in real data.

## Step 1: block

Comparing every pair of 43,000 records is 926 million comparisons. Instead, five blocking passes each give a record a key, and only records sharing a key in some pass are compared:

| pass | there for |
|---|---|
| last name Soundex + birth year | the workhorse; typos usually keep the Soundex code |
| full DOB + first initial | a last name that changed with marriage |
| phone | a birth date typed wrong |
| first name Soundex + zip | a changed last name and a wrong DOB, at one address |
| last 4 of SSN + birth month | everything else wrong, SSN right |

Together they find 99.13% of true pairs with about 5.7 candidates per record. A block with more than 200 records is skipped: a shared clinic phone number says nothing and costs thousands of comparisons. That skip is also why recall drops at a million records, when more blocks pass 200.

## Step 2: compare

Each candidate pair gets a level per field: first name exact, nickname or middle name, close (Jaro-Winkler at least 0.88), or different; DOB exact, day and month swapped, one part off, or different; and so on. A field missing on either side is no evidence either way.

## Step 3: weigh (Fellegi-Sunter)

Every level has two probabilities: M, how often the same person's two records land there, and U, how often two random records do. A pair's weight is the sum of log2(M/U) over its fields, in bits.

- **U comes from random pairs.** Almost every random pair is two people, so how often random pairs agree is how often non-matches agree.
- **M comes from EM over the candidate pairs.** Nobody labels which pairs match, so expectation maximization alternates between "how likely is each pair a match" and "given that, what are M and U". The U this fits is thrown away: candidate pairs share a blocking key and agree far more often than random pairs, so the random-pair U is the honest one.
- **The prior** is the chance two random records are one person. A pair merges when the posterior is at least 0.99. In an index that keeps growing, the prior shrinks with its size, since each record has about the same number of true partners however many records there are.

## Step 4: what plain Fellegi-Sunter gets wrong

On these records, Fellegi-Sunter with EM finds 97.7% of true pairs and puts 29 people per 1000 on someone else's chart. It merges 66% of twin pairs and 60% of fathers and sons.

The reason is that Fellegi-Sunter assumes fields are independent given match status, and a household breaks that. Two siblings agree on last name, street, zip and home phone. The model adds all four as separate evidence, and together they outweigh what the different first names and birth dates take away. EM makes it worse: with households in the candidate pairs, its "match" class learns to look like a household.

## Step 5: count the household once

Street, zip and phone are counted as one piece of evidence: only the largest positive weight among them counts, though a disagreement still counts in full. This is the biggest single fix. Without it, overlays go from 0.22 to 7.4 per 1000.

## Step 6: weigh names by rarity

Agreeing on SMITH is weak evidence; agreeing on KOVALENKO is strong. For an exact name match, the weight uses the share of records with that name instead of U. This finds more people (duplicates 38.7 to 34.7 per 1000) and lets a few more wrong ones through (overlays 0.16 to 0.22). It's a trade, not a pure win, and the ablation table shows both sides.

## Step 7: the chart guard

Charts grow one pair at a time, strongest first. Before two charts are joined, every pair of records across them is checked for a conflict:

- **Twins:** the same birthday (or one typo away) and first names more than two edits apart, even going through nicknames. Two edits, because both records can carry a typo: LIAS and LSIA are both LISA.
- **Father and son:** JR against SR, or birth years at least 12 apart on a different day of the year. Some desks type the current year for the birth year, which keeps month and day, so a different year alone looks like that mistake.
- **SSN:** two last-4s that differ by more than one keying mistake.

One conflict anywhere vetoes the merge. Checking every cross pair, not just the pair being merged, matters because charts chain: A links to B and B links to C, though A and C are twins. Without the guard, 43% of twin pairs merge.

My first version also had the twin and father-son rules inside the weights, zeroing the evidence twins share. The ablation showed the guard covered everything that did, while the zeroing pushed people who go by their middle name below the review threshold. They look exactly like twins: same birthday, same home, different first name. So I took it out. With the guard alone, those pairs keep their high probability, the guard stops the merge, and the pair goes to review.

## Step 8: review

A pair between 0.5 and 0.99, or above 0.99 but stopped by the guard, goes to a review queue, one entry per pair of charts. On the default population that's about 24 reviews per 1000 records, and 54% of them turn out to be one person. A reviewer working through the queue brings duplicates from 34.7 to 7.1 per 1000 and people who go by their middle name from 33% split to 1.9%.

## What it costs

Batch linking of 43,000 records takes 0.78 s. Linking one record into an index takes 6 µs at the median with 43,000 records and 25 µs with a million. Memory is about 500 bytes per record. See `results/bench.md`.

## What this design doesn't do

- **The data is synthetic.** The error rates are my guesses at what registration desks do, not measured from a real hospital. A real deployment should fit its model on its own records and check a labelled sample.
- **One reviewer who is always right** is the best case for the queue.
- **Twins named ANA and ANNA merge.** Names within two edits look like typos.
- **Charts are never split.** Once records are merged, neither Link nor Index takes them apart again, even if later records show the merge was wrong.
- **Blocks over 200 records are skipped.** At a million records that costs recall: duplicates rise from 31 to 46 per 1000. Bigger populations need finer keys.
