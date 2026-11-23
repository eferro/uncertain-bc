# Modelling guide

How to estimate well with `ubc`. The [tutorial](TUTORIAL.md) teaches the language; this guide is about the numbers you type into it: what a range claims, which shape to choose, how ranges combine, the mistakes that make a model lie, how to keep units straight without units in the language, how much to trust the simulation itself, and how to model something that may or may not happen.

Like the tutorial, every output below is checked against the binary by the test suite. Run the commands from the root of the repository.

## 1. What a 90 % interval claims

`a to b` says "I am 90 % sure the value is between `a` and `b`". `ubc` takes you at your word: 5 % of the samples fall below `a` and 5 % above `b`.

```console
$ ubc -e 'x = 10 to 20' -e 'P(x < 10)' -e 'P(x > 20)'
0.05
0.05
```

So a range is a claim you can check. If you give ten 90 % intervals for things you later find out, about nine should contain the true value. Most people, asked for the first time, get far fewer: their ranges are too narrow, because they write the range they *expect*, not the one they are 90 % *sure* of. Douglas Hubbard (*How to Measure Anything*) calls this being overconfident, and shows that it can be trained away.

How unlucky would a well-calibrated estimator have to be to hit only 7 of 10? `ubc` can tell you:

```console
$ ubc -e 'P(binomial(10, 0.9) <= 7)' -e 'P(binomial(10, 0.9) >= 9)'
0.07
0.74
```

Seven hits or fewer happens 7 % of the time to someone who is calibrated; nine or ten happens 74 % of the time. If you score 7 or less, widen your ranges.

Hubbard's tricks to get there:

- **Test yourself.** Answer ten trivia questions with 90 % intervals (the year a book came out, the length of a river) and count the hits. Repeat until you hit about nine.
- **The equivalent bet.** Would you rather win 1000 euros if the true value is inside your range, or by spinning a wheel that pays 90 % of the time? If you prefer the wheel, your range is too narrow; if you prefer your range, it may be too wide.
- **Start absurd.** Begin with ends that are certainly wrong (a house in Madrid costs between 1 euro and 1 billion) and move each end in until you start to feel unsure. Anchoring on a single best guess and adding a margin makes ranges too narrow.
- **Think about each end on its own.** "Is there a 5 % chance it is below `a`?" and "a 5 % chance it is above `b`?" are two separate questions.

## 2. Choosing a shape

`a to b` picks the shape from its ends, and that choice is usually right.

- **Both ends positive: lognormal.** The value cannot go below zero and the right tail is long. That suits durations, prices, sizes, counts and anything you would describe as "between 2 and 10 times".
- **An end at zero or below: normal.** Symmetric around the middle, which suits changes, differences and errors that can go either way.

Watch out for a zero at the lower end. `0 to 10` is a normal range, so it allows negative values:

```console
$ ubc -e 'cost = 0 to 10' -e 'P(cost < 0)'
0.05
```

A cost cannot be negative, so write the smallest value you actually believe in, such as `0.5 to 10`, and get a lognormal. For a growth rate or a change in percent, which can be negative, keep the normal and add it to 1:

```console
$ ubc -e 'growth = 1 + (-5% to 15%)' -e 'growth'
≈ 1.05   90% [0.95 … 1.15]   98% [0.9 … 1.2]   ▂▂▃▄▅▆▇██▇▆▅▄▃▂▂
```

When `a to b` is not enough:

- `pert(min, likely, max)` when you can name a most likely value and hard limits, as in three-point task estimates.
- `uniform(lo, hi)` only when every value in between really is equally likely, which is rare: a uniform says the ends are as likely as the middle.
- `[v1, v2, …]` or `--data` when you have past observations: let the data be the shape instead of guessing one.
- `poisson(lambda)` for a count of independent events in a period (incidents a month), `binomial(n, p)` for how many of `n` tries succeed.

## 3. Ranges are independent; variables are not

Every `a to b` draws its own samples, so two ranges never move together, even when they are written the same. A variable keeps its samples, so every use of it is the same quantity. That is how you say "these two things depend on the same unknown": name the unknown once and reuse it.

```ubc
# Two jobs billed at the same uncertain hourly rate, in euros.
rate = 40 to 90
shared = 100 * rate + 300 * rate                 # one rate, used twice
separate = 100 * (40 to 90) + 300 * (40 to 90)   # two independent guesses
shared
separate
P(shared > 30k)
P(separate > 30k)
```

```console
$ ubc examples/shared-rate.ubc
shared = ≈ 24k   90% [16k … 35.6k]   98% [13k … 43k]   ▂▄▅▆██▇▆▅▄▃▃▂▂▁▁
separate = ≈ 24.2k   90% [17.7k … 33.7k]   98% [15k … 39k]   ▂▄▅▇██▇▇▆▅▃▃▂▂▁▁
P(shared > 30k) = 0.18
P(separate > 30k) = 0.13
```

Both have the same middle, but with one shared rate a high rate makes both jobs expensive at once, so the total is wider and blowing a 30k budget is more likely. Treating related quantities as independent lets their errors cancel out, which makes the model look more certain than it is. Before writing the same range twice, ask whether it is really two separate unknowns.

A function body draws fresh samples on every call (`task(5) + task(3)` gets two independent errors). If the error is shared, for example because the whole team is optimistic, pass it in as a variable instead.

## 4. Typical mistakes

**Precedence of `to`.** `to` binds looser than arithmetic, so `8 to 25 * 12` is `8 to 300`, not a year of a monthly bill:

```console
$ ubc -e '8 to 25 * 12' -e '(8 to 25) * 12'
≈ 49.2   90% [7.99 … 286]   98% [3.7 … 630]   █▅▃▂▂▁▁▁▁▁▁▁▁▁▁▁
≈ 169   90% [96.9 … 298]   98% [76 … 380]   ▃▅▇██▇▆▅▄▃▃▂▂▂▁▁
```

Put parentheses around every range you compute with.

**Double counting.** Each uncertainty should enter the model once. If every task range already includes a margin for surprises, adding a 30 % contingency on top counts the surprises twice; if your estimates are clean, the contingency is a separate range (`forgotten = 10% to 40%` in the tutorial's project estimate). The same happens when two lines of a cost model both include, say, the hosting of the database.

**Mixing units.** `ubc` has no units, so it adds hours to minutes without complaint:

```console
$ ubc -e 'hours = 2 to 6' -e 'minutes = 30 to 90' -e 'hours + minutes' -e 'hours + minutes / 60'
≈ 55.4   90% [33.6 … 93.1]   98% [28 … 116]   ▃▅▇██▇▇▅▄▃▃▂▂▂▁▁
≈ 4.4   90% [2.83 … 6.82]   98% [2.3 … 8.5]   ▂▄▆▇██▇▅▄▃▃▂▂▁▁▁
```

Pick one unit per quantity, convert on the way in, and say it in the variable's comment or name (`hours`, `cost_eur`, `latency_ms`); [section 5](#5-units) shows how.

**Mean vs median.** The `≈` value is the median: half the outcomes are below it. For a lognormal the mean is higher, and for a wide one much higher:

```console
$ ubc -e 'x = 1k to 1M' -e 'x' -e 'median(x)' -e 'mean(x)' -e 'P(x > mean(x))'
≈ 31.9k   90% [997 … 916k]   98% [230 … 4.1M]   █▂▁▁▁▁▁▁▁▁▁▁▁▁▁▁
31913.40476
265343.5222
0.15
```

The mean is eight times the median, and only 15 % of outcomes reach it. Use the median for "what will probably happen" and a percentile, such as `p(x, 0.9)`, for a commitment. Use the mean when the same bet repeats many times and the outcomes add up (the yearly cost of many incidents, the budget of many projects): over many draws, totals approach the mean, not the median.

**Heavy tails.** A range whose ends are a factor of a thousand apart, like `1k to 1M`, says you do not know the order of magnitude; its 98 % interval reaches past 4M. Multiplying several wide ranges makes the tail longer still. That may be true, and then it is the most important thing the model tells you: look at the 98 % interval and at `P(x > limit)`, not only at the middle. If it is not true, the range is too wide: find out more about the quantity that drives it (`--sens` ranks them).

**A probability where you meant an event.** `P(...)` is one number, the same in every sample. Multiplying an impact by it spreads a tenth of the cost over every sample, which is the expected cost, not what happens:

```console
$ ubc --echo -e 'p_hit = P(uniform(0, 1) < 0.1)' -e 'loss = p_hit * (50k to 500k)' -e 'P(loss > 0)'
p_hit = 0.10
loss = ≈ 15.5k   90% [5010 … 49k]   98% [3100 … 79k]   ▆█▇▅▄▃▂▂▂▂▁▁▁▁▁▁
1
```

`loss` is never 0, and never as large as a real incident. To say "10 % of the time it happens, and then it costs 50k to 500k", keep the event: `hit = uniform(0, 1) < 0.1`, without `P`. [Section 7](#7-events-that-may-or-may-not-happen) shows how.

## 5. Units

`ubc` has no units, by decision ([ADR 0065](adr/0065-no-units-in-the-language.md)): every number is a plain number, and keeping track of what it measures is up to the model. Three habits do most of the work.

**Name the unit.** Put it in the variable's name or in its comment: `delay_weeks`, `price_eur_kwh`, `latency_ms`, `rate = 40 to 90   # euros an hour`. A reader, and you in three months, should not have to guess.

**One unit per quantity, one period per model.** Convert once, where the number comes in, not inside every formula. Keep rates and durations in the same period: a monthly cost with a number of months, a weekly throughput with a deadline in weeks.

**Convert with named constants, not magic numbers.** The shared formulas that every scenario loads (`--scenario formulas`, see the [scenario catalogue](../scenarios/README.md#reusable-formulas)) start with the usual conversions. Each one says whether it counts working time, what people put in, or calendar time, what a clock, a bill or an SLO window counts:

| Constant | Value | Counts |
|---|---|---|
| `minutes_per_hour` | 60 | both |
| `hours_per_day` | 8 | working hours in a working day |
| `hours_per_calendar_day` | 24 | calendar |
| `days_per_week` | 5 | working days |
| `days_per_month` | 30 | calendar, the month of a bill or an SLO |
| `working_days_per_year` | 220 | after weekends, public holidays and leave |
| `working_weeks_per_year` | 44 | `working_days_per_year / days_per_week` |
| `weeks_per_year` | 52 | calendar |
| `months_per_year` | 12 | both |

A daily stand-up of six people, in hours a year:

```console
$ ubc --scenario formulas -e 'standup_minutes = 10 to 25' -e 'people = 6' -e 'yearly_hours = standup_minutes / minutes_per_hour * people * working_days_per_year' -e 'yearly_hours'
≈ 348   90% [220 … 544]   98% [180 … 670]   ▂▄▆▇██▇▅▅▄▃▂▂▂▁▁
```

Working and calendar time are far apart, so pick the one the quantity really counts. A year of a server that is always on against a year of one engineer:

```console
$ ubc --scenario formulas -e 'months_per_year * days_per_month * hours_per_calendar_day' -e 'working_days_per_year * hours_per_day'
8640
1760
```

The constants are ordinary variables, so `--set` adapts them to your organisation, and the ones computed from them follow:

```console
$ ubc --set 'working_days_per_year = 200' --scenario formulas -e 'working_weeks_per_year'
40
```

## 6. Monte Carlo noise

`ubc` computes with 5000 random samples, so every result carries a little noise of its own, on top of the uncertainty you typed. The fixed seed hides it: the same model prints the same numbers on every run. Changing the seed shows how much of a result is noise:

```console
$ ubc --seed 1 -e 'p(20 to 45, 0.99)'
53.07508159
$ ubc --seed 2 -e 'p(20 to 45, 0.99)'
51.91173463
$ ubc --seed 3 -e 'p(20 to 45, 0.99)'
53.3506105
```

The middle and the 90 % interval are stable with 5000 samples; the far tails rest on few samples (the 99th percentile on about 50), so they move. That is why the 98 % interval is printed with fewer digits. Raise `--samples` when you need a tail, or the probability of a rare event, to be steadier:

```console
$ ubc -n 100000 --seed 1 -e 'p(20 to 45, 0.99)'
53.37301437
$ ubc -n 100000 --seed 2 -e 'p(20 to 45, 0.99)'
53.19856213
$ ubc -n 100000 --seed 3 -e 'p(20 to 45, 0.99)'
52.91262469
```

Twenty times more samples makes the noise about four and a half times smaller (it shrinks with the square root of the count). Raise the count when:

- a decision depends on a tail, such as `p(x, 0.99)` or a 98 % end;
- a probability is small, like `P(x > limit)` around 0.01, where 5000 samples hold only about 50 hits;
- two options are close and you want to know whether the difference is real or noise.

A quick check before trusting a close call: run with two or three seeds (or `--random`). If the answer changes, raise `--samples` or decide that the options are, as far as you know, equivalent. No sample count can fix a range you got wrong: the noise of the simulation is almost always far smaller than the uncertainty in your inputs.

## 7. Events that may or may not happen

Many questions are "if it happens, it costs X": an incident, a missed date, a risk, an experiment that wins. Name the event with a comparison, then multiply by what it costs when it happens:

```console
$ ubc --echo -e 'hit = uniform(0, 1) < 0.1' -e 'loss = hit * (50k to 500k)' -e 'P(loss > 0)' -e 'mean(loss)' -e 'p(loss, 0.95)'
hit = 0.10
loss = ≈ 0   90% [0 … 155k]   98% [0 … 410k]   █▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁
0.10
21032.27932
155428.0518
```

`uniform(0, 1) < 0.1` holds in about 10 % of the samples. A comparison of uncertain values is an *event*: `hit` keeps its 0/1 outcome in every sample and only prints as the probability that it happens, `0.10`. So `loss` is 0 in the 90 % of samples where nothing happened and the whole impact in the rest, and `P(loss > 0)` is about `0.1`. The median of `loss` is 0, because most of the time nothing happens; the mean, about 21k, is the expected loss, the figure to budget when the same risk repeats over many years or many projects; the 95th percentile, about 155k, is a bad year.

Three ways to write the event give the same samples:

- **A named comparison**, `hit = uniform(0, 1) < 0.1`, as above. It reads as the question you are asking, prints as a probability, and the chance can be uncertain itself: `uniform(0, 1) < (5% to 25%)`.
- **`binomial(1, p)`**, `hit = binomial(1, 0.1)`: one trial, 1 with probability `p`. It is an ordinary uncertain number, so it prints as a summary (`≈ 0   90% [0 … 1]`) instead of a probability.
- **Inline**, `(uniform(0, 1) < 0.1) * (50k to 500k)`, when the event has no name of its own.

```console
$ ubc -e 'chance = 5% to 25%' -e 'hit = uniform(0, 1) < chance' -e 'loss = hit * (100k to 1M)' -e 'P(loss > 0)' -e 'mean(loss)' -e 'P(loss > 500k)'
0.13
51353.26035
0.03
```

What not to write is `P(...) * impact`: `P` turns the event into one number, and the loss becomes a fraction of the impact in every sample ([typical mistakes](#4-typical-mistakes)). Use the probability when you want the probability (`P(loss > 0)`), and the event when you compute what happens.

For something that happens several times a year, an event that happens or not is the wrong shape: count the occurrences with `poisson(rate)`, or multiply by the expected count as the [incident cost scenario](../scenarios/reliability/incident-cost.ubc) does.

## Where next

- The [tutorial](TUTORIAL.md) covers the whole language step by step.
- The [scenario catalogue](../scenarios/README.md) has ready-to-use models; each one says where its ranges come from.
- [docs/SPEC.md](SPEC.md) describes every flag and function.
