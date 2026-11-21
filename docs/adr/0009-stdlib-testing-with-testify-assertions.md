# ADR 0009: stdlib `testing` with testify assertions

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

Quality and testing matter a lot in this project, and TDD drives every change. We want tests that read as specifications and failures that say exactly what went wrong, without leaving the standard toolchain. Two other Go projects by the author show the options: project A uses Ginkgo/Gomega, and project B moved from Ginkgo to stdlib `testing` plus testify. The project rule is "standard library unless a dependency clearly pays off".

## Decision

- Tests use the standard `testing` package and run with `go test`.
- Assertions use `github.com/stretchr/testify`: `require` for preconditions that make the rest of the test meaningless, `assert` for the verifications.
- Do not call `t.Error*` or `t.Fatal*` directly in tests. A check in `make validate` will enforce this.
- Table-driven tests for homogeneous cases (lexer tokens, operator results). Write separate test functions when the intent differs.
- Test names describe behaviour: `TestRange_IsLognormal_WhenBothEndsPositive`.
- Errors are checked against sentinel or typed errors with `require.ErrorIs` / `require.ErrorAs`, not by matching substrings.
- Use hand-written fakes. Mock libraries and generated mocks are not allowed, including `testify/mock`.

## Consequences

- One small, widely known test dependency. It never ships in the binary because it is only imported from `_test.go` files.
- Failure messages show the expected and actual values with no extra code.
- The same style as project B, so habits and agent prompts transfer.

## Alternatives considered

- **stdlib only:** zero dependencies, but it means verbose `if got != want` blocks and inconsistent failure messages.
- **Ginkgo/Gomega:** expressive BDD style, but it adds two dependencies and a separate runner, and project B already moved away from it.
