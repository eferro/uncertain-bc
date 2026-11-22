# ADR 0039: Round printed numbers by exact powers of ten

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0006](0006-bc-like-cli-with-honest-summary-output.md)

## Context

Issue #17: `p99 11.599999999999998M` instead of `p99 11.6M`. The output rounds a number to N significant figures as `round(x · 10^d) / 10^d`. When the rounding step is ten or more, `d` is negative and `10^d` (`1e-5`, `1e-15`…) is not exact in binary, so the division lands next to the round value instead of on it: `11560770` rounds to `11599999.999999998`, and scaling it to millions keeps the noise. The issue blamed the suffix division, but `11600000 / 1e6` is already `11.6`; the noise comes from the rounding. The same flaw shows in scalars with more than ten digits (`100000010000000` printed `100000009999999.98`).

## Decision

- **Round with a step that is exact.** When the step is ten or more, divide by `10^-d`, round, and multiply back. Whole powers of ten up to `1e22` are exact doubles, and a single division or multiplication of exact values is correctly rounded, so the result is the double closest to the round decimal, which prints without noise. Steps below one keep `round(x · 10^d) / 10^d`, whose last operation already divides by an exact power of ten.
- **Fix it at the root, not in the suffix code.** Every printed number goes through the same rounding (summaries, `--quantiles`, `:vars`, `:plot` and `:hist` labels, scalars), so one fix covers all of them.

## Consequences

- Printed numbers carry no binary noise; the values themselves are unchanged.
- Rounding to tens is exact with either formula for the magnitudes we print, so the mutant that leaves that case on the old path survives and is listed as equivalent in `docs/mutation-testing.md`.

## Alternatives considered

- **Divide by the suffix scale before rounding:** fixes `M`/`B`/`T`, but not four-digit or scalar output, and moves the choice of suffix away from the rounded value (SPEC §6).
- **Round again after scaling:** hides the symptom in one place and leaves the noisy value elsewhere.
- **Format with `strconv.FormatFloat(x, 'g', 3, 64)`:** switches to exponent notation for large numbers and changes trailing-zero handling.
