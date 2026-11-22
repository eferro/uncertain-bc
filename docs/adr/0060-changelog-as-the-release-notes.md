# ADR 0060: A changelog that is also the release notes

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0015](0015-github-actions-ci-and-goreleaser-releases.md)

## Context

Issue #38 asks for a user-facing `CHANGELOG.md` in the Keep a Changelog format, starting with everything unreleased since v0.7.2 and summarising the earlier versions, and for the GitHub release to say the same as the changelog. Until now goreleaser wrote the release notes from the Conventional Commit history (ADR 0015): a list of commit subjects written for developers, with refactorings filtered out but not grouped or explained, and nothing a user could read in the repository.

## Decision

- **`CHANGELOG.md` at the root**, in [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/) form: `## [Unreleased]` first, then `## [X.Y.Z] - YYYY-MM-DD` from newest to oldest, each with `### Added`, `### Changed`, `### Fixed` (and the other standard headings when needed), and a compare link at the bottom for every heading. Entries are written for users, one line per change, not copied from the commits. Versions 0.1.0 to 0.7.2 are summarised from the git history.
- **The release notes are the version's section.** A small Go program, `tools/releasenotes`, prints the body of one version's section (`go run ./tools/releasenotes CHANGELOG.md v1.2.0`, through `make -s release-notes VERSION=v1.2.0`); the leading `v` is optional and link definitions are left out. The release workflow writes it to a file and passes it to goreleaser with `--release-notes`, so the `changelog` section of `.goreleaser.yaml` is gone. A version with no section or an empty one makes the program fail, and the release stops before anything is published.
- **The test suite checks the shape** of `CHANGELOG.md`: the title, `Unreleased` first, released versions as `X.Y.Z` with an ISO date and non-empty notes, newest first without repeats, the oldest being 0.1.0, and a link for every heading. The program has its own testscript scenario and unit tests.
- **Development tools live in `tools/`**, outside `internal/` and `cmd/`: they are not part of `ubc`, so they stay out of the hexagonal layout and are not installed by `go install ./cmd/...`. Their tests run with `make test` and `make validate` like the rest.
- **Who writes it:** AGENTS.md asks every user-visible change (a `feat:` or `fix:`, or a guide users read) to add its entry under `Unreleased` in the same commit. Releasing means renaming that section to the version and date, opening a new empty `Unreleased`, updating the links, committing, and then pushing the tag.
- The archives and the Linux packages ship `CHANGELOG.md` next to `README.md`, and the README links to it.

## Consequences

- The changelog, the GitHub release and the packaged file say the same, and forgetting to move `Unreleased` before tagging fails the release instead of publishing empty notes.
- Writing the entry is a manual step for each change; the test suite checks the format, not that every change has an entry.
- The Homebrew and Scoop manifests keep no notes of their own; they point to the release.

## Alternatives considered

- **Keep goreleaser's notes from the commits and generate `CHANGELOG.md` from them** (git-cliff, goreleaser's changelog): no hand-written entries, but the text would still be commit subjects, and generating the file would need another tool and a commit at release time.
- **Extract the section with `awk` in the workflow:** shorter, but untested and only runnable where a POSIX shell exists; the Go program is tested on every platform CI runs.
- **goreleaser's `--release-notes-tmpl`:** templates cannot select a Markdown section of a file.
