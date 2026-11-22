# ADR 0049: Cloud and variable cost scenarios, in dollars

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0042](0042-scenario-catalogue.md), [ADR 0043](0043-reusable-formulas-for-the-scenarios.md), [ADR 0044](0044-team-scenarios.md), [ADR 0045](0045-delivery-scenarios.md)

## Context

Issue #27 asks for scenarios about cloud and other strongly variable costs: a serverless bill (linking the existing `examples/cloud-cost.ubc`), an LLM or API billed by the token, observability billed by the gigabyte of logs, storage with compound growth, reserved capacity or savings plans against on-demand under an uncertain usage, the infrastructure cost per customer, per-seat SaaS as the team grows, and ephemeral environments per pull request. The other areas count money in euros, but cloud providers, LLM APIs and SaaS vendors publish their prices in dollars, and the tutorial's cloud bill is already in dollars.

## Decision

- Seven scenarios open a new `cloud` area, one file each, following the catalogue template and the testing rule of ADR 0044 (one testscript per scenario, `testdata/script/scenarios-cloud-<name>.txtar`, pinning medians, the answer line, the `--set` answer and one `--sens` driver): `llm-cost`, `observability`, `storage-growth`, `reserved-capacity`, `cost-per-customer`, `saas-seats` and `ephemeral-environments`.
- Amounts in the `cloud` area are in dollars, the currency the price lists are published in, so a reader can copy a price without converting it; every header says so.
- The serverless bill is not copied: as with the project estimate in ADR 0045, its index row (`serverless-bill`) links `../examples/cloud-cost.ubc`, which the tutorial tests already pin and the catalogue test checks exists.
- `storage-growth` and `saas-seats` reuse the shared `compound` formula for monthly growth over a 24-month horizon. `storage-growth` also adds up every month's bill as a geometric series, `(compound(...) - today) / growth`, since the language has no loops; its header explains why compounding makes small monthly rates matter. `saas-seats` adds a share of unused seats on top of the people, because licences are rarely given back on time.
- `reserved-capacity` measures usage and the commitment at on-demand prices: the commitment is paid at a discount whether it is used or not, and usage above it is paid on demand (`max(usage - commitment, 0)`). It prints a year of savings, negative when the commitment is wasted, the chance of paying for unused capacity and the chance that committing is cheaper, which are different questions: a commitment can be partly unused and still save money.
- `cost-per-customer` divides the shared platform among the customers and adds what each one adds, and answers with the share of the price that goes to infrastructure against the 30% that a 70-80% SaaS gross margin leaves.
- `ephemeral-environments` prices environments a month as pull requests × hours alive × cost per hour, and also prints how many are up at once, by Little's law, because that is the number the cluster has to fit.
- `llm-cost` counts input and output tokens separately, with their own prices, because the prompt and its context usually outnumber the answer and output tokens cost several times more.
- Assumptions are tuned so that each answer sits away from 0 and 1 with the defaults, and the `--set` example in each header moves it clearly, which is what the testscript pins.

## Consequences

- Readers get models for the bills that move with usage and can see which assumption dominates with `--sens`: with the defaults, how long environments stay alive drives their cost, the monthly growth rate drives storage more than today's size, and the seat bill as much as the price per seat.
- Changing a cloud scenario's numbers or header means updating its testscript, as with the other areas.
- A manual mutation pass over the seven models (arithmetic operators, removed terms and factors, constants, `max` to `min` or removed, `compound` replaced by linear growth, the answer's comparison) killed 42 of 49 mutants. The survivors are the boundary flips `>` to `>=` (and `<` to `<=`) in the seven `P(...)` answers, which are equivalent: a continuous sample equals the threshold with probability zero.

## Alternatives considered

- **Euros, like the other areas:** consistent across the catalogue, but every price copied from a provider would need converting, and the linked tutorial model is in dollars.
- **A per-month loop or a list of months for the total storage cost:** the language has no loops; the closed-form geometric series gives the same total in one line.
- **Model reserved capacity in instance-hours:** closer to reserved instances, but savings plans commit to spend, and money at on-demand prices covers both with one formula.
