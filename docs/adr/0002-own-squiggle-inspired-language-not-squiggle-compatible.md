# ADR 0002: Implement our own Squiggle-inspired language, not Squiggle itself

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

Squiggle (QURI, MIT licence, TypeScript) is the closest existing language to what we want (`1k to 5k * 12`). It has no Go implementation. On 2026-11-20 we installed `@quri/squiggle-lang` 0.10.0 locally and ran a model through its CLI (see [research/2026-11-21-origin.md](../research/2026-11-21-origin.md)). It works, but the CLI shipped with missing dependencies, only accepts ASCII identifiers and prints distributions as `Sample Set Distribution` without any summary or chart.

Options evaluated to use Squiggle from Go:

| Option | Verdict |
|---|---|
| Embed a bundled squiggle-lang in goja (pure-Go JS engine) | Unverified compatibility, 0.3–1 s startup parsing a multi-MB bundle, tied to upstream releases |
| Run Node as a sidecar process | Breaks the single-binary goal (ADR 0001) |
| Embed V8 via v8go | Needs cgo, heavy and fragile cross-compilation |
| Write our own interpreter in Go | More work, full control, instant startup |

## Decision

Write our own small language in Go, **inspired by** Squiggle but not compatible with it. Keep Squiggle's ideas where they are good (`a to b` as a 90 % interval, number suffixes, sample-set semantics) and diverge freely where a `bc`-like calculator needs something else (REPL commands, `last`, output format, operator precedence).

## Consequences

- No JavaScript, no build-time Node; startup is instant.
- We own a lexer, a parser and an evaluator. The main risk becomes scope creep in the language; the MVP grammar is fixed in [SPEC.md](../SPEC.md).
- Squiggle users may expect identical results; every deliberate difference must be documented.
- Squiggle stays a reference for semantics and naming of functions, not a dependency.
