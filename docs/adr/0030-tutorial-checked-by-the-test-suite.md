# ADR 0030: A tutorial whose examples are checked by the test suite

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro
- Complements: [ADR 0014](0014-cli-acceptance-tests-with-testscript.md)

## Context

Issue #4 asks for a step-by-step tutorial, from plain arithmetic to models with ranges, functions and probabilities, with real-world examples. Every example must live in `examples/` and its output must be checked by a testscript scenario, so the tutorial cannot drift from the binary. A hand-written `.txtar` next to the tutorial would check the binary, but nothing would check that the Markdown shows the same output: the drift would move from "binary vs. docs" to "scenario vs. docs".

## Decision

- The tutorial is `docs/TUTORIAL.md`, linked from the README. Each step has one model in `examples/` (`arithmetic.ubc` … `probability.ubc`), followed by three real-world models: a project estimate, a Fermi problem (piano tuners in Madrid) and a monthly cloud bill that also shows `--set` and `--assert`.
- **The tutorial itself is the scenario.** `TestTutorial_EveryConsoleBlockMatchesTheBinary` (in `cmd/ubc`) reads `docs/TUTORIAL.md`, turns it into a testscript archive at test time and runs it with the same setup as `testdata/script/`. There is no second copy of the expected output.
- Conventions the converter enforces:
  - In a ```` ```console ```` block, a line starting with `$ ` is a command and the lines after it, up to the next command, are its exact stdout. A command without output must print nothing. Every command must start with `ubc ` and exit with 0; stderr is not checked. A block that is empty or starts with output is an error.
  - A ```` ```ubc ```` block shows a model file: it must be followed directly by a console block whose first command runs that file (the last `.ubc` argument), and it must match the file byte for byte. A `ubc` block that is never run is an error.
  - Other fenced blocks are ignored, and fence-like lines inside a block are content.
- Every scenario sees the example files under `examples/` in its working directory, so tutorial and README commands are exactly what a reader types from the repository root (`ubc examples/homelab.ubc`).
- Commands that fail (a failing `--assert`, exit 4) are described in prose instead of shown, to keep the converter small.

## Consequences

- Changing the language, the RNG or the output format breaks the tutorial test until the tutorial is updated, which is the point: the docs change in the same commit as the behaviour.
- Writing a tutorial example is just Markdown plus a file in `examples/`; the expected output is copied from a real run.
- The converter is test code with its own unit tests; it lives in `cmd/ubc/tutorial_script_test.go` and is not part of the binary.

## Alternatives considered

- **A hand-written `tutorial.txtar` with the same outputs:** checks the binary but not the Markdown, so the two copies can drift apart.
- **Generate the Markdown from the examples:** no duplication either, but it needs a generator, a template and a "regenerate" step that is easy to forget.
- **Check only that every example runs without error:** cheaper, but the outputs shown to readers could still be wrong.
