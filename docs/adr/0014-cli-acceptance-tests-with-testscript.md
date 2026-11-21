# ADR 0014: CLI acceptance tests with testscript

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

`ubc`'s contract with users is its command line: flags, stdin, files, stdout, stderr and exit codes (SPEC §1, §6, §9). The fixed seed (ADR 0005) makes output deterministic, so end-to-end tests can assert exact text. We want those tests to be readable as examples, cheap to add, and part of `go test` and coverage. Project B runs a bash script against the built binary, which tests the real artifact but sits outside `go test`.

## Decision

- CLI acceptance tests use `testscript` from `github.com/rogpeppe/go-internal/testscript` (the format the Go team uses for the `go` command).
- Scenarios are `.txtar` files in `cmd/ubc/testdata/script/`. Each one holds commands, stdin and fixture files, the expected stdout/stderr (regex or exact via `cmp`), and the exit status (`! exec` for failures).
- The `ubc` command runs in-process through `testscript.Main`, so scenarios count towards coverage and need no build step.
- Every acceptance example in SPEC §9 has a scenario. Adding a user-visible behaviour starts with a failing scenario (outer TDD loop), then unit tests inside the hexagon (inner loop).
- `main` stays a thin shell, `os.Exit(run(args, stdin, stdout, stderr))`, so `run` can also be unit-tested directly.

## Consequences

- Scenarios double as executable documentation of the CLI.
- One more test-only dependency. Like testify, it never ships in the binary.
- Exact output assertions will need updating when formatting or sampling changes. That is intended: such changes are user-visible.

## Alternatives considered

- **Hand-rolled golden files:** no dependency, but we would reimplement stdin, exit-code and file-fixture handling, and the tests read worse.
- **A bash script against the binary:** tests the real artifact, but it is outside `go test` and coverage, and less portable (Windows).
