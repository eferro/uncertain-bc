# ADR 0008: Keep units out of the MVP

- Status: Superseded by [ADR 0065](0065-no-units-in-the-language.md)
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

Unit-aware arithmetic (`€`, `h`, `kWh`, conversions) would catch mistakes in estimates and is attractive. It also touches the lexer, the type system of values and the output format, and multiplies the scope of the first version.

## Decision

No units in the MVP. Numbers are dimensionless; users put units in variable names and comments (`price_eur_kwh`, `# €/kWh`).

## Consequences

- Smaller grammar and evaluator for the first release.
- Revisit after the REPL phase if real usage shows unit errors are frequent.
