# ADR 0004: `a to b` is a 90 % interval, lognormal when both ends are positive

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

Guesstimate interprets a range as a 90 % confidence interval of a normal distribution by default. Squiggle interprets `a to b` the same way but uses a lognormal when both ends are positive. With a normal, wide ranges of positive quantities (costs, durations, counts) produce negative samples, which are nonsense.

## Decision

`a to b` means P(X < a) = 5 % and P(X > b) = 5 %. If `a > 0` and `b > 0` the distribution is lognormal; otherwise it is normal. `a >= b` is an error. Formulas are in [SPEC.md](../SPEC.md) §4.

## Consequences

- Matches Squiggle's behaviour, the closest reference for users.
- Positive quantities never go negative, and their medians are geometric means (`8 to 25` → ≈ 14.1).
- Users who want a symmetric distribution use `normal(mean, sd)` explicitly.
- Heavy right tails are common; histograms clip to P1–P99.
