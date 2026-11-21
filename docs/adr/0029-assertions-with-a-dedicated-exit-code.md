# ADR 0029: Assertions with `--assert` and a dedicated exit code

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro
- Complements: [ADR 0006](0006-bc-like-cli-with-honest-summary-output.md), [ADR 0028](0028-override-variables-with-set.md)

## Context

Issue #2 asks to use ubc as a check in scripts and CI: `ubc --assert 'P(cost > 100k) < 0.1' model.ubc`. A script needs to tell "the check failed" from "the model is broken", a failure has to explain itself, and with uncertain values "holds" needs a precise meaning.

## Decision

- `--assert EXPR` is repeatable. Assertions run after the whole program, in order, in the same session, so they see every variable, `--set` overrides included. They print nothing and bind nothing, not even `last`.
- **Exit code 4** means "an assertion failed". It is distinct from 1 (evaluation error), 2 (usage) and 3 (parse), so a CI job can tell a failed check from a broken model.
- **An assertion holds only when its value is exactly 1.** The value must be a condition: a comparison (a probability at top level), `P(...)`, or a plain 0/1 scalar. For uncertain operands a comparison must hold in every sample; a probability threshold is written explicitly (`P(x > 0) > 0.95`). Anything else (a number such as 2 or 0.5, an uncertain value) is an evaluation error, exit 1: it is a mistake in the check, not a failed check.
- **The message shows what was compared** (chosen by the user over printing only the expression's value, which for `P(cost > 100k) < 0.1` would be a useless `0`). For a top-level comparison it prints both operands and the operator, `ubc: --assert fixed > 200 failed: 150 > 200`; uncertain operands appear as their one-line summary without sparkline, and when the comparison held in some samples the probability follows, `(P = 0.61)`. Anything else prints its value. The operands are evaluated once and the probability is computed from those same samples, so the message is consistent with the verdict.
- **Every failure is reported**, then the run exits with 4, so one CI run shows all broken checks. An assertion that cannot be evaluated (bad syntax, unknown variable, not a condition) stops at once with its usual code, reported as `ubc: --assert EXPR: line 1, col C: message`, like `--set`.
- If the program fails, assertions are not evaluated and the program's error and exit code stand.
- `--assert` with the REPL is a usage error (`--assert needs a program: -e, files or stdin`): an interactive session has no "after". With `--watch`, assertions run after every run and the last run's code is returned.
- Messages go to stderr and ignore `--json` and `--plain`, which shape stdout only; `--ascii` is honoured.
- `domain.Assert` evaluates the condition and returns a `Check` (value, operands, operator); the CLI decides how to word it.

## Consequences

- A model doubles as a test: `ubc --set team=5 --assert 'P(cost > 100k) < 0.1' model.ubc` in CI.
- "Certainly holds" is strict: `--assert 'x > 0'` fails if a single sample is not positive. That keeps checks from passing by luck and pushes the tolerance into the expression, where it is visible.

## Alternatives considered

- **Reuse exit code 1:** a script could not tell a failed check from a broken model.
- **Hold when the probability is above 0.5 (or any fixed threshold):** a hidden, arbitrary tolerance.
- **Stop at the first failed assertion:** simpler, but CI would show one problem per run.
- **Print only the expression's value:** for the commonest form, `P(...) < t`, it is always `0`.
