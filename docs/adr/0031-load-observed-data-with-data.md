# ADR 0031: Load observed data from a file with `--data`

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro
- Complements: [ADR 0028](0028-override-variables-with-set.md)

## Context

Issue #5 asks to load observed data into a variable, for example lead times exported from a tracker: `ubc --data lead_time=times.csv model.ubc`. The language already has empirical sample sets (`[3, 5, 8, 13]`, resampled with replacement), but typing hundreds of observations into a list is impractical. The open questions were the file format and whether this is a flag or a language function.

## Decision

- **A flag, not a language function.** `--data NAME=PATH` is repeatable. The language stays free of I/O: a model reads the same everywhere, the domain needs no file port (ADR 0010), and a function such as `data("times.csv")` would make a model depend on the working directory. The path is split at the first `=` and the name is checked like a `--set` name.
- **One plain number per line.** The format is what a spreadsheet column, `cut` or a SQL client exports. Numbers are parsed as Go floats (`12`, `-3.5`, `1e3`); ubc suffixes (`1k`, `5%`) are not accepted, because exported data never uses them. Surrounding spaces, blank lines and `\r\n` line ends are ignored. Anything else, including a header line, `NaN` or an infinity, is an error that names the file and line, so a stray header never becomes silent data. CSV with several columns is not supported; `cut -d, -f2` already turns it into this format.
- **Same value as a list literal.** The numbers become an empirical SampleSet resampled with replacement from the session RNG, exactly what `[v1, v2, ...]` gives, so `--data x=f -e x` prints the same as `-e '[...]'` with the same numbers. The resampling code is shared.
- **An override, like `--set`.** A model can keep a placeholder (`lead = 2 to 10`) for when no data is given; `--data` replaces it with the rules of ADR 0028 (later assignments skipped, the last `--data` for a name wins and a `--set` for it wins over `--data`, forgotten by `:samples`, applies to the REPL and every `--watch` run).
- **Order: `--data` before `--set`.** Data files are loaded first, so an override can derive from data (`--set 'worst=max(lead)'`). Go's `flag` does not keep the interleaving of different flags, and a fixed order is easier to explain.
- **Errors.** A malformed flag value (no `=`, bad name, empty path) is a usage error, exit 2, checked with the other flags. An unreadable file, a bad line or a file without numbers is an evaluation error, exit 1, like any unreadable program file, reported as `ubc: --data NAME: ...`; the first failing file stops the run before any program.
- **Reading files is the CLI adapter's job**; the domain gets the numbers through `OverrideData(name, values, env)`, which also rejects an empty list.
- `--watch` reads the data files again on every run but does not watch them; only program files trigger a run.

## Consequences

- Real observations can drive a model without editing it, and a model can ship with an estimate to fall back on.
- Resampling the data draws from the RNG before any program runs, so adding `--data` changes the samples of later ranges; output stays deterministic for the same command line.
- Data files with headers or several columns need a preprocessing step (`tail -n +2`, `cut`).

## Alternatives considered

- **A `data("file")` function in the language:** puts I/O in the domain and makes models depend on where they run.
- **CSV with a column selector (`--data lead=times.csv:2`):** more syntax for a case shell tools already cover; can be added later without breaking this format.
- **Accept ubc number syntax per line (`1k`, `5%`):** would let a model file and a data file share notation, but exported data never needs it, and it would tie the data format to the language's lexer.
- **Bind the data as a plain variable that the model can reassign:** inconsistent with `--set`, and the model's placeholder would silently win.
