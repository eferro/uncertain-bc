# Uncertain bc — Vision

> `ubc` is `bc` for numbers you are not sure about. You type ranges instead of numbers, it runs a Monte Carlo simulation behind the scenes and answers with a distribution: a central value, an interval and a small chart, right in the terminal.

Status: phase 4 done (expressiveness, v0.5.0). Date: 2026-11-21. Detailed MVP behaviour: [SPEC.md](SPEC.md). Decisions: [adr/](adr/README.md).

## 1. The problem

Most of the numbers we use to make decisions are guesses: how long a migration will take, how much a server costs per year, how many users will try a feature. Spreadsheets and calculators force us to pick a single value for each guess, and the final answer inherits a precision it never had.

Tools that model uncertainty properly exist, but none of them fits the "quick calculation in the terminal" habit:

- **Guesstimate** is a spreadsheet for uncertainty (ranges, formulas between cells, Monte Carlo with 5,000 samples, histograms). It is great, but it is a web app with a canvas, accounts and a backend.
- **Squiggle** is a small language for estimation (`1k to 5k * 12`) with a rich standard library. It runs in the browser, in VS Code and from Node, but its CLI only prints numbers and it needs a JavaScript runtime.
- **bc** is the calculator many of us reach for in the terminal: instant, scriptable, pipe-friendly. It only knows exact numbers.

`ubc` sits at the intersection: the ergonomics of `bc`, the expressiveness of Squiggle-style ranges, and a single Go binary you can share with anyone.

## 2. Goals

1. **Feel like bc.** Start instantly, read expressions from the REPL, from `-e`, from a file or from stdin, and print a short answer. Usable inside shell scripts and pipes.
2. **Ranges are first-class numbers.** `10 to 20` means "90 % confident it is between 10 and 20". Every operator and function works on ranges exactly as on plain numbers.
3. **Honest output.** Show a central value, an interval and the shape of the distribution with few significant figures. Never print `96.98057045115692` for a guess.
4. **Charts while you calculate.** A sparkline next to every result, a full histogram on demand, and export to SVG/PNG for sharing.
5. **One static binary.** No runtime, no network, no account. `go install` or download and run on Linux, macOS and Windows.
6. **Reproducible by default.** Same input, same output: the random seed is fixed unless you ask otherwise.

## 3. Non-goals

- **Not a Squiggle implementation.** The language is inspired by Squiggle and stays close to it where that makes sense, but compatibility is not a goal and `ubc` will diverge freely (bc-isms, REPL commands, simpler semantics).
- **Not a spreadsheet or canvas.** No grid of cells, no GUI. Models are text.
- **Not a statistics package.** No fitting, no inference, no symbolic distributions (at least at first). Everything is sample-based Monte Carlo.
- **Not a server or a collaboration tool.** Sharing means sharing a text file or an exported chart.

## 4. What using it looks like

These examples are a design sketch, not a final syntax.

One-shot, like `bc -e`:

```console
$ ubc -e '(8 to 25) * 12'
≈ 170   90% [96 … 300]   ▁▃▇█▅▃▂▁
```

Pipes and scripts:

```console
$ echo '(0.12 to 0.25) * 24 * 365 * (12 to 30) / 1000' | ubc
≈ 37    90% [19 … 70]     ▂▆█▆▃▂▁▁

$ ubc --json -e '1k to 5k' | jq .median
2236
```

Interactive REPL with variables, kept between lines (no re-sampling, so `a` keeps its value):

```console
$ ubc
ubc> watts = 12 to 30
≈ 19    90% [12 … 30]     ▂▇█▆▃▂▁▁
ubc> price = 0.12 to 0.25
ubc> yearly = watts / 1000 * 24 * 365 * price
≈ 33    90% [17 … 62]     ▂▆█▆▃▂▁▁
ubc> P(yearly > 50)
0.16
ubc> :hist yearly
   15 ┤▇▇▇▇▇
   25 ┤▇▇▇▇▇▇▇▇▇▇▇▇▇▇
   35 ┤▇▇▇▇▇▇▇▇▇▇
   45 ┤▇▇▇▇▇
   55 ┤▇▇▇
   65 ┤▇
ubc> :plot yearly yearly.svg
```

A model file:

```console
$ cat homelab.ubc
# Does my homelab beat a VPS?
watts    = 12 to 30
price    = 0.12 to 0.25          # €/kWh
hardware = 120 to 200
life     = 3 to 6                # years
homelab  = watts / 1000 * 24 * 365 * price + hardware / life
vps      = (8 to 25) * 12
saving   = vps - homelab
P(saving < 0)

$ ubc homelab.ubc
```

## 5. Language sketch

Inspired by Squiggle, trimmed down for a calculator:

| Feature | Example | Notes |
|---|---|---|
| Plain numbers | `42`, `1.5e3` | |
| Suffixes | `1k`, `2.5M`, `3B`, `5%` | `5%` = `0.05` |
| 90 % range | `10 to 20` | Lognormal when both ends are positive, normal otherwise (same rule as Squiggle) |
| Explicit distributions | `normal(100, 15)`, `lognormal(...)`, `uniform(1, 6)`, `beta(2, 8)`, `triangular(1, 3, 10)` | |
| Empirical data | `[3, 5, 8, 13]` | Resampled |
| Arithmetic | `+ - * / ^`, parentheses | Sample by sample |
| Variables | `x = 1 to 5` | REPL state persists |
| Last result | `last` | Like bc's `last` |
| Summaries | `mean(x)`, `median(x)`, `stdev(x)`, `p(x, 0.9)` | |
| Probabilities | `P(x > 50)`, `P(a < b)` | |
| Functions | `min`, `max`, `sqrt`, `log`, `exp`, `round`, … | |
| User functions | `cost(n) = n * (2 to 4)` | Later phase |
| Comments | `# …` | |

REPL commands start with `:` so they never clash with the language: `:hist x`, `:plot x file.svg`, `:samples 10000`, `:seed 42`, `:vars`, `:save model.ubc`, `:help`.

## 6. Output design

- **Default line:** `≈ median   90% [p5 … p95]   sparkline`, three significant figures.
- **Plain numbers stay plain:** `2 + 2` prints `4`, as in bc.
- **Probabilities:** a single number with two decimals.
- **Machine output:** `--json` (summary stats plus optional samples) and `--plain` (just the median, for `$(ubc -e ...)`).
- **Unicode** always, `--ascii` on request (ADR 0022); **colour** only when stdout is a TTY.

## 7. Architecture

Everything is Go, no embedded JavaScript (ADR 0001, ADR 0002). Hexagonal layout (ADR 0010):

```
cmd/ubc/                composition root: flags, wiring, exit codes
internal/domain/        pure: lexer, parser, AST, values, distributions, statistics, evaluator
internal/application/   Session use case (variables, `last`, RNG) and the ports it needs
internal/adapters/      CLI/REPL, text and JSON output, later history, model files and chart export
```

Key decisions:

- **Sample sets everywhere.** A range becomes N samples (default 5,000) at definition time; operators are applied element-wise. Correlation between uses of the same variable is preserved for free, because `x - x` reuses the same samples and gives exactly 0.
- **Seeded RNG per session.** Deterministic by default; `--seed` and `:seed` to change it, `--random` to opt out.
- **Scalars stay scalars.** Only values that touch a range become sample sets, so `ubc` is also a normal calculator.
- **Libraries:** standard library first. Candidates: a readline library for the REPL history, gonum (`gonum.org/v1/gonum/stat`, `gonum.org/v1/plot`) for stats and image export, and a terminal chart library if hand-rolled sparklines are not enough.

## 8. Risks and open questions

1. **Scope creep in the language.** Writing an interpreter is the fun part and the easy place to lose months. Mitigation: a fixed MVP grammar (section 5 minus user functions) and a test corpus of real estimates.
2. **Semantics that look like Squiggle but are not.** People coming from Squiggle will expect the same results. Mitigation: document every deliberate difference; reuse Squiggle's `to` rule.
3. **Monte Carlo noise.** Percentiles move between seeds. Mitigation: fixed seed, three significant figures, a `:samples` knob.
4. **Correlations and independence.** Two different ranges are always independent; users may not expect that. Document it clearly.
5. **Heavy tails.** Lognormals with wide ranges produce huge outliers that break histograms and means. Mitigation: clip histogram axes to p1–p99 and say so.
6. **Name searchability.** "ubc" alone is dominated by the University of British Columbia. The repo name "uncertain-bc" carries discovery; see section 10.
7. **Decided:** `to` is lognormal for positive ranges and normal otherwise (ADR 0004).
8. **Decided:** units (`€`, `h`, `kWh`) are out of the MVP (ADR 0008).

## 9. Roadmap

| Phase | Outcome |
|---|---|
| 0. Walking skeleton ✅ (v0.1.0) | `ubc -e '1 + 2'` and `ubc -e '10 to 20'` print a result; CI, release binaries with goreleaser |
| 1. Calculator MVP ✅ (v0.2.0) | Lexer, parser, scalar + sample-set evaluation, `to`, suffixes, arithmetic, summaries, `P(...)`, default output line with sparkline, stdin and file modes |
| 2. REPL ✅ (v0.3.0) | History, persistent variables, `last`, `:` commands, `:hist` |
| 3. Sharing ✅ (v0.4.0) | `--json`, `--plain` (landed in phase 1), `:plot` to SVG/PNG, model files |
| 4. Expressiveness ✅ (v0.5.0) | User functions, more distributions, empirical data |
| 5. Maybe | Units, interactive sliders (Guesstimate-style "calculators") in a TUI |

## 10. Name

- **Binary:** `ubc`, short to type and an obvious nod to `bc`.
- **Repository and project:** `uncertain-bc` (Uncertain bc), module path `github.com/eferro/uncertain-bc`, installed with `go install github.com/eferro/uncertain-bc/cmd/ubc@latest`. This is the same pattern as ripgrep → `rg` or fd-find → `fd`.
- **Collision check (2026-11-21):** no `ubc` package in Homebrew, apt or crates.io; npm `ubc` is an abandoned placeholder and PyPI `ubc` is a utility bill calculator, both unrelated ecosystems. No GitHub repository matches "uncertain bc", "uncertainty calculator cli" or similar.
- Rejected alternatives: `napkin` (very crowded on GitHub), `circa`, `ballpark`, `sqc`.

## 11. References

- Guesstimate: https://www.getguesstimate.com and https://github.com/getguesstimate/guesstimate-app (MIT)
- Squiggle language: https://www.squiggle-language.com and https://github.com/quantified-uncertainty/squiggle (MIT)
- Quantified Uncertainty Research Institute (QURI): https://quantifieduncertainty.org
- GNU bc: https://www.gnu.org/software/bc/
- Napkin math, estimation from first principles: https://github.com/sirupsen/napkin-math
- Fermi problem: https://en.wikipedia.org/wiki/Fermi_problem
- Monte Carlo method: https://en.wikipedia.org/wiki/Monte_Carlo_method
- Douglas W. Hubbard, *How to Measure Anything* (calibrated 90 % confidence intervals)
- Pratt parsing: https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html
- Gonum (stats, plot): https://www.gonum.org
- GoReleaser: https://goreleaser.com
