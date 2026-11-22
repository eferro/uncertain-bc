# ADR 0046: Product scenarios around cost of delay, with ranges in place of confidence factors

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0042](0042-scenario-catalogue.md), [ADR 0043](0043-reusable-formulas-for-the-scenarios.md), [ADR 0044](0044-team-scenarios.md)

## Context

Issue #24 asks for twelve product strategy scenarios with cost of delay as the central theme: the cost of delay of a feature, CD3 over a small backlog, the urgency profiles, delaying a launch vs a reduced scope, the cost of queues, a feature's business case, RICE with ranges, TAM/SAM/SOM, LTV and CAC payback, a price change, an A/B experiment, and buy vs build. Some of them need choices the issue leaves open: the language has no conditionals, so a fixed-date profile cannot be written as "if late then penalty"; RICE multiplies by a confidence percentage, which overlaps with ranges; and "a likely 1–3% uplift" says nothing about variants that do not win.

## Decision

- Twelve scenarios open the `product` area, one file each, following the catalogue template and the testing rule of ADR 0044 (one testscript per scenario, `testdata/script/scenarios-product-<name>.txtar`, pinning medians, the answer line, the `--set` answer and one `--sens` driver): `cost-of-delay`, `cd3`, `urgency-profiles`, `delay-vs-reduced-scope`, `queue-cost`, `business-case`, `rice`, `tam-sam-som`, `ltv-cac`, `price-change`, `ab-experiment` and `buy-vs-build`.
- The cost of delay scenarios reuse the shared formulas: `delay_cost` in `cost-of-delay`, `urgency-profiles`, `delay-vs-reduced-scope` and `queue-cost`; `cd3` in `cd3`; `lead_time` in `queue-cost` to show how long items wait; `ltv` and `cac_payback` in `ltv-cac`. The others need none and run alone.
- Costs of delay are weekly and durations in weeks throughout the area, so the scenarios can be compared and combined.
- `urgency-profiles` models each profile by the shape of its cost over the same delay: standard and expedite are linear (`delay_cost`) with different weekly costs; fixed date is `(delay_weeks > slack_weeks) * missed_date`, because a comparison is 0 or 1 in every sample and stands in for the missing conditional; intangible is a weekly cost that grows geometrically, summed in closed form.
- `delay-vs-reduced-scope` measures both options by the value lost against shipping everything today: the delay loses the whole weekly value for the extra weeks; the reduced scope loses only the missing share, for the longer time until the rest ships, plus the cost of a second release.
- `rice` drops the confidence factor: a shaky guess is a wide range, and the answer is how likely one feature is to score higher than the other, so multiplying by a confidence percentage as well would discount the same doubt twice. The header says so.
- `ab-experiment` multiplies the uplift by an uncertain chance that the variant wins, an expected value like ADR 0044's attrition, because most experiments do not win and the issue's 1–3% only describes the winning case.
- `buy-vs-build` compares undiscounted money over three years; the opportunity cost of the team building it is named in the header and left out, as `in-house-vs-contractors` leaves out knowledge.

## Consequences

- Readers get a coherent set of cost of delay models in the same units, most of them using the shared formulas.
- Changing a product scenario's numbers or header means updating its testscript, as with the team and delivery scenarios.
- Readers used to RICE's confidence percentage will not find it; the header explains why.
- A manual mutation pass over the twelve models killed every mutant but one, `P(payback <= 12)` turned into `P(payback < 12)` in `ltv-cac`, which is equivalent: a continuous sample is exactly 12 with probability zero.

## Alternatives considered

- **Keep RICE's confidence factor next to the ranges:** familiar, but it counts the same uncertainty twice and the probabilistic answer already says how sure the ranking is.
- **Model the A/B experiment's success with a Bernoulli draw per sample:** shows the all-or-nothing outcome, but the expected value is the formula product teams use to decide and reads more simply.
- **Discount buy vs build with `npv` or `present_value`:** more correct for long horizons, but over three years it changes little and adds a rate the reader must guess.
- **One file with the four urgency profiles plotted against several delays:** shows the curves, but the language has no loops and the scenario would grow long; one uncertain delay compares the profiles in a few lines.
