# ADR 0027: Track work in GitHub issues

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

After phase 4 and the installation channels, the pending work lived only in conversation hand-off notes: feature ideas (`--set`, `--assert`, `mixture`, a tutorial, …), checks that were never automated, and a token that must be renewed. Nothing in the repository said what was left or in which order, and that list was easy to lose between sessions.

## Decision

- Pending work is tracked as issues in the public GitHub repository, one issue per feature, check or chore.
- Labels: `enhancement` for features, `documentation`, `testing` for verification that is missing or only manual, `maintenance` for housekeeping, and `question` for ideas that need a decision before any work.
- The issue states the goal and the acceptance criteria; design choices still go to an ADR and behaviour to `docs/SPEC.md`, as before. The issue links to them, not the other way round.
- The commit that finishes an issue closes it with `Closes #N` in its message body.
- Issue numbers reflect the order agreed at creation; the next item to work on is the lowest open `enhancement` or `documentation` issue unless the user says otherwise.

## Consequences

- What is left, and why, is visible to anyone who opens the repository, and survives the end of a working session.
- Issues are public: they must not mention employer projects or secrets, the same rule as commits and ADRs.
- One more place to keep honest: an issue that no longer applies is closed with a comment explaining why.

## Alternatives considered

- **A `docs/TODO.md` file:** versioned with the code, but no way to link commits to items or close them automatically.
- **GitHub Projects board:** more structure than a one-person project needs; it can be added later on top of the same issues.
