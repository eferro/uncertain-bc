# ADR 0026: Installation channels

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

Releases publish archives for six platforms (ADR 0015), but installing them by hand means picking the right archive, extracting it, moving the binary onto the `PATH` and, on macOS and Windows, getting past warnings about unsigned binaries. `go install` works only for people with Go, and it printed `dev` as the version. The user wants `ubc` to be very easy to install, with or without Go.

## Decision

Offer one obvious command per platform, all built from the same goreleaser release:

| Channel | Command | Notes |
|---|---|---|
| Install script (Linux, macOS) | `curl -fsSL …/install.sh \| sh` | POSIX `sh` so it runs where `sh` is dash. Detects OS and architecture, verifies the SHA-256 checksum, installs to `/usr/local/bin` if writable or `~/.local/bin`. A file fetched by `curl` gets no macOS quarantine attribute, so Gatekeeper does not block it. |
| Install script (Windows) | `irm …/install.ps1 \| iex` | Verifies the checksum, installs to `%LOCALAPPDATA%\Programs\ubc`, adds it to the user `PATH` and runs `Unblock-File`, so SmartScreen does not warn. |
| Homebrew (macOS, Linux) | `brew install eferro/tap/ubc` | goreleaser publishes the formula to `eferro/homebrew-tap`. |
| Scoop (Windows) | `scoop bucket add eferro https://github.com/eferro/scoop-bucket` then `scoop install ubc` | goreleaser publishes the manifest to `eferro/scoop-bucket`. |
| Linux packages | `.deb`, `.rpm`, `.apk` attached to each release | goreleaser `nfpms`. No package repository: download and install with `dpkg`/`rpm`/`apk`. |
| Go | `go install github.com/eferro/uncertain-bc/cmd/ubc@latest` | `ubc --version` reads the module version from the build info. |

Both install scripts accept `UBC_VERSION` and `UBC_INSTALL_DIR`. The POSIX script is tested offline against a fake GitHub (`install_test.go`), and CI runs `shellcheck` on every shell script. After each release, CI installs the new version through the scripts on Linux, macOS and Windows, and through `go install`, and checks the reported version.

Publishing to the tap and the bucket needs a token that can push to those two repositories. It is stored as the `TAP_GITHUB_TOKEN` secret, because the workflow's own `GITHUB_TOKEN` cannot write to other repositories. The secret is a fine-grained token with Contents read and write on those two repositories only; it was set on 2026-11-21 and must be renewed before it expires, or Homebrew and Scoop stop receiving new versions (releases still succeed, as both are skipped without the token).

**Not done: code signing.** Signing and notarising for macOS needs a paid Apple Developer account, and Windows Authenticode needs a code-signing certificate. The channels above avoid the warnings instead (no quarantine with `curl`, `Unblock-File` on Windows, and package managers). If signing becomes necessary, a new ADR will add it to the goreleaser pipeline.

## Consequences

- Each platform has a one-line install with integrity checks, with no Go and no manual `PATH` editing.
- Two extra public repositories (tap and bucket) and one secret to maintain.
- The install scripts become part of the public contract: their URLs and environment variables must stay stable.

## Alternatives considered

- **Only manual downloads:** works, but it was the friction this ADR removes.
- **Hosting apt/yum repositories:** automatic updates through the system package manager, but more infrastructure than a small tool needs. Packages attached to the release cover most uses.
- **Signing the binaries now:** the cleanest experience on macOS and Windows, but it has ongoing costs. Deferred.
