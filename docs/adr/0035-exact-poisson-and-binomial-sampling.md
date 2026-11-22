# ADR 0035: Exact `poisson` and `binomial` sampling by recursive reduction

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0003](0003-monte-carlo-sample-sets-as-the-only-distribution-representation.md), [ADR 0005](0005-deterministic-seed-by-default.md)

## Context

Issue #9 asks for the discrete distributions `poisson(lambda)` and `binomial(n, p)`. They were left out because the textbook methods are slow for large parameters: Knuth's product of uniforms needs about λ draws per sample, and counting Bernoulli trials needs n. The fast exact methods (PTRS, BTPE) are long, full of tuned constants, and hard to test. Approximations (a normal for large λ or n) would be wrong in the tails and not whole numbers. Whatever we choose must draw only from the session RNG, so the fixed seed keeps output reproducible (ADR 0005).

## Decision

- **Small parameters, direct methods:** `poisson` with λ < 16 multiplies uniforms until the product drops below e^−λ; `binomial` with n < 16 counts n Bernoulli trials.
- **Large parameters, exact reduction** (Knuth, TAOCP vol. 2, 3.4.1 F), reusing the gamma and beta samplers we already have:
  - Poisson(λ): draw X ~ Gamma(m) with m = ⌊7λ/8⌋, the m-th arrival time of a rate-1 process. If X < λ, the result is m + Poisson(λ − X); otherwise it is Binomial(m − 1, λ / X).
  - Binomial(n, p): draw X ~ Beta(a, n + 1 − a) with a = ⌊n/2⌋ + 1, the a-th smallest of n uniforms. If X ≥ p, the result is Binomial(a − 1, p / X); otherwise it is a + Binomial(n − a, (p − X) / (1 − X)).
  Each step is exact, and the recursion shrinks the parameter geometrically, so even `poisson(1e12)` with 100 000 samples takes a fraction of a second. The split points and the limit of 16 only affect speed, not the distribution.
- **Parameters:** `lambda` must be positive and finite; `n` a positive whole number; `p` between 0 and 1, both included (`p = 0` and `p = 1` give a certain result). Anything else is an evaluation error naming the parameter, like the other distributions. A NaN or infinite `lambda` would make the gamma draw loop forever, so it is rejected explicitly.
- Samples are whole numbers stored as float64, like every other value (ADR 0003); there is no integer type.

## Consequences

- Exact draws for every parameter in about 60 lines, with no tables or tuned constants, and only the session RNG, so output stays deterministic per seed.
- The tests compare the empirical CDF with the exact one on both sides of the limits, which kills every mutant except the equivalent ones listed in `docs/mutation-testing.md` (split points, limits, boundaries of continuous draws).
- The sparkline histogram of a discrete result over a narrow range has gaps (`▂▃▁▆▁█▁█`), because 16 cells span fewer than 16 distinct values. That is honest about the support, and it is left as it is.

## Alternatives considered

- **Direct methods only, with an upper limit on λ and n:** simpler, but it arbitrarily refuses common models (requests per month, visitors per year).
- **Normal approximation above a threshold:** fast, but not exact, not whole numbers, and wrong in the tails that a risk calculator cares about.
- **PTRS (Hörmann 1993) / BTPE (Kachitvichyanukul & Schmeiser 1988):** faster per draw, but much more code with constants that are hard to verify, for a speed we do not need.
