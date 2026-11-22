# ADR 0059: The README as a short entry point, checked like the guides

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0056](0056-modelling-guide-checked-like-the-tutorial.md)

## Context

Issue #37 asks to rewrite the README as a short entry point now that the tutorial, the modelling guide, the how-to guides, the scenario catalogue and the reference exist. The old README had grown into a summary of every flag and the whole language in long paragraphs, which repeated those documents and drifted from them. The structure was agreed with the maintainer: an opening example, why use it, install, get started in three commands, how to read the output line, a short flag table and links; the inspiration from `bc`, Guesstimate and Squiggle stays explicit. Its console examples were not checked by the test suite, except through the scenarios of ADR 0030 that copied some of them by hand.

## Decision

- The README has, in this order: the opening example (`(8 to 25) * 12`), "Why" (ranges instead of numbers; `bc` for the terminal habit and the name, Guesstimate and Squiggle for modelling with ranges), "Install", "Get started in three commands" (`-e` with `P(...)`, a model file, a built-in scenario with `--set`), "Reading the output" (the annotated line of the reference), "Flags you will use most" (eight rows; `ubc --help` and the reference have the rest) and "Learn more" (tutorial, modelling guide, how-to guides, scenario catalogue, reference), then a short paragraph for contributors and the license. Detail lives in the linked documents, not in the README.
- **The README is checked like the guides.** It joins the list of `TestGuides_EveryConsoleBlockMatchesTheBinary`, which now names each guide by its path from the repository root, so every `console` block runs against the binary. Commands that are not `ubc` (install lines) go in `sh` or `powershell` blocks, which the converter ignores, and the annotated output line in a plain block. The interactive session is described in prose, since the converter cannot drive the REPL.
- The link test also checks that the README links to the tutorial and the scenario catalogue, so it now covers all five documents the README points readers to.

## Consequences

- The README's examples cannot drift from the binary: a change in the RNG, the output or a scenario breaks the test until it is updated.
- The README no longer lists every flag, REPL command and built-in function; readers who want them follow the links. A new feature does not need a README change unless it belongs among the most used flags.
- The install lines are not checked, as before.

## Alternatives considered

- **Keep the README unchecked:** less to maintain, but it is the first page readers see and the one most likely to go stale.
- **Keep the long feature summary below the links:** the README would still repeat the reference, which the issue asks to stop.
- **Show a REPL session in a `text` block:** readable, but unchecked output in the first page is the drift this decision avoids; one sentence and `:help` are enough.
