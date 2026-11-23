# ADR 0063: Every scenario ends with its output, checked by the test suite

- Status: Accepted
- Effective date: 2026-11-23
- Deciders: eduardo.ferro
- Complements: [ADR 0030](0030-tutorial-checked-by-the-test-suite.md), [ADR 0042](0042-scenario-catalogue.md), [ADR 0044](0044-team-scenarios.md)

## Context

Issue #40 asks every scenario in `scenarios/` to end with the output its run command prints, so a reader sees the answer (and what an edited assumption changes) without installing `ubc`. A hand-written output comment would go stale with the next formatting change, so it has to be checked like the tutorial (ADR 0030) and regenerated rather than edited. Until now each scenario also had a `scenarios-<area>-<name>.txtar` that pinned the medians and the answer of its default run: a second copy of the same output.

## Decision

- Every `scenarios/<area>/<name>.ubc` ends with a blank line and an **output block**: the line `# Output (seed 1, 5000 samples):` followed by one comment line per output line, `#   ` plus the line exactly as `ubc --scenario <area>/<name>` prints it (default seed and sample count, shared formulas first, no `--echo`). The block is the last thing in the file; `formulas.ubc` prints nothing and has none.
- **The scenario is the expectation.** `TestScenarios_EveryScenarioEndsWithItsOutput` (in `cmd/ubc`) runs each scenario the way `--scenario` does, rebuilds the file with that output in its block (everything after the header line is replaced) and fails, per scenario, with a line diff and the hint `run make scenario-outputs` when the file differs.
- **Regenerated, not hand-edited.** `make scenario-outputs` runs the same test with `-update-scenario-outputs`, which writes the rebuilt files instead of failing. After an intended output change, the diff of `scenarios/` shows that change across the whole catalogue and is reviewed in the same commit. The logic lives in one test file, so checking and regenerating cannot disagree.
- **One source of truth.** The default run is no longer pinned in the `.txtar` files: each `scenarios-<area>-<name>.txtar` keeps only the runs the block does not cover (`--set`, `--sens`, `--data`, `--echo`), and `embedded_scenarios.txtar` only checks that the embedded scenario runs from another directory. Empty stderr on the default run stays checked by `TestScenarios_EveryFileEvaluates`.
- The block is a comment, so it changes neither the evaluation nor `--scenarios`; `--scenario-source` prints it with the rest of the file, and a copied and edited scenario keeps the block of the original until its author removes or updates it.

## Consequences

- Readers of the catalogue on GitHub or the website see each scenario's answer next to its assumptions.
- Changing the RNG, the formatting or a scenario breaks the test until `make scenario-outputs` is run, which is the point: the scenarios change in the same commit as the behaviour.
- The test runs every scenario a second time (they already run in `TestScenarios_EveryFileEvaluates`), a few seconds in total.

## Alternatives considered

- **Keep the `.txtar` expectations as well:** two copies of the same output to update on every format change, for no extra protection.
- **A generator under `tools/`:** a second program with the same rules as the test; the test flag keeps one implementation.
- **Line labels as in the issue (`L41: …`):** results are already labelled by name or expression (§6), so the block shows exactly what `ubc` prints.
