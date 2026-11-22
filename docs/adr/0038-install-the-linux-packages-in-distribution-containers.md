# ADR 0038: Install the `.rpm` and `.apk` packages in distribution containers after each release

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0026](0026-installation-channels.md)

## Context

Issue #13: every release attaches `.deb`, `.rpm` and `.apk` packages (ADR 0026), but only the `.deb` was ever installed in CI. The GitHub runners are Ubuntu, so neither `rpm` nor `apk` packages can be installed on them directly.

Installing the v0.7.2 packages by hand showed a real bug: on Fedora 45 (RPM 6.1) `rpm -i ubc-0.7.2-1.x86_64.rpm` fails with `package ubc-0.7.2-1.x86_64 does not verify: no signature`, because RPM now enforces signatures and the packages are not signed (signing is deferred by ADR 0026). The README told Fedora users to run exactly that command. `dnf install ./ubc-….rpm` installs the same package with a warning that it skipped the OpenPGP check. On Alpine, `apk add --allow-untrusted` works.

## Decision

1. **Document `dnf` for the `.rpm`.** The README says `sudo dnf install ./ubc-*.rpm` instead of `sudo rpm -i`. It is also the command Fedora and RHEL users expect for a downloaded package. Signing the packages stays out, as in ADR 0026.
2. **A `linux-package-install-check` job in `release.yml`** for the three packages. After the release job, a matrix downloads the `x86_64` `.deb`, `.rpm` and `.apk` of the tag, installs each one in an `ubuntu:latest`, a `fedora:latest` and an `alpine:latest` container (`docker run` on an Ubuntu runner, so `gh` stays outside the container) with the documented command, and checks that `ubc --version` reports the tag.
3. **`latest` images, not pinned versions.** The point is to catch what breaks on current distributions, as RPM 6 just did. A pinned image would hide that until someone bumps it.

## Consequences

- Every documented installation channel is now exercised after each release. The `.deb` check, a separate job until now, moved into the same matrix, so the three packages are checked the same way.
- A new Fedora or Alpine release can turn the job red without any change on our side. That is the signal we want: the README must keep working for current users.
- The job only runs on release tags, so it was verified locally by running the same script against v0.7.2 (and checked that it fails with `rpm -i` or a wrong version); the first tagged release runs it in CI.
- Only `x86_64` is installed; the `aarch64` packages are built by the same goreleaser config and are not checked.

## Alternatives considered

- **`container:` jobs** running directly in Fedora and Alpine: `gh` is not in those images, so the download would need another tool or an install step per distribution.
- **Building a snapshot and installing it on every push:** catches packaging changes before a release, but adds a goreleaser build and two container runs to every push for a part of the project that rarely changes.
- **Signing the `.rpm`:** would make `rpm -i` work, but needs a key to manage and publish. Deferred with the rest of signing (ADR 0026).
