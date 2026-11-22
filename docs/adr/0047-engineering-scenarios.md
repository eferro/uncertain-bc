# ADR 0047: Engineering scenarios for technical debt and productivity, counted in recoverable hours

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0042](0042-scenario-catalogue.md), [ADR 0043](0043-reusable-formulas-for-the-scenarios.md), [ADR 0044](0044-team-scenarios.md)

## Context

Issue #25 asks for scenarios about technical debt and engineering productivity: the cost of not fixing a piece of debt against the cost of fixing it, whether a task is worth automating (the xkcd "Is It Worth the Time?" table, with ranges), a slow CI pipeline, flaky tests, and a migration's cost and probability of exceeding its budget. None of the existing areas fits them: they are neither team costs nor delivery forecasts nor product strategy. All of them turn engineers' time into a cost, and the issue's formulas (hours lost × people × cost) count every lost hour as fully lost and every fix as fully effective, which overstates both sides.

## Decision

- Five scenarios open a new `engineering` area, one file each, following the catalogue template and the testing rule of ADR 0044 (one testscript per scenario, `testdata/script/scenarios-engineering-<name>.txtar`, pinning medians, the answer line, the `--set` answer and one `--sens` driver): `debt-cost`, `automation`, `slow-ci`, `flaky-tests` and `migration`. The area is named `engineering` rather than `debt` because CI, flaky tests and automation are productivity, not debt.
- Only the hours that are really lost or really recovered count, each through its own assumption: `debt-cost` credits the fix with a `recovered` share of the lost hours, and `slow-ci` counts an `idle_share` of the wait, because people switch to other work while the pipeline runs. `flaky-tests` counts the wait for a rerun whole, because a red build blocks the merge; its header says so.
- `debt-cost`, `slow-ci` and `migration` price hours with the shared `loaded_cost` formula and show the two-file command. `automation` and `flaky-tests` answer in hours and run alone: in `automation` the same people's time is on both sides, so the cost of an hour cancels out.
- `automation` keeps the xkcd table's five-year horizon and adds what the table leaves out: the time spent keeping the automation working every year.
- `migration` gives the unlisted work its own `forgotten` share, as the tutorial's project estimate does, and adds the months of running the old and the new systems side by side; the budget is a plain number the reader replaces.
- Assumptions are tuned so that each answer sits away from 0 and 1 with the defaults, and the `--set` example in each header moves it clearly, which is what the testscript pins.

## Consequences

- Readers get debt and productivity models that do not overstate the case for fixing or automating, and can see which share they would argue about.
- Changing an engineering scenario's numbers or header means updating its testscript, as with the other areas.
- A manual mutation pass over the five models (arithmetic operators, removed factors, the answer's comparison, renamed `--set` variables) killed every mutant but the boundary flips `P(fix_cost < yearly_saving)` to `<=` in `debt-cost`, `P(net_hours > 0)` to `>=` in `automation` and `P(cost > budget)` to `>=` in `migration`, which are equivalent: a continuous sample equals the threshold with probability zero.

## Alternatives considered

- **Count every lost hour and every fix at 100%, as the issue's formulas read:** shorter, but it inflates the cost of waiting and the payback of every fix, the usual way such cases are oversold.
- **Model flaky failures with `binomial` per run:** shows the discrete count, but the expected value is the formula the issue names and reads more simply, as with attrition in ADR 0044.
- **Price `automation` in euros:** adds two assumptions that cancel out and hides the xkcd table's point that it is about time.
- **Put these scenarios under `delivery`:** they affect delivery, but the questions are about where engineering time goes, and the area would mix forecasts with cost models.
