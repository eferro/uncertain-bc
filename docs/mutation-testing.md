# Mutation testing

`make test-mutation` runs [gremlins](https://github.com/go-gremlins/gremlins) on `internal/domain` and fails below 85% test efficacy (ADR 0012, ADR 0020). CI runs it on every push. While developing, mutate the changed code by hand with the `mutation-testing` skill before each commit (ADR 0021).

## Equivalent mutants

Surviving mutants that no test can kill with the current behaviour. Revisit this list when the grammar or the distributions change. Each entry should either disappear or still hold.

| Location | Mutation | Why it survives |
|---|---|---|
| `distribution.go` `normal`, `mean + sd·N` | `mean - sd·N` | The standard normal is symmetric, so both draw the same distribution. Only the order of samples changes. |
| `distribution.go` `lognormal`, `mu + sigma·N` | `mu - sigma·N` | Same symmetry, in log space. |
