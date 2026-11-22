# ADR 0044: Team scenarios, self-contained and pinned by their answers

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0042](0042-scenario-catalogue.md), [ADR 0043](0043-reusable-formulas-for-the-scenarios.md)

## Context

Issue #22 asks for scenarios about people and team cost: hiring, attrition, onboarding, real team capacity, the yearly cost of a recurring meeting, in-house vs contractors vs outsourcing, and the bus factor. The cost of replacing a person already exists (`team/replace-a-person.ubc`, ADR 0042). The catalogue test only checks that a scenario evaluates and is indexed, and ADR 0042 pinned only the reference scenario's output. A manual mutation pass showed what that leaves open: an operator flipped in a model, or a variable renamed so that the `--set` example in its header silently overrides nothing, still passes.

## Decision

- Seven scenarios join the `team` area, one file each, following the catalogue template: `hiring-cost`, `attrition-cost`, `onboarding`, `team-capacity`, `meeting-cost`, `in-house-vs-contractors` and `bus-factor`, each with an index row.
- They reuse the shared formulas where one fits, and then show the two-file command: `annual_loss` in `bus-factor`, `real_capacity` in `team-capacity`, `loaded_cost` in `meeting-cost` and `in-house-vs-contractors`. `hiring-cost`, `onboarding` and `attrition-cost` need none and run alone; `attrition-cost` multiplies its expected leavers by the replacement cost so the number of leavers is computed in one place.
- Every scenario is self-contained: a value that another scenario computes is an assumption with that scenario named as its source (`attrition-cost` takes `replacement_cost = 40k to 100k` from `replace-a-person`'s 90% interval), instead of chaining scenario files, so each one can be read, copied and run on its own.
- Attrition and the bus factor are modelled as expected values, `people × P(leave) × cost`, with the replacement cost as the uncertain average cost of one departure; counting discrete leavers with `binomial` would add little for the reader and hide the formula the issue names.
- Results print the intermediate values a reader wants to see (time to fill, hours, costs per option) and end with one `P(...)` line that answers the question; the comparison of three options answers it with `P(in_house < min(contractors, outsourcing))`.
- Each scenario has its own testscript, `testdata/script/scenarios-team-<name>.txtar`, so one broken model does not hide the others. It runs every command the scenario's header shows and pins, with the fixed seed, the median of every printed result and the whole answer line of the plain run, the answer line of the `--set` run (which proves the override reaches a variable the model uses) and one expected driver in the `--sens` ranking. Intervals and sparklines are not pinned.

## Consequences

- A wrong operator, a renamed variable or a broken header command in a team scenario turns the build red; all the mutants tried in the manual pass were killed.
- Changing a scenario's numbers or its header length means updating its testscript, as with the reference scenario; changes to the RNG or to how medians and probabilities print do too.
- ADR 0042's lighter rule (evaluate and index) still applies to future scenarios that do not add a testscript; new areas are encouraged to follow this one.

## Alternatives considered

- **Only the catalogue test, as ADR 0042 allows:** cheapest, but the mutation pass showed it misses wrong arithmetic and `--set` examples that override nothing.
- **Pin each scenario's whole output with `cmp`:** strongest check, but intervals and sparklines would make every assumption tweak a hand-copied diff, for little more detection than the medians and the answer.
- **Chain scenarios (`replace-a-person` before `attrition-cost`):** reuses the replacement cost exactly, but it would print the other scenario's results and make the files depend on each other's variable names.
- **One testscript for all seven scenarios:** fewer files, but testscript stops at the first failure, so one broken model would hide the rest.
