# ADR 0051: Round the ends of the 98% interval outwards

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0040](0040-nested-98-percent-interval.md)

## Context

Issue #29: ADR 0040 rounds P1 and P99 to the *nearest* multiple of two significant figures of their distance. In a right-skewed distribution the distance is set by the long tail, so the lower end is rounded far too coarsely: `(1 to 10) * 1000` printed `90% [999 … 9710]   98% [1000 … 16k]`, with a real P1 of 615.9. The outer interval started inside the inner one, and with seeds 1 to 40 the nesting broke in 27 of them.

## Decision

- **Round outwards:** P1 down and P99 up. The printed 98% interval always contains the true P1 and P99.
- **Never coarser than two significant figures of the end's own value:** the step is the finer of the two (distance and own value), so `615.9` prints `610`, not `1000`. The distance rule of ADR 0040 still applies where it is finer, as in `1000 to 1010`, which now prints `98% [997 … 1013]`.
- **Never inside the printed inner interval:** rounding P1 down can still land above P5 rounded to three significant figures (P1 = P5 = 1004 prints `1004` against `1000`). The outer ends are clamped to the printed inner ends, so the 98% interval always contains the 90% one as printed.
- **Exact steps:** the end is first rounded to the nearest step and then moved one step out only if that crossed it, so a value already on a step (`0.29`, `0.07`) stays there instead of picking up the float error of x·10^k (ADR 0039).
- The `-0` guard of ADR 0040 goes away: an end never rounds to zero, because its step is always finer than its own value.
- JSON keeps the raw `p1` and `p99`; only the text summary changes.

## Consequences

- The outer interval is an honest bound and always nests the inner one, whatever the skew.
- The lower end of skewed results shows more digits (`610`, `75` in place of `1000`, `80`), and ends crossing zero may show a small negative value (`-1.5` in place of `0`), which is the real P1.
- Every documented 98% interval was regenerated: scenarios, tutorial, README, VISION and SPEC.

## Alternatives considered

- **Outwards only, with the distance step:** keeps P1 below P5 in value but still prints `0` for P1 = 615.9 in a 15k spread, losing the lower tail.
- **Two significant figures of each end** without the distance rule: rejected in ADR 0040 for narrow distributions far from zero, where the outer interval looked narrower than the inner one.
- **Round the inner interval outwards too:** would change every 90% interval for a problem that only concerns the outer one.
