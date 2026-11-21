# ADR 0007: Name the binary `ubc` and the project `uncertain-bc`

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

Candidates considered: `circa`, `napkin`, `ballpark`, `ubc`, `sqc`. The short name must be quick to type like `bc`; the project name must be findable.

Checks done on 2026-11-21:

- `ubc`: no package in Homebrew, apt or crates.io. npm `ubc` is an abandoned placeholder (last modified 2022); PyPI `ubc` is a "Utility Bill Calculator". Search results for "ubc" are dominated by the University of British Columbia.
- "uncertain bc": zero GitHub repositories; "uncertainty calculator cli": zero.
- `napkin`: no Homebrew formula, but a crowded GitHub namespace (sirupsen/napkin-math 5.9k stars, vEnhance/napkin 2k, several agent tools named napkin).

## Decision

Binary `ubc`; repository and project name `uncertain-bc` ("Uncertain bc"); Go module `github.com/eferro/uncertain-bc`; entry point `cmd/ubc`. Same pattern as ripgrep → `rg` and fd-find → `fd`.

## Consequences

- `go install github.com/eferro/uncertain-bc/cmd/ubc@latest` installs a binary called `ubc`.
- Discoverability relies on the long name in titles, descriptions and topics; never market the tool as just "ubc".
