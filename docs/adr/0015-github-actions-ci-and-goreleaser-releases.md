# ADR 0015: GitHub Actions CI and goreleaser releases

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

The project is public, and its module path is `github.com/eferro/uncertain-bc`. Goal 5 in VISION asks for one static binary for Linux, macOS and Windows, installable with `go install` or by downloading it. Local hooks (ADR 0011) can be skipped, so a shared gate is needed. Mutation testing (ADR 0012) is too slow for the pre-commit hook and needs a place to run.

## Decision

- The code is hosted on GitHub. Work is trunk-based on `main` (no branches unless asked).
- **CI** (`.github/workflows/ci.yml`) runs on every push to `main` and on pull requests. It executes `make validate` and `make test-mutation`, and checks that every pushed commit message follows Conventional Commits. A red pipeline is fixed before anything else lands.
- No automatic retries for failing tests. A flaky test is a bug.
- **Releases** (`.github/workflows/release.yml`) run when a `vX.Y.Z` tag is pushed. goreleaser builds `CGO_ENABLED=0` binaries for linux, darwin and windows on amd64 and arm64, plus checksums, and publishes a GitHub Release. Release notes come from the Conventional Commit history.
- **Versioning** follows SemVer, starting at `v0.x` until the MVP spec is complete. The version is injected with `-ldflags "-X main.version=..."` and shown by `ubc --version`.
- GitHub Actions and tools are pinned to exact versions, actions by commit SHA.

## Consequences

- Every push is checked by the same `make validate` that runs locally, plus mutation testing.
- Users get prebuilt binaries for six targets, and `go install github.com/eferro/uncertain-bc/cmd/ubc@latest` works.
- The GitHub repository must exist and be configured before the first push.

## Alternatives considered

- **CI only, deciding on releases later:** postpones a decision the roadmap already makes in phase 0.
- **Local only:** hooks can be bypassed and mutation testing would never run.
- **GitLab CI with CalVer scripts (project B):** fits an internal monorepo, not a public Go tool on GitHub.
