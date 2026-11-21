# Agent rules — uncertain-bc

`ubc` is a `bc`-like command-line calculator for uncertain numbers: a small Squiggle-inspired language evaluated with Monte Carlo sample sets, shipped as a single static Go binary.

## Read first

- [docs/VISION.md](docs/VISION.md): why, goals, non-goals, roadmap.
- [docs/SPEC.md](docs/SPEC.md): MVP behaviour (CLI, grammar, semantics, output, REPL). Source of truth for what to build.
- [docs/adr/](docs/adr/README.md): decisions already taken. Do not reopen them silently; propose a new ADR that supersedes one if needed.
- [docs/research/](docs/research/): origin and early experiments.

## How we work

- **Baby steps.** One small, verifiable step at a time; follow the roadmap phases in order, starting with the walking skeleton.
- **TDD.** Start every behaviour with a failing test. Table-driven tests for the lexer, parser and evaluator; golden-output tests for the CLI (the fixed seed makes output deterministic, ADR 0005).
- **Test style** (ADR 0009): stdlib `testing` + testify (`require` for preconditions, `assert` for checks), behaviour names like `TestX_Y_WhenZ`, hand-written fakes, no mocks.
- **Simplicity first.** Standard library unless a dependency clearly pays off. No cgo (ADR 0001).
- **Keep the spec honest.** If behaviour changes, update `docs/SPEC.md` in the same commit; if a decision changes, add an ADR.
- **Small functions, clear names,** no comments that restate the code.
- **English** for code, docs, commits and ADRs. Conversation with the user can be in Spanish.
- **Ask when in doubt,** one question at a time.

## Go conventions

- Module: `github.com/eferro/uncertain-bc`. Binary entry point: `cmd/ubc`. Internals under `internal/` (see VISION §7).
- Before every commit: `gofmt`, `go vet ./...`, `go test ./...` must pass. Add `staticcheck` / `golangci-lint` when the Makefile exists.
- Conventional commits (`feat:`, `fix:`, `docs:`, `test:`, `refactor:`, `chore:`).
- Commit directly to `main`; no branches unless asked.
