# ADR 0034: Configurable interval width with `--interval`

- Status: Accepted, amended by [ADR 0040](0040-nested-98-percent-interval.md)
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0006](0006-bc-like-cli-with-honest-summary-output.md), [ADR 0025](0025-command-line-ux-additions.md)

## Context

Issue #8 asks to choose the width of the printed interval, `--interval 80` printing `80% [p10 … p90]` instead of the fixed 90 % of ADR 0006. `--quantiles 10,90` (ADR 0025) already shows the same two numbers, so the issue asked whether a separate flag is worth it.

## Decision

- **A separate flag, `--interval W`.** `--quantiles` prints a list of labelled percentiles (`p10 10.8   p90 18.4`), a different shape from the one-line summary people learn to read. Many users simply prefer an 80 % or 50 % interval; keeping the familiar `W% [lo … hi]` form for them is worth one flag. `--quantiles` stays for arbitrary percentile lists.
- **Central interval.** W is the share of samples inside; the ends are the (50 − W/2)th and (50 + W/2)th percentiles, as `a to b` reads its bounds (ADR 0004). W is a percentage strictly between 0 and 100, decimals allowed (`99.5`); anything else is a usage error. Not given means 90, so `--interval 90` changes nothing.
- **Everywhere a summary appears.** Results, REPL and `--echo` assignments, `:vars`, `:plot` titles and the operands of a failed `--assert`: the summary is one concept, so its width is one setting. The sparkline and `:hist` keep spanning P1–P99.
- **JSON keeps its schema.** `p5` and `p95` stay so scripts do not break; when `--interval` is given its ends are added to the `quantiles` object that `--quantiles` already defines (`--interval 80` gives the same JSON as `--quantiles 10,90`).
- **Not with `--quantiles`.** Both decide what replaces the default interval, so giving both is a usage error rather than a silent precedence rule.

## Consequences

- An 80 % or 50 % reading is one flag away, in the familiar format.
- The default output, and every existing script, is unchanged.
- Two flags touch the same part of the line; the usage error makes the overlap explicit.

## Alternatives considered

- **No new flag, document `--quantiles 10,90`:** no code, but a different output shape and no `80%` label.
- **Make `--quantiles` with two values print an interval:** overloads one flag with two output formats depending on the list length.
- **Rename `p5`/`p95` in JSON after the width:** more faithful, but breaks every script that reads the current fields.
