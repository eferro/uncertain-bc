# ADR 0018: MIT license

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

The repository is public on GitHub (ADR 0015). Without a license, nobody can legally use, modify or redistribute the code. The projects that inspired `ubc`, Guesstimate and Squiggle, are both MIT. We borrow their ideas, not their code.

## Decision

`ubc` is released under the MIT license (`LICENSE` at the repository root).

## Consequences

- Anyone can use, modify, embed and redistribute `ubc`, including in closed-source products, as long as they keep the copyright notice.
- The project uses the same license as its inspirations, and goreleaser archives include the `LICENSE` file.
- Dependencies must stay compatible with MIT. testify (MIT) and go-internal (BSD-3-Clause) are only used in tests.

## Alternatives considered

- **Apache-2.0:** permissive, with an explicit patent grant, but longer and heavier than a small calculator needs.
- **GPL-3.0 (like GNU bc):** copyleft, but it limits embedding `ubc` in other tools.
