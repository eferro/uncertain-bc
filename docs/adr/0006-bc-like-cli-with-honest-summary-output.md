# ADR 0006: bc-like CLI with an honest one-line summary

- Status: Accepted, amended by [ADR 0022](0022-unicode-output-everywhere.md)
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

The tool replaces `bc` in daily use, so it must accept input the same ways (REPL, `-e`, file, stdin) and print short answers. A distribution has no single right answer to print, and Squiggle's CLI shows how unhelpful raw output is (`96.98057045115692`, or just `Sample Set Distribution`).

## Decision

- Modes: REPL on a TTY, `-e EXPR`, file arguments, stdin.
- Scalars print plainly (`2 + 2` → `4`). Sample sets print `≈ MEDIAN   90% [P5 … P95]   SPARKLINE` with 3 significant figures. Probabilities print with two decimals.
- `--json` and `--plain` for scripts. Colour and Unicode only on a TTY.
- REPL commands start with `:` (`:hist`, `:plot`, `:seed`, `:vars`, …) so they never collide with the language.
- Charts: sparkline always, `:hist` for a terminal histogram, `:plot` exports SVG/PNG.

Details in [SPEC.md](../SPEC.md) §1, §6 and §7.

## Consequences

- Usable in pipes and `$(...)` like `bc`.
- No false precision on screen.
- Inline terminal images (kitty, sixel) are not needed for the MVP; they may come later as an option.
