# ADR 0050: Risk and security scenarios around the annual loss expectancy

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0042](0042-scenario-catalogue.md), [ADR 0043](0043-reusable-formulas-for-the-scenarios.md), [ADR 0044](0044-team-scenarios.md), [ADR 0048](0048-reliability-scenarios.md)

## Context

Issue #28 asks for scenarios about risk and security: the annual loss expectancy (probability per year × impact) and whether a mitigation pays off, and the cost of a data breach as records × cost per record + a fine. The shared `annual_loss` formula already exists (ADR 0043). A mitigation's effectiveness is a share between 0 and 1, and an `a to b` range of a reduction close to 100% (`50% to 90%`) is a lognormal whose tail goes above 1, which would make the loss with the mitigation negative.

## Decision

- Two scenarios open a new `risk` area, one file each, following the catalogue template and the testing rule of ADR 0044 (one testscript per scenario, `testdata/script/scenarios-risk-<name>.txtar`, pinning medians, the answer line, the `--set` answer and one `--sens` driver): `annual-loss` and `data-breach`. Amounts are in euros, like every area except `cloud`.
- `annual-loss` computes the expected loss a year before and after a mitigation with `annual_loss`, and answers with the chance that the loss it avoids exceeds its yearly cost (`P(net_benefit > 0)`). The mitigation is described by the share of the probability that remains (`residual = 20% to 50%`), as availability is by its downtime share in ADR 0048, so no sample turns the loss negative. Its header says that the annual loss expectancy is an average over years, not what one year costs, and points to `reliability/incident-cost` for events that happen several times a year.
- `data-breach` prices one breach as `records × cost_per_record + fine`, the fine being a share of the yearly turnover, and turns it into a yearly loss with `annual_loss` and a breach probability. It answers with the chance that one breach costs more than 5M euros. Its header cites the order of magnitude of IBM's Cost of a Data Breach per-record cost, warns that it falls for very large breaches, and gives the GDPR cap (4% of the worldwide turnover or 20M euros, whichever is higher). The cap is stated, not modelled: with realistic fine shares it never binds, so a `min` would be dead code no test could tell apart.
- Assumptions are tuned so that each answer sits away from 0 and 1 with the defaults, and the `--set` example in each header moves it clearly, which is what the testscript pins.

## Consequences

- Readers get the classic quantitative risk calculation with honest ranges and can see with `--sens` what drives it: with the defaults, the impact decides whether the mitigation pays off, and the number of records exposed decides the cost of a breach.
- Changing a risk scenario's numbers or header means updating its testscript, as with the other areas.
- A manual mutation pass over the two models (arithmetic operators, removed terms and factors, `residual` dropped from the loss after, the shared formula bypassed, constants, the answer's comparison and threshold, the assumptions' ranges) killed 29 of 31 mutants. The survivors are the boundary flips `>` to `>=` in `P(net_benefit > 0)` and `P(breach_cost > 5M)`, which are equivalent: a continuous sample equals the threshold with probability zero.

## Alternatives considered

- **Describe the mitigation by the reduction it brings (`50% to 90%`):** reads naturally, but the lognormal tail goes above 100% and some samples make the residual loss negative.
- **Draw whether the incident happens with `binomial(1, probability)`:** shows that most years cost nothing, but the issue asks for the expected loss, and the header explains the difference instead.
- **Model the GDPR cap with `min`:** faithful to the law, but it never binds with the defaults, so it adds a line no test can check.
