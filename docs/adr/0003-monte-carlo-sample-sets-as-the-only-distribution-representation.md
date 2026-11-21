# ADR 0003: Monte Carlo sample sets as the only distribution representation

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

Uncertain values can be represented symbolically (closed-form distributions), as point sets (PDF approximations) or as samples. Squiggle mixes the three; Guesstimate uses only samples (5,000 per metric).

## Decision

Every uncertain value is a SampleSet of N float64 samples (default N = 5,000, configurable). Operations are applied element-wise. Values that never touch a distribution stay plain scalars.

## Consequences

- Any function works on distributions for free: apply it sample by sample.
- Reusing a variable reuses its samples, so `x - x` is exactly 0 and dependencies between derived values are preserved.
- Two different distribution literals are always independent; correlation between distinct ranges is out of scope.
- Results carry Monte Carlo noise; output uses 3 significant figures (ADR 0006).
- Memory is O(N) per variable, negligible for a calculator.

## Alternatives considered

- **Symbolic distributions first** (Squiggle style): exact for simple cases but much more code; rejected for the MVP.
