# Architecture Decision Records

Long-lived decisions about `uncertain-bc`. One file per decision, numbered, never renumbered. To change a decision, add a new ADR that supersedes the old one and update the status of both.

Template: Status, Effective date, Deciders, then Context, Decision, Consequences and, when useful, Alternatives considered.

## Index

| ADR | Title | Status | Effective date |
| --- | --- | --- | --- |
| [0001](0001-ship-a-single-static-go-binary.md) | Ship a single static Go binary | Accepted | 2026-11-21 |
| [0002](0002-own-squiggle-inspired-language-not-squiggle-compatible.md) | Implement our own Squiggle-inspired language, not Squiggle itself | Accepted | 2026-11-21 |
| [0003](0003-monte-carlo-sample-sets-as-the-only-distribution-representation.md) | Monte Carlo sample sets as the only distribution representation | Accepted | 2026-11-21 |
| [0004](0004-ranges-are-90-percent-intervals-lognormal-when-positive.md) | `a to b` is a 90 % interval, lognormal when both ends are positive | Accepted | 2026-11-21 |
| [0005](0005-deterministic-seed-by-default.md) | Deterministic RNG seed by default | Accepted | 2026-11-21 |
| [0006](0006-bc-like-cli-with-honest-summary-output.md) | bc-like CLI with an honest one-line summary | Accepted | 2026-11-21 |
| [0007](0007-name-the-binary-ubc-and-the-project-uncertain-bc.md) | Name the binary `ubc` and the project `uncertain-bc` | Accepted | 2026-11-21 |
| [0008](0008-units-out-of-the-mvp.md) | Keep units out of the MVP | Accepted | 2026-11-21 |
| [0009](0009-stdlib-testing-with-testify-assertions.md) | stdlib `testing` with testify assertions | Accepted | 2026-11-21 |
| [0010](0010-hexagonal-architecture.md) | Hexagonal architecture | Accepted | 2026-11-21 |
| [0011](0011-makefile-as-single-entry-point-and-git-hooks.md) | Makefile as the single entry point, enforced by git hooks | Accepted | 2026-11-21 |
| [0012](0012-coverage-and-mutation-gates.md) | Coverage and mutation testing gates | Accepted | 2026-11-21 |
| [0013](0013-strict-golangci-lint-with-complexity-limits.md) | Strict golangci-lint with complexity limits | Accepted | 2026-11-21 |
| [0014](0014-cli-acceptance-tests-with-testscript.md) | CLI acceptance tests with testscript | Accepted | 2026-11-21 |
| [0015](0015-github-actions-ci-and-goreleaser-releases.md) | GitHub Actions CI and goreleaser releases | Accepted | 2026-11-21 |
| [0016](0016-stdlib-flag-for-cli-parsing.md) | Standard library `flag` for CLI parsing | Accepted | 2026-11-21 |
| [0017](0017-comments-explain-why-not-what.md) | Comments explain why, not what | Accepted | 2026-11-21 |
| [0018](0018-mit-license.md) | MIT license | Accepted | 2026-11-21 |
| [0019](0019-gitleaks-on-pre-push-and-ci.md) | gitleaks on pre-push and in CI | Accepted | 2026-11-21 |
| [0020](0020-reliable-mutation-gate.md) | A reliable mutation gate (amends 0012) | Accepted | 2026-11-21 |
| [0021](0021-manual-mutation-in-the-tdd-loop.md) | Manual mutation in the TDD loop, gremlins as the CI safety net | Accepted | 2026-11-21 |
| [0022](0022-unicode-output-everywhere.md) | Unicode output everywhere, ASCII on request (amends 0006) | Accepted | 2026-11-21 |
| [0023](0023-line-editing-with-x-term.md) | Line editing with golang.org/x/term | Accepted | 2026-11-21 |
| [0024](0024-png-charts-with-x-image.md) | PNG charts with golang.org/x/image | Accepted | 2026-11-21 |
| [0025](0025-command-line-ux-additions.md) | Command-line UX additions | Accepted | 2026-11-21 |
| [0026](0026-installation-channels.md) | Installation channels | Accepted | 2026-11-21 |
| [0027](0027-track-work-in-github-issues.md) | Track work in GitHub issues | Accepted | 2026-11-21 |
| [0028](0028-override-variables-with-set.md) | Override variables from the command line with `--set` | Accepted | 2026-11-21 |
| [0029](0029-assertions-with-a-dedicated-exit-code.md) | Assertions with `--assert` and a dedicated exit code | Accepted | 2026-11-21 |
| [0030](0030-tutorial-checked-by-the-test-suite.md) | A tutorial whose examples are checked by the test suite | Accepted | 2026-11-21 |
| [0031](0031-load-observed-data-with-data.md) | Load observed data from a file with `--data` | Accepted | 2026-11-21 |
| [0032](0032-export-samples-with-samples-out.md) | Export raw samples with `--samples-out` | Accepted | 2026-11-22 |
| [0033](0033-sensitivity-by-rank-correlation.md) | Sensitivity analysis by rank correlation, with `:sens` and `--sens` | Accepted | 2026-11-22 |
| [0034](0034-configurable-interval-width.md) | Configurable interval width with `--interval` | Accepted | 2026-11-22 |
| [0035](0035-exact-poisson-and-binomial-sampling.md) | Exact `poisson` and `binomial` sampling by recursive reduction | Accepted | 2026-11-22 |
| [0036](0036-raw-mode-repl-tested-through-a-pseudo-terminal.md) | Raw-mode REPL tested through testscript's pseudo-terminal | Accepted | 2026-11-22 |
| [0037](0037-repl-on-windows-consoles.md) | REPL on Windows consoles: escape processing on, checked in CI by reading the code | Accepted | 2026-11-22 |
| [0038](0038-install-the-linux-packages-in-distribution-containers.md) | Install the `.rpm` and `.apk` packages in distribution containers after each release (amends 0026) | Accepted | 2026-11-22 |
| [0039](0039-round-by-exact-powers-of-ten.md) | Round printed numbers by exact powers of ten | Accepted | 2026-11-22 |
| [0040](0040-nested-98-percent-interval.md) | A nested 98% interval in every summary (amends 0006, 0034) | Accepted | 2026-11-22 |
| [0041](0041-whole-test-suite-on-windows.md) | The whole test suite on Windows: portable tests and LF checkouts (amends 0037) | Accepted | 2026-11-22 |
| [0042](0042-scenario-catalogue.md) | A scenario catalogue checked by the test suite | Accepted | 2026-11-22 |
| [0043](0043-reusable-formulas-for-the-scenarios.md) | Reusable formulas for the scenarios in one shared file | Accepted | 2026-11-22 |
| [0044](0044-team-scenarios.md) | Team scenarios, self-contained and pinned by their answers | Accepted | 2026-11-22 |
| [0045](0045-delivery-scenarios.md) | Delivery scenarios, with observed data as a replaceable default | Accepted | 2026-11-22 |
| [0046](0046-product-scenarios.md) | Product scenarios around cost of delay, with ranges in place of confidence factors | Accepted | 2026-11-22 |
| [0047](0047-engineering-scenarios.md) | Engineering scenarios for technical debt and productivity, counted in recoverable hours | Accepted | 2026-11-22 |
| [0048](0048-reliability-scenarios.md) | Reliability scenarios, with availability described by its downtime share | Accepted | 2026-11-22 |
| [0049](0049-cloud-cost-scenarios.md) | Cloud and variable cost scenarios, in dollars | Accepted | 2026-11-22 |
| [0050](0050-risk-scenarios.md) | Risk and security scenarios around the annual loss expectancy | Accepted | 2026-11-22 |
| [0051](0051-round-the-98-percent-ends-outwards.md) | Round the ends of the 98% interval outwards (amends 0040) | Accepted | 2026-11-22 |
| [0052](0052-keep-enough-digits-for-narrow-intervals.md) | Keep enough digits to tell the ends of a narrow interval apart (amends 0006) | Accepted | 2026-11-22 |
| [0053](0053-label-results-with-their-expression.md) | Label results with their variable name or short expression (amends 0025) | Accepted | 2026-11-22 |
| [0054](0054-scenarios-built-into-the-binary.md) | Scenarios built into the binary, with `--scenarios`, `--scenario` and `--scenario-source` (amends 0042, 0043) | Accepted | 2026-11-22 |
| [0055](0055-catalogue-index-by-area-and-sources.md) | Catalogue index by area, a currency note and sources for every assumption (amends 0042) | Accepted | 2026-11-22 |
| [0056](0056-modelling-guide-checked-like-the-tutorial.md) | A modelling guide checked like the tutorial (amends 0030) | Accepted | 2026-11-22 |
| [0057](0057-user-reference-checked-by-the-test-suite.md) | A user reference checked by the test suite, with errors shown (amends 0030, 0056) | Accepted | 2026-11-22 |
| [0058](0058-how-to-guides-checked-by-the-test-suite.md) | How-to guides checked by the test suite, with data files and `cat` (amends 0030, 0057) | Accepted | 2026-11-22 |
| [0059](0059-readme-as-a-short-entry-point.md) | The README as a short entry point, checked like the guides (amends 0056) | Accepted | 2026-11-22 |
| [0060](0060-changelog-as-the-release-notes.md) | A changelog that is also the release notes (amends 0015) | Accepted | 2026-11-22 |
| [0061](0061-documentation-website-with-mkdocs.md) | A documentation website built by MkDocs and deployed to GitHub Pages | Accepted | 2026-11-22 |
| [0062](0062-version-1-and-its-stability-promise.md) | Version 1.0.0 and its stability promise (amends 0015) | Accepted | 2026-11-23 |
