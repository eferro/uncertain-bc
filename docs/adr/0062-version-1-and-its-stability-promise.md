# ADR 0062: Version 1.0.0 and its stability promise (amends 0015)

- Status: Accepted
- Effective date: 2026-11-23
- Deciders: eduardo.ferro

## Context

ADR 0015 started versioning at `v0.x` "until the MVP spec is complete". Since then `ubc` has grown into an application people can rely on: the full language, the REPL, every installation channel, scenario and assertion flags, a built-in scenario catalogue, and user documentation (tutorial, modelling guide, how-tos, reference, website). The documentation no longer calls it an MVP, and the maintainer decided to release it as `1.0.0`.

Under SemVer, `0.x` promises nothing between minor versions; `1.0.0` defines a public interface that only a major version may break.

## Decision

- The release after `v0.7.2` is `v1.0.0`.
- From `1.0.0` on, the **public interface** is what [REFERENCE.md](../REFERENCE.md) and [SPEC.md](../SPEC.md) describe:
  - the language: syntax, operators, built-in functions and their meaning;
  - command-line flags and REPL commands;
  - machine output: `--json` fields and `--plain`;
  - exit codes;
  - scenario names in the built-in catalogue (`--scenario NAME`).
- Breaking any of them needs a new major version. New functions, flags, fields or scenarios are minor versions; fixes are patch versions.
- **Not** part of the promise: the exact digits of a Monte Carlo result (they may change when sampling improves, and the changelog says so), the layout of the human-readable text output (spacing, sparklines, labels), and the assumptions inside a scenario.
- Every release keeps its section in `CHANGELOG.md` (ADR 0060); a breaking change is listed under a **Breaking** heading.

## Consequences

- Users and scripts can upgrade within `1.x` without surprises in flags, JSON or exit codes.
- Changing a flag or a JSON field now costs a major version or a deprecation period; that pushes design decisions to be settled before they ship.
- Model results can still move slightly between versions, as with any Monte Carlo tool.

## Alternatives considered

- **`v0.8.0`**: keeps freedom to break things, but undersells a tool that is documented and installable everywhere, and leaves users without any compatibility promise.
- **Promise identical numeric output**: would freeze the sampling algorithms forever; the fixed seed already makes each version reproducible.
