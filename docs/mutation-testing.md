# Mutation testing

`make test-mutation` runs [gremlins](https://github.com/go-gremlins/gremlins) on `internal/domain` and fails below 85% test efficacy (ADR 0012, ADR 0020). CI runs it on every push.

## Equivalent mutants

Surviving mutants that no test can kill with the current behaviour. Revisit this list when the grammar or the distributions change. Each entry should either disappear or still hold.

| Location | Mutation | Why it survives |
|---|---|---|
| `distribution.go` `ninetyPercentRange`, `to > 0` | `to >= 0` | Without negative literals, a range never ends at 0 with a start below it. Killable once unary minus exists (`-5 to 0`). |
| `distribution.go` `ninetyPercentRange`, `from + to` | `from - to` | The normal branch only runs with `from = 0` today, so both give the same mean. Killable with unary minus (`-10 to 20`). |
| `distribution.go` `ninetyPercentRange`, `to - from` | `to + from` | Same reason: `from = 0` in every reachable normal range. |
| `distribution.go` `normal`, `mean + sd·N` | `mean - sd·N` | The standard normal is symmetric, so both draw the same distribution. Only the order of samples changes. |
| `distribution.go` `lognormal`, `mu + sigma·N` | `mu - sigma·N` | Same symmetry, in log space. |
