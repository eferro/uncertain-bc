# Uncertain bc — MVP Specification

Scope: phases 0–3 of the roadmap in [VISION.md](VISION.md). Everything here is a proposal until it is implemented and covered by tests; deviations must update this file or add an ADR in [adr/](adr/README.md).

Status: draft. Date: 2026-11-21.

## 1. Invocation

```
ubc                         # interactive REPL when stdin is a TTY
ubc FILE...                 # evaluate files in order, print the result of each top-level expression
ubc -e EXPR [-e EXPR...]    # evaluate expressions, print results
cmd | ubc                   # read program from stdin (non-TTY), print results
```

| Flag | Meaning | Default |
|---|---|---|
| `-e, --eval EXPR` | Evaluate an expression (repeatable) | |
| `-n, --samples N` | Monte Carlo sample count | `5000` |
| `-s, --seed N` | RNG seed | `1` |
| `--random` | Use a time-based seed | off |
| `--json` | One JSON object per printed result | off |
| `--plain` | Print only the central value (median) | off |
| `--no-color` | Disable ANSI colour | auto (TTY only) |
| `--ascii` | ASCII instead of Unicode (sparklines, `≈`, `…`) | auto |
| `-h, --help`, `--version` | | |

Flags go before files (`ubc -e 'x = 2' model.ubc`); anything after the first file name is read as another file. `-e` expressions are evaluated before files, all in one session, so variables flow from one source to the next. Without `-e` or files, the program is read from stdin.

Exit codes: `0` success, `1` evaluation error, `2` usage error (bad flags), `3` parse error.

Errors go to stderr as `ubc: line L, col C: message`, with the file name for errors inside a file: `ubc: model.ubc: line L, col C: message`. A file or stdin that cannot be read is an evaluation error (exit 1). In the REPL an error prints and the session continues.

## 2. Lexical structure

- **Numbers:** `42`, `3.14`, `.5`, `1e6`, `2.5e-3`.
- **Suffixes** (no space between number and suffix): `k` = 1e3, `M` = 1e6, `B` = 1e9, `T` = 1e12, `%` = 1e-2.
- **Identifiers:** `[A-Za-z_][A-Za-z0-9_]*`, case-sensitive. Unicode identifiers are a non-goal for the MVP.
- **Keywords:** `to`.
- **Operators:** `+ - * / ^ ( ) , = < > <= >= == !=`.
- **Comments:** `#` to end of line.
- **Statements** are separated by newlines or `;`.

## 3. Grammar (MVP)

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
6. Division by zero yields `±Inf`/`NaN` samples; summaries ignore non-finite samples and the output reports how many were dropped.

### `a to b`

- 90 % confidence interval: P(X < a) = 5 %, P(X > b) = 5 %.
- If `a > 0` and `b > 0`: lognormal with `μ = (ln a + ln b) / 2`, `σ = (ln b − ln a) / (2 · 1.6449)`.
- Otherwise: normal with `μ = (a + b) / 2`, `σ = (b − a) / (2 · 1.6449)`.
- `a >= b` is an evaluation error.
- Both bounds must be plain numbers: `(1 to 2) to 3` is an evaluation error (uncertain bounds are out of the MVP).

## 5. Built-in functions (MVP)

| Function | Returns |
|---|---|
| `normal(mean, sd)`, `lognormal(mu, sigma)`, `uniform(lo, hi)`, `triangular(lo, mode, hi)`, `beta(a, b)` | SampleSet |
| `mean(x)`, `median(x)`, `stdev(x)`, `min(x)`, `max(x)` | Scalar (identity on scalars); `stdev` is the population standard deviation |
| `p(x, q)` | q-quantile, `q` in [0, 1] |
| `P(cond)` | Probability |
| `sqrt, exp, log, log10, abs, round, floor, ceil` | Element-wise |
| `min(a, b, ...)`, `max(a, b, ...)` with 2+ args | Element-wise |
| `[v1, v2, ...]` | Empirical SampleSet, resampled with replacement |

Built-in variables: `last` (result of the previous top-level expression), `pi`, `e`.

## 6. Output

One line per top-level expression (assignments print too, as in the REPL examples; files and `-e` print only expressions unless `--echo` is added later).

| Value | Default output |
|---|---|
| Scalar | Shortest decimal representation, max 10 significant figures: `4`, `0.1`, `3.141592654` |
| SampleSet | `≈ MEDIAN   90% [P5 … P95]   SPARKLINE` |
| Probability | `0.16` (two decimals; `<0.01` / `>0.99` at the extremes) |

- A SampleSet whose samples are all equal prints as a scalar: `x = 1 to 5; x - x` prints `0`.
- SampleSet numbers use **3 significant figures** without trailing zeros (`≈ 14.2   90% [10 … 19.8]`), and suffixes from 10k up, chosen number by number after rounding: `≈ 22.4k   90% [9990 … 49k]`.
- The sparkline is 16 cells wide, built from a histogram over [P1, P99]; Unicode block elements `▁▂▃▄▅▆▇█`, ASCII fallback ` .:-=+*#`.
- `--plain` prints the median (scalar as is) with full float precision, for `$(ubc -e ...)`.
- `--json` prints, per result: `{"kind":"samples","median":…,"mean":…,"stdev":…,"p5":…,"p95":…,"n":5000,"dropped":0}` or `{"kind":"scalar","value":…}` or `{"kind":"probability","value":…}`.

## 7. REPL

- Prompt `ubc> `, line editing and persistent history in `$XDG_STATE_HOME/ubc/history`.
- State (variables, `last`, RNG) persists for the session.
- Commands start with `:` and are not part of the language:

| Command | Effect |
|---|---|
| `:help` | List commands and functions |
| `:vars` | Show variables with their one-line summary |
| `:hist EXPR` | Full terminal histogram (≈20 rows, P1–P99, axis labels) |
| `:plot EXPR FILE` | Export histogram as SVG (`.svg`) or PNG (`.png`) |
| `:samples N` | Change sample count (clears variables, with a warning) |
| `:seed N` | Reseed the RNG |
| `:save FILE` | Write the session's successful statements to a `.ubc` file |
| `:load FILE` | Evaluate a file into the session |
| `:quit` / Ctrl-D | Exit |

## 8. Files

- Extension `.ubc`, UTF-8 text, same syntax as the REPL (no `:` commands).
- `ubc model.ubc` prints the result of every top-level expression, prefixed with the source line when there is more than one (`L12: ≈ 97   90% [22 … 237]   ▁▂▅█▆▃▁`).

## 9. Acceptance examples

These must hold with the default seed and sample count (exact digits are fixed once the RNG is implemented; until then, check the shape).

| Input | Expected |
|---|---|
| `2 + 2` | `4` |
| `1k * 3` | `3000` |
| `5% * 200` | `10` |
| `(8 to 25) * 12` | median ≈ 170, 90 % interval ≈ [96 … 300] |
| `x = 1 to 5` then `x - x` | `0` |
| `P(normal(0, 1) > 0)` | ≈ `0.50` |
| `10 to 5` | evaluation error, exit 1 |
| `-10 to 10` | normal, median ≈ 0 |
| `echo '1 + 2' \| ubc --plain` | `3` |

## 10. Out of scope for the MVP

User-defined functions, units, symbolic distributions, correlation between distinct ranges, a TUI with sliders, Squiggle compatibility, Unicode identifiers.
