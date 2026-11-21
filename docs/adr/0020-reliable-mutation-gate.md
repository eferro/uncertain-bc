# ADR 0020: A reliable mutation gate

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro
- Amends: [ADR 0012](0012-coverage-and-mutation-gates.md)

## Context

When the mutation gate from ADR 0012 was first implemented, three problems showed up:

1. gremlins v0.6.0 exits with `0` even when efficacy is below `--threshold-efficacy`, so the gate never failed.
2. Some mutants turn the lexer loop into an endless `append`. With a generous timeout they ate all memory and the run was killed. With a short timeout, the number of timed-out mutants varied wildly between runs (9 vs 47), and so did the efficacy, which ignores timed-out mutants.
3. ADR 0012 says equivalent mutants are recorded "in the test file with a short reason", but ADR 0017 forbids comments in tests.

## Decision

- `make test-mutation` runs `scripts/checks/mutation.sh`. The script runs gremlins, reads `Test efficacy` from the report, and enforces the minimum itself.
- gremlins runs with `--workers 2 --timeout-coefficient 20` under `ulimit -v` (4 GB). A mutant that exhausts memory crashes its test binary and counts as killed. Two consecutive runs gave the same result (94.44%, 1 timeout).
- Equivalent mutants are listed in [docs/mutation-testing.md](../mutation-testing.md) with the reason each one survives. This replaces the "in the test file" clause of ADR 0012. The list is reviewed whenever the code it mentions changes.

## Consequences

- The gate fails when efficacy drops, and its result is repeatable.
- One more small script to maintain. If a future gremlins release fixes the exit code, the script can shrink to a plain call.
- Equivalent mutants are visible in one place without breaking the no-comments-in-tests rule.

## Alternatives considered

- **Trust `--threshold-efficacy`:** it does not fail the build in v0.6.0.
- **A short timeout without a memory cap:** fast, but the efficacy jumps between runs.
- **Comments in test files for equivalent mutants:** contradicts ADR 0017.
