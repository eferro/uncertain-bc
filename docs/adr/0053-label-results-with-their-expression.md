# ADR 0053: Label results with their variable name or short expression

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0025](0025-command-line-ux-additions.md)

## Context

Issue #31: a model's results were hard to read. `L30: ≈ 61.1k …` and `L31: 0.75` made the reader open the file to learn what each line was. ADR 0025 named assignments with `--echo` and rejected naming expressions after their source text because long expressions would make unreadable labels. A length limit removes that objection.

## Decision

- **A label replaces `L<n>:` exactly where the prefix used to appear**: when one program prints more than one expression. A single result still prints bare, so `$(ubc -e ...)` and `$(ubc model.ubc)` capture just the value. The maintainer chose this over labelling every result.
- **`-e` follows the same rule as files and stdin**: one `-e` that prints several expressions (`-e 'x = 2; x; x * 3'`) labels them; before, `-e` never had prefixes. Each `-e`, file, stdin and `:load`ed file is still its own program, so separate `-e` flags with one expression each stay bare.
- **The label is the expression's source text as written**, from its first token to its last, without the trailing comment and without reformatting spaces: `total = ≈ 61.1k …`, `P(total > 50k) = 0.75`. A bare variable is a special case of this, labelled with its name whatever its length.
- **The limit is 30 characters, inclusive.** It fits typical questions (`P(monthly_cost > 1.5k)`, `round(p(days, 0.85))`) and keeps the summary in view on an 80-column terminal. A longer expression keeps `L<n>:`.
- **`--json` carries the label as `"name"`**, the field `--echo` assignments already use; a longer expression has no `"name"`, as JSON never had line prefixes. JSON output stops escaping `<`, `>` and `&` as `<`, `>`, `&`, so a label reads as written; the output is never embedded in HTML, and every JSON parser reads both forms.
- **`--plain` stays unlabelled**: it exists to print bare values for scripts.
- **The REPL is unchanged**: one input line prints one result, which the user has just typed.
- The domain's `Result` carries the expression's `Source`, sliced from the program by token offsets; the CLI adapter decides whether it is short enough to be a label.

## Consequences

- Models read as reports without repeating variable names or using `--echo`.
- Programs with several expressions in one `-e` now print labels, a visible change for anyone who captured that output; a single expression per `-e`, `--plain` and `--json` parsers are unaffected (JSON only gains a field).
- Documented outputs (SPEC, tutorial, README, scenario testscripts) changed from `L<n>:` to labels.

## Alternatives considered

- **Label every result, even a single one:** breaks `$(ubc -e ...)` captures, the reason ADR 0025 kept `--echo` off by default.
- **Normalise spaces in the label:** would print something the user did not write; the source text is what they will search for in the file.
- **A limit in terminal columns or relative to the summary width:** harder to explain than a fixed number, for no practical gain.
