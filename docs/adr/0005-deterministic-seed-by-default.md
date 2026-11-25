# ADR 0005: Deterministic RNG seed by default

- Status: Accepted, amended by [ADR 0069](0069-machine-output-rounded-to-ten-significant-figures.md)
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

With a random seed, the same command prints slightly different percentiles each time. That is confusing in a calculator, breaks golden-output tests and makes shared results impossible to reproduce.

## Decision

Seed the RNG once per session with a fixed default (`1`). `--seed N` and `:seed N` change it; `--random` opts into a time-based seed. The REPL keeps one RNG stream for the whole session and never re-samples existing variables.

## Consequences

- Same program, same flags, same output, on every machine.
- Tests can assert exact output.
- Output can still change between releases if the sampling algorithms change; that is acceptable and noted in release notes.
