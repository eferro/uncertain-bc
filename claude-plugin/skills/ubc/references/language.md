# ubc language and flags

A compact reference for writing and running models. The full reference, with every error message, is at <https://eferro.github.io/uncertain-bc/docs/REFERENCE/>. Every console block below is checked against the binary by the project's test suite.

## Contents

- Running ubc
- Numbers, operators and statements
- Ranges, variables and functions
- Distributions and observed data
- Comparisons, events and probabilities
- Summaries and maths
- Flags
- Output and exit codes
- Errors worth knowing

## Running ubc

- `ubc model.ubc` runs a file; several files run in order in one session, so variables and functions flow from one to the next. Flags go before files.
- `ubc -e 'EXPR'` runs an expression; repeat `-e` for more. Every `-e`, file or stdin is its own program: a program that prints one result prints it bare, one that prints several labels each line with the variable name, the expression when it is 30 characters or shorter, or `L<line>:`.
- `echo 'EXPR' | ubc` reads stdin. `ubc` alone starts the REPL; avoid it in automation.
- The seed is fixed (1), so the same program prints the same numbers every time.

## Numbers, operators and statements

- Numbers: `42`, `3.14`, `.5`, `1e6`. Suffixes right after the number: `k`, `M`, `B`, `T`, `%` (`5%` is `0.05`).
- Operators: `+ - * / ^` and parentheses. From loosest to tightest: comparisons, `to`, `+ -`, `* /`, leading `-`, `^` (right-associative, so `-2 ^ 2` is `-4`).
- Statements are separated by new lines or `;`. `#` starts a comment.
- Names: letters, digits and `_`, starting with a letter or `_`; case-sensitive; ASCII only. Built in: `pi`, `e`, `last` (the previous result).

## Ranges, variables and functions

- `a to b`: 90% sure the value is between `a` and `b`. Lognormal when both ends are positive, normal otherwise. `a` must be below `b`, and both must be plain numbers.
- `name = expr` binds a variable that keeps its samples: `x - x` is exactly 0, while `(1 to 5) - (1 to 5)` is two independent draws.
- `name(p1, p2) = expr` defines a function. Parameters shadow variables; ranges in the body are drawn again on every call; no conditionals and no recursion (more than 64 nested calls is an error).

```console
$ ubc -e 'task(days) = days * (0.8 to 1.6); task(5) + task(3) + task(8)'
≈ 18.4   90% [14.8 … 22.7]   98% [13 … 26]   ▂▃▄▆████▇▅▄▃▂▂▁▁
```

## Distributions and observed data

Arguments are plain numbers; shift or scale a fixed distribution to combine it with a range, as in `(10 to 20) + normal(0, 1)`.

- `normal(mean, sd)`, `lognormal(mu, sigma)` (of the logarithm), `uniform(lo, hi)`.
- `triangular(lo, mode, hi)`, `pert(lo, mode, hi)` for three-point estimates (most likely value plus hard limits).
- `beta(a, b)` for shares between 0 and 1; `exponential(rate)` for waiting times; `gamma(shape, scale)`.
- `poisson(lambda)` for a count of independent events in a period; `binomial(n, p)` for successes in `n` tries; `binomial(1, p)` is a single yes/no event.
- `[v1, v2, ...]`: observed values, resampled with replacement. `--data name=file.txt` does the same from a file with one plain number per line.

## Comparisons, events and probabilities

- `< > <= >= == !=`. Between plain numbers they give `1` or `0`. With an uncertain value they give an event: `1` or `0` in every sample, kept wherever it goes (assigned, returned by a function, given to `--set`) and shown as the probability that it holds.
- `P(cond)` is that probability as a plain number. `P(x)` of something that is not a comparison is an error.
- `event * value` is 0 in the samples where the event did not happen. `P(event) * value` is a different question: the expected value spread over every sample.

```console
$ ubc -e 'x = 10 to 20; slow = x > 15; slow; slow * 100; P(slow)'
slow = 0.39
slow * 100 = ≈ 0   90% [0 … 100]   98% [0 … 100]   █▁▁▁▁▁▁▁▁▁▁▁▁▁▁▅
P(slow) = 0.39
```

## Summaries and maths

- To a plain number: `mean(x)`, `median(x)`, `stdev(x)`, `min(x)`, `max(x)`, `p(x, q)` (the `q` quantile, `p(x, 0.9)` is the 90th percentile).
- Sample by sample: `sqrt`, `exp`, `log` (natural), `log10`, `abs`, `round`, `floor`, `ceil`; `min(a, b, ...)` and `max(a, b, ...)` with two or more arguments, so `max(0, x)` cuts the negative part.

## Flags

- `--echo` prints assignments too, so every step of a model shows its uncertainty.
- `--set 'name = expr'` replaces the model's value of `name`; everything computed from it follows. Repeatable. `--data name=file` does the same with observed data.
- `--sens expr` ranks the uncertain variables by how strongly they move with `expr` (rank correlation, strongest first). Not with `--json` or `--plain`.
- `--assert 'cond'` exits 4 unless `cond` holds in every sample; state looser checks as probabilities: `--assert 'P(cost > 100k) < 0.1'`.
- `--scenarios` lists the built-in catalogue; `--scenario NAME` runs one (after the shared formulas, before `-e` and files, so `--set`, `--sens` and the rest apply); `--scenario-source NAME` prints its source; `--scenario formulas` loads only the shared formulas and conversion constants.
- `--json` prints one object per result: `kind` (`samples`, `scalar`, `probability`), `name` when labelled, and for samples `median`, `mean`, `stdev`, `p1`, `p5`, `p95`, `p99`, `n`, `dropped`. `--plain` prints only the median.
- `--samples-out file.csv` (or `.json`) writes every variable's samples, row *i* being sample *i* of each, for pandas or a spreadsheet.
- `-n N` sample count (default 5000), `-s N` seed, `--random` clock seed.
- `--interval 80` prints an 80% interval instead of the 90% one; `--quantiles 10,50,90` prints those percentiles instead; `--hist` draws a histogram; `--ascii` avoids Unicode.

```console
$ ubc --set 'vps = 5 to 10' examples/homelab.ubc
homelab = ≈ 67   90% [46.6 … 96.6]   98% [40 … 115]   ▂▄▆▇██▇▆▅▄▃▂▂▂▁▁
vps = ≈ 7.08   90% [5 … 9.91]   98% [4.3 … 11.6]   ▂▄▅▆▇██▇▅▅▃▃▂▂▁▁
saving = ≈ -59.7   90% [-89.3 … -39.1]   98% [-108 … -33]   ▁▁▂▂▂▃▄▅▇▇██▇▅▄▂
```

## Output and exit codes

- Uncertain: `≈ median   90% [p5 … p95]   98% [p1 … p99]   sparkline`, 3 significant figures (more when the interval is narrow), `k`/`M`/`B`/`T` from 10k up. Probability: two decimals, `<0.01` and `>0.99` at the ends. Plain number: up to 10 significant figures.
- Non-finite samples (a division by zero) are left out and counted: `(N non-finite dropped)`.
- Exit codes: `0` success, `1` evaluation error, `2` usage error, `3` parse error, `4` an `--assert` failed. Errors go to stderr as `ubc: [file: ]line L, col C: message`.

## Errors worth knowing

```console
$ ubc -e '(1 to 2) to 3'
ubc: line 1, col 10: range bounds must be plain numbers
$ ubc -e 'normal(10 to 20, 1)'
ubc: line 1, col 1: normal expects plain numbers
$ ubc -e '10 to 5'
ubc: line 1, col 4: range start 10 must be below its end 5
$ ubc -e 'P(10 to 20)'
ubc: line 1, col 1: P expects a comparison
```
