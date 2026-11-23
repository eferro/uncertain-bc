# ADR 0065: No units in the language; conversion constants in the shared formulas

- Status: Accepted
- Effective date: 2026-11-23
- Deciders: eduardo.ferro
- Supersedes: [ADR 0008](0008-units-out-of-the-mvp.md)

## Context

ADR 0008 kept units (`€`, `h`, `kWh`) out of the MVP and promised to revisit them after the REPL phase if real use showed unit errors to be frequent. The roadmap kept them as phase 5, "Maybe" (issue #10). Since then the language, the tutorial, three guides and the 43 scenarios of the catalogue have been written without a single unit error being reported or found. What did show up is a weaker form of the problem: conversions written as magic numbers (`/ 60`, `* 12`, `30 * 24`, `43200`), and the same assumption (220 working days a year, 44 working weeks, 8 hours a day) repeated in several scenarios under different names, with `weeks_per_year` meaning 44 working weeks in one scenario where a reader would expect 52.

Dimensional analysis would have to reach almost everything: the lexer (a unit `M` clashes with the `M` suffix, `h` and `d` are common identifiers, `%` is already a suffix), every value (a sample set would carry a dimension), every built-in function (what is the unit of `log`, `P`, `p`?), and every output format (`--json`, `--plain`, `--samples-out`, the summaries). It is the kind of scope the vision warns about (risk 1).

## Decision

- **No units in the language**, as a lasting decision, not an MVP cut. Numbers are dimensionless; the unit of a quantity lives in its name and its comment (`delay_weeks`, `price_eur_kwh`, `# € per kWh`). This can be reopened with evidence of real unit errors that naming and the constants below do not prevent.
- **Conversion constants in `scenarios/formulas.ubc`**, in a "Conversions" section before the formulas: `minutes_per_hour = 60`, `hours_per_day = 8`, `hours_per_calendar_day = 24`, `days_per_week = 5`, `days_per_month = 30`, `working_days_per_year = 220`, `working_weeks_per_year = working_days_per_year / days_per_week` (44), `weeks_per_year = 52`, `months_per_year = 12`.
  - Each says whether it counts **working time** (what people put in: 8 hours, 5 days, 220 days, 44 weeks) or **calendar time** (what a clock, a bill or an SLO window counts: 24 hours, 30 days, 52 weeks). The unqualified names follow the usual meaning in estimates (`hours_per_day` is a working day; `weeks_per_year` is the calendar year); the other meaning carries the qualifier (`hours_per_calendar_day`, `working_weeks_per_year`).
  - `working_weeks_per_year` is computed from the other two, so a `--set` of `working_days_per_year` or `days_per_week` moves it too.
  - They are ordinary assignments: `--set` overrides them like any variable (`--set 'hours_per_day = 7.5'`), `--echo` prints them first because `--scenario` loads the formulas like a file (§1), and `--sens` leaves them out because they are plain numbers.
- **The scenarios use them** wherever a magic number did a conversion, and lose their local copies (`hours_per_day = 8`, `days_per_year = 220`, `weeks_per_year = 44`, `hours_per_week = 40`, `days = 30`, `minutes_per_month = 43200` becomes `days_per_month * hours_per_calendar_day * minutes_per_hour`). Assumptions specific to one scenario stay local (`working_days = 21` a month in ephemeral environments, `working_days = 63` a quarter in team capacity, `meetings_per_year = 46`). Every scenario's output is unchanged.
- **The modelling guide** gains a section on units: name them, keep one unit per quantity, convert with the constants; its examples are checked by the test suite like the rest of the guide (ADR 0056).

## Consequences

- The grammar, the values and every output format stay as they are; the stability promise of 1.0 (ADR 0062) is untouched.
- A conversion now reads as what it is (`/ minutes_per_hour`), and changing an organisation's working year is one `--set`.
- `ubc --echo --scenario NAME` and `ubc --echo --scenario formulas model.ubc` start with nine lines of constants.
- A model that loads the formulas and assigns one of these names (for example its own `hours_per_day = 7`) simply replaces the shared value for the rest of the run, as any reassignment does.
- No tool catches `hours + minutes`; the guide shows the mistake and the fix.

## Alternatives considered

- **Dimensional analysis in the language** (as in Frink or GNU units): catches the error, but touches the lexer, every value, function and output format, for an error not seen in practice.
- **Units as comments checked by a linter:** a second language to specify and maintain inside comments.
- **Constants as zero-argument functions** (`hours_per_day()`): would hide them from `--echo`, but cannot be overridden with `--set` and read worse in a formula.
- **Hide the shared constants from `--echo`:** cleaner output, but a special case in the CLI for one file; the maintainer chose to echo them, which also shows the conversion assumptions a run used.
