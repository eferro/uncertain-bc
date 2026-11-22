# Scenario catalogue

Ready-to-use `ubc` models for questions that come up in software engineering and product strategy. Pick the one closest to your question, change the assumptions you know better, and run it. You do not need to learn the whole language first: every scenario explains itself in its comments. If you want to learn the language, start with the [tutorial](../docs/TUTORIAL.md).

Every scenario is built into `ubc`, so it runs wherever `ubc` is installed, with no clone of this repository ([ADR 0054](../docs/adr/0054-scenarios-built-into-the-binary.md)). List them, run one, or print one to copy and edit it:

```console
$ ubc --scenarios
$ ubc --scenario team/replace-a-person
$ ubc --set 'vacancy_months = 1 to 3' --scenario team/replace-a-person
$ ubc --scenario-source team/replace-a-person > replace-a-person.ubc
```

`--scenario` combines with every other flag (`--set`, `--data`, `--sens`, `--echo`, `--assert`, …), runs before any `-e` or file, and loads the [shared formulas](#reusable-formulas) first. From a clone, the files run as any model does: `ubc scenarios/formulas.ubc scenarios/team/replace-a-person.ubc`.

Every scenario is checked by the test suite: it must evaluate without errors, be listed in its area's section of the index below, have a `# Sources:` line and a comment on every assumption, and every model the index links to must exist ([ADR 0042](../docs/adr/0042-scenario-catalogue.md), [ADR 0045](../docs/adr/0045-delivery-scenarios.md), [ADR 0055](../docs/adr/0055-catalogue-index-by-area-and-sources.md)).

## Index

`ubc` has no units: amounts are in euros in most scenarios and in dollars in the cloud ones, the currency providers bill in. To use your own currency, change the numbers; the calculation stays the same. Every assumption says in its comment where its range comes from, and each scenario's `# Sources:` line names the published references it relies on, or says that the ranges come from typical experience when there is none.

### Start here

Five scenarios, one per area, that show what the catalogue does:

| Scenario | Question it answers |
|----------|---------------------|
| [team/replace-a-person](team/replace-a-person.ubc) | How much does it cost to replace an engineer who leaves, and how likely is it to exceed 50k euros? |
| [delivery/delivery-date](delivery/delivery-date.ubc) | With the throughput the team has really had, how many weeks will the backlog take, and how likely is it to be done within 12 weeks? |
| [product/cost-of-delay](product/cost-of-delay.ubc) | How much value do we lose if a feature ships later than planned, and how likely is the delay to cost more than 100k euros? |
| [engineering/debt-cost](engineering/debt-cost.ubc) | What does living with a piece of technical debt cost a year compared with fixing it, and how likely is the fix to pay for itself within a year? |
| [reliability/incident-cost](reliability/incident-cost.ubc) | What do incidents cost a year, and how likely is it to exceed 500k euros? |

### Team

| Scenario | Question |
|----------|----------|
| [replace-a-person](team/replace-a-person.ubc) | What does replacing an engineer who leaves cost? |
| [hiring-cost](team/hiring-cost.ubc) | How long does filling a position take, and how many hours of hiring? |
| [attrition-cost](team/attrition-cost.ubc) | What does attrition cost in a year? |
| [onboarding](team/onboarding.ubc) | How much work does onboarding a new engineer take? |
| [team-capacity](team/team-capacity.ubc) | How many person-days can a team really deliver in a quarter? |
| [meeting-cost](team/meeting-cost.ubc) | What does a recurring weekly meeting cost in a year? |
| [in-house-vs-contractors](team/in-house-vs-contractors.ubc) | Is hiring, contractors or outsourcing cheapest for a year of work? |
| [bus-factor](team/bus-factor.ubc) | What does depending on the only person who knows a system cost? |

### Delivery

| Scenario | Question |
|----------|----------|
| [delivery-date](delivery/delivery-date.ubc) | How many weeks will the backlog take at our real throughput (`--data throughput=weeks.txt`)? |
| [project-estimate](../examples/project-estimate.ubc) | How long will a feature estimated as task ranges take? (From the [tutorial](../docs/TUTORIAL.md).) |
| [littles-law](delivery/littles-law.ubc) | How long does an item take, and how much would a WIP limit help? |
| [scope-creep](delivery/scope-creep.ubc) | How much later do we finish if the scope grows while we build? |
| [context-switching](delivery/context-switching.ubc) | How much of the day and of the budget do interruptions eat? |

### Product

| Scenario | Question |
|----------|----------|
| [cost-of-delay](product/cost-of-delay.ubc) | How much value do we lose if a feature ships late? |
| [cd3](product/cd3.ubc) | Which of three features goes first by CD3? |
| [urgency-profiles](product/urgency-profiles.ubc) | What does the same delay cost under each urgency profile? |
| [delay-vs-reduced-scope](product/delay-vs-reduced-scope.ubc) | Delay the launch for the whole feature, or ship a reduced scope now? |
| [queue-cost](product/queue-cost.ubc) | What does the work waiting in queues cost? |
| [business-case](product/business-case.ubc) | Does a feature pay for itself in two years? |
| [rice](product/rice.ubc) | Scored by RICE with ranges, which feature scores higher? |
| [tam-sam-som](product/tam-sam-som.ubc) | How big are the TAM, SAM and SOM of a B2B product? |
| [ltv-cac](product/ltv-cac.ubc) | What is a customer worth, and when do they pay back their acquisition? |
| [price-change](product/price-change.ubc) | Does a price rise make up for the customers it loses? |
| [ab-experiment](product/ab-experiment.ubc) | Is an A/B experiment worth running? |
| [buy-vs-build](product/buy-vs-build.ubc) | Over three years, is buying or building a tool cheaper? |

### Engineering

| Scenario | Question |
|----------|----------|
| [debt-cost](engineering/debt-cost.ubc) | Living with a piece of technical debt or fixing it: which costs more? |
| [automation](engineering/automation.ubc) | Is automating a routine task worth the time? |
| [slow-ci](engineering/slow-ci.ubc) | What does a slow CI pipeline cost a year? |
| [flaky-tests](engineering/flaky-tests.ubc) | How much time do flaky tests waste a year? |
| [migration](engineering/migration.ubc) | What will a migration cost, and will it exceed the budget? |

### Reliability

| Scenario | Question |
|----------|----------|
| [incident-cost](reliability/incident-cost.ubc) | What do incidents cost a year? |
| [availability](reliability/availability.ubc) | What availability do pieces in series and in parallel give? |
| [error-budget](reliability/error-budget.ubc) | Will we overspend an SLO's error budget this month? |
| [capacity](reliability/capacity.ubc) | How many nodes does next year's peak traffic need? |
| [bug-cost](reliability/bug-cost.ubc) | What does a bug cost when caught in review, in QA or in production? |

### Cloud

| Scenario | Question |
|----------|----------|
| [serverless-bill](../examples/cloud-cost.ubc) | What will a small serverless web app cost a month? (From the [tutorial](../docs/TUTORIAL.md).) |
| [llm-cost](cloud/llm-cost.ubc) | What does an LLM feature billed by the token cost a month? |
| [observability](cloud/observability.ubc) | What do logs billed by the gigabyte cost a month? |
| [storage-growth](cloud/storage-growth.ubc) | What will growing storage cost in two years? |
| [reserved-capacity](cloud/reserved-capacity.ubc) | Is a year of reserved capacity cheaper than on demand? |
| [cost-per-customer](cloud/cost-per-customer.ubc) | What does the infrastructure for one customer cost a month? |
| [saas-seats](cloud/saas-seats.ubc) | What will per-seat SaaS tools cost as the team grows? |
| [ephemeral-environments](cloud/ephemeral-environments.ubc) | What do ephemeral environments per pull request cost a month? |

### Risk

| Scenario | Question |
|----------|----------|
| [annual-loss](risk/annual-loss.ubc) | What does a risk cost a year, and does a mitigation pay off? |
| [data-breach](risk/data-breach.ubc) | What does a data breach cost, and what loss should we expect a year? |

## Reusable formulas

[`formulas.ubc`](formulas.ubc) defines formulas that come up again and again as `ubc` user functions, each with a comment explaining it and where it comes from. `ubc --scenario <area>/<name>` loads them before the scenario by itself. To use them in your own model, load them with `--scenario formulas`: everything runs in one session, so the functions flow into the files that follow.

```console
$ ubc --scenario formulas my-model.ubc
```

Inside the model, call them with plain or uncertain arguments, for example `lead_time(15 to 30, 3 to 6)` or `npv(50k, 2k to 4k, 0.01, 36)`. The test suite runs every scenario after `formulas.ubc`, as `--scenario` does, so a scenario may use them.

## Template

Every scenario lives in `scenarios/<area>/<name>.ubc` and follows the same shape, written as `#` comments in the file itself:

1. **The question** it answers, in one or two lines.
2. **How to run it**: the plain command, a `--set` example that changes one assumption, and a `--sens` example that shows which assumption matters most.
3. **Sources**: a `# Sources:` line in the header that names the published references the ranges rely on (a book, a report, a public price list), or says that they come from typical experience when there is none. Never invent a reference.
4. **Assumptions**: one variable per line, each with its range (`a to b` is a 90% confidence interval) and a comment on where the range comes from: typical experience, an industry rule of thumb or a named source.
5. **Calculation**: each step with a short comment explaining it.
6. **Result**: the final value, and a `P(...)` line that answers the question.

```
# Question: <what does this answer?>
#
# Run it:
#   ubc --scenario <area>/<name>
# Change an assumption:
#   ubc --set '<variable> = <a> to <b>' --scenario <area>/<name>
# See which assumption matters most:
#   ubc --sens <result> --scenario <area>/<name>
#
# Sources: <the references the ranges rely on, or "typical experience">

# Assumptions
<variable> = <a> to <b>      # where this range comes from

# Calculation
<step> = <expression>        # what this step adds up

# Result
<result>
P(<result> > <threshold>)    # the answer to the question
```

The first line, `# Question:`, and the comment lines that follow it up to the first empty `#` are what `ubc --scenarios` prints next to the name.

To add a scenario: copy the template, put the file under its area (one level deep: `<area>/<name>.ubc` is what the binary embeds), add a row to its area's table in the index, and run `make validate`.
