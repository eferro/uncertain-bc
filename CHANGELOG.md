# Changelog

All notable changes to `ubc` are documented in this file. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html). Each GitHub release shows the section of its version.

## [Unreleased]

## [1.0.0] - 2026-11-23

First stable version: from now on the language, flags, REPL commands, `--json` and `--plain` output, exit codes and scenario names only change in a backwards-incompatible way with a new major version ([ADR 0062](https://github.com/eferro/uncertain-bc/blob/main/docs/adr/0062-version-1-and-its-stability-promise.md)).

### Added

- `--set name=expr` overrides a model's variables from the command line.
- `--assert COND` checks a condition after the run and exits with code 4 when it does not hold, for checks in CI.
- `--data NAME=FILE` sets a variable to the numbers in a file, one per line, resampled like a list literal.
- `--samples-out FILE.csv` exports the raw samples of every variable, for pandas or a spreadsheet.
- `:sens EXPR` and `--sens` rank the variables that drive a result's uncertainty.
- `--interval W` chooses the width of the printed interval.
- Exact `poisson` and `binomial` distributions.
- Every summary shows a nested 98% interval next to the 90% one.
- A scenario catalogue of ready-to-use models for teams, delivery, product, engineering, reliability, cloud costs and risk, with reusable formulas and a source for every assumption. It is built into the binary: `--scenarios` lists it, `--scenario NAME` runs one and `--scenario-source NAME` prints its source.
- A tutorial, a modelling guide, a user reference and how-to guides, whose examples the test suite runs against the binary.
- This changelog.
- A documentation website at <https://eferro.github.io/uncertain-bc/> with the README, the guides, the reference, the scenario catalogue and this changelog, and a search box.

### Changed

- Results are labelled with their variable name or a short form of their expression.
- Narrow intervals keep enough digits to tell their ends apart.
- The README is a short entry point that links to the guides.

### Fixed

- Printed numbers no longer show binary floating-point noise.
- The ends of the 98% interval are rounded outwards, so it always contains the 90% one.
- The REPL handles escape sequences in the Windows console.

## [0.7.2] - 2026-11-21

### Fixed

- The Homebrew cask and the Scoop manifest are published with every release.

## [0.7.1] - 2026-11-21

### Fixed

- Release republished to configure publishing to the Homebrew tap; `ubc` itself is unchanged.

## [0.7.0] - 2026-11-21

### Added

- One-command installers, `install.sh` and `install.ps1`.
- `.deb`, `.rpm` and `.apk` packages with every release.
- A Homebrew cask and a Scoop manifest.
- Install instructions for every platform, with or without Go.

### Fixed

- `ubc --version` reports the module version when installed with `go install`.

## [0.6.0] - 2026-11-21

### Added

- `--echo` prints assignments, and the REPL names them.
- `--hist` draws histograms and `--quantiles` chooses the percentiles printed.
- `--watch` evaluates files again whenever they change.
- `--completion` prints bash, zsh and fish completion scripts.

## [0.5.0] - 2026-11-21

### Added

- Empirical data written as `[v1, v2, ...]`, resampled like a distribution.
- User functions defined with `name(params) = expr`.
- `exponential`, `gamma` and `pert` (Beta-PERT) distributions.

## [0.4.0] - 2026-11-21

### Added

- `--json` describes every result as a JSON object.
- `:plot EXPR FILE` exports a histogram to SVG or PNG.

## [0.3.0] - 2026-11-21

### Added

- An interactive session (REPL) when stdin is a terminal, with line editing and history across sessions.
- `:help`, `:vars`, `:seed`, `:samples`, `:quit`, `:save`, `:load` and `:hist EXPR` in the session.

## [0.2.0] - 2026-11-21

### Added

- Programs of statements with variables and comments, read from repeated `-e` expressions, files and stdin in one session.
- Arithmetic on sample sets: `+ - * / ^`, negation and parentheses.
- `normal`, `lognormal`, `uniform`, `triangular` and `beta` distributions.
- Comparisons and probabilities with `P(...)`.
- `mean`, `median`, `stdev`, `min`, `max`, `p` and per-sample `sqrt`, `exp`, `log`, `log10`, `abs`, `round`, `floor` and `ceil`.
- `pi`, `e`, and `last` for the previous result.
- `k`, `M`, `B`, `T` and `%` suffixes for numbers, and in large summaries.
- A sparkline after every summary, `--ascii` output, `--plain` for the central value only.
- `-n`, `--seed` and `--random` to choose the sample count and the seed.
- `-h`/`--help`.
- Non-finite samples are left out of summaries and counted.

## [0.1.0] - 2026-11-21

### Added

- `ubc -e EXPR` evaluates numbers and `a to b` ranges as Monte Carlo sample sets.
- Syntax errors report their position.
- `ubc --version`, and release binaries for Linux, macOS and Windows.

[Unreleased]: https://github.com/eferro/uncertain-bc/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/eferro/uncertain-bc/compare/v0.7.2...v1.0.0
[0.7.2]: https://github.com/eferro/uncertain-bc/compare/v0.7.1...v0.7.2
[0.7.1]: https://github.com/eferro/uncertain-bc/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/eferro/uncertain-bc/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/eferro/uncertain-bc/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/eferro/uncertain-bc/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/eferro/uncertain-bc/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/eferro/uncertain-bc/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/eferro/uncertain-bc/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/eferro/uncertain-bc/releases/tag/v0.1.0
