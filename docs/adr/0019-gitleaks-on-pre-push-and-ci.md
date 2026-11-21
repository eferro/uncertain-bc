# ADR 0019: gitleaks on pre-push and in CI

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

The repository is public (ADR 0015). Once a token reaches GitHub it must be treated as compromised: it gets indexed and scraped within minutes, and rewriting history does not undo that. The last moment to stop a leak cheaply is before `git push`. Agents and humans both commit to this repo, so the check must not rely on anyone remembering it.

## Decision

- Use [gitleaks](https://github.com/gitleaks/gitleaks), pinned (`GITLEAKS_VERSION` in the Makefile) and installed into `.tools/` by `make local-setup`.
- `make secrets-scan` scans git history, with secrets redacted in its output. `LOG_OPTS="A..B"` narrows the scan to a range.
- The versioned `scripts/hooks/pre-push` hook scans exactly the commits being pushed: `remote..local`, or for a new branch everything not yet on any remote. Any finding aborts the push.
- CI runs `make secrets-scan` over the full history on every push, because hooks can be skipped with `--no-verify`.
- A false positive is handled with a `.gitleaksignore` entry (finding fingerprint) or an inline `gitleaks:allow` comment, in a commit that explains why. The scan is never disabled.
- If a real secret ever reaches the remote, revoke it first. Cleaning history is secondary.

## Consequences

- A token committed by mistake is caught locally before it leaves the machine, and again in CI.
- Pushes take a fraction of a second longer, which is negligible for a repository this size.
- One more pinned tool to install. `make local-setup` handles it.

## Alternatives considered

- **gitleaks in pre-commit (as in project A):** catches leaks even earlier, but it adds latency to every commit. Pre-push already blocks before anything is public.
- **Only GitHub secret scanning / push protection:** useful as an extra layer, but it covers known provider patterns only, and it acts after the data has reached GitHub.
