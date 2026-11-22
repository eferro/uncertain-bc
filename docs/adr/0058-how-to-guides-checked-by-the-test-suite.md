# ADR 0058: How-to guides checked by the test suite, with data files and `cat`

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0030](0030-tutorial-checked-by-the-test-suite.md), [ADR 0057](0057-user-reference-checked-by-the-test-suite.md)

## Context

Issue #36 asks for task-oriented guides in one file, `docs/HOWTO.md`, one short section per task, each with a real command checked by the test suite: using your own numbers (`--set`, `--data`), comparing two options, finding the assumption to research first (`--sens`), a check in CI (`--assert`, exit 4), taking the samples to pandas or a spreadsheet (`--samples-out`), forecasting a delivery date from real throughput, and running and adapting a catalogue scenario. Two of them could not be checked by the ADR 0030 converter: `--data` needs a data file, and only `.ubc` files reached the test's working directory (ADR 0057 left `--data` in prose for that reason); `--samples-out` prints nothing, and what matters is the file it writes.

## Decision

- The guides are `docs/HOWTO.md`, one section per task of the issue, in that order, with a list of the tasks at the top. They are checked by the guide converter like the other guides (they join the list of ADR 0056) and linked from the README, the tutorial, the reference and the catalogue index; the link test covers the four links.
- **Data files live in `examples/`.** The test setup now copies `.txt` files from `examples/` and `scenarios/` as well as `.ubc` files, so a guide can run `--data NAME=examples/FILE.txt` exactly as a reader types it. `examples/throughput.txt` holds twelve weeks of items done, used by both the `--data` recipe and the delivery forecast with `delivery/delivery-date`.
- **`cat FILE` in console blocks.** A console command that is exactly `cat` and one file shows that file's content: the lines after it are compared with the file, byte for byte (testscript's `cmp`), instead of running a process. It shows the CSV that `--samples-out` wrote and the data file before it is loaded. `cat` with no content shown is an error, as is `cat` with more than one file; every other command must still start with `ubc `.
- The samples recipe uses `-n 4`, so the CSV shown is short; the pandas snippet and the GitHub Actions step are in unchecked `python` and `yaml` blocks, since they need tools the test suite does not have.
- Adapting a scenario uses `--echo` to show the names of its assumptions and `--set` to change them; copying it with `--scenario-source NAME > mine.ubc` stays in prose, because the converter has no redirection and the edit itself cannot be shown as a command.

## Consequences

- Every command and its output in the how-to guides, the samples file and the data file included, are compared with the binary, so a change in the RNG, the output or a scenario breaks the guides until they are updated.
- The tutorial, the modelling guide and the reference can now show `--data` and the files `ubc` writes too.
- The prose about the numbers (for example "building is cheaper half of the time") is not checked, as in the other guides.

## Alternatives considered

- **One guide per task in `docs/howto/`:** easier to link one by one, but the issue asks for one file, and seven short sections read well together.
- **Show the samples in an unchecked block:** simpler, but the file shown could drift from what `ubc` writes, the drift ADR 0030 rejected.
- **`--samples-out /dev/stdout`:** shows the samples without a new command, but it is not a path on Windows and the extension rule of ADR 0032 rejects it.
- **Support shell redirection (`> file`) in the converter:** would check `--scenario-source` into a file, but it needs a second translation for a step whose interesting part, the edit, still could not be shown.
