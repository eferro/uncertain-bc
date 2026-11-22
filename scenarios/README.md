# Scenario catalogue

Ready-to-use `ubc` models for questions that come up in software engineering and product strategy. Pick the one closest to your question, change the assumptions you know better, and run it. You do not need to learn the whole language first: every scenario explains itself in its comments. If you want to learn the language, start with the [tutorial](../docs/TUTORIAL.md).

Run a scenario from the repository root:

```console
$ ubc scenarios/team/replace-a-person.ubc
```

Every scenario is checked by the test suite: it must evaluate without errors and be listed in the index below, and every model the index links to must exist ([ADR 0042](../docs/adr/0042-scenario-catalogue.md), [ADR 0045](../docs/adr/0045-delivery-scenarios.md)).

## Index

| Area | Scenario | Question it answers |
|------|----------|---------------------|
| team | [replace-a-person](team/replace-a-person.ubc) | How much does it cost to replace an engineer who leaves, and how likely is it to exceed 50k euros? |
| team | [hiring-cost](team/hiring-cost.ubc) | How long does it take to fill a position, how many hours does the hiring process take, and how likely is the seat to stay empty for more than three months? |
| team | [attrition-cost](team/attrition-cost.ubc) | How much does attrition cost in a year (people × chance of leaving × replacement cost), and how likely is it to exceed 500k euros? |
| team | [onboarding](team/onboarding.ubc) | How much work does onboarding a new engineer take, ramp-up and the team's help included, and how likely is it to exceed three months of an engineer? |
| team | [team-capacity](team/team-capacity.ubc) | How many person-days can a team really deliver in a quarter, after holidays, sick leave, on-call and the focus factor, and how likely is it to fall short of the roadmap? |
| team | [meeting-cost](team/meeting-cost.ubc) | How much does a recurring weekly meeting cost in a year, and how likely is it to exceed 30k euros? |
| team | [in-house-vs-contractors](team/in-house-vs-contractors.ubc) | For a year of work, what do hiring, contractors and outsourcing cost, and how likely is hiring to be the cheapest? |
| team | [bus-factor](team/bus-factor.ubc) | What does it cost a year to depend on the only person who knows a system, and is spreading the knowledge worth it? |
| delivery | [delivery-date](delivery/delivery-date.ubc) | With the throughput the team has really had (`--data throughput=weeks.txt`), how many weeks will the backlog take, and how likely is it to be done within 12 weeks? |
| delivery | [project-estimate](../examples/project-estimate.ubc) | How long will a feature take when it is estimated as a sum of task ranges plus forgotten work, and how likely is it to miss the deadline? (A worked model from the [tutorial](../docs/TUTORIAL.md).) |
| delivery | [littles-law](delivery/littles-law.ubc) | How long does an item take, given the work in progress and the throughput (Little's law), and how much shorter would a WIP limit make it? |
| delivery | [scope-creep](delivery/scope-creep.ubc) | If the scope grows while we build it, how much later do we finish than the plan says, and how likely are we to miss the date? |
| delivery | [context-switching](delivery/context-switching.ubc) | How much of the working day do interruptions eat (interruptions a day × minutes lost × people), what does it cost a year, and how likely is it to be more than a quarter of the day? |
| product | [cost-of-delay](product/cost-of-delay.ubc) | How much value do we lose if a feature ships later than planned (weekly value × weeks of delay), and how likely is the delay to cost more than 100k euros? |
| product | [cd3](product/cd3.ubc) | With three features competing for the same team, which one should go first by CD3 (cost of delay divided by duration), and how likely is the small one to be the right first pick? |
| product | [urgency-profiles](product/urgency-profiles.ubc) | For the same delay, what do the standard, expedite, fixed-date and intangible urgency profiles cost, and how likely is the fixed-date work to lose more than the standard one? |
| product | [delay-vs-reduced-scope](product/delay-vs-reduced-scope.ubc) | Should we delay the launch until the whole feature is ready or ship now with a reduced scope, and how likely is the reduced scope to lose less value? |
| product | [queue-cost](product/queue-cost.ubc) | What does the work waiting in queues cost (items waiting × weekly cost of delay), and how likely is it to exceed 1M euros a year? |
| product | [business-case](product/business-case.ubc) | Does a feature pay for itself in two years (reach × adoption × conversion × ARPU − build and run cost), and how likely is it to make money? |
| product | [rice](product/rice.ubc) | Scored by RICE with ranges instead of false precision, how likely is one feature to score higher than another? |
| product | [tam-sam-som](product/tam-sam-som.ubc) | Fermi style, how big are the TAM, SAM and SOM of a B2B product, and how likely is the SOM to exceed 10M euros a year? |
| product | [ltv-cac](product/ltv-cac.ubc) | What is a customer worth over their lifetime, how many months until they pay back their acquisition cost, and how likely is the payback to stay within 12 months? |
| product | [price-change](product/price-change.ubc) | If we raise the price, does the extra revenue from the customers who stay make up for the ones we lose, and how likely is revenue to grow? |
| product | [ab-experiment](product/ab-experiment.ubc) | Is an A/B experiment with a likely 1–3% uplift worth running, knowing most variants do not win, and how likely is it to pay back its cost? |
| product | [buy-vs-build](product/buy-vs-build.ubc) | Over three years, is a SaaS licence or building and maintaining our own tool cheaper, and how likely is building to be the cheaper option? |
| engineering | [debt-cost](engineering/debt-cost.ubc) | What does living with a piece of technical debt cost a year (hours lost a week × people × cost per hour) compared with fixing it, and how likely is the fix to pay for itself within a year? |
| engineering | [automation](engineering/automation.ubc) | Is it worth automating a routine task (the xkcd "Is It Worth the Time?" table, with ranges and maintenance), and how likely is it to save more time than it takes? |
| engineering | [slow-ci](engineering/slow-ci.ubc) | What does a slow CI pipeline cost a year (minutes waiting × builds a day × people × cost per hour), and how likely is it to cost more than 100k euros? |
| engineering | [flaky-tests](engineering/flaky-tests.ubc) | How much time do flaky tests waste a year (failed runs × retries × pipeline duration, plus checking each failure), and how likely is it to be more than 300 hours? |
| engineering | [migration](engineering/migration.ubc) | What will a migration to a new framework, language or cloud cost, forgotten work and running both systems included, and how likely is it to exceed the budget? |
| reliability | [incident-cost](reliability/incident-cost.ubc) | What do incidents cost a year (incidents × hours to restore × cost per hour of downtime, plus the responders' time), and how likely is it to exceed 500k euros? |
| reliability | [availability](reliability/availability.ubc) | What availability does a service get from pieces in series and in parallel, how many minutes a month is it down, and how likely is it to meet a 99.9% target? |
| reliability | [error-budget](reliability/error-budget.ubc) | How many minutes of downtime does an SLO allow in a month, how many can we expect, and how likely are we to overspend the error budget? |
| reliability | [capacity](reliability/capacity.ubc) | How many nodes does next year's peak traffic need, with headroom and one spare, and how likely is the current cluster to fall short? |
| reliability | [bug-cost](reliability/bug-cost.ubc) | What does a bug cost when it is caught in code review, in QA or in production, and how likely is a production bug to cost more than ten times one caught in QA? |
| all | [formulas](formulas.ubc) | Reusable formulas shared by every scenario: loaded cost, real capacity, Little's law, cost of delay and CD3, LTV and CAC payback, annual loss expectancy, availability in series and in parallel, compound growth, present value and NPV. |

## Reusable formulas

[`formulas.ubc`](formulas.ubc) defines formulas that come up again and again as `ubc` user functions, each with a comment explaining it and where it comes from. Run it before a model, in the same command: files run in one session, so the functions flow from one file to the next.

```console
$ ubc scenarios/formulas.ubc scenarios/<area>/<name>.ubc
$ ubc scenarios/formulas.ubc my-model.ubc
```

Inside the model, call them with plain or uncertain arguments, for example `lead_time(15 to 30, 3 to 6)` or `npv(50k, 2k to 4k, 0.01, 36)`. The test suite runs every scenario after `formulas.ubc`, so a scenario may use them; when it does, its "Run it" lines show the two-file command.

## Template

Every scenario lives in `scenarios/<area>/<name>.ubc` and follows the same shape, written as `#` comments in the file itself:

1. **The question** it answers, in one or two lines.
2. **How to run it**: the plain command, a `--set` example that changes one assumption, and a `--sens` example that shows which assumption matters most.
3. **Assumptions**: one variable per line, each with its range (`a to b` is a 90% confidence interval) and a comment on where the numbers come from.
4. **Calculation**: each step with a short comment explaining it.
5. **Result**: the final value, and a `P(...)` line that answers the question.

```
# Question: <what does this answer?>
#
# Run it:
#   ubc scenarios/<area>/<name>.ubc
# Change an assumption:
#   ubc --set '<variable> = <a> to <b>' scenarios/<area>/<name>.ubc
# See which assumption matters most:
#   ubc --sens <result> scenarios/<area>/<name>.ubc

# Assumptions
<variable> = <a> to <b>      # where these numbers come from

# Calculation
<step> = <expression>        # what this step adds up

# Result
<result>
P(<result> > <threshold>)    # the answer to the question
```

To add a scenario: copy the template, put the file under its area, add a row to the index, and run `make validate`.
