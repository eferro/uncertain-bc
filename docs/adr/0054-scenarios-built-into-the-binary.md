# ADR 0054: Scenarios built into the binary, with `--scenarios`, `--scenario` and `--scenario-source`

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0042](0042-scenario-catalogue.md), [ADR 0043](0043-reusable-formulas-for-the-scenarios.md)

## Context

The scenario catalogue (ADR 0042) and its shared formulas (ADR 0043) live in `scenarios/` and are run by path, `ubc scenarios/formulas.ubc scenarios/<area>/<name>.ubc`. People who install `ubc` with Homebrew, Scoop, the Linux packages or the install scripts (ADR 0026) get a binary and no `scenarios/` directory, so the catalogue only worked from a clone. Issue #32 agreed with the maintainer to embed the catalogue in the binary, list it with `--scenarios`, run one scenario with `--scenario NAME` combined with every other flag, load the shared formulas automatically, and leave the exact flag for copying a scenario to this ADR.

## Decision

- **Embedding.** `scenarios/scenarios.go` (package `scenarios`, at the root of the module because `go:embed` cannot reach a parent directory) embeds `*.ubc` and `*/*.ubc`: `formulas.ubc` and one level of areas. `cmd/ubc` hands that `fs.FS` to `cli.App.Scenarios`; the adapter only sees an `fs.FS`, so its tests use `fstest.MapFS`. An `App` without a catalogue behaves as an empty one. A test in `cmd/ubc` fails when a `.ubc` file under `scenarios/` is not embedded, so a deeper directory cannot be forgotten silently.
- **Names.** A scenario is named by its path without the extension, `team/replace-a-person`; `formulas` names the shared file. A name that does not match an embedded file exactly (with `.ubc`, an area alone, `..`) is unknown.
- **`ubc --scenarios`** prints one line per scenario, sorted by name: the name, padded to the longest one, three spaces, and the question. The question is read from the template (ADR 0042): the text after `# Question:` and the comment lines that follow it up to the first empty `#` or non-comment line, joined on one line, its first letter capitalised. A file without that line prints its name alone. The name and the question share a line, so `grep` and `cut -d' ' -f1` work on it; the area is the first part of the name. `formulas` is not listed: it answers no question, and `--help`, the README and the formulas' own header show `--scenario formulas`.
- **`ubc --scenario NAME`** adds the shared formulas and then the scenario as the first sources of the program, before `-e` expressions and files, all in one session: flags such as `--set`, `--data`, `--sens`, `--echo`, `--assert` and `--samples-out` apply exactly as they do to a file, and `-e` can ask a follow-up question about the scenario's variables. `--scenario formulas` adds only the formulas, which is how a reader's own model uses them: `ubc --scenario formulas my-model.ubc`. With `--scenario` there is a program, so stdin is not read and the REPL does not start. An error inside a scenario names it as a file would be named (`ubc: team/replace-a-person: line L, col C: …`). The flag takes one name; the last one given wins.
- **`ubc --scenario-source NAME`** prints the embedded file unchanged to stdout and exits, so a reader can copy it and edit it: `ubc --scenario-source team/replace-a-person > replace.ubc`, then `ubc --scenario formulas replace.ubc`. It accepts `formulas` too. Like `--version` and `--completion`, `--scenarios` and `--scenario-source` print and exit without running a program.
- An unknown name, with `--scenario` or `--scenario-source`, is a usage error (exit 2): `ubc: unknown scenario 'NAME'; ubc --scenarios lists them`.
- **Completion.** The flags come from the real flag set, as before; after `--scenario` and `--scenario-source` the bash, zsh and fish scripts complete the scenario names, `formulas` included, taken from the catalogue when the script is generated.
- **Docs.** Every scenario's "Run it" lines, the template in `scenarios/README.md` and their testscripts use `--scenario`, so the tests check the commands a reader copies. Running the files by path still works from a clone and is mentioned once in `scenarios/README.md`.

## Consequences

- Every installation channel gets the whole catalogue; the binary grows by the size of the `.ubc` files (tens of kilobytes).
- A scenario that calls a formula needs nothing special in its commands; forgetting the formulas is no longer possible with `--scenario`.
- A scenario fixed after a release reaches users only with the next release, like any other behaviour.
- A completion script lists the scenarios of the binary that generated it; a newer binary needs the script regenerated to complete new names.

## Alternatives considered

- **`--scenario NAME --source` or a `--show` flag:** a modifier flag that changes what `--scenario` does is harder to discover and to combine; a separate `--scenario-source` reads as what it does and completes like `--scenario`.
- **`ubc --scenario-copy NAME DIR` writing files:** saves a redirection, but writing files is one more failure mode, and stdout already lets the reader choose the name.
- **Columns for area and name in `--scenarios`:** closer to the README table, but the reader would have to join them to type `--scenario`; the full name is what the flag takes.
- **Listing `formulas` with a description:** it would mix a library into a list of questions.
- **Run the scenario after `-e` and files, as written on the command line:** files are positional and flags are not, so their relative order is lost; running the scenario first is what both `--scenario formulas my-model.ubc` and follow-up `-e` questions need.
- **Download the catalogue on demand:** needs the network and a versioning scheme for a few kilobytes of text.
