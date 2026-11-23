# ADR 0066: An assigned comparison keeps its samples (amends 0028)

- Status: Accepted
- Effective date: 2026-11-23
- Deciders: eduardo.ferro
- Amends: [ADR 0028](0028-override-variables-with-set.md)

## Context

Issue #41: `hit = uniform(0, 1) < 0.1` bound the probability `0.1` to `hit`, not the sample-by-sample outcomes, so `loss = hit * (50k to 500k)` was 10 % of the impact in every sample and `P(loss > 0)` was `1` instead of about `0.1`. Nothing warned about it. The same comparison written inline, `(uniform(0, 1) < 0.1) * (50k to 500k)`, kept its samples and gave the expected result.

SPEC §4 rule 5 said that a bare comparison at top level *prints* as a probability, which reads as display only. The evaluator turned the whole value into a probability for any statement whose expression was a comparison, assignments included, and ADR 0028 copied that rule for `--set` ("a bare comparison is a probability, as at top level"). "If it happens, it costs X" is a common way to model risks, incidents, missed dates and wins, and naming the event first is the natural way to write it.

## Decision

- **A comparison of uncertain values is an *event*:** the 0/1 outcome of the comparison in every sample. It keeps its samples wherever it goes, and it is *shown* as the probability that it happens: the mean of its samples. A comparison of two plain numbers stays a plain probability, `0` or `1`, as before.
- **Where an event is shown, it prints as a probability**, exactly as a top-level comparison did: the result line, `name = 0.10` for an assignment in the REPL and with `--echo`, a bare variable (`hit`), `:vars`, the skipped assignment of a `--set` override with `--echo`, `--plain` (the probability with full precision) and `--json` (`"kind":"probability"`). No histogram is drawn for it with `--hist`, and `:hist`/`:plot` still reject it as not uncertain, as before.
- **Where an event is used, it is its samples:** assigned to a variable, returned by a function whose body is a comparison, bound to `last` by a bare comparison, given to `--set`, passed to arithmetic or a function. Arithmetic turns it into an ordinary SampleSet (`hit * (50k to 500k)`, `last * 1`), which prints as a summary. `P(event)` is its probability; summaries (`mean`, `median`, `min`, `max`, `p`) treat it as its 0/1 samples.
- **Analysis sees the samples:** `--samples-out` writes a 0/1 column for an event; `--sens` counts an event variable as a driver and accepts a comparison as its target, ranked by rank correlation of the 0/1 samples.
- **`--assert` is unchanged in meaning:** a variable bound to an event is a condition and holds when its probability is exactly 1, the rule ADR 0029 already applies to a comparison; the failure message prints that probability.
- This amends ADR 0028: `--set 'hit=uniform(0, 1) < 0.1'` binds the event, like an assignment in the program, instead of its probability.
- **It is a bug fix under ADR 0062**, released as a patch version (1.1.1): the previous behaviour contradicted SPEC §4 rule 5, which only described how a comparison prints, and silently produced wrong models. The printed output of every program that does not compute with an assigned comparison is identical. A model that did compute with one (`hit * impact`) changes its result to what it meant.
- In the domain an `Event` is its own value type next to `Scalar`, `SampleSet` and `Probability`. `domain.Samples` returns the samples of a SampleSet or an Event, so code that needs samples does not care which one it gets.

## Consequences

- `hit = uniform(0, 1) < 0.1; loss = hit * (50k to 500k); P(loss > 0)` gives `0.10`; `binomial(1, p)` and the inline comparison keep working as before.
- `last` after a bare comparison of uncertain values is the event, so `last * 100` is a SampleSet of 0 and 100, not the probability times 100. `P(x > 5) * 100` remains the way to compute with the probability.
- A function whose body is a comparison now prints a probability instead of the summary of its 0/1 samples.
- The scenarios are unchanged: those that model an event already use the inline form or an expected count; `risk/annual-loss` models the expected annual loss, a different question.

## Alternatives considered

- **Document the old behaviour** and recommend `binomial(1, p)` or the inline comparison: no change, but the most natural way to write the model would stay silently wrong.
- **Keep the samples only in assignments**, with a top-level comparison still a plain probability and `:vars`/`--set` echo looking the binding up separately: smaller, but a variable would print as a probability when assigned and as a 0/1 summary when shown later, and every display path would need to know which names hold events.
- **Bind `last` to the probability** while assignments bind the event: keeps `last * 100` as before, but makes `last` the one variable that behaves differently from the rest; `P(...)` already covers that use.
