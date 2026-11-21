# ADR 0021: Manual mutation in the TDD loop, gremlins as the CI safety net

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro
- Complements: [ADR 0012](0012-coverage-and-mutation-gates.md), [ADR 0020](0020-reliable-mutation-gate.md)

## Context

gremlins (ADR 0012, 0020) checks the whole domain in about 30 seconds, and it will slow down as the domain grows. It only applies mechanical mutations: boundary changes, operator swaps, negated conditions. Some of the most plausible bugs in a numerical interpreter are beyond it: a deleted statement, a wrong constant, a default return value. Running it after every change would also slow down development.

Manual mutation testing (the `mutation-testing` skill) targets only the code that just changed. A handful of chosen mutants takes seconds, and it forces thinking about behaviour rather than chasing a percentage. But it depends on discipline, it is not repeatable, and it does not guard against regressions.

## Decision

1. **Inner loop, every increment:** after green and before committing, apply the `mutation-testing` skill by hand to the functions changed in the increment. Try 3–6 meaningful mutants (delete a statement, change a constant or operator, flip a condition, return a default value), run `make test-pkg`, and kill every survivor with a test or justify it as equivalent in `docs/mutation-testing.md`. Revert every mutant before committing. Then run the Test Desiderata review.
2. **CI safety net:** `make test-mutation` (gremlins, 85% efficacy) keeps running in CI on every push, and it stays out of the git hooks and `validate`.
3. **When gremlins gets slow** (more than about 3 minutes in CI), pushes switch to `--diff` (mutate only what changed since the previous commit). Full runs stay on release tags.

## Consequences

- Fast, targeted feedback while developing, including mutations gremlins cannot generate.
- An automatic, repeatable gate still catches regressions and skipped manual reviews.
- Each increment spends a few extra minutes on manual mutants. That is accepted as part of the definition of done.

## Alternatives considered

- **gremlins only:** automatic, but slower to give feedback and blind to non-mechanical mutations.
- **Manual mutation only:** fast, but nothing protects against regressions or a skipped step.
- **gremlins in the pre-commit hook:** it would make every commit wait 30+ seconds, and growing.
