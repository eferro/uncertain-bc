# ADR 0023: Line editing with golang.org/x/term

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

SPEC §7 asks for line editing and persistent history in the REPL (`$XDG_STATE_HOME/ubc/history`). Reading lines with `bufio` gives neither: arrow keys print escape codes and nothing is remembered between sessions. The standard library has no terminal line editor, so a dependency is needed (ADR 0001 rules out cgo).

## Decision

- Use `golang.org/x/term`: its `Terminal` provides line editing (cursor movement, word deletion, in-memory history navigation) on any `io.ReadWriter`, and its pluggable `History` interface lets us persist history ourselves. `term.MakeRaw` puts the TTY in raw mode.
- The REPL reads lines through a small interface. On a real terminal it uses an `x/term` editor; elsewhere (tests, pipes) it keeps the `bufio` reader.
- History lives in `$XDG_STATE_HOME/ubc/history`, falling back to `~/.local/state/ubc/history`. `cmd/ubc` resolves the path, because only `cmd/ubc` reads the environment (ADR 0010). Each entered line is appended to the file. A history file that cannot be written never breaks the session.

## Consequences

- One small dependency, maintained by the Go team, pure Go, with no transitive dependencies beyond `golang.org/x/sys`.
- `x/term` has no tab completion or reverse search. Those would need a bigger library and are not in the MVP.
- Raw-mode setup cannot be unit-tested without a pseudo-terminal. It stays a few lines at the edge, and the editor and history are tested through in-memory readers and temporary files.

## Alternatives considered

- **chzyer/readline:** feature-rich (completion, reverse search), but unmaintained since 2022.
- **ergochat/readline / peterh/liner:** maintained, but bigger third-party code for features we do not need yet.
- **No line editing:** simplest, but it fails SPEC §7 and makes the REPL painful to use.
