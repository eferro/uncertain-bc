# ADR 0052: Keep enough digits to tell the ends of a narrow interval apart

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0006](0006-bc-like-cli-with-honest-summary-output.md)

## Context

Issue #30: three significant figures per number (ADR 0006) are not enough when the spread is small compared with the value. `1M to 1.002M` printed `≈ 1M   90% [1M … 1M]`, and the availability scenario printed `balancer = ≈ 1   90% [1 … 1]`, while the 98% interval of the same lines already followed its spread (ADR 0040, ADR 0051). For availability, success rates or large values with a small spread, the hidden digits are exactly the information that matters.

## Decision

- **Every number of a summary but the 98% ends (median, inner interval ends, `--quantiles`) is rounded to the finer of 3 significant figures of its own value and 2 significant figures of the inner interval's width.** `1M to 1.002M` prints `≈ 1.001M   90% [1M … 1.002M]`; `0.999 to 0.9999` prints `≈ 0.99945   90% [0.999 … 0.99989]`.
- **The width is the one of the inner interval as configured**, P95 − P5 by default and the `--interval` width otherwise, because those are the ends that must stay apart. `--quantiles` uses the same width even though it hides the interval, so one rule covers every summary and the median prints the same with or without the flag.
- **A zero width (P5 = P95) leaves only the 3 significant figures**, as before.
- **The rule only adds digits**, so ordinary results do not change: `(8 to 25) * 12` still prints `≈ 170   90% [96 … 296]`. Of the documented outputs, only the availability scenario and `1000 to 1010` (`≈ 1005   90% [1000 … 1009.9]`) changed.
- **The 98% ends keep their own rule** (ADR 0051). Their step is never finer than the inner one, since P99 − P1 is at least P95 − P5, so they still nest the printed inner interval; the clamp of ADR 0051 stays for the zero-width case.
- Histogram labels, `:plot` ticks, `--plain` and `--json` are unchanged.

## Consequences

- Narrow intervals show distinct ends, and the 90% interval is never coarser than the 98% one around it.
- Numbers of a narrow result carry more digits (`1009.9` in place of `1010`). Those digits are within the Monte Carlo noise of the median, but not of the spread being shown, which is what the reader compares.

## Alternatives considered

- **More significant figures for every number:** would clutter every ordinary result to fix the narrow ones.
- **Two significant figures of the 98% spread:** coarser than needed for the inner ends, since the 98% spread is the wider one.
- **Keep 3 figures and only warn:** the reader would still not see the ends.
