# ADR 0037: REPL on Windows consoles: escape processing on, checked in CI by reading the code

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0023](0023-line-editing-with-x-term.md), [ADR 0036](0036-raw-mode-repl-tested-through-a-pseudo-terminal.md)
- Amended by: [ADR 0041](0041-whole-test-suite-on-windows.md) (the Windows job now runs the whole suite)

## Context

Issue #12: the raw-mode REPL (line editing and history with `golang.org/x/term`, ADR 0023) had never run on Windows. Nobody on the project has a Windows machine at hand, so the check was done by reading `x/term`, the Go runtime and our code for Windows behaviour:

- `term.MakeRaw` on Windows sets the console **input** mode only: it turns off echo, line input and Ctrl-C processing and turns on `ENABLE_VIRTUAL_TERMINAL_INPUT`, so arrow keys arrive as the same escape sequences as on Unix and `term.Terminal` decodes them. Enter (`\r`), Backspace (DEL or `^H`), Ctrl-A/E/P/N, Ctrl-C and Ctrl-D behave as on Linux.
- Go reads console input with `ReadConsoleW` and treats Ctrl-Z at the start of a read as end of file, so Ctrl-Z ends the session, as Windows users expect (in raw mode every key is read at once, no Enter needed).
- The history file and its directory are created with ordinary file calls; the default path is `%USERPROFILE%\.local\state\ubc\history` (SPEC §7). It works, although it is not where Windows programs usually keep state.
- **The bug:** `term.Terminal` redraws the line with ANSI escape sequences (cursor left/right, erase), but nothing turned on `ENABLE_VIRTUAL_TERMINAL_PROCESSING` on the **output** console. A console that starts without it (the classic console host behind Windows PowerShell 5.1 and `cmd.exe`) prints them as text, such as `←[D`, as soon as a line is edited. Windows Terminal and PowerShell 7 usually leave the mode on, which hides the bug there.

## Decision

- While the raw-mode REPL runs, `ubc` turns on `ENABLE_VIRTUAL_TERMINAL_PROCESSING` for stdout when it is a Windows console, and restores the previous mode on exit (`internal/adapters/cli/console_windows.go`). If stdout is not a console (redirected) or the mode cannot be set, nothing changes. On other systems the function does nothing (`console_other.go`): terminals there always interpret escape sequences.
- It uses `golang.org/x/sys/windows`, already in the module through `x/term`; it is now a direct dependency. No cgo (ADR 0001).
- A new CI job on `windows-latest` runs the REPL tests of `internal/adapters/cli` (`TestRepl_*`, `TestHistory_*`, `TestLineEditor_*` and the console tests) through the Makefile. The rest of the package is not run there yet: some of its tests compare Unix error texts (`no such file or directory`) or run bash, and fail on Windows for reasons unrelated to the REPL. The Windows-only unit tests open the real console (`CONOUT$`, allocating one if the test process has none) and check that the mode gains exactly the escape-processing flag and gets its old value back on restore, and that a file that is not a console is left alone.
- The rest of the REPL on Windows (keys, Ctrl-Z, history) is accepted from the reading above and is not run on a real console: no CI runner offers an interactive console that a test could type into, and testscript's `ttyin` is Linux and macOS only (ADR 0036).

## Consequences

- Line editing on Windows PowerShell 5.1 and `cmd.exe` should no longer print raw escape codes.
- The REPL on Windows is still not exercised end to end: whether `terminalRepl` calls `enableVirtualTerminal` at all is only visible on a real Windows console, so that surviving mutant is listed in `docs/mutation-testing.md`. A report from a Windows user may still find problems this reading missed.
- Only the REPL tests run on Windows; making the other tests portable and widening the job to the whole module is separate work.

## Alternatives considered

- **Document the gap only:** cheaper, but leaves a likely visible bug in the most common Windows shell.
- **Wait for a manual check on Windows before closing the issue:** more certain, but blocks on hardware the project does not have; the user chose to fix and record instead.
- **Turn on escape processing once at start-up for every mode:** `ubc` only prints escape sequences in the line editor, so a wider change would alter the console for nothing.
- **Move the default history to `%LOCALAPPDATA%` on Windows:** more idiomatic, but the current path works and changing it is a product decision outside this issue.
