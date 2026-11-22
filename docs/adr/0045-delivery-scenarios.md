# ADR 0045: Delivery scenarios, with observed data as a replaceable default

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0042](0042-scenario-catalogue.md), [ADR 0043](0043-reusable-formulas-for-the-scenarios.md), [ADR 0044](0044-team-scenarios.md)

## Context

Issue #23 asks for scenarios about delivery forecasting and estimation: a delivery date from historical throughput read with `--data throughput=weeks.txt`, an estimate as a sum of task ranges (linking the existing `examples/project-estimate.ubc`), Little's law and the effect of lowering WIP, scope creep, and context switching. The catalogue test (ADR 0042) runs every scenario with no flags, so a model that only gets its throughput from `--data` would fail with an unknown variable. Testscripts copy only `.ubc` files into their working directory, so a data file shipped next to the scenario would not reach them. The tutorial's project estimate is a model, but not a scenario: it lives in `examples/`, where the tutorial tests pin its output, and it does not follow the catalogue template.

## Decision

- Four scenarios open the `delivery` area, one file each, following the catalogue template and the testing rule of ADR 0044 (one testscript per scenario, `testdata/script/scenarios-delivery-<name>.txtar`, pinning medians, the answer line, the `--set` answer and one `--sens` driver): `delivery-date`, `littles-law`, `scope-creep` and `context-switching`.
- `delivery-date` assigns `throughput` an example list of twelve weekly counts. Since `--data` overrides assignments like `--set`, `ubc --data throughput=weeks.txt` replaces it with the reader's own history, and the scenario runs on its own for the catalogue test and for a first look. No data file is shipped: the header names the reader's own `weeks.txt`, and the testscript supplies one inline to prove the override reaches the model.
- The forecast does not divide the backlog by one week's throughput, which would treat a single bad week as the pace for the whole project. The pace over the coming weeks is the average of many weeks, so it is drawn as `normal(mean(throughput), stdev(throughput) / sqrt(horizon))`, with `horizon` the rough number of weeks ahead (`mean(backlog) / mean(throughput)`). The language cannot sum a variable number of resampled weeks, and this central-limit approximation keeps the formula short and visible. As a consequence `--sens weeks` ranks `weekly_rate`, not `throughput`, which only feeds the model through its mean and spread.
- `littles-law` reuses `lead_time` from `formulas.ubc` for the current WIP and for a WIP limit, with the same throughput, because Little's law says limiting WIP shortens the wait without changing what the team finishes per week.
- `context-switching` reuses `loaded_cost` and answers with the share of the working day lost; `scope-creep` needs no shared formula and reports the delay that the growth alone adds.
- The estimate as a sum of task ranges is not copied into the catalogue: its index row links `../examples/project-estimate.ubc`, keeping one copy that the tutorial tests already check. Because the index can now point outside `scenarios/`, the catalogue test also fails when any `.ubc` link in the index points to a file that does not exist.

## Consequences

- Readers can forecast with their own throughput file without editing the model, and the catalogue test still runs every scenario with no flags.
- Renaming or moving a linked example breaks the build instead of leaving a dead link in the index.
- Changing a delivery scenario's numbers or header means updating its testscript, as with the team scenarios.

## Alternatives considered

- **Ship `scenarios/delivery/weeks.txt` and make the catalogue test pass `--data`:** the header command would work as written, but the test would need per-scenario flags and testscripts would need to copy data files too.
- **`weeks = backlog / throughput` with one resampled week:** simplest, but it makes the forecast far wider than the history justifies.
- **Copy `project-estimate.ubc` into `scenarios/delivery/`:** it would follow the template, but there would be two copies of the same model to keep in step.
