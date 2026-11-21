# ADR 0013: Strict golangci-lint with complexity limits

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

Simplicity is one of the project's mottos, and the rules ask for small functions with clear names. Rules that only live in a document erode over time, so we want the linter to enforce what it can. The other Go projects by the author use golangci-lint v2: project B has a minimal set that excludes tests, and project A adds gocritic and misspell, lints tests, and sets `max-issues-per-linter: 0`.

## Decision

Use golangci-lint v2 with a version pinned in the Makefile. `.golangci.yml` enables:

| Linter | Why |
|---|---|
| errcheck, govet, staticcheck, unused, ineffassign | Correctness baseline |
| gocritic | Idiomatic and simpler code |
| misspell | English docs and identifiers |
| gocyclo (max 10), gocognit (max 10) | Keep functions easy to follow |
| funlen (20 lines, 15 statements) | Small functions |

- The formatter is **gofumpt**, a stricter superset of gofmt, run through `golangci-lint fmt`. `make format` applies it and `make lint` checks it.
- Test files are linted too. The one exception: funlen is disabled for `_test.go`, because table-driven tests are long by nature (one table literal holds many cases). gocyclo and gocognit still apply to tests, so test logic stays flat.
- `max-issues-per-linter: 0` and `max-same-issues: 0`: every finding is reported, none is hidden.
- `//nolint` is allowed only with the linter's name and a reason (`//nolint:gocognit // reason`). nolintlint enforces this.
- After any config change, run `golangci-lint config verify` (v1 keys are silently ignored by v2).

## Consequences

- Size and complexity limits are checked on every commit, not left to review.
- Some code (the Pratt parser loop, the formatting of numbers) may need extra helper functions to stay under the limits. That pressure is intentional.
- gofumpt can reformat code that plain gofmt accepts. Editors should be set to gofumpt to avoid noisy diffs.

## Alternatives considered

- **The project A set without size limits:** relies on review to keep functions small.
- **The minimal project B set:** least noise, but too weak for the stated priorities, and tests go unlinted.
