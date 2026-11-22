# ADR 0025: Command-line UX additions

- Status: Accepted, amended by [ADR 0053](0053-label-results-with-their-expression.md)
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

After phase 4 the user asked for every improvement that makes `ubc` nicer to use from the command line, as opposed to new modelling capabilities. Output only named results in `:vars`: a model's output said `L18: ≈ 101 …` and the reader had to look up what line 18 was, and the REPL echoed only the value of an assignment.

## Decision

**Names in the output**

- In the REPL, an assignment prints `name = summary` (`watts = ≈ 19 …`). Expressions print as before.
- `--echo` makes files, stdin and `-e` print assignments too, as `name = summary`. Expressions keep their `L<n>:` prefix when a program prints more than one. With `--json`, an echoed assignment carries a `"name"` field. With `--plain`, it prints `name = value`.
- `--echo` is off by default, so `$(ubc -e ...)` and existing scripts keep their output.

**Charts and percentiles outside the REPL**

- `--hist` prints the 20-row histogram (the one `:hist` draws) under every uncertain result, in files, stdin, `-e` and the REPL. It cannot be combined with `--json` or `--plain`, which are for machines.
- `--quantiles 10,50,90` replaces `90% [p5 … p95]` with the listed percentiles (`p10 … p50 … p90 …`). In `--json` it adds a `"quantiles"` object and keeps the standard fields. An empty list means the default.

**Watch mode**

- `--watch FILE...` evaluates the files, then evaluates them again whenever one changes, printing `--- FILE changed ---` before each run. It needs at least one file.
- Changes are detected by polling modification time and size every 300 ms, with the standard library only (no fsnotify dependency). Each run uses a fresh session with the same seed, so the output changes only when the model does. Errors are reported and watching continues.
- The wait between polls is injected into the CLI adapter, so tests drive it synchronously.

**Shell completion**

- `--completion bash|zsh|fish` prints a completion script for flags and file names. The list of flags is read from the same `FlagSet` the CLI parses, so completion cannot drift from the real flags.
- Tests check every script's syntax with the shell itself when it is installed, and run the bash completion function to see what it offers.

## Consequences

- A model reads as a report (`ubc --echo model.ubc`) without repeating variable names at the end just to print them.
- The REPL transcript says which variable each line is.
- With an editor on one side and `ubc --watch` on the other, a model gets a fast feedback loop, much like a spreadsheet.

## Alternatives considered

- **Always print assignments in files:** noisy for models with many intermediate variables, and it changes the output of existing scripts.
- **Name expressions after the source text:** long expressions would make unreadable labels.
