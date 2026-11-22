# ADR 0055: Catalogue index by area, a currency note and sources for every assumption

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0042](0042-scenario-catalogue.md)

## Context

The scenario catalogue (ADR 0042) grew to 43 scenarios in seven areas (ADRs 0044 to 0050), indexed in one long table in `scenarios/README.md` whose question column repeated each scenario's full question. A newcomer had no obvious place to start, some readers asked why amounts are in euros or dollars, and the assumptions said what each variable is more often than where its range comes from. Issue #33 agreed with the maintainer: a "Start here" section with five featured scenarios, then one section per area with its own table and shorter questions; keep euros and dollars and add a note that `ubc` has no units; every assumption explains where its range comes from and names a known source when there is one, never an invented one.

## Decision

- **Index.** `scenarios/README.md` opens its index with a note on currency and sources, then `### Start here` with five scenarios, one from each of the areas a newcomer is most likely to ask about (`team/replace-a-person`, `delivery/delivery-date`, `product/cost-of-delay`, `engineering/debt-cost`, `reliability/incident-cost`), with their full question. Then one `### <Area>` section per area (`### Team`, `### Delivery`, `### Product`, `### Engineering`, `### Reliability`, `### Cloud`, `### Risk`), each with a table of the scenario name and a short question; the tutorial's examples stay in their area's table. The `formulas` row goes away: the "Reusable formulas" section already links it.
- **Currency.** Amounts stay in euros, and in dollars in the cloud area, where providers bill in dollars (ADR 0049). The index says that `ubc` has no units and that a reader changes the numbers to their currency; the calculation does not change.
- **Sources.** Every scenario's header ends with a `# Sources:` paragraph that names the published references its ranges rely on (author, title and year, a report or a public price list), or says that they come from typical experience when nothing published backs them. Every assignment under `# Assumptions` keeps a comment on its own line that says where its range comes from: typical experience, a rule of thumb, the reader's own data, a price list or a named source. A reference is cited only when it exists and says what the comment attributes to it; figures are attributed only when the source gives them, otherwise the comment points to the source's current edition. Only comments change: the numbers and the outputs the testscripts pin stay the same.
- **Tests.** No CLI behaviour changes, so the acceptance checks are tests over the catalogue in `cmd/ubc/scenarios_test.go`, next to the existing index tests: the "Start here" section links exactly five scenarios, each area's section links every scenario of that area, and every scenario has a `# Sources:` line and a comment on every assumption.

## Consequences

- A reader finds a first scenario in the first screen of the index and the rest grouped by area, and can trust that a cited source exists.
- A new scenario needs a `# Sources:` line, a comment per assumption and a row in its area's table, or `make validate` fails; a new area needs its own `### <Area>` section.
- Results without a short label are named by their line (ADR 0053), so the `# Sources:` lines moved the `L<n>` labels of five scenarios; their testscripts were updated with the new line numbers and the values did not change.
- Whether a comment really explains the range, and whether a citation is accurate, stays a matter of review: the tests check the shape, not the content.

## Alternatives considered

- **Keep one table with an area column:** shorter to maintain, but long questions made it hard to scan and there was no starting point.
- **Convert every amount to one currency, or to no currency:** cloud prices are quoted in dollars and the rest of the examples in euros; a neutral "units" would read worse than a note.
- **A link or footnote per assumption:** URLs rot and most ranges come from experience, not a document; one `# Sources:` paragraph per scenario plus a short comment per line keeps the files readable.
- **A testscript for the change:** the CLI output is unchanged, so a testscript would only repeat the file tests.
