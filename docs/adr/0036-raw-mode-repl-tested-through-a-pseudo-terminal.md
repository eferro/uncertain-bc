# ADR 0036: Raw-mode REPL tested through testscript's pseudo-terminal

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0014](0014-cli-acceptance-tests-with-testscript.md), [ADR 0023](0023-line-editing-with-x-term.md)

## Context

Issue #11: the raw-mode REPL (`terminalRepl`, ADR 0023) only ran when stdin was a real terminal, so no test reached it; it had only been checked by hand with `script`. Testing it needs a pseudo-terminal, and ADR 0001 rules out cgo.

## Decision

- Drive the REPL from a testscript scenario (`cmd/ubc/testdata/script/repl_terminal.txtar`) with testscript's built-in `ttyin -stdin`. It opens a pseudo-terminal in pure Go and attaches it to the standard input of `ubc`, so `ubc` takes the raw-mode path. There is no new dependency, no cgo and no test-only seam in production code.
- The scenario checks what a user would see: the input echoed by the editor (proof that raw mode and `x/term` are active), a line edited with Ctrl-A, a history entry recalled with Ctrl-P from a seeded history file, and history saved across sessions. It sets `XDG_STATE_HOME` to the work directory so it never touches the real history.
- Keystrokes are written as Go escapes (`\x01`, `\n`) and turned into raw bytes by a small `unescape` script command in `cmd/ubc/main_test.go`, so key files stay readable and diffable.
- Every line ends in `\n`, not `\r`, and sessions end with `:quit`, not Ctrl-D. testscript writes all the keys as soon as `ubc` starts, before it switches to raw mode, and the terminal is still in canonical mode then: it turns `\r` into `\n` and swallows Ctrl-D as end of file, which made a Ctrl-D scenario hang. Keys that canonical mode treats specially (backspace, Ctrl-U, Ctrl-W, Ctrl-D) are avoided for the same reason.
- `ttyin` only works on Linux and macOS, so the scenario is skipped elsewhere (the Windows CI jobs).

## Consequences

- The raw-mode setup, the editor wiring and history persistence are covered by CI on Linux and macOS, and count towards coverage.
- Ctrl-D and restoring the terminal state on exit stay untested: the pseudo-terminal is discarded when `ubc` exits, so no script can observe whether it was restored. The surviving mutant is listed in `docs/mutation-testing.md`.

## Alternatives considered

- **`github.com/creack/pty` in a Go test:** pure Go too, but a new dependency for what testscript already provides.
- **Opening `/dev/ptmx` by hand with `golang.org/x/sys/unix`:** no new dependency, but platform code to maintain in tests.
- **Raw control bytes in the `.txtar` files:** no extra command, but invisible in editors and diffs.
