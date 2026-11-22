# ADR 0043: Reusable formulas for the scenarios in one shared file

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0042](0042-scenario-catalogue.md)

## Context

Issue #21 asks for formulas that many scenarios need (loaded cost of a person, Little's law, LTV and CAC payback, annual loss expectancy, availability in series and in parallel, compound growth, net present value, cost of delay and CD3, focus factor) written once as `ubc` user functions, so scenarios can share them. Files given on the command line already run in one session, so a function defined in one file can be called from the next (`ubc scenarios/formulas.ubc scenarios/<area>/<name>.ubc`). The catalogue test (ADR 0042) runs every scenario on its own, which would reject a scenario that calls a shared formula.

## Decision

- The formulas live in one file, `scenarios/formulas.ubc`, at the top of the catalogue, not in an area: they serve every area. It holds only definitions, so running it prints nothing.
- Each formula is one user function with a comment that explains it, gives a rule-of-thumb range for its parameters when one exists, and names where it comes from. The functions are grouped by topic (people, flow, product, risk, money over time):
  - `loaded_cost(salary, overhead)`, `real_capacity(people, days, focus)`;
  - `lead_time(wip, throughput)` (Little's law), `delay_cost(cost_of_delay, periods)`, `cd3(cost_of_delay, duration)`;
  - `ltv(revenue_per_period, gross_margin, churn)`, `cac_payback(cac, revenue_per_period, gross_margin)`;
  - `annual_loss(probability, impact)`, `availability_series(a, b)`, `availability_parallel(a, b)`;
  - `compound(value, rate, periods)`, `present_value(amount, rate, periods)`, `npv(investment, cash_per_period, rate, periods)`.
- Names are snake_case like the scenarios' variables. Functions and variables live in different namespaces, so a model may still have a variable called `ltv`.
- The language has no variadic user functions, so the availability formulas take two services; more are chained (`availability_series(a, availability_series(b, c))`). NPV takes a constant cash flow, the present value of an annuity minus the investment, because lists cannot be folded in the language; `present_value` covers one-off amounts. Rates and periods are in whatever unit the caller keeps consistent.
- The catalogue test runs every scenario after `formulas.ubc`, in one session, exactly as the README tells readers to; `formulas.ubc` itself runs alone. A scenario that calls a formula shows the two-file command in its "Run it" lines. `formulas.ubc` is indexed like any scenario (`](formulas.ubc)`), and `scenarios/README.md` gains a section on how to use it.
- A testscript (`testdata/script/scenarios-formulas.txtar`) checks that the file prints nothing and pins the value of every formula with plain numbers, so a wrong formula turns the build red; it also runs the uncertain example shown in `scenarios/README.md`.

## Consequences

- Scenarios can share a formula instead of repeating it, and readers can reuse the file in their own models.
- A scenario that uses a formula no longer runs on its own: forgetting `scenarios/formulas.ubc` in the command is an evaluation error (`unknown function`).
- Renaming or changing a formula's parameters breaks the scenarios that call it, which the catalogue test catches.

## Alternatives considered

- **Built-in functions in the language:** checked by unit tests and always available, but they would grow the language with domain formulas and hide the arithmetic that the comments make visible and editable.
- **One formulas file per area:** keeps them closer to their scenarios, but several formulas (loaded cost, NPV) serve more than one area, and readers would have to know which file to load.
- **An `import` or `:load` statement in files:** would make scenarios self-contained, but it is a language change bigger than this issue needs; passing two files already works.
