# Scenario catalogue

Ready-to-use `ubc` models for questions that come up in software engineering and product strategy. Pick the one closest to your question, change the assumptions you know better, and run it. You do not need to learn the whole language first: every scenario explains itself in its comments. If you want to learn the language, start with the [tutorial](../docs/TUTORIAL.md).

Run a scenario from the repository root:

```console
$ ubc scenarios/team/replace-a-person.ubc
```

Every scenario is checked by the test suite: it must evaluate without errors and be listed in the index below ([ADR 0042](../docs/adr/0042-scenario-catalogue.md)).

## Index

| Area | Scenario | Question it answers |
|------|----------|---------------------|
| team | [replace-a-person](team/replace-a-person.ubc) | How much does it cost to replace an engineer who leaves, and how likely is it to exceed 50k euros? |

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
