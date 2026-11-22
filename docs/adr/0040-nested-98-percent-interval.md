# ADR 0040: A nested 98% interval in every summary

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0006](0006-bc-like-cli-with-honest-summary-output.md), [ADR 0034](0034-configurable-interval-width.md)

## Context

Issue #18 asks to always show Guesstimate's detail percentiles (P1, P5, P50, P95, P99) without a flag, keeping one line per result. The median and P5/P95 are already on the line, so only P1 and P99 are new. An LLM council session (2026-11-22) chose a second, nested interval on the same line and two significant figures for its ends, and left open how it combines with `--interval`, `--quantiles`, `--json` and the other places that print a summary.

The prototype on the tutorial examples showed that two significant figures *of each end* break on narrow distributions far from zero: `1000 to 1010` printed `90% [1000 … 1010]   98% [1000 … 1000]`, an outer interval that looks narrower than the inner one, and `11.5M to 11.7M` printed `98% [11M … 12M]`.

## Decision

- **Same line, nested interval:** `≈ MEDIAN   90% [P5 … P95]   98% [P1 … P99]   SPARKLINE`. About 16 more characters, no repeated numbers, and still one line per result for pipes.
- **The 98% ends keep two significant figures of their distance**, P99 − P1, not of their own size (chosen by the user after the prototype). P1 and P99 rest on few samples; with 5000 samples P99 moved about 30 % between seeds, and that noise grows with the spread, not with how far the values are from zero. `(8 to 25) * 12` prints `98% [80 … 380]`, `1000 to 1010` prints `98% [998 … 1012]`, `11.5M to 11.7M` prints `98% [11.46M … 11.74M]`. When P1 equals P99 there is no distance to measure, so the ends use the usual three significant figures. A small negative end that rounds to zero prints `0`, never `-0`.
- **`--interval W` sets the inner interval only; the outer one stays at 98%.** When W is 98 or more the outer interval is left out, since it would no longer be the wider one and would repeat or contradict the inner one.
- **`--quantiles` replaces both intervals.** Whoever lists percentiles has said exactly which ones to see.
- **`--json` adds `p1` and `p99` as standard fields**, next to `p5` and `p95`. Additive, so existing scripts keep working; they are raw values, not rounded.
- **Everywhere a summary appears:** results, REPL and `--echo` assignments, `:vars`, `:plot` titles and the operands of a failed `--assert`. The summary is one concept (ADR 0034), so it is built in one place.

## Consequences

- The tails are visible without a flag, and a model's spread can be read at two levels of confidence.
- Every documented output changes: the `.txtar` scenarios, the tutorial, README, VISION and SPEC examples were regenerated in the same commit.
- Discrete or empirical distributions can print the same numbers twice (`90% [3 … 13]   98% [3 … 13]`). That is a true reading, so it stays.
- With few samples (`-n 100`), P1 and P99 come from a single sample each; the coarse rounding is the only protection.

## Alternatives considered

- **The five percentiles in parentheses on the same line:** repeats P5, P50 and P95 and reaches about 100 characters.
- **A second line per result:** doubles the output and breaks one line per result.
- **Output that changes between a terminal and a pipe:** the same command would print different things.
- **More samples only for the tails:** session variables have a fixed sample count.
- **Two significant figures of each end** (the council's first choice): rejected after the prototype for the artefacts above.
- **Three significant figures like the rest:** consistent, but claims a precision the tails do not have (ADR 0006).
