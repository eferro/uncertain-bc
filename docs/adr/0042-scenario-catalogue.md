# ADR 0042: A scenario catalogue checked by the test suite

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0030](0030-tutorial-checked-by-the-test-suite.md)

## Context

Issue #20 asks for ready-to-use scenarios close to software engineering and product strategy, so people can start using `ubc` by changing a few assumptions instead of learning the whole language. The tutorial (ADR 0030) teaches the language step by step; a catalogue serves a different reader, one who has a question and wants a model to adapt. Like the tutorial, the catalogue must not rot when the language or the output changes.

## Decision

- The catalogue lives in this repository, under `scenarios/<area>/<name>.ubc`, one directory per area (`team`, and later others such as `product` or `delivery`). The first, reference scenario is `scenarios/team/replace-a-person.ubc`: interview hours × interviewer cost, agency or job-ad fees, months of vacancy × value lost, and the new hire's ramp-up.
- Every scenario follows one template, written as `#` comments in the file itself, so the file is self-explanatory wherever it is copied: the question; how to run it, with a `--set` and a `--sens` example; an **Assumptions** block, one variable per line with its range and the source of the numbers; the calculation, one comment per step; the result and a `P(...)` line that answers the question.
- `scenarios/README.md` explains the catalogue, shows the template and holds the index table (area, scenario, question). An entry counts as indexed when the README links to its path relative to `scenarios/` (`](team/replace-a-person.ubc)`).
- The test suite (`cmd/ubc/scenarios_test.go`) finds every `.ubc` under `scenarios/` at any depth and fails when one does not evaluate (non-zero exit or anything on stderr, with the fixed seed) or is missing from the index. It runs the files in-process through `cli.App`, so adding a scenario needs no new test code.
- The reference scenario also has a testscript (`testdata/script/scenarios.txtar`) that checks its output and the `--set` and `--sens` commands its header shows. testscript scenarios see `scenarios/` in their working directory, as they already see `examples/`.
- Only evaluation and indexing are enforced. The template is a convention checked in review: parsing comments to check it would make the test brittle for little gain.
- The catalogue is linked from the README and from the tutorial's "Where next".

## Consequences

- A new scenario is one `.ubc` file and one index row; forgetting the row or breaking the model turns the build red.
- A change to the language that breaks a scenario must fix the scenario in the same commit.
- Printed results of scenarios other than the reference one are not pinned, so their numbers can change silently with the RNG or the output format; that is acceptable because readers run them to get their own numbers.

## Alternatives considered

- **Put the scenarios in `examples/`:** mixes the tutorial's teaching models with question-driven ones and loses the area grouping.
- **A separate repository or a wiki:** easier to contribute to, but nothing would check that the models still run.
- **Pin the output of every scenario in a testscript:** strongest check, but every new scenario would need a hand-copied expected output, and the reader-facing value is the model, not the numbers.
