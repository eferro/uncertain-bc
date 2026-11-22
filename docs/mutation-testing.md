# Mutation testing

`make test-mutation` runs [gremlins](https://github.com/go-gremlins/gremlins) on `internal/domain` and fails below 85% test efficacy (ADR 0012, ADR 0020). CI runs it on every push. While developing, mutate the changed code by hand with the `mutation-testing` skill before each commit (ADR 0021).

## Equivalent mutants

Surviving mutants that no test can kill with the current behaviour. Revisit this list when the grammar or the distributions change. Each entry should either disappear or still hold.

| Location | Mutation | Why it survives |
|---|---|---|
| `distribution.go` `normal`, `mean + sd·N` | `mean - sd·N` | The standard normal is symmetric, so both draw the same distribution. Only the order of samples changes. |
| `distribution.go` `lognormal`, `mu + sigma·N` | `mu - sigma·N` | Same symmetry, in log space. |
| `distribution.go` `gamma`, `v := 1 + c·x` | `1 - c·x` | Same symmetry: `x` is a standard normal draw. |
| `distribution.go` `gamma`, `c := 1 / √(9d)` | other values of `c` | In Marsaglia & Tsang the density of the accepted `x` does not depend on `c`; `c` only shapes the envelope, so it changes efficiency, not the distribution. |
| `distribution.go` `triangular`, `u < split`; `gamma`, `v <= 0` and `log(u) < …` | `<=` / `<` | The two sides differ only when a continuous draw hits a boundary exactly, which has probability zero. |
| `summary.go` `stdev`, `(x - m)·(x - m)` | `(x + m)·(x - m)` | Summed over the samples, `Σ(x+m)(x−m) = Σx² − n·m² = Σ(x−m)²`. |
| `rank_correlation.go` `ranks`, `cmp.Compare(values[a], values[b])` | `cmp.Compare(values[b], values[a])` | Both sides are ranked the same way, so descending order reverses both rankings and Pearson's coefficient of the ranks stays the same. |
| `rank_correlation.go` `ranks`, `(start+end+1) / 2` | `(start+end) / 2` | Every rank moves by the same 0.5, and a correlation does not change when a variable is shifted. |
| `rank_correlation.go` `pearson`, `my := mean(y)` | `mean(x)` | `pearson` only ever gets ranks, and average ranks of n values always sum to n(n+1)/2, so both means are (n+1)/2. |
| `discrete_distribution.go` `poissonDirectLimit`, `binomialDirectLimit` | other limits | Both paths draw exactly the same distribution; the limit only decides which one is faster (ADR 0035). |
| `discrete_distribution.go` `poissonDraw`, `m := ⌊7λ/8⌋`; `binomialDraw`, `a := ⌊n/2⌋ + 1` | other split points | The reduction is exact for any split point; it only changes how many steps the recursion takes. |
| `discrete_distribution.go` `poissonDraw`, `x < lambda`; `binomialDraw`, `x >= p`; `poissonByProduct`, `product > limit`; `bernoulliTrials`, `u < p` | `<=` / `<` / `>=` | The two sides differ only when a continuous draw hits a boundary exactly, which has probability zero. |
| `adapters/cli/format.go` `roundSignificant`, `decimals < 0` | `decimals <= 0` / `decimals < -1` | A step of one changes nothing, and rounding to tens is exact with either formula: an integer below 2^30 divided by `0.1` is always exact, and the ten-digit range stays within half an ulp (ADR 0039). |
| `adapters/cli/terminal.go` `enterRawMode`, `term.Restore(fd, state)` in `restore` | deleted | The pseudo-terminal of the scenario is discarded when `ubc` exits, so no script can observe whether its state was restored (ADR 0036). |
| `adapters/cli/terminal.go` `enterRawMode`, `restoreOutput = enableVirtualTerminal(out)` | deleted | It does nothing outside Windows, and no CI runner has an interactive Windows console to run the REPL on; `enableVirtualTerminal` itself is tested on the Windows CI job (ADR 0037). |
