# ADR 0064: No TUI with sliders

- Status: Accepted
- Effective date: 2026-11-23
- Deciders: eduardo.ferro

## Context

Roadmap phase 5 in [VISION.md](../VISION.md) listed, as a "maybe", Guesstimate-style interactive sliders in a terminal UI: the inputs of a model as sliders and the outputs updating live. Issue #16 asked whether to do it.

It would need a TUI library (a new dependency), a second interaction model to design, test and document, and terminal handling on three operating systems, for a tool whose strength is being a small `bc`-like calculator. Meanwhile `--watch` with an editor already gives a live loop, `--set` runs scenarios from the command line, and `--sens` shows which input matters most.

## Decision

`ubc` will not have a TUI with sliders. Issue #16 is closed as not planned. The roadmap drops it from phase 5; units (ADR 0008, issue #10) remain the only open idea there.

## Consequences

- No new dependency and no second user interface to maintain.
- Exploring a model stays text-based: edit and `--watch`, or `--set` and `--sens`.
- The idea can come back with a new ADR if real users ask for it.

## Alternatives considered

- **Build it now**: too much complexity for the value it adds today (maintainer's decision).
- **Keep it as a "maybe"**: leaves an open question nobody plans to answer; closing it makes the roadmap honest.
