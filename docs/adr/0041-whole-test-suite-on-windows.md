# ADR 0041: The whole test suite on Windows: portable tests and LF checkouts

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0037](0037-repl-on-windows-consoles.md) (the Windows CI job ran only the REPL tests)

## Context

Issue #19: the Windows CI job ran only the REPL, history, line editor and console tests (ADR 0037), because the rest of the suite failed there. Nobody on the project has a Windows machine, so the Windows job is the only place to check. Reading the tests showed problems in the tests, not in `ubc`:

- Tests for `--data`, `--samples-out`, `:save`, `:load` and `:plot` compared the whole error message, including the operating system's reason (`no such file or directory`). Windows words it differently (`The system cannot find the file specified.`); `ubc` passes the reason through unchanged on every system.
- The bash completion tests gave bash absolute temporary paths. On Windows these are `C:\…` paths, which the bash shipped with Git for Windows reads as escapes.
- The tests of `install.sh` ran it under Git for Windows' `sh`, where it correctly refuses to run and points to `install.ps1` (ADR 0026).
- Git for Windows checks text files out with CRLF line ends by default. The testscript archives, the golden outputs and the tutorial are compared byte for byte with what `ubc` prints, and `ubc` prints LF on every system.

## Decision

- Unit tests build the expected file error from the operating system itself: a helper opens the same missing path, checks that the error is `fs.ErrNotExist` and returns its text. The tests still check every byte `ubc` adds around it, and they check the kind of error instead of its wording.
- Testscript scenarios check the part of a file error that `ubc` writes and accept any reason after it (`open missing.txt: .+`).
- The bash completion tests run bash in the temporary directory with relative file names; they still skip when bash is missing.
- The `install.sh` tests skip on Windows, as they already do when `sh`, `curl` or `tar` is missing.
- A `.gitattributes` file checks out every text file with LF line ends on every system. It changes nothing on Linux or macOS.
- The Windows CI job runs `make test`, the whole module, through the Makefile (ADR 0011). Scenarios that need a pseudo-terminal keep skipping there (ADR 0036).

## Consequences

- Every test except the pseudo-terminal scenario and the `install.sh` tests runs on Windows on every push, so a change that breaks `ubc` on Windows shows up in CI.
- SPEC §5 states that a file error ends with the operating system's reason, which differs between systems.
- Contributors on Windows get LF files whatever their Git settings; editors there handle LF.
- The coverage, lint and mutation gates still run on Linux only.

## Alternatives considered

- **Skip the failing tests on Windows:** cheaper, but leaves file errors and completion untested there.
- **Match only the exit code of file errors:** portable, but loses the check on how `ubc` frames the message.
- **Translate the operating system's reasons into `ubc`'s own words:** one message everywhere, but a product change that hides useful detail and needs a mapping for every error; not needed to make the tests portable.
- **Turn off `core.autocrlf` in the CI job instead of `.gitattributes`:** fixes CI only, and Windows contributors would still get failing tests.
