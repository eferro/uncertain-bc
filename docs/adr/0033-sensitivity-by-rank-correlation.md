# ADR 0033: Sensitivity analysis by rank correlation, with `:sens` and `--sens`

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Complements: [ADR 0032](0032-export-samples-with-samples-out.md)

## Context

Issue #7 asks to show which inputs drive most of the uncertainty of a result, as a tornado-style list. The open questions were the method (rank correlation or variance contribution) and the surface (a REPL command or a flag). A complexity review of the options came first, as the issue asked.

## Decision

- **Spearman's rank correlation of samples.** Sample *i* of every variable belongs to the same simulated scenario (SPEC §4 rule 1), so correlating each variable's samples with the result's samples pairs every input with the outcome it produced. Ranks instead of raw values make the coefficient catch any monotonic relation (a lognormal input multiplied in is not linear) and keep a few extreme samples from dominating. Tied samples, common with `[...]` and `--data`, share their average rank. Pairs where either sample is not finite are skipped, as summaries skip them (SPEC §4 rule 6). The coefficient keeps its sign, so the list also says in which direction each variable pushes. The domain computes it (`RankCorrelation`); the application ranks the session's variables (`Session.Sensitivity`); the CLI draws the list.
- **Every uncertain variable is a candidate.** ubc keeps no dependency graph, so the list cannot tell an input from an intermediate variable built from it; both appear, and the coefficient of an unrelated variable is near zero. Plain numbers, probabilities, variables that do not vary (their correlation is undefined) and the variable named exactly like the target are left out. Order is by the absolute coefficient, strongest first, ties in assignment order.
- **A tornado list in text.** A header `sensitivity of EXPR (rank correlation):`, then one line per variable with its name padded to a common width, the coefficient with sign and two decimals and a bar of `█` (`#` with `--ascii`, ADR 0022) whose length is |coefficient| × 20 cells, rounded. Bars share one fixed scale so a weak top driver looks weak.
- **Both surfaces.** `:sens EXPR` in the REPL and `--sens EXPR` for programs, printed after the program and `--samples-out` and before the assertions, like `--samples-out` (ADR 0032), so a failing check still shows what drove the result. The user chose both over a REPL command or a flag alone. EXPR is evaluated without binding anything, not even `last`, as assertions are (ADR 0029): asking a question does not change the session. `--sens` takes one expression; it is not repeatable.
- **Errors.** A plain-number EXPR, or no other uncertain variable to compare with, is an evaluation error (exit 1), as is any evaluation error of EXPR; a parse error is exit 3. `--sens` with `--json` or `--plain` is a usage error, as `--hist` is (there is no machine-readable form yet), and so is `--sens` in the REPL, which has `:sens`.

## Consequences

- A model answers "where should I reduce uncertainty first?" without leaving ubc.
- The coefficient measures how strongly a variable moves with the result, not the share of variance it explains; the SPEC says so to avoid reading the bars as percentages.
- Intermediate variables can crowd the list in large models; naming fewer intermediates, or a later dependency-aware filter, addresses it.

## Alternatives considered

- **Variance contribution (squared standardized regression coefficients, Sobol indices):** reads as "% of variance", but needs a regression or many extra model runs and assumes linearity or additivity; much more code and math for a calculator. It can be added later behind the same surfaces if rank correlation proves too coarse.
- **Pearson correlation of raw values:** simpler, but misses non-linear monotonic effects and is dominated by heavy lognormal tails, the common case in ubc models.
- **Tracking which variables an expression depends on:** would list true inputs only, but needs a dependency graph through assignments and user functions; postponed until intermediate variables prove a real nuisance.
- **JSON output and a repeatable flag:** no need yet; both fit the current design when asked for.
