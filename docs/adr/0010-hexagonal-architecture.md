# ADR 0010: Hexagonal architecture

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

`ubc` reads programs from several sources (the REPL, `-e`, files and stdin) and writes results in several formats (a text line, `--plain`, `--json`, histograms, and later SVG/PNG files). The core is a language and a Monte Carlo evaluator, and it should stay easy to test exhaustively without any I/O. The other Go projects by the author (project A, project B) use hexagonal architecture, and their habits transfer to this one.

VISION §7 sketched a flat layout with one package per pipeline stage. We prefer an explicit boundary between the pure core and the outside world, enforced by tooling.

## Decision

Three layers, with dependencies pointing inward only:

```
cmd/ubc/                  composition root: flags, wiring, exit codes; the only place that reads os.Args, env or the clock
internal/domain/          pure: language (lexer, parser, AST), values (scalar, sample set, probability),
                          distributions, statistics, evaluator
internal/application/     use cases: a Session that evaluates statements and keeps variables, `last` and the RNG;
                          defines ports for the I/O it needs
internal/adapters/        thin translations to the outside world: CLI/REPL and text/JSON output,
                          and later the history file, model files and chart export
```

Rules:

1. **Inside or outside?** Anything that does I/O, or depends on something that does, is an adapter. A seeded pseudo-random generator is deterministic computation, so it lives in the domain. The time-based seed for `--random` is read in `cmd/ubc`.
2. `domain` imports only the standard library and nothing from `application` or `adapters`. `application` imports `domain`. `adapters` import `application` and `domain`. Nothing under `internal/` imports `os`, except adapters.
3. **Ports are discovered, not designed up front.** A port (Go interface) appears when a use case first needs I/O, and it is defined in `application` next to its consumer. Ports stay small.
4. Sub-packages inside each layer emerge through refactoring when a package grows. The starting point is one package per layer.
5. Wiring is manual in `cmd/ubc/main.go`. There is no DI framework.
6. `make validate` will include a dependency check (`go list` based) that fails when an import breaks rule 2.

Testing per layer: the domain gets unit tests with no doubles. The application gets unit tests with hand-written fakes for its ports. Adapters get tests against real files, buffers and the built binary.

## Consequences

- The language and evaluator can be tested exhaustively in memory and fast.
- New front ends (a TUI, a WASM playground) and new outputs only add adapters.
- There is more structure than a calculator strictly needs. Rules 3 and 4 keep the overhead proportional: no speculative ports or packages.
- VISION §7 is updated to point to this layout.

## Alternatives considered

- **Flat pipeline packages (VISION §7):** simpler, but the boundary between the pure core and I/O would be only a convention.
- **A single package, split later:** fastest to start, but it postpones a decision we already know the answer to.
