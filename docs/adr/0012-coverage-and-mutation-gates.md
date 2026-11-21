# ADR 0012: Coverage and mutation testing gates

- Status: Accepted, amended by [ADR 0020](0020-reliable-mutation-gate.md)
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

Testing quality is a core value of the project. The domain (language, evaluator, statistics) is pure logic, so near-complete testing is realistic and cheap. A coverage percentage alone can be gamed by tests that execute code without checking it. Mutation testing measures whether the tests actually catch changes, but it is too slow for a pre-commit hook. In the other Go projects by the author, project A gates coverage at 60% and project B gates gremlins mutation efficacy at 80–90% per package.

## Decision

- `make validate` fails when statement coverage of `internal/` is below **90%**. `cmd/` and test helpers are excluded from the measurement.
- `make test-mutation` runs [gremlins](https://github.com/go-gremlins/gremlins) on `internal/domain` and fails below **85% test efficacy**. It is pinned like every other tool.
- `make test-mutation` is not part of `validate` or the pre-commit hook. CI runs it on every push, and a failure turns the pipeline red.
- Thresholds only go up. Lowering one needs a new ADR.
- Surviving mutants are fixed with better tests, or recorded as equivalent mutants in the test file with a short reason.

## Consequences

- Untested code cannot land, and tests that don't check anything are detected.
- Some low-value code paths will need tests to keep coverage at 90%. If a threshold proves wrong, it is revisited through an ADR, not quietly bypassed.
- CI runs longer because of mutation testing. The domain is small and fast, so this is acceptable.

## Alternatives considered

- **Coverage only, with mutation testing on demand:** simpler, but nothing stops tests that execute code without checking it.
- **Mutation testing only:** a strong signal, but feedback is slow and nothing gates locally.
- **No numeric gate:** relies entirely on discipline, against the stated priority on testing.
