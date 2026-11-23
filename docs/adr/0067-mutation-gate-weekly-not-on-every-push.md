# ADR 0067: The mutation gate runs weekly, not on every push (amends 0012, 0020, 0021)

- Status: Accepted
- Effective date: 2026-11-23
- Deciders: eduardo.ferro

## Context

ADR 0012 and ADR 0020 set an 85 % gremlins efficacy gate on `internal/domain`, and ADR 0021 kept gremlins in CI as the safety net behind manual mutation in the TDD loop. As the domain grew, `make test-mutation` came to take about 10 of the 11 minutes of the CI `validate` job; everything else in the job (lint, tests with coverage, script checks, build, goreleaser check) takes about 20 seconds, and the Windows job already runs in parallel in about a minute. Every push waited 8 to 12 minutes for feedback, mostly for a check that rarely changes its verdict, because each increment is already mutated by hand before it is committed (ADR 0021).

Splitting the remaining work into parallel jobs would not help: each job pays about 30 seconds to set up Go to save a few seconds of work.

## Decision

- Push and pull request CI (`ci.yml`) no longer runs `make test-mutation`.
- A separate workflow, `mutation.yml`, runs `make test-mutation` every week and on demand (`workflow_dispatch`). It keeps the same 85 % minimum and the same reliability settings (ADR 0020); a red run is fixed first, like any red pipeline.
- Manual mutation of the changed code stays part of the definition of done (ADR 0021), and `make test-mutation` stays available locally.
- Before tagging a release, run the mutation workflow (or `make test-mutation`) if it has not run since the last change to `internal/domain`.

## Consequences

- Push feedback drops from about 11 minutes to about 1.
- A mutation regression can now reach `main` and is caught within a week instead of on the same push. The manual mutation step and the pre-release check limit that window.
- The gate is not weakened: the threshold is the same and still blocks, only less often.

## Alternatives considered

- **Keep gremlins on every push**: the slowest feedback for the least new information.
- **Remove gremlins entirely**: loses the only automated check that the manual mutation step is actually being done.
- **Run it only in the release workflow**: would make every release 10 minutes slower and fail late, after the tag is pushed.
- **Parallel jobs or sharding gremlins**: more workflow complexity for a check that does not need per-push feedback.
