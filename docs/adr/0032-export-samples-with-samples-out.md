# ADR 0032: Export raw samples with `--samples-out`

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0031](0031-load-observed-data-with-data.md)

## Context

Issue #6 asks to write the raw samples of a result to a file (`--samples-out FILE`) so that they can be analysed with other tools: pandas, R, a spreadsheet, a plotting library. The summary line and `--json` only give statistics. The open questions were which results to export and in which format.

## Decision

- **Every variable, one column each, plus `last`.** The file holds the variables in the order they were first assigned, `--data` and `--set` overrides included, then `last` when the program printed an expression and did not assign `last` itself. Row *i* is sample *i* of every column, so the correlations the model builds (`revenue = price * units`) survive the export: this is what makes the raw samples worth more than the summaries. Plain numbers and probabilities repeat on every row, so all columns have `--samples` rows and the table stays rectangular. The user chose this over exporting only `last` or only named variables.
- **Format from the extension.** `.csv` writes a header row with the names and one row per sample; `.json` writes one object `{"name": [samples...], ...}` with keys in variable order. Any other extension is a usage error (exit 2), checked with the other flags, like `:plot` with `.svg`/`.png`. The user chose this over CSV only or a separate `--samples-format` flag. The extension check is case-sensitive, as `:plot`'s is.
- **Numbers in their shortest exact form** (`strconv` `'g'`, precision `-1`), so a file read back gives the same float64 values. Non-finite samples are an empty cell in CSV (the usual missing value) and `null` in JSON, as in `--json` (SPEC §6).
- **When it is written.** After the whole program succeeds and before the assertions are checked, so a failing `--assert` still leaves the samples behind to look into. A program that fails writes no file, as a half-run model would export a misleading state. With `--watch` the file is written again after every successful run.
- **Errors.** Nothing to export (no variable and no `last`) and a file that cannot be written are evaluation errors (exit 1), reported as `ubc: --samples-out FILE: ...`, and the assertions are not checked. `--samples-out` in the REPL is a usage error, like `--assert`: there is no end of program to write at; `:plot` covers one-off exports there.
- **Writing the file is the CLI adapter's job.** The application session gives the values through `VariablesAndLast()`; the domain is unchanged.

## Consequences

- A model's samples can be analysed with any tool, keeping the pairing between variables.
- Files grow with `--samples` × variables; at the default 5000 samples this is small.
- Built-ins `pi` and `e` and user functions are not exported; they are not results of the model.

## Alternatives considered

- **Only `last`, one number per line:** would round-trip with `--data`, but loses every intermediate variable and the correlations between them.
- **Named variables only:** a model that ends with an expression would export nothing about its result.
- **A `--samples-format` flag:** one more flag for a choice the file name already makes.
- **Writing `NaN`/`Inf` in CSV:** readers disagree on their spelling; an empty cell is understood everywhere as missing.
