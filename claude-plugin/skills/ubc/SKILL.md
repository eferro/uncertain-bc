---
name: ubc
description: Estimates uncertain quantities with ubc, a Monte Carlo calculator where `a to b` is a 90% confidence interval. Use whenever an answer rests on guessed or uncertain inputs (back-of-the-envelope and Fermi estimates, costs, effort and delivery dates, capacity, risk, revenue, cloud bills, "is it worth it", "how likely are we to exceed X"), even if ubc is not mentioned, instead of multiplying point guesses. Not for exact arithmetic or fitting statistics to a dataset.
allowed-tools: Bash(ubc *) Bash(${CLAUDE_SKILL_DIR}/scripts/ensure-ubc.sh)
---

# Estimating with ubc

`ubc` is `bc` for numbers nobody knows exactly. Every uncertain input is a range, `a to b`, meaning "90% sure it is between a and b"; `ubc` draws 5000 samples of each, computes with all of them and answers with the median, a 90% and a 98% interval, and `P(...)` for how likely something is. A product of "most likely" values hides how wrong it can be and never says how likely the bad case is; a ubc model says both, and the user can rerun it with their own numbers.

## 1. Make sure ubc is installed

Run `${CLAUDE_SKILL_DIR}/scripts/ensure-ubc.sh`. It prints `ubc X ready: PATH`; call ubc by that path when it is not plain `ubc`.

Exit 1 means ubc is missing or older than 1.1.1. Tell the user it will install the latest release into `~/.local/bin` with the project's installer (checksum verified, no sudo), then run `${CLAUDE_SKILL_DIR}/scripts/ensure-ubc.sh --install` and use the path it prints. On Windows, or without network access, give the user the install line for their system from <https://github.com/eferro/uncertain-bc#install> instead (PowerShell: `irm https://raw.githubusercontent.com/eferro/uncertain-bc/main/install.ps1 | iex`).

## 2. Build the model

1. **Frame a decision, not a number.** Name the quantity, its unit and period, and the threshold that matters: the budget, the deadline, "worth it" as a net gain above zero. The answer the user needs is usually `P(cost > budget)`, not the cost.
2. **Look for a ready model.** `ubc --scenarios` lists the built-in catalogue (team, delivery, product, engineering, reliability, cloud, risk). When one fits, read it with `ubc --scenario-source NAME` and run it with the user's numbers: `ubc --set 'name = a to b' --scenario NAME`. To change its structure, save the source to a file and edit it.
3. **Otherwise write a `.ubc` file**, in the user's project when it is worth keeping, else in a scratch directory. Follow the catalogue's shape: a comment with the question; the assumptions, one per line, each with its range and a comment on where it comes from; the calculation, a step per line; the result and the `P(...)` that answers the question.
4. **Decompose until every input can be bounded** (Fermi style): users per month × requests per user × price per request, not "traffic cost".
5. **Ask for the ranges the user knows** (their team, prices, history) and propose the rest with the reasoning behind each, labelled as your assumption. Observed history beats a guess: `--data name=file.txt` reads one number per line, and `[3, 5, 8]` resamples a few values.
6. **Run, then probe.** `ubc --echo model.ubc` prints every step; `ubc --sens result model.ubc` ranks the inputs that drive the result's uncertainty; `--set` runs the alternatives (a cheaper plan, a smaller scope) without editing the file.

## 3. Rules that keep a model honest

- **Ranges must be wide enough to be 90% sure.** People write the range they expect, which is too narrow. Start from ends that are certainly wrong and move each one in until unsure.
- **Put parentheses around every range you compute with.** `to` binds looser than arithmetic:

```console
$ ubc -e '8 to 25 * 12; (8 to 25) * 12'
8 to 25 * 12 = ≈ 49.2   90% [7.99 … 286]   98% [3.7 … 630]   █▅▃▂▂▁▁▁▁▁▁▁▁▁▁▁
(8 to 25) * 12 = ≈ 169   90% [96.9 … 298]   98% [76 … 380]   ▃▅▇██▇▆▅▄▃▃▂▂▂▁▁
```

- **The ends choose the shape.** Both positive is lognormal (durations, prices, sizes: never negative, long right tail); an end at zero or below is normal, so `0 to 10` allows negative costs. Write the smallest value you believe in, `0.5 to 10`. Changes that can go either way: `1 + (-5% to 15%)`.
- **Different ranges are independent; a variable is one quantity.** Two costs driven by the same unknown rate must reuse one `rate` variable, or their errors cancel and the model looks more certain than it is. A range inside a user function is drawn again on every call.
- **Range ends and distribution arguments are plain numbers.** `(1 to 2) to 3` is an error; scale instead: `base * (0.8 to 1.2)`.
- **Events replace conditionals.** A comparison of uncertain values keeps a 0/1 outcome per sample, so "if it happens, it costs X" is a product; `max(0, x)` and `min` clamp:

```console
$ ubc -e 'hit = uniform(0, 1) < 0.1; loss = hit * (50k to 500k); loss; P(loss > 0)'
loss = ≈ 0   90% [0 … 155k]   98% [0 … 410k]   █▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁
P(loss > 0) = 0.10
```

- **There are no units.** Keep one unit and one period per model, name them (`cost_eur_month`, `delay_weeks`), and convert with the constants that `--scenario formulas` loads (`hours_per_day`, `working_days_per_year`, `weeks_per_year`…) along with formulas such as `lead_time`, `cost_of_delay`, `ltv`, `annual_loss`, `npv`: `ubc --scenario formulas model.ubc`.
- **Count each uncertainty once.** No contingency on top of task ranges that already include surprises.
- **Median, percentile or mean.** `≈` is the median, what probably happens; a commitment is a percentile, `p(x, 0.85)`; the mean is for totals of many repeated bets (a year of incidents), and for a wide lognormal it is far above the median.
- **Check close calls for noise.** When two options are close or the answer is a small probability, rerun with `--seed 2` and `--seed 3`, or `-n 100000`. If the answer flips, say the options are equivalent as far as anyone knows.

## 4. Read and report

- An uncertain result prints `≈ median   90% [p5 … p95]   98% [p1 … p99]   sparkline`; a probability prints with two decimals.
- Each `-e` is a separate program, and a program that prints one result prints it unlabelled. Run a file, or one `-e` with statements separated by `;`, so every result carries its name. To parse results, add `--json`: one object per result, with `name`, `median`, `p5`, `p95` and the rest.
- `--sens` also lists intermediate variables built from the inputs; recommend researching the inputs, the assumptions at the top of the model.
- Tell the user, in their units: the answer to their question as a probability, the median and the 90% interval, the one or two assumptions that drive it and what would narrow them, and where the model file is so they can change it. Keep ubc's rounding (never `96.98057`), and say that the answer is only as good as the ranges.

## Reference

- [references/language.md](references/language.md): syntax, functions, distributions, flags, REPL-free usage and exit codes. Read it before writing anything beyond plain ranges and arithmetic.
- The guides (modelling, how-tos, full reference, scenario catalogue) are published at <https://eferro.github.io/uncertain-bc/>.
