# Uncertain bc — Specification

The complete behaviour of `ubc`: invocation, language, semantics, output and REPL. It describes what the released binary does and is kept in sync with the code: any change in behaviour updates this file in the same commit, and any change of decision adds an ADR in [adr/](adr/README.md).

## 1. Invocation

```
ubc                         # interactive REPL when stdin is a TTY
ubc FILE...                 # evaluate files in order, print the result of each top-level expression
ubc -e EXPR [-e EXPR...]    # evaluate expressions, print results
ubc --set NAME=EXPR FILE    # run a model with NAME overridden (scenarios)
ubc --data NAME=PATH FILE   # run a model with NAME bound to observed data
ubc --assert EXPR FILE      # run a model and fail with exit 4 unless EXPR holds (CI checks)
ubc --samples-out F.csv FILE # run a model and write every variable's samples to a file
ubc --sens EXPR FILE        # run a model and rank the variables that drive EXPR's uncertainty
ubc --scenarios             # list the built-in scenarios and the question each answers
ubc --scenario NAME         # run a built-in scenario, after the shared formulas
cmd | ubc                   # read program from stdin (non-TTY), print results
```

| Flag | Meaning | Default |
|---|---|---|
| `-e, --eval EXPR` | Evaluate an expression (repeatable) | |
| `--set NAME=EXPR` | Set a variable, overriding the program's own assignments to it (repeatable, ADR 0028) | |
| `--data NAME=PATH` | Set a variable to the numbers in a file, one per line, resampled like `[...]`; overrides like `--set` (repeatable, ADR 0031) | |
| `--assert EXPR` | After the program, exit 4 unless EXPR certainly holds (repeatable, ADR 0029) | |
| `--samples-out FILE` | After the program, write the samples of every variable and `last` to FILE, `.csv` or `.json` (ADR 0032) | |
| `--sens EXPR` | After the program, list the variables that drive EXPR's uncertainty, strongest first (ADR 0033) | |
| `-n, --samples N` | Monte Carlo sample count | `5000` |
| `-s, --seed N` | RNG seed | `1` |
| `--random` | Use a time-based seed | off |
| `--json` | One JSON object per printed result | off |
| `--echo` | Also print assignments, as `name = summary` (ADR 0025) | off |
| `--hist` | Draw the 20-row histogram under every uncertain result (not with `--json`/`--plain`) | off |
| `--interval W` | Print a W% interval instead of the 90% one, e.g. `80`; the 98% interval stays; adds its ends to JSON `quantiles` (ADR 0034, ADR 0040) | `90` |
| `--quantiles LIST` | Show these percentiles instead of the 90% and 98% intervals, e.g. `10,50,90`; adds `quantiles` to JSON | |
| `--watch` | Evaluate the files again whenever they change (polls every 300 ms) | off |
| `--scenarios` | List the built-in scenarios, one per line: name and question (ADR 0054) | |
| `--scenario NAME` | Run a built-in scenario, after the shared formulas and before `-e` and files; `formulas` loads only the formulas (ADR 0054) | |
| `--scenario-source NAME` | Print a built-in scenario's source, to copy and edit it (ADR 0054) | |
| `--completion SHELL` | Print a completion script for `bash`, `zsh` or `fish`; it also completes scenario names after `--scenario` and `--scenario-source` | |
| `--plain` | Print only the central value (median) | off |
| `--no-color` | Disable ANSI colour | auto (TTY only) |
| `--ascii` | ASCII instead of Unicode (sparklines, `≈`, `…`) | off (ADR 0022) |
| `-h, --help`, `--version` | | |

Flags go before files (`ubc -e 'x = 2' model.ubc`); anything after the first file name is read as another file. `-e` expressions are evaluated before files, all in one session, so variables flow from one source to the next. Without `-e` or files, the program is read from stdin.

`--set` overrides are evaluated first, in order, before `-e`, files, stdin or the REPL; each value is split at its first `=` and `EXPR` is one expression that may use earlier overrides. Any later assignment to an overridden name is skipped without evaluating its right-hand side, so everything computed from the name uses the override; with `--echo` the skipped assignment prints the override. The same name given twice: the last one wins. `--set` prints nothing itself. A value without `=`, or whose name is not an identifier, is a usage error (exit 2); a bad expression is a parse (3) or evaluation (1) error reported as `ubc: --set NAME: line 1, col C: message`.

`--data` files are read first, in order, before the `--set` overrides (so `--set 'worst=max(lead)'` can use the data) and before any program. Each file holds one plain number per line (`12`, `-3.5`, `1e3`; no suffixes); spaces around a number, blank lines and Windows line ends are ignored. The variable becomes an empirical SampleSet that resamples those numbers with replacement, exactly as the list literal `[v1, v2, ...]` would, and it is an override with the rules of `--set`: later assignments to it are skipped, the last `--data` for a name wins, a `--set` for the same name wins over `--data`, and `:samples` forgets it. With `--watch` the data files are read again on every run, but only changes to the program files start one. A value without `=`, a name that is not an identifier or an empty path is a usage error (exit 2). A file that cannot be read, a line that is not a finite number (`ubc: --data lead: times.csv: line 3: 'days' is not a number`) or a file with no numbers (`ubc: --data lead: times.csv: no values to resample`) is an evaluation error (exit 1) and stops the run before any program.

`--assert` expressions are evaluated after the whole program (`--set`, `-e`, files or stdin), in order, in the same session; they print nothing and change no variable, not even `last`. Each must be one expression whose value is a condition: a comparison, a probability (`P(...)`) or a plain `0`/`1`. It holds only when that value is exactly 1, so a comparison of uncertain values must hold in every sample (`P(x > 0) > 0.95` states a looser check). Every assertion that does not hold prints `ubc: --assert EXPR failed: ACTUAL` to stderr, where `ACTUAL` is, for a comparison, both sides and the operator (`ubc: --assert fixed > 200 failed: 150 > 200`), uncertain sides as their one-line summary without sparkline and followed by the probability that the comparison held (`≈ 142   90% [100 … 198]   98% [86 … 240] < 150 (P = 0.61)`), and for anything else its value (`0.25`). The run goes on to the next assertion and exits with `4` if any failed. An assertion that is not a condition (`ubc: --assert cost: line 1, col 1: an assertion must be a condition or a probability`), or that fails to parse or evaluate, stops the run with exit 1 or 3. If the program itself fails, the assertions are not evaluated. `--assert` without a program (the REPL) is a usage error; with `--watch`, failures are reported after every run.

`--samples-out FILE` writes, after the whole program has run without error and before the assertions are checked, the raw samples of every variable in the order they were first assigned (`--data` and `--set` overrides included), followed by `last` when the program printed an expression and did not assign `last` itself. Row *i* holds sample *i* of every variable, so correlations survive (`x` and `2 * x` stay paired); plain numbers and probabilities repeat their value on every row, so every column has `--samples` rows. The format comes from the extension, which is case-sensitive: `.csv` writes a header row with the names and one row per sample; `.json` writes one object whose keys are the names, in the same order, and whose values are arrays of samples. Numbers are written in their shortest exact form (`2.5`, `1e+21`); non-finite samples are an empty cell in CSV and `null` in JSON. Nothing goes to stdout. Any other extension is a usage error (exit 2: `ubc: --samples-out writes .csv or .json files, got 'out.txt'`), and so is `--samples-out` without a program (the REPL). A program that fails writes no file. A program that leaves no variable and no `last` (`ubc: --samples-out out.csv: nothing to export`) or a file that cannot be written (`ubc: --samples-out out.csv: open out.csv: permission denied`) is an evaluation error (exit 1); when it happens the assertions are not checked. With `--watch` the file is written again after every run that succeeds.

`--sens EXPR` prints, after the whole program has run without error, after `--samples-out` and before the assertions are checked, which variables drive the uncertainty of EXPR (ADR 0033). EXPR is one expression, evaluated like an assertion: it sees the session and binds nothing, not even `last`. Each variable whose value is uncertain is compared with EXPR by Spearman's rank correlation of their samples, pairing sample *i* with sample *i* and skipping pairs where either is not finite; plain numbers, probabilities, variables whose samples do not vary and the variable named exactly EXPR are left out. Variables are listed by the absolute value of the coefficient, strongest first, ties in assignment order, one per line under a header, the names padded to the same width, the coefficient with its sign and two decimals and a bar of `█` (`#` with `--ascii`) whose length is the coefficient's absolute value times 20, rounded:

```
$ ubc --sens profit model.ubc
sensitivity of profit (rank correlation):
  units  +0.98  ████████████████████
  price  +0.16  ███
  cost   -0.01
```

A positive coefficient means EXPR tends to grow with the variable, a negative one that it shrinks. The coefficient measures how strongly they move together, any monotonic relation counting, not how much variance each one explains, and intermediate variables built from EXPR's inputs appear too. EXPR whose value is a plain number (`ubc: --sens 2: needs an uncertain value, got a plain number`) or no other uncertain variable to compare with (`ubc: --sens x: needs other uncertain variables to compare with`) is an evaluation error (exit 1); a broken EXPR is a parse (3) or evaluation (1) error reported as `ubc: --sens EXPR: line 1, col C: message`; then the assertions are not checked. A program that fails prints no ranking. `--sens` with `--json` or `--plain`, or without a program (the REPL, which has `:sens`), is a usage error. With `--watch` the ranking is printed after every run that succeeds.

`--scenarios`, `--scenario NAME` and `--scenario-source NAME` use the scenario catalogue built into the binary (§8, ADR 0054), so it works wherever `ubc` is installed. A scenario is named by its path under `scenarios/` without `.ubc` (`team/replace-a-person`); `formulas` names the shared formulas. `--scenarios` prints one line per scenario, sorted by name: the name padded to the longest one, three spaces and the question, which is the scenario's `# Question:` comment and the comment lines after it up to the first empty `#`, joined on one line and capitalised (`team/replace-a-person   How much does it cost to replace an engineer who leaves the team, and how likely is it to cost more than 50k euros?`); `formulas` is not listed. `--scenario NAME` makes the shared formulas and then the scenario the first sources of the program, before `-e` and files, in the same session, so every other flag (`--set`, `--data`, `--echo`, `--assert`, `--samples-out`, `--sens`, `--watch` on the files) applies as it does to a file; stdin is not read and the REPL does not start. `--scenario formulas` loads only the formulas, for your own model: `ubc --scenario formulas my-model.ubc`. Errors inside a scenario carry its name like a file's (`ubc: team/replace-a-person: line L, col C: message`). `--scenario-source NAME` prints the file unchanged and exits, like `--version`: `ubc --scenario-source team/replace-a-person > mine.ubc`. An unknown name, with either flag, is a usage error (exit 2): `ubc: unknown scenario 'NAME'; ubc --scenarios lists them`.

`--interval W` sets the width of the interval in every summary: results, assignments in the REPL and with `--echo`, `:vars`, `:plot` titles and failed assertions. W is a percentage strictly between 0 and 100, decimals allowed; the interval is central, leaving (100 − W) / 2 % of the samples out on each side, so `--interval 80` prints `≈ 170   80% [109 … 262]   98% [75 … 380]`, the 10th and 90th percentiles. The nested 98% interval (§6) stays as it is, and is left out when W is 98 or more, because it would no longer be the wider one. `--interval 90` is the default output. With `--json` the standard fields (`p5`, `p95`) stay and the interval's ends are added to the `quantiles` object (`"quantiles":{"p10":…,"p90":…}`); the text form of `:hist` and the sparkline are unchanged (they always span P1–P99). A value that is not such a percentage (`ubc: --interval expects a percentage between 0 and 100, got '100'`) and `--interval` together with `--quantiles`, which replaces the interval, are usage errors (exit 2).

Exit codes: `0` success, `1` evaluation error, `2` usage error (bad flags), `3` parse error, `4` assertion failed (ADR 0029).

Errors go to stderr as `ubc: line L, col C: message`, with the file name for errors inside a file: `ubc: model.ubc: line L, col C: message`. A file or stdin that cannot be read is an evaluation error (exit 1); a file error ends with the operating system's own reason, which differs between systems (`open model.ubc: no such file or directory` on Linux and macOS, `open model.ubc: The system cannot find the file specified.` on Windows; ADR 0041). In the REPL an error prints and the session continues.

## 2. Lexical structure

- **Numbers:** `42`, `3.14`, `.5`, `1e6`, `2.5e-3`.
- **Suffixes** (no space between number and suffix): `k` = 1e3, `M` = 1e6, `B` = 1e9, `T` = 1e12, `%` = 1e-2.
- **Identifiers:** `[A-Za-z_][A-Za-z0-9_]*`, case-sensitive. Unicode identifiers are not supported.
- **Keywords:** `to`.
- **Operators:** `+ - * / ^ ( ) , = < > <= >= == !=`.
- **Comments:** `#` to end of line.
- **Statements** are separated by newlines or `;`.

## 3. Grammar

```
program     = { statement ( NEWLINE | ";" ) } ;
statement   = assignment | expression ;
assignment  = IDENT "=" expression ;
expression  = comparison ;
comparison  = range [ ( "<" | ">" | "<=" | ">=" | "==" | "!=" ) range ] ;
range       = additive [ "to" additive ] ;
additive    = multiplicative { ( "+" | "-" ) multiplicative } ;
multiplicative = unary { ( "*" | "/" ) unary } ;
unary       = [ "-" ] power ;
power       = primary [ "^" unary ] ;              (* right-associative *)
primary     = NUMBER | IDENT | call | "(" expression ")" | list ;
call        = IDENT "(" [ expression { "," expression } ] ")" ;
list        = "[" [ expression { "," expression } ] "]" ;
```

Precedence, lowest to highest: comparison, `to`, `+ -`, `* /`, unary `-`, `^`. Note that `to` binds looser than arithmetic: `1 + 2 to 10` is `(1 + 2) to 10`. When in doubt, users write parentheses: `(8 to 25) * 12`.

## 4. Values and semantics

Every value is one of:

- **Scalar:** a float64.
- **SampleSet:** N float64 samples (N = `--samples`), produced by any distribution and propagated by any operation that touches one.
- **Probability result:** a scalar in [0, 1] produced by `P(...)`.

Rules:

1. Scalar ⊕ scalar → scalar. Scalar ⊕ SampleSet and SampleSet ⊕ SampleSet → SampleSet, applied element-wise (sample *i* with sample *i*).
2. A variable bound to a SampleSet keeps its samples for the whole session. `x - x` is exactly 0; `x * y` with two independent variables is the product of independent draws.
3. Every distribution literal (`to`, `normal(...)`, …) draws fresh samples from the session RNG at evaluation time.
4. The RNG is seeded once per session (default seed 1), so the same program gives the same output.
5. Comparisons between SampleSets yield a SampleSet of 0/1; `P(expr)` returns its mean. A bare comparison at top level prints as a probability. `P` of anything that is not 0/1 in every sample is an evaluation error. A probability behaves as a scalar in arithmetic (`P(x > 50) * 100`).
6. Division by zero yields `±Inf`/`NaN` samples; summaries ignore non-finite samples and the output reports how many were dropped: `≈ 0.702   90% [0.517 … 1.08]   98% [0.5 … 1.3]   █▇▇▆▆▆▅▄▃▃▂▂▂▁▁▁   (2486 non-finite dropped)`, or `no finite samples (5000 dropped)`.

### `a to b`

- 90 % confidence interval: P(X < a) = 5 %, P(X > b) = 5 %.
- If `a > 0` and `b > 0`: lognormal with `μ = (ln a + ln b) / 2`, `σ = (ln b − ln a) / (2 · 1.6449)`.
- Otherwise: normal with `μ = (a + b) / 2`, `σ = (b − a) / (2 · 1.6449)`.
- `a >= b` is an evaluation error.
- Both bounds must be plain numbers: `(1 to 2) to 3` is an evaluation error (uncertain bounds are not supported).

## 5. Built-in functions

| Function | Returns |
|---|---|
| `normal(mean, sd)`, `lognormal(mu, sigma)`, `uniform(lo, hi)`, `triangular(lo, mode, hi)`, `beta(a, b)` | SampleSet |
| `exponential(rate)`, `gamma(shape, scale)`, `pert(lo, mode, hi)` (phase 4) | SampleSet; `pert` is Beta-PERT with mean `(lo + 4·mode + hi) / 6` |
| `poisson(lambda)`, `binomial(n, p)` | SampleSet of whole numbers, drawn exactly for any size (ADR 0035); `lambda` positive and finite, `n` a positive whole number, `p` in [0, 1] |
| `mean(x)`, `median(x)`, `stdev(x)`, `min(x)`, `max(x)` | Scalar (identity on scalars); `stdev` is the population standard deviation |
| `p(x, q)` | q-quantile, `q` in [0, 1] |
| `P(cond)` | Probability |
| `sqrt, exp, log, log10, abs, round, floor, ceil` | Element-wise |
| `min(a, b, ...)`, `max(a, b, ...)` with 2+ args | Element-wise |
| `[v1, v2, ...]` | Empirical SampleSet, resampled with replacement; elements must be plain numbers and the list cannot be empty |

Built-in variables: `last` (result of the previous top-level expression), `pi`, `e`.

### User functions (phase 4)

```
cost(n) = n * (2 to 4)
cost(10)
```

- `name(p1, p2, ...) = expr` defines a function; the definition prints nothing. Redefining a function replaces it.
- A call evaluates the body with the parameters bound to the arguments and the global variables as they are at call time. Parameters shadow globals without changing them.
- Ranges and distributions in the body draw fresh samples on every call, so `cost(1) - cost(1)` is not 0.
- Redefining a built-in function, repeating a parameter or calling with the wrong number of arguments is an evaluation error.
- There are no conditionals, so recursion cannot end: calls nested more than 64 deep are an evaluation error.
- Functions survive `:samples`, which only clears variables.

## 6. Output

One line per top-level expression. In the REPL, assignments print too, as `name = summary`; files, stdin and `-e` print only expressions unless `--echo` is given (ADR 0025).

When one program (a file, stdin, one `-e` or a `:load`ed file) prints more than one expression, each says what it is (ADR 0053): an expression that is a variable name prints as `name = summary`, any other expression whose source text, as written up to any comment, is at most **30 characters** prints as `text = summary` (`P(total > 50k) = 0.75`), and a longer one keeps its line number, `L12: summary`. A program that prints a single expression prints it bare, so `$(ubc -e ...)` captures just the value; the REPL prints expressions bare too.

| Value | Default output |
|---|---|
| Scalar | Shortest decimal representation, max 10 significant figures: `4`, `0.1`, `3.141592654` |
| SampleSet | `≈ MEDIAN   90% [P5 … P95]   98% [P1 … P99]   SPARKLINE` (the inner width follows `--interval`, §1) |
| Probability | `0.16` (two decimals; `<0.01` / `>0.99` at the extremes) |

- A SampleSet whose samples are all equal prints as a scalar: `x = 1 to 5; x - x` prints `0`.
- Every summary nests a 98% interval, `[P1 … P99]`, around the inner one, so the tails show without a flag and one result stays one line (ADR 0040). It appears wherever the summary does: results, assignments in the REPL and with `--echo`, `:vars`, `:plot` titles and failed assertions. `--quantiles` replaces both intervals.
- SampleSet numbers use **3 significant figures** without trailing zeros (`≈ 14.2   90% [10 … 19.8]`), and suffixes from 10k up, chosen number by number after rounding: `≈ 22.4k   90% [9990 … 49k]`.
- When the interval is narrow compared with its values, a summary number keeps more digits: the median, the inner interval ends and `--quantiles` use the finer of 3 significant figures of their own value and **2 significant figures of the inner interval's width** (P95 − P5, or the `--interval` width), so the ends never collapse into one number (ADR 0052): `1M to 1.002M` prints `≈ 1.001M   90% [1M … 1.002M]`, `1000 to 1010` prints `≈ 1005   90% [1000 … 1009.9]`. Ordinary results are unchanged.
- The ends of the 98% interval are rounded **outwards**, P1 down and P99 up, to **2 significant figures of their distance** (P99 − P1), because P1 and P99 rest on few samples and move with the spread between seeds, but never coarser than 2 significant figures of their own value, and never inside the printed inner interval. So the 98% interval always contains the true P1 and P99 and the printed inner interval (ADR 0040, ADR 0051): `(8 to 25) * 12` prints `98% [75 … 380]`, `1000 to 1010` prints `98% [997 … 1013]` and `(1 to 10) * 1000` prints `90% [999 … 9710]   98% [610 … 16k]`. When P1 equals P99 they use 3 significant figures like the rest.
- Rounding is exact in decimal, so no printed number shows binary floating-point noise: `11.5M to 11.7M` prints `≈ 11.6M   90% [11.5M … 11.7M]   98% [11.45M … 11.75M]`, never `11.599999999999998M` (ADR 0039).
- The sparkline is 16 cells wide, built from a histogram over [P1, P99]; Unicode block elements `▁▂▃▄▅▆▇█`, ASCII fallback ` .:-=+*#`.
- `--plain` prints the median (scalar as is) with full float precision and no labels or `L<n>:` prefixes, for `$(ubc -e ...)`.
- `--json` prints, per result: `{"kind":"samples","median":…,"mean":…,"stdev":…,"p1":…,"p5":…,"p95":…,"p99":…,"n":5000,"dropped":0}` or `{"kind":"scalar","value":…}` or `{"kind":"probability","value":…}`. One object per line, no `L<n>:` prefixes; a labelled expression carries its label as `"name"`, like an assignment with `--echo` (`{"name":"P(total > 50k)","kind":"probability","value":0.75}`), with `<`, `>` and `&` written as they are; statistics ignore non-finite samples, and values that are not finite (`1 / 0`, or no finite samples at all) are `null`. `--json` and `--plain` cannot be combined.

## 7. REPL

- Prompt `ubc> `, line editing and persistent history in `$XDG_STATE_HOME/ubc/history` (default `~/.local/state/ubc/history`), using `golang.org/x/term` (ADR 0023); a testscript scenario drives it through a pseudo-terminal (ADR 0036). On Windows the REPL turns on escape-sequence processing for the console while it runs and restores it on exit; Ctrl-Z also ends the session there (ADR 0037). The REPL starts when stdin is a terminal and there are no `-e` expressions or files.
- Every line prints its result, assignments included. `:` commands work only in the REPL.
- State (variables, `last`, RNG) persists for the session.
- Commands start with `:` and are not part of the language:

| Command | Effect |
|---|---|
| `:help` | List commands and functions |
| `:vars` | Show variables with their one-line summary |
| `:hist EXPR` | Full terminal histogram (≈20 rows, P1–P99, axis labels) |
| `:sens EXPR` | Rank the variables that drive the uncertainty of EXPR, as `--sens` does (§1); it changes nothing, not even `last` |
| `:plot EXPR FILE` | Export a 40-bin histogram over P1–P99 as SVG (`.svg`) or PNG (`.png`, ASCII text, ADR 0024), titled with EXPR and its summary |
| `:samples N` | Change sample count (clears variables, with a warning) |
| `:seed N` | Reseed the RNG |
| `:save FILE` | Write the session's successful statements to a `.ubc` file |
| `:load FILE` | Evaluate a file into the session |
| `:quit` / Ctrl-D | Exit |

## 8. Files

- Extension `.ubc`, UTF-8 text, same syntax as the REPL (no `:` commands).
- `ubc model.ubc` prints the result of every top-level expression, labelled as §6 describes when there is more than one (`saving = ≈ 101   90% [21.5 … 229]   98% [-1.5 … 310]   ▃▅▆▇█▇▇▆▄▃▃▂▂▂▁▁`).
- `examples/` holds the worked models of [TUTORIAL.md](TUTORIAL.md); the test suite runs every command shown in the tutorial and compares its output (ADR 0030).
- `scenarios/<area>/<name>.ubc` is a catalogue of ready-to-use models that follow one template, indexed in `scenarios/README.md`; the test suite fails when a scenario does not evaluate, is missing from the index or is not embedded in the binary (ADR 0042, ADR 0054). The index opens with a note that `ubc` has no units (amounts are in euros, and in dollars in the cloud area; change the numbers to your currency), then a "Start here" section with five featured scenarios and one section per area with short questions; every scenario's header has a `# Sources:` line naming the published references its ranges rely on or saying they come from typical experience, and every assumption's comment says where its range comes from; the test suite checks that shape (ADR 0055). The binary embeds `formulas.ubc` and `<area>/<name>.ubc`, which `--scenarios`, `--scenario` and `--scenario-source` use (§1), and each scenario's "Run it" lines use `--scenario`. The `team` area covers replacing a person, hiring, attrition, onboarding, real team capacity, the cost of a recurring meeting, in-house vs contractors vs outsourcing and the bus factor; each team scenario has a testscript that runs the commands in its header and pins its medians and its answer (ADR 0044). The `delivery` area covers the delivery date from historical throughput (a default list that `--data` replaces), Little's law and a WIP limit, scope creep and context switching, with testscripts like the team ones; its index also links `examples/project-estimate.ubc`, and the test suite fails when a model the index links to does not exist (ADR 0045). The `product` area centres on cost of delay: the cost of delay of a feature, CD3 prioritisation, urgency profiles, delaying a launch vs a reduced scope and the cost of queues, plus a feature's business case, RICE with ranges, TAM/SAM/SOM, LTV and CAC payback, a price change, an A/B experiment and buy vs build, with testscripts like the team ones (ADR 0046). The `engineering` area covers technical debt and productivity: the cost of not fixing a piece of debt against fixing it, whether automating a task is worth the time, a slow CI pipeline, flaky tests, and a migration's cost against its budget, with testscripts like the team ones (ADR 0047). The `reliability` area covers operations: the yearly cost of incidents, the availability of services in series and in parallel, an SLO's error budget, capacity planning for peak traffic, and the cost of a bug by the stage that catches it, with testscripts like the team ones (ADR 0048). The `cloud` area covers cloud and other variable costs, in dollars: an LLM feature billed by the token, observability billed by the gigabyte of logs, storage with compound growth, a reserved capacity commitment against on-demand under an uncertain usage, the infrastructure cost per customer, per-seat SaaS as the team grows and ephemeral environments per pull request, with testscripts like the team ones; its index also links `examples/cloud-cost.ubc` as the serverless bill (ADR 0049). The `risk` area covers risk and security with the shared `annual_loss` formula: the annual loss expectancy of a risk and whether a mitigation pays off, and the cost of a data breach (records × cost per record + the regulator's fine) with the loss to expect a year, with testscripts like the team ones (ADR 0050).
- `scenarios/formulas.ubc` defines reusable formulas as user functions (loaded cost, real capacity, Little's law, cost of delay, CD3, LTV, CAC payback, annual loss expectancy, availability in series and in parallel, compound growth, present value, NPV) and prints nothing; `--scenario` loads it before every scenario, and `ubc --scenario formulas model.ubc` loads it before your own model (from a clone, `ubc scenarios/formulas.ubc model.ubc` does the same). The test suite runs every scenario after it and checks each formula with plain numbers (ADR 0043).

## 9. Acceptance examples

These hold with the default seed and sample count.

| Input | Expected |
|---|---|
| `2 + 2` | `4` |
| `1k * 3` | `3000` |
| `5% * 200` | `10` |
| `(8 to 25) * 12` | median ≈ 170, 90 % interval ≈ [96 … 300], 98 % interval ≈ [80 … 380] |
| `x = 1 to 5` then `x - x` | `0` |
| `P(normal(0, 1) > 0)` | ≈ `0.50` |
| `10 to 5` | evaluation error, exit 1 |
| `-10 to 10` | normal, median ≈ 0 |
| `echo '1 + 2' \| ubc --plain` | `3` |

## 10. Out of scope

Not supported today, by decision: units (ADR 0008), symbolic distributions, correlation between distinct ranges, a TUI with sliders, Squiggle compatibility, Unicode identifiers.
