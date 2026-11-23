# ADR 0028: Override variables from the command line with `--set`

- Status: Accepted, amended by [ADR 0066](0066-an-assigned-comparison-keeps-its-samples.md)
- Effective date: 2026-11-21
- Deciders: eduardo.ferro
- Complements: [ADR 0025](0025-command-line-ux-additions.md)

## Context

Issue #1 asks for scenarios without editing the model: `ubc --set team=5 --set rate=0.8 model.ubc`. A model assigns its own value to `team`, so the command line needs a rule for which value wins, and the variables computed from `team` must see the winning value.

## Decision

- `--set NAME=EXPR` is repeatable. The value is split at its first `=`, so `EXPR` may contain comparisons (`--set 'ok=a==b'`); spaces around `NAME` are ignored. `EXPR` is a single ubc expression (`--set 'effort=10 to 20'`); a bare comparison is a probability, as at top level.
- **The model's assignment is skipped.** The overrides are evaluated first, in the order given, in the same session as `-e`, files and stdin. After that, an assignment to an overridden name keeps the override and does not evaluate its right-hand side, so an invalid default in the model cannot fail an overridden run. Variables computed after it use the override. With `--echo`, the skipped assignment prints `name = <override>` at its usual place, so the report shows the scenario.
- Evaluating `--set` last was rejected: every variable the model derives from the name would keep the model's value, which defeats the purpose.
- Overrides are evaluated before any program, so they may use earlier overrides and built-in variables (`--set rate=20 --set 'team=rate / 10'`), but not the model's variables. Giving the same name twice: the last one wins.
- `--set` itself prints nothing.
- A value without `=` or whose name is not an identifier (or is the keyword `to`) is a usage error, exit 2, checked with the other flags before anything is evaluated. A bad expression keeps the usual codes: parse error 3, evaluation error 1, reported as `ubc: --set NAME: line 1, col C: message`. The first failing override stops the run.
- Overrides apply to `-e`, files, stdin, every `--watch` run (each run starts a fresh session with the same overrides) and the REPL, with the same rule: typing `team = 3` in a REPL started with `--set team=5` prints `team = 5`. `:samples` clears variables and their overrides with them. Overrides are not part of the `:save` transcript, which holds only what was typed.
- The override lives in the domain (`Env` remembers overridden names) because skipping an assignment is evaluation semantics; the CLI only parses the flag and reports errors.

## Consequences

- One model file serves many scenarios: `ubc --set team=5 model.ubc`, `ubc --set 'team=4 to 8' model.ubc`.
- Overriding a name the model never assigns simply defines it, so a model can also leave parameters to the command line.
- Skipping the right-hand side means an overridden assignment draws no random samples, so the samples drawn after it differ from a run without `--set`. Results stay deterministic for the same command line.

## Alternatives considered

- **Evaluate `--set` after the model:** derived variables would ignore the override.
- **Re-evaluate the model's right-hand side and discard it:** keeps the random stream aligned with the run without overrides, but fails on defaults that only make sense with the override, and costs a draw for nothing.
- **Reject a REPL assignment to an overridden name with an error:** louder, but a different rule from files; printing the kept value already shows what happened.
