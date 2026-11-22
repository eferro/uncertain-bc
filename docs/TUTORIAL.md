# Tutorial

A step-by-step guide to `ubc`, from plain arithmetic to models with ranges, functions and probabilities, ending with three real-world estimates.

Every example is a file in [`examples/`](../examples/), and every output below is checked against the binary by the test suite, so what you read is what `ubc` prints. Run the commands from the root of the repository (or download the files) after [installing `ubc`](../README.md#install). Results are the same on every run because `ubc` uses a fixed random seed.

## 1. Plain arithmetic

`ubc` works like `bc`: give it expressions and it prints their values.

```ubc
# Plain arithmetic, like bc.
2 + 3 * 4
(2 + 3) * 4
2 ^ 10
1.5k * 12          # suffixes: k, M, B, T
15% * 340          # and percentages
```

```console
$ ubc examples/arithmetic.ubc
2 + 3 * 4 = 14
(2 + 3) * 4 = 20
2 ^ 10 = 1024
1.5k * 12 = 18000
15% * 340 = 51
```

A file prints one line per expression, labelled with the expression itself: its variable name, or its text when that is short (a longer one says its line instead, `L12:`). `#` starts a comment. For a one-off calculation, skip the file: `ubc -e '2 ^ 10'`.

## 2. Ranges

When you do not know a number, give a range instead. `a to b` means "I am 90 % sure the value is between `a` and `b`".

```ubc
# "a to b" means: 90% sure the value is between a and b.
8 to 25            # a monthly VPS bill, in euros
(8 to 25) * 12     # the same bill over a year
-5 to 5            # a change that can go either way
```

```console
$ ubc examples/ranges.ubc
8 to 25 = ≈ 14.2   90% [8 … 24.6]   98% [6.2 … 32]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
(8 to 25) * 12 = ≈ 169   90% [96.9 … 298]   98% [76 … 380]   ▃▅▇██▇▆▅▄▃▃▂▂▂▁▁
-5 to 5 = ≈ -0.0367   90% [-4.94 … 5.2]   98% [-6.9 … 7.3]   ▂▂▃▄▅▆▇█▇▇▆▅▃▃▂▂
```

`ubc` draws 5000 random samples for every range and computes with all of them (a Monte Carlo simulation). Each uncertain result is one line:

- `≈ 169`: the median, the value in the middle.
- `90% [96.9 … 298]`: 90 % of the samples fall in this interval.
- `98% [76 … 380]`: a wider interval that holds 98 % of the samples, so you can see how far the tails reach. Its ends rest on few samples, so they are rounded more coarsely and outwards: the 98% interval never ends inside the 90% one.
- the sparkline: the shape of the distribution, from low to high values.

A range with both ends positive is lognormal: values cannot go below zero and the right tail is longer, which suits prices, durations and sizes. A range that includes zero or negative numbers is normal and symmetric. Write parentheses around a range before using it: `to` binds looser than arithmetic, so `8 to 25 * 12` means `8 to 300`.

## 3. Variables

Name a value to use it later. A variable keeps its samples, so it stays the same uncertain quantity wherever you use it.

```ubc
# A variable keeps its samples, so it stays one uncertain quantity.
price = 8 to 25
yearly = price * 12
yearly
yearly - price * 12     # exactly 0: same samples on both sides
(8 to 25) - (8 to 25)   # two independent guesses do not cancel
```

```console
$ ubc --echo examples/variables.ubc
price = ≈ 14.2   90% [8 … 24.6]   98% [6.2 … 32]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
yearly = ≈ 170   90% [96 … 296]   98% [75 … 380]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
yearly = ≈ 170   90% [96 … 296]   98% [75 … 380]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
yearly - price * 12 = 0
(8 to 25) - (8 to 25) = ≈ -0.0758   90% [-12.6 … 12]   98% [-21 … 19]   ▁▂▂▃▄▅▇██▇▆▄▃▂▂▁
```

Assignments print nothing by default; `--echo` prints them too, so a model reads like a report. Every range draws new samples: the two `8 to 25` on the last line are independent guesses, and their difference is anything between about -12 and 12.

## 4. Distributions

`a to b` covers most guesses. When you know more about the shape, use a distribution.

```ubc
# When a 90% interval is not the right shape, pick a distribution.
normal(100, 15)               # mean and standard deviation
uniform(10, 20)               # anything in between is equally likely
pert(3, 5, 12)                # minimum, most likely, maximum
[12, 15, 15, 18, 31]          # past observations, resampled
```

```console
$ ubc examples/distributions.ubc
normal(100, 15) = ≈ 100   90% [75.3 … 124]   98% [64 … 135]   ▂▂▃▄▅▆▇██▇▆▅▄▃▂▂
uniform(10, 20) = ≈ 14.9   90% [10.5 … 19.5]   98% [10 … 20]   ██▇█▇█▇█▇▇▇▇▇█▇█
pert(3, 5, 12) = ≈ 5.65   90% [3.59 … 8.64]   98% [3.2 … 9.9]   ▄▆▇▇█▇▇▇▆▅▅▄▃▃▂▂
[12, 15, 15, 18, 31] = ≈ 15   90% [12 … 31]   98% [12 … 31]   ▅▁█▁▁▅▁▁▁▁▁▁▁▁▁▅
```

`pert` is the classic three-point estimate: it leans towards the most likely value. A list of numbers is empirical data: each sample picks one of them at random. The [specification](SPEC.md#5-built-in-functions) lists every distribution (`lognormal`, `triangular`, `beta`, `exponential`, `gamma`).

## 5. Functions

Define a function when the same calculation appears more than once.

```ubc
# A function names a calculation you repeat.
task(days) = days * (0.8 to 1.6)   # every guess is a bit off, either way
task(5)
task(5) + task(3) + task(8)
```

```console
$ ubc examples/functions.ubc
task(5) = ≈ 5.66   90% [4 … 7.93]   98% [3.4 … 9.3]   ▂▄▅▆▇██▇▅▅▃▃▂▂▁▁
task(5) + task(3) + task(8) = ≈ 18.3   90% [14.8 … 22.8]   98% [13 … 25]   ▂▃▄▅▆██▇▆▆▄▃▃▂▂▂
```

A range inside a function draws fresh samples on every call, so each task gets its own error. That is why the sum of three tasks (≈ 18.3) is relatively narrower than a single one: the errors partly cancel out.

## 6. Probabilities and summaries

A range answers "how much?". To decide, you usually need "how likely?" and "what number do I commit to?".

```ubc
# Ask questions about an uncertain value.
delivery = 20 to 45         # days
P(delivery > 40)            # how likely is it to slip past 40 days?
delivery < 25               # a bare comparison is a probability too
median(delivery)
round(p(delivery, 0.9))     # 90th percentile, rounded: a safe commitment
mean(delivery)
```

```console
$ ubc examples/probability.ubc
P(delivery > 40) = 0.12
delivery < 25 = 0.23
median(delivery) = 30.03223521
round(p(delivery, 0.9)) = 41
mean(delivery) = 30.87175467
```

`P(...)` is the share of samples where the condition holds. The summaries `median`, `mean`, `stdev`, `min`, `max` and `p(x, q)` (the `q` quantile) turn a distribution into a plain number, ready for `round`, `floor` or `ceil`. Note how the mean is above the median: the right tail of a lognormal pulls it up.

## 7. A project estimate

Now everything together. A team estimates each task of a feature, adds a share for the work nobody listed, and asks whether the deadline is safe.

```ubc
# How long will the new billing feature take? Working days.
# Estimate each task as a 90% interval, then add a share for the
# things nobody listed.

design = 2 to 5
backend = 8 to 20
frontend = 5 to 12
migration = 1 to 6
testing = 3 to 8
work = design + backend + frontend + migration + testing

forgotten = 10% to 40%     # unplanned work, as a share of the plan
days = work * (1 + forgotten)

days
P(days > 50)               # the deadline is 50 working days away
round(p(days, 0.85))       # days we can promise with 85% confidence
```

```console
$ ubc examples/project-estimate.ubc
days = ≈ 38.8   90% [29.6 … 51.8]   98% [26 … 60]   ▂▄▅▇███▇▆▅▃▃▂▂▂▁
P(days > 50) = 0.07
round(p(days, 0.85)) = 46
```

The sum of the most likely values would say about 30 days. The model says 39 is the middle, there is a 7 % chance of missing the 50-day deadline, and 46 days is a promise the team can keep 85 % of the time.

## 8. A Fermi problem

A Fermi problem estimates something you cannot look up by chaining quantities you can guess. The classic one: how many piano tuners work in a city?

```ubc
# Fermi problem: how many piano tuners work in Madrid?

people = 3M to 3.6M
per_household = 2.3 to 2.7
households = people / per_household
with_piano = 1% to 5%           # households that own a piano
pianos = households * with_piano + (500 to 3k)  # plus schools, bars, halls

tunings = 0.5 to 1.5            # per piano and year
per_day = 2 to 4                # tunings a tuner can do in a day
workdays = 200 to 240
per_tuner = per_day * workdays

tuners = pianos * tunings / per_tuner
tuners
```

```console
$ ubc --echo examples/fermi.ubc
people = ≈ 3.29M   90% [3M … 3.59M]   98% [2.88M … 3.74M]   ▂▂▃▄▅▆██▇▆▆▄▃▂▂▂
per_household = ≈ 2.49   90% [2.3 … 2.7]   98% [2.22 … 2.79]   ▂▂▄▅▆▇████▆▅▃▃▂▂
households = ≈ 1.32M   90% [1.17M … 1.49M]   98% [1.11M … 1.57M]   ▂▂▄▄▆▇█▇▆▆▅▄▃▂▂▂
with_piano = ≈ 0.0222   90% [0.0101 … 0.0516]   98% [0.0074 … 0.072]   ▄▇██▆▅▄▃▃▂▂▂▁▁▁▁
pianos = ≈ 30.9k   90% [14.6k … 69.3k]   98% [10k … 101k]   ▄▇█▇▆▅▄▃▂▂▂▁▁▁▁▁
tunings = ≈ 0.863   90% [0.493 … 1.48]   98% [0.39 … 1.9]   ▃▅▆▇██▆▅▅▃▃▂▂▂▁▁
per_day = ≈ 2.84   90% [1.98 … 4.07]   98% [1.7 … 4.7]   ▂▄▅▇███▇▆▅▃▃▂▂▂▁
workdays = ≈ 220   90% [200 … 241]   98% [192 … 251]   ▂▂▄▅▇▇███▇▆▅▃▃▂▂
per_tuner = ≈ 622   90% [431 … 901]   98% [370 … 1050]   ▃▄▆▇██▇▇▅▄▃▃▂▂▂▁
tuners = ≈ 42.8   90% [15.8 … 118]   98% [10 … 190]   ▅██▆▅▃▃▂▂▂▁▁▁▁▁▁
≈ 42.8   90% [15.8 … 118]   98% [10 … 190]   ▅██▆▅▃▃▂▂▂▁▁▁▁▁▁
```

With `--echo`, every step shows its own uncertainty, so you can see which guess drives the result: the share of households with a piano spans a factor of five, more than anything else. That is the number worth researching first. (A file with a single expression prints it without a label.)

## 9. A cloud bill, with scenarios and checks

The last model estimates a monthly cloud bill and asks how often it would exceed a 300-dollar budget.

```ubc
# Monthly cloud bill for a small web app, in dollars.

users = 20k to 60k              # monthly active users
requests = users * (50 to 200)  # requests per user and month
per_million = 0.2 to 0.6        # serverless price per million requests
compute = requests / 1M * per_million

storage_gb = 200 to 500
storage = storage_gb * 0.023

egress_gb = users * (20 to 80) / 1k   # megabytes per user, in gigabytes
egress = egress_gb * 0.09

database = 50 to 120            # managed database, depends on the plan
bill = compute + storage + egress + database

bill
P(bill > 300)                   # how often would we blow the budget?
```

```console
$ ubc examples/cloud-cost.ubc
bill = ≈ 216   90% [130 … 387]   98% [100 … 500]   ▃▅▇█▇▆▅▄▃▂▂▂▂▁▁▁
P(bill > 300) = 0.17
```

What if the launch goes well? `--set` runs a scenario without editing the file: the model's own assignment to `users` is skipped and everything computed from it uses the new value.

```console
$ ubc --set 'users=100k to 200k' examples/cloud-cost.ubc
bill = ≈ 600   90% [321 … 1170]   98% [250 … 1600]   ▄▆██▇▆▅▄▃▃▂▂▂▁▁▁
P(bill > 300) = 0.96
```

`--assert` turns the model into a check, for example in CI when someone changes a price. It prints nothing when the condition holds, and otherwise explains the failure on stderr and exits with code 4.

```console
$ ubc --assert 'P(bill > 300) < 0.2' examples/cloud-cost.ubc
bill = ≈ 216   90% [130 … 387]   98% [100 … 500]   ▃▅▇█▇▆▅▄▃▂▂▂▂▁▁▁
P(bill > 300) = 0.17
```

With the `--set` scenario above, the same assertion fails: the budget is blown 96 % of the time.

## Where next

- The [scenario catalogue](../scenarios/README.md) has ready-to-use models for software engineering and product questions, such as the cost of replacing a person: change the assumptions and run them.
- Run `ubc` with no arguments for an interactive session: assignments print as you type, and `:hist`, `:plot` and `:save` draw histograms, export charts and save the session as a model.
- `--json` and `--plain` make the results easy to use from scripts; `--hist` draws a histogram under every result; `--interval 80` prints an 80 % interval instead of the 90 % one; `--quantiles 10,50,90` shows other percentiles.
- [docs/SPEC.md](SPEC.md) describes the whole language and every flag.
