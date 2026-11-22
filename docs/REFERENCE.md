# Reference

Everything you can type into `ubc`: the language, the built-in functions, the flags, the REPL commands, the output formats, the exit codes and the common errors with how to fix them. To learn the language step by step, read the [tutorial](TUTORIAL.md); to estimate well, the [modelling guide](MODELLING.md).

Every command shown in a console block here is run by the test suite, which compares its output with this page, so what you read is what `ubc` prints. Lines that start with `ubc: ` go to stderr; the rest go to stdout. The examples use the default seed and sample count, so you get the same numbers.

## Running `ubc`

```
ubc                    interactive session (REPL), when stdin is a terminal
ubc FILE...            run the files in order
ubc -e EXPR...         run the expressions in order
cmd | ubc              run the program read from stdin
```

Flags go before files: anything after the first file name is read as another file. `-e` expressions run before files, all in one session, so a variable defined by one is seen by the next:

```console
$ ubc -e 'x = 3' -e 'x * 2'
6
```

Files use the extension `.ubc` and the same language as the REPL, without the `:` commands. The [examples](../examples/) are models you can run from the repository root:

```console
$ ubc examples/homelab.ubc
homelab = ≈ 66.5   90% [47 … 95.7]   98% [40 … 111]   ▂▃▅▆▇█▇▆▅▄▃▃▂▂▂▁
vps = ≈ 169   90% [94.7 … 295]   98% [75 … 380]   ▃▅▇▇█▇▆▅▄▃▃▂▂▂▁▁
saving = ≈ 101   90% [21.5 … 229]   98% [-1.5 … 310]   ▃▅▆▇█▇▇▆▄▃▃▂▂▂▁▁
```

## Language

### Numbers

Numbers are written `42`, `3.14`, `.5`, `1e6` or `2.5e-3`. A suffix right after the number, without a space, scales it: `k` (thousand), `M` (million), `B` (billion), `T` (trillion) and `%` (hundredth).

```console
$ ubc -e '1.5e3' -e '2.5M' -e '5% * 200' -e '.5k'
1500
2500000
10
500
```

### Operators

`+ - * / ^` and parentheses, with the usual precedence. `^` is right-associative and binds tighter than a leading minus, so `-2 ^ 2` is `-(2 ^ 2)`:

```console
$ ubc -e '2 * 3 ^ 2' -e '-2 ^ 2' -e '2 ^ 3 ^ 2'
18
-4
512
```

From loosest to tightest: comparisons, `to`, `+ -`, `* /`, leading `-`, `^`.

### Ranges: `a to b`

`a to b` means "90 % sure the value is between `a` and `b`". When both ends are positive the range is lognormal, skewed to the right like most costs and durations; otherwise it is normal and symmetric:

```console
$ ubc -e '10 to 20'
≈ 14.2   90% [10 … 19.8]   98% [8.6 … 24]   ▂▄▅▆▇██▇▅▅▃▃▂▂▁▁
$ ubc -e '-10 to 10'
≈ 0.0265   90% [-10 … 9.74]   98% [-15 … 15]   ▂▂▃▄▅▆▇██▇▆▅▄▃▂▂
```

The start must be below the end, and both ends must be plain numbers, not other ranges. `to` binds looser than arithmetic, so `1 + 2 to 10` is the range `3 to 10`; write parentheses to add to a range:

```console
$ ubc -e '1 + 2 to 10'
≈ 5.49   90% [3 … 9.85]   98% [2.3 … 13]   ▃▅▇██▇▆▅▄▃▂▂▂▁▁▁
$ ubc -e '1 + (2 to 10)'
≈ 5.48   90% [3 … 10.8]   98% [2.4 … 15]   ▄▇██▇▅▄▄▃▂▂▂▁▁▁▁
```

### Variables

`name = expression` binds a name; names are letters, digits and `_`, start with a letter or `_`, and are case-sensitive. A variable keeps its samples, so it is correlated with itself, while every range is a fresh, independent one:

```console
$ ubc -e 'x = 1 to 5' -e 'x - x'
0
$ ubc -e '(1 to 5) - (1 to 5)'
≈ -0.0123   90% [-2.88 … 2.89]   98% [-4.9 … 5]   ▁▁▂▂▃▄▆█▇▆▄▃▂▂▁▁
```

Three names are built in: `pi`, `e` and `last`, the result of the previous expression.

```console
$ ubc -e 'pi' -e 'e'
3.141592654
2.718281828
$ ubc -e '2 + 3' -e 'last * 10'
5
50
```

### Statements and comments

Statements are separated by new lines or `;`. `#` starts a comment that runs to the end of the line.

```console
$ ubc -e 'a = 2; b = 3; a * b   # six'
6
```

### Functions of your own

`name(p1, p2, ...) = expression` defines a function. Parameters hide variables with the same name; other variables are read when the function is called. A range in the body draws fresh samples on every call.

```console
$ ubc -e 'cost(n) = n * (2 to 4)' -e 'cost(10)'
≈ 28.3   90% [20 … 39.6]   98% [17 … 47]   ▂▄▅▆▇██▇▅▅▃▃▂▂▁▁
```

You cannot redefine a built-in function, and since there are no conditionals a function that calls itself never ends: calls nested more than 64 deep are an error.

### Observed data: `[...]`

A list of plain numbers is resampled with replacement, so it behaves like any other uncertain value built from what you saw:

```console
$ ubc -e '[3, 5, 8, 13]'
≈ 5   90% [3 … 13]   98% [3 … 13]   █▁▁█▁▁▁▁█▁▁▁▁▁▁█
```

### Comparisons and probabilities

`< > <= >= == !=` compare two values. Between plain numbers they give `1` (true) or `0` (false); with an uncertain value they give, sample by sample, how often the comparison holds, printed as a probability. `P(...)` turns a comparison into a number you can compute with:

```console
$ ubc -e '2 > 1' -e '1 == 2'
1
0
$ ubc -e 'x = 10 to 20' -e 'x > 15'
0.39
$ ubc -e 'x = 10 to 20' -e 'P(x > 15) * 100'
38.54
```

## Built-in functions

Function names are case-sensitive. Arguments of distributions must be plain numbers.

### Distributions

| Function | Draws |
|---|---|
| `normal(mean, sd)` | Normal; `sd` positive |
| `lognormal(mu, sigma)` | Lognormal with `mu` and `sigma` of the logarithm; median `exp(mu)` |
| `uniform(lo, hi)` | Every value between `lo` and `hi` equally likely |
| `triangular(lo, mode, hi)` | Triangle with its peak at `mode` |
| `pert(lo, mode, hi)` | Beta-PERT, smoother than the triangle; mean `(lo + 4·mode + hi) / 6` |
| `beta(a, b)` | Beta, between 0 and 1, for rates and shares |
| `exponential(rate)` | Waiting times; mean `1 / rate` |
| `gamma(shape, scale)` | Gamma; mean `shape · scale` |
| `poisson(lambda)` | Whole counts of events with mean `lambda` |
| `binomial(n, p)` | Successes in `n` trials of probability `p` |

```console
$ ubc -e 'normal(100, 15)'
≈ 100   90% [75.3 … 124]   98% [64 … 135]   ▂▂▃▄▅▆▇██▇▆▅▄▃▂▂
$ ubc -e 'pert(1, 2, 6)'
≈ 2.37   90% [1.3 … 4.09]   98% [1.1 … 4.8]   ▅▇▇▇█▇▆▇▅▅▄▃▃▂▂▂
$ ubc -e 'poisson(4)'
≈ 4   90% [1 … 8]   98% [0 … 9]   ▂▃▁▆▁█▁█▇▁▅▁▃▁▂▂
```

### Summaries

They turn an uncertain value into a plain number, and return a plain number unchanged.

| Function | Returns |
|---|---|
| `mean(x)`, `median(x)` | The mean and the median of the samples |
| `stdev(x)` | The population standard deviation |
| `min(x)`, `max(x)` | The smallest and the largest sample |
| `p(x, q)` | The `q` quantile, `q` between 0 and 1: `p(x, 0.9)` is the 90th percentile |
| `P(cond)` | The probability that `cond`, a comparison, holds |

```console
$ ubc -e 'x = 10 to 20' -e 'mean(x)' -e 'median(x)' -e 'p(x, 0.9)'
14.43920894
14.15512332
18.40096995
$ ubc -e 'P(normal(0, 1) > 1)'
0.15
```

### Mathematics

Applied sample by sample: `sqrt(x)`, `exp(x)`, `log(x)` (natural), `log10(x)`, `abs(x)`, `round(x)`, `floor(x)`, `ceil(x)`. With two or more arguments, `min(a, b, ...)` and `max(a, b, ...)` also work sample by sample, so `max(0, x)` cuts the negative part of `x`:

```console
$ ubc -e 'round(p(20 to 40, 0.85))'
35
$ ubc -e 'min(3, 7, 5)'
3
$ ubc -e 'max(0, -5 to 5)'
≈ 0.0132   90% [0 … 4.87]   98% [0 … 7.1]   █▂▂▂▂▂▂▁▁▁▁▁▁▁▁▁
```

## Flags

```console
$ ubc --help
Usage:
  ubc [flags] [FILE...]

Evaluates expressions with uncertain numbers. Without -e or files, it reads
the program from stdin. Flags go before files.

Flags:
  -e, --eval EXPR     evaluate EXPR (repeatable)
      --set NAME=EXPR set NAME to EXPR, overriding the model's value (repeatable)
      --data NAME=FILE
                      set NAME to the numbers in FILE, one per line (repeatable)
      --assert EXPR   after the program, exit 4 unless EXPR certainly holds (repeatable)
      --samples-out FILE
                      after the program, write every variable's samples to FILE (.csv or .json)
      --sens EXPR     after the program, rank the variables that drive EXPR's uncertainty
      --scenarios     list the built-in scenarios with the question each one answers
      --scenario NAME run a built-in scenario, e.g. team/replace-a-person, after the
                      shared formulas and before -e and files; 'formulas' loads only them
      --scenario-source NAME
                      print a built-in scenario's source, to copy and edit it
  -n, --samples N     Monte Carlo sample count (default 5000)
  -s, --seed N        random seed (default 1)
      --random        seed from the clock
      --ascii         print ASCII instead of Unicode
      --plain         print only the central value (median), for scripts
      --json          print one JSON object per result
      --echo          also print assignments, with their names
      --hist          draw a histogram under every uncertain result
      --watch         evaluate the files again whenever they change
      --interval W    print a W% interval instead of the 90% one (default 90)
      --quantiles L   show these percentiles instead of the intervals (10,50,90)
  -h, --help          show this help
      --version       print the version
      --completion S  print a completion script for bash, zsh or fish
```

### Scenarios: `--set` and `--data`

`--set NAME=EXPR` replaces the model's own value of `NAME`, and everything computed from it follows. Quote the value when it has spaces:

```console
$ ubc --set 'vps = 5 to 10' examples/homelab.ubc
homelab = ≈ 67   90% [46.6 … 96.6]   98% [40 … 115]   ▂▄▆▇██▇▆▅▄▃▂▂▂▁▁
vps = ≈ 7.08   90% [5 … 9.91]   98% [4.3 … 11.6]   ▂▄▅▆▇██▇▅▅▃▃▂▂▁▁
saving = ≈ -59.7   90% [-89.3 … -39.1]   98% [-108 … -33]   ▁▁▂▂▂▃▄▅▇▇██▇▅▄▂
```

`--data NAME=FILE` does the same with observed data: `FILE` holds one plain number per line (no suffixes; blank lines are ignored), resampled like `[...]`. For example, `ubc --data lead=lead_times.txt model.ubc`. When both name the same variable, `--set` wins.

### Checks for CI: `--assert`

`--assert EXPR` checks a condition after the program. It prints nothing when the condition holds in every sample; otherwise it explains the failure on stderr and `ubc` exits with `4`. Use a probability to state a looser check:

```console
$ ubc --assert 'P(x > 3) < 0.5' -e 'x = 1 to 5'
$ ubc --assert 'P(x > 3) < 0.1' -e 'x = 1 to 5'
ubc: --assert P(x > 3) < 0.1 failed: 0.27 < 0.1
```

### Analysis: `--sens` and `--samples-out`

`--sens EXPR` ranks the uncertain variables by how strongly they move with `EXPR` (Spearman's rank correlation), strongest first. A positive number means `EXPR` grows with the variable:

```console
$ ubc --sens saving examples/homelab.ubc
homelab = ≈ 66.5   90% [47 … 95.7]   98% [40 … 111]   ▂▃▅▆▇█▇▆▅▄▃▃▂▂▂▁
vps = ≈ 169   90% [94.7 … 295]   98% [75 … 380]   ▃▅▇▇█▇▆▅▄▃▃▂▂▂▁▁
saving = ≈ 101   90% [21.5 … 229]   98% [-1.5 … 310]   ▃▅▆▇█▇▇▆▄▃▃▂▂▂▁▁
sensitivity of saving (rank correlation):
  vps       +0.96  ███████████████████
  homelab   -0.26  █████
  power     -0.19  ████
  watts     -0.14  ███
  life      +0.13  ███
  price     -0.13  ███
  hardware  -0.10  ██
```

`--samples-out FILE` writes the raw samples of every variable to `FILE`, one column per variable (`.csv`) or one array per variable (`.json`), for a spreadsheet, pandas or R. Row *i* holds sample *i* of every variable, so correlations survive. The results still go to stdout.

### Built-in scenarios: `--scenarios`, `--scenario`, `--scenario-source`

`ubc --scenarios` lists the [scenario catalogue](../scenarios/README.md) built into the binary, one per line with the question it answers. `--scenario NAME` runs one, after the shared formulas and before any `-e` or file, so every other flag applies to it:

```console
$ ubc --scenario team/replace-a-person
total = ≈ 61.1k   90% [37.8k … 100k]   98% [30k … 122k]   ▃▅▆██▇▇▆▅▄▃▂▂▂▁▁
P(total > 50k) = 0.75
```

`ubc --scenario-source team/replace-a-person > mine.ubc` prints a scenario to copy and edit; `ubc --scenario formulas mine.ubc` loads only the shared formulas before your own model.

### Simulation: `--samples`, `--seed`, `--random`

`ubc` draws 5000 samples from a random stream seeded with `1`, so the same program always prints the same numbers. More samples make the tails steadier; another seed shows how much the numbers move by chance:

```console
$ ubc -e '(8 to 25) * 12'
≈ 170   90% [96 … 296]   98% [75 … 380]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
$ ubc -n 100000 -e '(8 to 25) * 12'
≈ 170   90% [96 … 300]   98% [75 … 390]   ▃▅▇██▇▆▅▄▃▃▂▂▂▁▁
$ ubc --seed 2 -e '(8 to 25) * 12'
≈ 169   90% [95.8 … 296]   98% [76 … 370]   ▃▅▇█▇▇▇▅▄▃▃▂▂▂▁▁
```

`--random` seeds from the clock, so every run differs.

### Output: `--interval`, `--quantiles`, `--hist`, `--ascii`

`--interval W` prints a W % interval in place of the 90 % one; `--quantiles` prints the percentiles you list in place of both intervals:

```console
$ ubc --interval 80 -e '(8 to 25) * 12'
≈ 170   80% [109 … 262]   98% [75 … 380]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
$ ubc --quantiles 10,50,90 -e '(8 to 25) * 12'
≈ 170   p10 109   p50 170   p90 262   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
```

`--hist` draws a 20-row histogram from the 1st to the 99th percentile under every uncertain result:

```console
$ ubc --hist -e '(8 to 25) * 12'
≈ 170   90% [96 … 296]   98% [75 … 380]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
   83.1 ┤███████████
   98.2 ┤█████████████████████████
    113 ┤███████████████████████████████████
    129 ┤████████████████████████████████████████
    144 ┤████████████████████████████████████████████████
    159 ┤██████████████████████████████████████████████████
    174 ┤███████████████████████████████████████████████
    189 ┤█████████████████████████████████████████
    204 ┤███████████████████████████████
    219 ┤███████████████████████████
    234 ┤████████████████████████
    250 ┤████████████████
    265 ┤█████████████
    280 ┤█████████
    295 ┤████████
    310 ┤█████
    325 ┤████
    340 ┤██
    356 ┤███
    371 ┤█
(P1 to P99 shown, 100 samples outside)
```

`--ascii` replaces `≈`, `…` and the block characters with `~`, `...` and ` .:-=+*#`, for terminals and files that do not show Unicode. `--plain`, `--json` and `--echo` are described under [Output](#output).

### Other flags

- `--watch` runs the files again every time one of them changes, until you press Ctrl-C.
- `--completion bash` (or `zsh`, `fish`) prints a completion script, which also completes scenario names: `source <(ubc --completion bash)`.
- `--version` prints the version; `-h`, `--help` the list above.

## REPL

`ubc` with no arguments, in a terminal, starts an interactive session with line editing and a history kept in `~/.local/state/ubc/history` (or `$XDG_STATE_HOME/ubc/history`). Every line prints its result, assignments included, and variables live until you leave:

```
$ ubc
ubc> watts = 12 to 30
watts = ≈ 19   90% [12 … 29.7]   98% [9.8 … 37]   ▂▄▆▇██▇▅▅▄▃▂▂▂▁▁
ubc> P(watts > 25)
0.13
```

Lines that start with `:` are commands; they work only in the REPL:

| Command | Does |
|---|---|
| `:help` | List the commands and the language |
| `:vars` | Show every variable with its summary |
| `:hist EXPR` | Draw the histogram of `EXPR`, as `--hist` does |
| `:plot EXPR FILE` | Write the histogram of `EXPR` to `FILE`, an `.svg` image or a `.png` |
| `:sens EXPR` | Rank the variables that drive `EXPR`, as `--sens` does |
| `:seed N` | Restart the random stream with seed `N` |
| `:samples N` | Change the sample count; clears the variables, since their samples no longer fit |
| `:save FILE` | Write the statements that worked to `FILE`, to run it later as a model |
| `:load FILE` | Run `FILE` in the session |
| `:quit` | Leave, as Ctrl-D does (and Ctrl-Z on Windows) |

## Output

### Text

Each printed expression gives one line. A plain number prints with up to 10 significant figures, a probability with two decimals (`<0.01` and `>0.99` at the ends), and an uncertain value as a summary:

```
≈ 170   90% [96 … 296]   98% [75 … 380]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
  │       │                │                └ sparkline: the shape from P1 to P99
  │       │                └ 98% interval: the 1st and 99th percentiles, rounded outwards
  │       └ 90% interval: the 5th and 95th percentiles
  └ median
```

Summary numbers have 3 significant figures, more when the interval is narrow, and use the `k`, `M`, `B` and `T` suffixes from 10k up. An uncertain value whose samples are all equal prints as a plain number. Samples that are not finite, such as a division by zero, are left out of the summary, which says how many it dropped.

A program that prints a single expression prints it bare, so `$(ubc -e ...)` captures just the value. When a program (one file, one `-e` or stdin) prints more than one, each line says what it is: the variable name, the expression when it is 30 characters or shorter, or its line number otherwise:

```console
$ ubc -e 'x = 10 to 20; P(x > 15); x * 2; round(p(x, 0.9) * 100) / 100 + 1000'
P(x > 15) = 0.39
x * 2 = ≈ 28.3   90% [20 … 39.6]   98% [17 … 47]   ▂▄▅▆▇██▇▅▅▃▃▂▂▁▁
L1: 1018.4
```

`--echo` also prints assignments, so a model reads like a report:

```console
$ ubc --echo -e 'x = 10 to 20' -e 'double = x * 2'
x = ≈ 14.2   90% [10 … 19.8]   98% [8.6 … 24]   ▂▄▅▆▇██▇▅▅▃▃▂▂▁▁
double = ≈ 28.3   90% [20 … 39.6]   98% [17 … 47]   ▂▄▅▆▇██▇▅▅▃▃▂▂▁▁
```

### Plain

`--plain` prints only the median, with full precision and no labels, for scripts:

```console
$ ubc --plain examples/homelab.ubc
66.4979572018464
169.18163356927604
101.46083104992495
```

### JSON

`--json` prints one JSON object per result, on its own line. `kind` is `samples`, `scalar` or `probability`; a labelled result carries its label as `name`. Statistics leave out non-finite samples, and values that are not finite are `null`:

```console
$ ubc --json -e 'P(normal(0, 1) > 1)' -e '2 + 2'
{"kind":"probability","value":0.151}
{"kind":"scalar","value":4}
$ ubc --json --echo -e 'x = 1 to 5'
{"name":"x","kind":"samples","median":2.2408390648919774,"mean":2.5085321845122825,"stdev":1.2886668459505415,"p1":0.7126492324906979,"p5":0.9992998523855567,"p95":4.898313091927286,"p99":6.937904093115295,"n":5000,"dropped":0}
```

`n` is the sample count and `dropped` the number of non-finite samples. With `--interval` or `--quantiles` a `quantiles` object adds the requested percentiles. `--json` and `--plain` cannot be combined.

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | Evaluation error: the program is well written but cannot be computed (an unknown name, a bad range, a file that cannot be read) |
| `2` | Usage error: a bad flag or flag value |
| `3` | Parse error: the program is not well written |
| `4` | An `--assert` did not hold |

## Common errors

Errors go to stderr as `ubc: line L, col C: message`, with the file name in front for errors inside a file (`ubc: model.ubc: line L, col C: message`). In the REPL an error prints and the session goes on.

**Unexpected input** (exit 3). A parenthesis is not closed or an operator lacks an operand. Look at the column the message points to.

```console
$ ubc -e '(10 to 20'
ubc: line 1, col 10: unexpected end of input
$ ubc -e '2 +* 3'
ubc: line 1, col 4: unexpected '*'
```

**Range start must be below its end** (exit 1). Write the lower end first: `5 to 10`.

```console
$ ubc -e '10 to 5'
ubc: line 1, col 4: range start 10 must be below its end 5
```

**Range bounds must be plain numbers** (exit 1). The ends of a range cannot be uncertain. Pick plain ends that cover both uncertainties, or add the uncertainty on top: `(1 to 2) + (0 to 1)`.

```console
$ ubc -e '(1 to 2) to 3'
ubc: line 1, col 10: range bounds must be plain numbers
```

**Expects plain numbers** (exit 1). The arguments of a distribution cannot be uncertain either. Scale or shift a fixed distribution instead: `(10 to 20) + normal(0, 1)`.

```console
$ ubc -e 'normal(10 to 20, 1)'
ubc: line 1, col 1: normal expects plain numbers
```

**Unknown variable** (exit 1). The name is misspelled, or is used before the line that defines it; names are case-sensitive.

```console
$ ubc -e 'price * 2'
ubc: line 1, col 1: unknown variable 'price'
```

**Unknown function** (exit 1). Function names are lower case, except `P`; [Built-in functions](#built-in-functions) lists them.

```console
$ ubc -e 'Normal(0, 1)'
ubc: line 1, col 1: unknown function 'Normal'
```

**Wrong arguments** (exit 1). Check the number and the range of the arguments in [Built-in functions](#built-in-functions).

```console
$ ubc -e 'normal(0, 1, 2)'
ubc: line 1, col 1: normal expects 2 arguments, got 3
$ ubc -e 'normal(0, -1)'
ubc: line 1, col 1: normal: sd must be positive, got -1
```

**P expects a comparison** (exit 1). `P` asks how likely something is, so it needs a condition: `P(x > 15)`, not `P(x)`.

```console
$ ubc -e 'P(10 to 20)'
ubc: line 1, col 1: P expects a comparison
```

**Function calls nest too deeply** (exit 1). A function calls itself, directly or through another; without conditionals that never ends.

```console
$ ubc -e 'f(x) = f(x)' -e 'f(1)'
ubc: line 1, col 8: function calls nest too deeply (does f call itself?)
```

**An empty list** (exit 1). `[...]` needs at least one number to resample.

```console
$ ubc -e '[]'
ubc: line 1, col 1: an empty list has no values to resample
```

**Flags that do not fit** (exit 2). The message names the flag and what it expects:

```console
$ ubc --json --plain -e '1'
ubc: --json and --plain cannot be combined
$ ubc --set vps examples/homelab.ubc
ubc: --set expects name=expr, got 'vps'
$ ubc --interval 100 -e '1'
ubc: --interval expects a percentage between 0 and 100, got '100'
$ ubc --scenario team/replace
ubc: unknown scenario 'team/replace'; ubc --scenarios lists them
```

A flag `ubc` does not know prints `flag provided but not defined` and the list of flags. A flag written after a file name is read as another file, so it fails as a file that cannot be read: put flags first, `ubc -n 100000 model.ubc`.

**A result you did not expect.** Not every surprise is an error:

- `1 + 2 to 10` is `3 to 10`, because `to` binds looser than arithmetic: write `1 + (2 to 10)`.
- `(1 to 5) - (1 to 5)` is not 0, because each range is independent: bind it to a variable once and reuse the variable.
- Numbers move a little with `--seed`; that is Monte Carlo noise, and `--samples` reduces it. The [modelling guide](MODELLING.md) explains both.
