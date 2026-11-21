# ADR 0017: Comments explain why, not what

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

Comments that restate code go stale and add noise. Names and tests should carry the intent. But a numerical tool has facts that code alone cannot express, like where a constant comes from, why an algorithm was chosen, or which SPEC rule a branch implements. Project B bans comments entirely, while project A keeps short "why" comments.

## Decision

- No comment may say *what* the code does. Rename, extract a function or write a test instead.
- Allowed comments:
  - Package doc comments (`// Package dist ...`).
  - Doc comments on exported identifiers, only when the name and signature are not enough.
  - Short *why* notes: the origin of non-obvious constants and math, a rejected obvious alternative, or a reference to the SPEC section or ADR a rule comes from.
- Test files have no comments. Test names and scenario file names describe the behaviour.
- No commented-out code and no `TODO` without context. Open work goes into the roadmap or an issue.

## Consequences

- Names and tests do most of the documenting.
- The few comments that remain are worth reading, and they record knowledge that would otherwise be lost (for example, why the 90 % interval uses `1.6449`).
- Whether a comment explains "why" or "what" is decided in review. No linter checks it.

## Alternatives considered

- **Zero comments (project B):** maximally strict, but it loses the reasoning behind numerical constants and choices.
