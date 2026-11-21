# ADR 0011: Makefile as the single entry point, enforced by git hooks

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

Every commit should leave the project a bit better than before, and it should never leave it broken. Humans and agents both work on the repo, so the way to build, test and check it must be one obvious command, identical locally and in CI. The other Go projects by the author (project A, project B) use a Makefile with canonical target names and the rule "never call the tools directly". Project A also versions git hooks so the gates run without relying on memory.

## Decision

- The `Makefile` is the only documented way to run tools. Do not call `go test`, `gofmt`, `go vet` or `golangci-lint` directly. If a new operation is needed, add a target.
- Canonical targets, each with a `##` description that the default `make help` target prints:
  - `help`: list targets.
  - `local-setup`: install pinned tools and git hooks.
  - `test`: fast unit tests.
  - `test-pkg PKG=...`: inner TDD loop on one package.
  - `format`: rewrite with `gofmt`.
  - `lint`: format check, `go vet`, `golangci-lint`.
  - `validate`: everything CI runs (lint, tests, project checks such as test style and the architecture rule).
  - `build`: static binary (`CGO_ENABLED=0`).
- Git hooks are versioned in `scripts/hooks/` and enabled with `git config core.hooksPath scripts/hooks` by `make local-setup`:
  - `pre-commit` runs `make validate`.
  - `commit-msg` rejects messages that do not follow Conventional Commits: `^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([a-z0-9 /-]+\))?!?: .+`.
- CI runs `make validate` and checks every commit message again, so `--no-verify` cannot sneak anything in.
- Never weaken or remove a gate to make a commit or a pipeline pass. Fix the cause.
- Tool versions are pinned, never `latest`.

## Consequences

- One command (`make validate`) answers "is it OK to commit?" for humans, agents and CI.
- Hooks make commits a bit slower. `validate` must stay fast; slow suites, if any appear, get their own target outside `validate`.
- Contributors must run `make local-setup` once after cloning.

## Alternatives considered

- **Makefile only, without hooks:** less machinery, but it relies on discipline and only CI catches mistakes.
- **Plain go tooling without a Makefile:** fewer files, but no single `validate` command, and local and CI checks drift apart.
