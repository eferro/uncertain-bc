# Agent rules — uncertain-bc

`ubc` is a `bc`-like command-line calculator for uncertain numbers: a small Squiggle-inspired language evaluated with Monte Carlo sample sets, shipped as a single static Go binary.

## Read first

- [docs/VISION.md](docs/VISION.md): why, goals, non-goals, roadmap.
- [docs/SPEC.md](docs/SPEC.md): MVP behaviour (CLI, grammar, semantics, output, REPL). Source of truth for what to build.
- [docs/adr/](docs/adr/README.md): decisions already taken. Do not reopen them silently; propose a new ADR that supersedes one if needed.
- [docs/research/](docs/research/): origin and early experiments.

## How we work

- **Baby steps.** One small, verifiable step at a time; follow the roadmap phases in order, starting with the walking skeleton.
- **TDD.** Start every behaviour with a failing test. Table-driven tests for the lexer, parser and evaluator; CLI acceptance scenarios as testscript `.txtar` files in `cmd/ubc/testdata/script/` (ADR 0014; the fixed seed makes output deterministic, ADR 0005). Outer loop: a failing scenario; inner loop: unit tests.
- **Definition of done per increment** (ADR 0021): green → manual mutation of the changed code with the `mutation-testing` skill (kill or document survivors, revert mutants) → Test Desiderata review → `make validate` → commit.
- **Test style** (ADR 0009): stdlib `testing` + testify (`require` for preconditions, `assert` for checks), behaviour names like `TestX_Y_WhenZ`, hand-written fakes, no mocks.
- **Quality gates** (ADR 0012): `make validate` requires 90% coverage of `internal/`; `make test-mutation` (gremlins, CI) requires 85% efficacy on `internal/domain`; equivalent mutants are listed in `docs/mutation-testing.md` (ADR 0020). Thresholds only go up.
- **Simplicity first.** Standard library unless a dependency clearly pays off. No cgo (ADR 0001).
- **Keep the spec honest.** If behaviour changes, update `docs/SPEC.md` in the same commit; if a decision changes, add an ADR.
- **Small functions, clear names.** Comments only explain *why* (package docs, non-obvious math, SPEC/ADR references); never *what*; none in tests (ADR 0017). Enforced by golangci-lint (ADR 0013): funlen 20 lines, gocyclo/gocognit 10, gofumpt; `//nolint` only with linter name and reason.
- **English** for code, docs, commits and ADRs. Conversation with the user can be in Spanish.
- **Ask when in doubt,** one question at a time.

## Go conventions

- Module: `github.com/eferro/uncertain-bc`. Binary entry point: `cmd/ubc`. Hexagonal layout under `internal/` (`domain`, `application`, `adapters`), dependencies point inward, ports only when a use case needs I/O (ADR 0010).
- Always go through the Makefile, never call `go test`, `gofmt`, `go vet` or `golangci-lint` directly (ADR 0011). `make help` lists targets; `make test-pkg PKG=...` is the TDD loop.
- `make validate` must pass before every commit; the versioned `pre-commit` hook runs it and `commit-msg` checks Conventional Commits; `pre-push` runs gitleaks on the pushed commits (ADR 0019). Never commit tokens or secrets. Never weaken a gate to get green.
- Conventional commits (`feat:`, `fix:`, `docs:`, `test:`, `refactor:`, `chore:`).
- Commit directly to `main`; no branches unless asked. CI is GitHub Actions; releases are SemVer tags built by goreleaser (ADR 0015). A red pipeline is fixed first.
