# ADR 0048: Reliability scenarios, with availability described by its downtime share

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0042](0042-scenario-catalogue.md), [ADR 0043](0043-reusable-formulas-for-the-scenarios.md), [ADR 0044](0044-team-scenarios.md)

## Context

Issue #26 asks for scenarios about reliability and operations: the yearly cost of incidents (frequency × MTTR × cost per hour of downtime), the composite availability of services in series and in parallel, an SLO's error budget, capacity planning from peak traffic and headroom, and the cost of a bug depending on where it is caught. None of the existing areas fits them. Two of them deal with availabilities a hair below 100%, where a `a to b` range of availabilities (99.9% to 99.99%) is a lognormal whose tail goes above 1, and a summary with three significant figures prints `0.998` for very different services.

## Decision

- Five scenarios open a new `reliability` area, one file each, following the catalogue template and the testing rule of ADR 0044 (one testscript per scenario, `testdata/script/scenarios-reliability-<name>.txtar`, pinning medians, the answer line, the `--set` answer and one `--sens` driver): `incident-cost`, `availability`, `error-budget`, `capacity` and `bug-cost`.
- Availability is described by the share of the time each piece is down (`0.05% to 0.3%`) and turned into an availability with `1 - down`, so no sample exceeds 100%. `availability` prints the number of nines (`-log10(1 - availability)`) and the minutes down a month instead of the availability itself, which three significant figures would flatten. It uses the shared `availability_series` and `availability_parallel` formulas and says in its header that they assume independent failures.
- `incident-cost` counts the hours of downtime and the team's time (responders during the incident, follow-up afterwards) and multiplies by the incidents a year with the shared `annual_loss` formula, its frequency playing the annual rate of occurrence; the expected count is used rather than `poisson`, as with attrition in ADR 0044.
- `error-budget` follows the Google SRE book: the budget is `1 - SLO` of a 30-day window, and an incident spends only the share of requests it fails, as a request-based SLO counts it; planned maintenance counts against it.
- `capacity` rounds up to whole nodes with `ceil`, plans for a target utilisation below 100% and adds one spare node (N+1), and prints `p(nodes, 0.95)` as the size to provision as well as the chance that today's cluster falls short.
- `bug-cost` uses hours and a customer impact the reader can estimate rather than the 1:10:100 multipliers, whose evidence is weak; its header says so.
- `incident-cost` and `bug-cost` price hours with the shared `loaded_cost`; `error-budget` and `capacity` need no shared formula and run alone.
- Assumptions are tuned so that each answer sits away from 0 and 1 with the defaults, and the `--set` example in each header moves it clearly, which is what the testscript pins.

## Consequences

- Readers get reliability models in the units operations people use (minutes of downtime, nines, nodes) and can see which piece or assumption dominates with `--sens`: with the defaults, the single database decides the availability.
- Changing a reliability scenario's numbers or header means updating its testscript, as with the other areas.
- A manual mutation pass over the five models (arithmetic operators, removed terms and factors, `ceil` to `floor`, the spare node, the quantile, the answer's comparison, the shared formulas swapped or bypassed) killed 40 of 44 mutants. The survivors are the boundary flips `>` to `>=` in `P(yearly_cost > 500k)`, `P(availability > 99.9%)`, `P(expected_minutes > allowed_minutes)` and `P(production_cost > 10 * qa_cost)`, which are equivalent: a continuous sample equals the threshold with probability zero. The same flip in `capacity`, where nodes are whole numbers, is killed.

## Alternatives considered

- **Ranges of availabilities (`99.9% to 99.99%`):** reads naturally, but some samples exceed 100% and the composite comes out wrong in the tail.
- **Print the availability as a percentage:** `99.8` still hides the difference between 99.75% and 99.84%, which is a factor in downtime; nines and minutes show it.
- **Draw the incident count with `poisson`:** shows the discrete count, but the expected value is the formula the issue names and reads more simply.
- **Fixed 1:10:100 multipliers for the cost of a bug:** the usual rule of thumb, but presenting it as data would contradict the catalogue's aim of honest ranges.
