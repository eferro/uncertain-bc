# ADR 0057: A user reference checked by the test suite, with errors shown

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0030](0030-tutorial-checked-by-the-test-suite.md), [ADR 0056](0056-modelling-guide-checked-like-the-tutorial.md)

## Context

Issue #35 asks for a user reference separate from the specification: the language, built-in functions, flags, REPL commands, output formats, exit codes and common errors with how to fix them. `docs/SPEC.md` stays the development contract and links to the reference, without duplicated text that can drift; the examples should be checked by the test suite where practical. The ADR 0030 converter, which already checks the tutorial and the modelling guide, only accepts commands that succeed, so it cannot check the error messages the reference has to show.

## Decision

- The reference is `docs/REFERENCE.md`, for users only, organised by what they look up: running `ubc`, the language, built-in functions, flags, the REPL, output (text, plain, JSON), exit codes and common errors, each with the message, its exit code and the fix. It is linked from the README, the tutorial and the specification.
- **The guide converter checks it** like the other guides (it joins the list of ADR 0056), so every command shown, its output and its error messages are compared with the binary.
- **Errors in console blocks.** In a console block, an output line that starts with `ubc: ` is stderr, since every message `ubc` writes there starts that way and no result does. A command that shows such a line must fail (`! exec`), and its stderr must match those lines exactly; its stdout is checked as before. Commands without such lines must still succeed. The exit code itself is not compared, because testscript's `exec` only tells success from failure; the codes are pinned by the scenarios in `cmd/ubc/testdata/script/` and stated in prose in the reference.
- **The flag list is the binary's.** The reference shows `ubc --help` in a console block, so the list of flags is checked byte for byte and cannot fall behind a new flag.
- **Completeness tests.** A test in `internal/domain` fails when a built-in function is not documented in the reference (written as `` `name( ``), and one in `internal/adapters/cli` when a REPL command is not (`` `:name ``). Both read the maps that define them, so a new function or command breaks the build until it is documented.
- What the converter cannot run stays in prose or in unchecked plain blocks: the REPL (it needs a terminal; ADR 0036 tests it separately), `--data` (it needs a data file outside `examples/`), `--watch`, `--random`, `--completion`, `--version`, `--ascii` (its sparkline ends in spaces that editors strip) and unknown flags (the flag package's message does not start with `ubc: `).
- `docs/SPEC.md` keeps its detail, because it is the contract the code is built against; it now points users to the reference. The reference describes behaviour through examples the suite checks, so the overlap cannot drift silently. The stale `--no-color` row is removed from the specification: `ubc` prints no colour and has no such flag, which the `--help` block in the reference made visible.

## Consequences

- A change in an error message breaks the reference until it is updated, as a change in output already did for the guides.
- The tutorial and the modelling guide can now show failing commands too, instead of describing them in prose as ADR 0030 required.
- Prose around the examples, such as the meaning of an exit code, is not checked.

## Alternatives considered

- **Move the user-facing parts of SPEC into the reference and leave SPEC with the internals:** removes the overlap, but the code and the ADRs are written against SPEC's precise rules (rounding, labels, override order), and a reader-first reference would lose them or become a second specification.
- **A separate marker for stderr (`2> ` or a `stderr` block):** explicit, but readers would see a format that is not what their terminal shows; the `ubc: ` prefix is what they see.
- **Check the exit code with a custom testscript command:** exact, but more converter code for codes the scenarios already pin.
- **A generated reference from `--help`, `:help` and the function tables:** no drift, but no room for the fixes and examples users need.
