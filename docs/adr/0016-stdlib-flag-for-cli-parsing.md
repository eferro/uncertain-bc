# ADR 0016: Standard library `flag` for CLI parsing

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

`ubc` has no subcommands. Its interface is a handful of flags (SPEC §1), each with a short and a long name, and `-e/--eval` can be repeated. The project rule is "standard library unless a dependency clearly pays off". Project B uses cobra, which suits a tool with many subcommands.

## Decision

- Parse flags with the standard library's `flag` package, using a `flag.NewFlagSet("ubc", flag.ContinueOnError)` created inside `run`. `run` never uses the global `flag.CommandLine`, so it stays testable.
- Register short and long names against the same variable (`-n` and `--samples`). Go's `flag` accepts one or two dashes for any name.
- A small `multiString` type implementing `flag.Value` collects repeated `-e` expressions in order.
- The usage text is written by hand to match SPEC §1. Flag errors print the usage to stderr and exit with code `2` (SPEC §1).
- Combined short flags (`-ne`) are not supported, and flags must come before file arguments (`flag` stops at the first positional argument). Both are documented as known differences from GNU tools.

## Consequences

- No dependency, and full control over help and error output.
- A few lines of glue code (double registration, `multiString`).
- If `ubc` ever grows subcommands or the flag glue becomes painful, revisit this with a new ADR (spf13/pflag would be the first candidate).

## Alternatives considered

- **spf13/pflag:** GNU-style flags with combined short options. One dependency for a small UX gain.
- **cobra:** a full framework, and overkill without subcommands.
