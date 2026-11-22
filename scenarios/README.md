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
