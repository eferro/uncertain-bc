# How-to guides

Short recipes for common tasks, one section each. They assume you know the basics of the language; if not, start with the [tutorial](TUTORIAL.md). The [reference](REFERENCE.md) has every flag in detail, and the [modelling guide](MODELLING.md) explains how to choose good ranges.

Every command shown in a console block is run by the test suite, which compares its output with this page, so you get the same numbers with the default seed. Run the commands from the root of the repository, where `examples/` lives. Lines that start with `ubc: ` go to stderr.

- [Use my own numbers](#use-my-own-numbers)
- [Compare two options](#compare-two-options)
- [Find which assumption to research first](#find-which-assumption-to-research-first)
- [Use ubc as a check in CI](#use-ubc-as-a-check-in-ci)
- [Take the samples to pandas or a spreadsheet](#take-the-samples-to-pandas-or-a-spreadsheet)
- [Forecast a delivery date from the team's real throughput](#forecast-a-delivery-date-from-the-teams-real-throughput)
- [Run and adapt a scenario from the catalogue](#run-and-adapt-a-scenario-from-the-catalogue)

## Use my own numbers

A model is written with someone's guesses. You do not need to edit it to try yours: `--set NAME=EXPR` replaces the model's value of `NAME` with any expression, a plain number or a range. This is the project estimate of the tutorial, first as written and then with a smaller share of forgotten work:

```console
$ ubc examples/project-estimate.ubc
days = ≈ 38.8   90% [29.6 … 51.8]   98% [26 … 60]   ▂▄▅▇███▇▆▅▃▃▂▂▂▁
P(days > 50) = 0.07
round(p(days, 0.85)) = 46
$ ubc --set 'forgotten = 10% to 20%' examples/project-estimate.ubc
days = ≈ 36.7   90% [28.4 … 47.6]   98% [25 … 54]   ▂▃▄▅▇▇█▇▆▅▄▃▂▂▂▂
P(days > 50) = 0.03
round(p(days, 0.85)) = 43
```

Quote the whole `NAME=EXPR` when it has spaces. Repeat `--set` to change several values; the model's own lines for those names are skipped, and everything that depends on them is computed again.

When you have measurements instead of a guess, put them in a file, one number per line, and load it with `--data NAME=FILE`. The variable becomes your data: every sample picks one of your numbers at random, as the list `[4, 6, 3, …]` would.

```console
$ cat examples/throughput.txt
4
6
3
5
5
2
6
4
7
3
5
4
$ ubc --data weekly=examples/throughput.txt -e weekly
≈ 5   90% [2 … 7]   98% [2 … 7]   ▃▁▁▅▁▁█▁▁█▁▁▅▁▁▃
```

`--data` overrides like `--set` does, so it replaces a guess in a model too; see [the delivery forecast](#forecast-a-delivery-date-from-the-teams-real-throughput) below.

## Compare two options

To choose between two options, ask how likely one is to beat the other: `P(a < b)` is the share of samples where `a` came out below `b`. Building a tool in-house might cost 40k to 120k; buying one, 60k to 80k:

```console
$ ubc -e 'build = 40k to 120k' -e 'buy = 60k to 80k' -e 'P(build < buy)' -e 'P(build > 100k)' -e 'P(buy > 100k)'
0.50
0.13
0
```

Building is cheaper half of the time, so on the central value it is a coin toss; but building has a 13 % chance of going over 100k and buying has none. Comparing the medians alone would have hidden that difference in risk.

`P(a < b)` pairs sample *i* of `a` with sample *i* of `b`. If both options depend on the same uncertain quantity, such as the number of users, assign it to a variable once and use the variable in both, so that the comparison is made under the same conditions; the [modelling guide](MODELLING.md#3-ranges-are-independent-variables-are-not) explains why. The [buy vs build scenario](../scenarios/product/buy-vs-build.ubc) is a fuller version of this comparison.

## Find which assumption to research first

Some ranges matter much more than others. `--sens EXPR` ranks the variables by how strongly the result moves with them, strongest first:

```console
$ ubc --sens days examples/project-estimate.ubc
days = ≈ 38.8   90% [29.6 … 51.8]   98% [26 … 60]   ▂▄▅▇███▇▆▅▃▃▂▂▂▁
P(days > 50) = 0.07
round(p(days, 0.85)) = 46
sensitivity of days (rank correlation):
  work       +0.88  ██████████████████
  backend    +0.63  █████████████
  forgotten  +0.43  █████████
  frontend   +0.39  ████████
  migration  +0.28  ██████
  testing    +0.25  █████
  design     +0.17  ███
```

`work` is the sum of the tasks, not an assumption, so skip it: the input to research first is `backend`. Spend an hour breaking it down, and the narrower range you get narrows the whole estimate:

```console
$ ubc --set 'backend = 10 to 14' examples/project-estimate.ubc
days = ≈ 37.4   90% [30.3 … 47.8]   98% [27 … 53]   ▂▃▅▆▇█▇▇▆▄▃▂▂▂▂▁
P(days > 50) = 0.02
round(p(days, 0.85)) = 43
```

Researching `design` instead would have changed almost nothing. In the REPL, `:sens days` prints the same ranking.

## Use ubc as a check in CI

`--assert EXPR` turns a model into a test: after the program runs, `ubc` exits with code 4 unless every assertion certainly holds, and it says which one failed. A probability is compared with its value, so you can state the risk you accept:

```console
$ ubc --assert 'P(days > 50) < 0.1' examples/project-estimate.ubc
days = ≈ 38.8   90% [29.6 … 51.8]   98% [26 … 60]   ▂▄▅▇███▇▆▅▃▃▂▂▂▁
P(days > 50) = 0.07
round(p(days, 0.85)) = 46
$ ubc --assert 'P(days > 50) < 0.05' examples/project-estimate.ubc
days = ≈ 38.8   90% [29.6 … 51.8]   98% [26 … 60]   ▂▄▅▇███▇▆▅▃▃▂▂▂▁
P(days > 50) = 0.07
round(p(days, 0.85)) = 46
ubc: --assert P(days > 50) < 0.05 failed: 0.07 < 0.05
```

The first run exits with 0 and the second with 4. Any other failure has its own code (1, 2 or 3; see the [reference](REFERENCE.md#exit-codes)), so a broken model is not mistaken for a broken assumption. An assertion on an uncertain value, such as `days < 60`, holds only when it is true for every sample.

Because the seed is fixed, the same model gives the same answer on every run, and a red pipeline means the model or its inputs changed. A GitHub Actions step, for example, installs `ubc` and runs the check:

```yaml
- name: Check the delivery risk
  run: |
    curl -fsSL https://raw.githubusercontent.com/eferro/uncertain-bc/main/install.sh | sh
    ubc --assert 'P(days > 50) < 0.1' model.ubc
```

## Take the samples to pandas or a spreadsheet

`--samples-out FILE` writes the raw samples of every variable after the program runs: a CSV file with one column per variable, or a JSON object with one array per variable when `FILE` ends in `.json`. Row *i* holds sample *i* of every variable, so you can plot one against another. With only four samples, to keep the file short here:

```console
$ ubc -n 4 --samples-out samples.csv -e 'users = 20k to 60k' -e 'cost = users * 0.01'
$ cat samples.csv
users,cost
46886.069381712114,468.86069381712116
40177.29526999524,401.7729526999524
26668.21390176827,266.68213901768274
43196.684311693854,431.96684311693855
```

The results still go to the terminal; this program assigns and prints nothing. Leave out `-n` to get the usual 5000 rows, then read the file with pandas:

```python
import pandas as pd

samples = pd.read_csv("samples.csv")
samples["cost"].describe()
samples.plot.scatter(x="users", y="cost")
```

A spreadsheet opens the CSV file as it is, one sample per row; R reads it with `read.csv("samples.csv")`.

## Forecast a delivery date from the team's real throughput

Instead of estimating every item, use what the team has really finished. Count the items done in each of the last weeks, one number per line, as in `examples/throughput.txt` above, and run the [delivery-date scenario](../scenarios/delivery/delivery-date.ubc) with them:

```console
$ ubc --data throughput=examples/throughput.txt --scenario delivery/delivery-date
weeks = ≈ 10.8   90% [8.45 … 14]   98% [7.5 … 15.5]   ▂▃▅▆▇██▇▆▆▄▃▃▂▂▁
P(weeks <= deadline) = 0.74
```

The scenario divides the backlog left (40 to 60 items, as written) by the team's weekly rate, and gives the weeks to finish and the chance of finishing within the 12 weeks we were asked for. Set your own backlog and deadline with `--set`, and add `-e` to ask your own question after the scenario, such as the number of weeks you can promise with 85 % confidence:

```console
$ ubc --data throughput=examples/throughput.txt --set 'backlog = 30 to 40' --set deadline=8 --scenario delivery/delivery-date -e 'p(weeks, 0.85)'
weeks = ≈ 7.67   90% [6.11 … 9.72]   98% [5.5 … 10.8]   ▂▃▄▆▇██▇▆▅▄▃▂▂▂▁
P(weeks <= deadline) = 0.62
8.891761291
```

Ten or more weeks of history give a fair picture; keep only weeks that look like the ones ahead, without holidays or a team half its size.

## Run and adapt a scenario from the catalogue

`ubc --scenarios` lists the [catalogue](../scenarios/README.md) of ready-made models built into `ubc`, one per line with the question it answers. Run one by name:

```console
$ ubc --scenario team/replace-a-person
total = ≈ 61.1k   90% [37.8k … 100k]   98% [30k … 122k]   ▃▅▆██▇▇▆▅▄▃▂▂▂▁▁
P(total > 50k) = 0.75
```

To see the assumptions you can change, add `--echo`, which also prints every assignment with its name:

```console
$ ubc --echo --scenario team/replace-a-person
candidates = ≈ 9.5   90% [6 … 14.8]   98% [4.9 … 19]   ▂▄▆▇██▇▅▅▄▃▂▂▂▁▁
hours_per_candidate = ≈ 4.88   90% [3.02 … 7.96]   98% [2.4 … 9.8]   ▃▅▇██▇▇▆▄▄▃▂▂▂▁▁
interviewer_cost = ≈ 66.9   90% [50.2 … 91]   98% [44 … 103]   ▂▃▅▆▇█▇▆▆▅▄▃▂▂▂▁
hiring_fees = ≈ 7750   90% [2930 … 20.4k]   98% [2000 … 30k]   ▅▇█▇▆▄▃▃▂▂▂▁▁▁▁▁
vacancy_months = ≈ 3.15   90% [1.98 … 4.94]   98% [1.6 … 6]   ▃▄▆▇▇█▇▆▅▄▃▂▂▂▁▁
value_per_month = ≈ 9550   90% [5930 … 15.3k]   98% [4900 … 19k]   ▃▅▆▇██▇▆▅▄▃▂▂▂▁▁
ramp_up_months = ≈ 4.28   90% [3.01 … 6.05]   98% [2.6 … 7.1]   ▂▄▅█████▆▅▄▃▂▂▂▁
ramp_up_gap = ≈ 0.425   90% [0.302 … 0.599]   98% [0.25 … 0.68]   ▂▃▅▆▇█▇▇▆▅▄▃▃▂▂▁
interviewing = ≈ 3100   90% [1500 … 6390]   98% [1100 … 8600]   ▃▆██▇▆▅▄▃▂▂▂▂▁▁▁
vacancy = ≈ 30.1k   90% [15.3k … 57.3k]   98% [11k … 75k]   ▃▅▇█▇▇▅▄▄▃▂▂▂▁▁▁
ramp_up = ≈ 17.3k   90% [8770 … 33.8k]   98% [6900 … 45k]   ▄▆███▆▅▄▃▃▂▂▂▁▁▁
total = ≈ 61.1k   90% [37.8k … 100k]   98% [30k … 122k]   ▃▅▆██▇▇▆▅▄▃▂▂▂▁▁
total = ≈ 61.1k   90% [37.8k … 100k]   98% [30k … 122k]   ▃▅▆██▇▇▆▅▄▃▂▂▂▁▁
P(total > 50k) = 0.75
```

The first eight are the assumptions, the rest are calculations. Replace the ones you know better with `--set`, as in [Use my own numbers](#use-my-own-numbers); a faster hiring process, for example:

```console
$ ubc --set 'vacancy_months = 1 to 3' --scenario team/replace-a-person
total = ≈ 47.2k   90% [29.4k … 76.8k]   98% [24k … 95k]   ▃▅▇██▇▇▆▄▄▃▂▂▂▂▁
P(total > 50k) = 0.43
```

`--sens`, `--data`, `--assert` and `--samples-out` work with `--scenario` as they do with a file. When the change is bigger than a value, such as a new cost or another question, copy the scenario and edit it: `ubc --scenario-source team/replace-a-person > mine.ubc` writes its source with all its comments, and `ubc --scenario formulas mine.ubc` runs your copy with the shared formulas the scenarios use.
