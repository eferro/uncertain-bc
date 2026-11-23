# ADR 0068: A Claude Code skill, shipped from this repository as a plugin

- Status: Accepted
- Effective date: 2026-11-23
- Deciders: eduardo.ferro

## Context

AI coding agents answer questions such as "what will this cost?", "when will we finish?" or "is this worth automating?" by multiplying point guesses, the habit `ubc` exists to replace. Claude Code loads *skills*, folders with a `SKILL.md` whose description tells the agent when to use them, and installs them from *plugin marketplaces*, which can be any git repository with a `.claude-plugin/marketplace.json`.

A skill for `ubc` has to stay true to the binary it drives: the language, the flags, the output and the scenario names. Kept in another repository, its examples would drift silently. The project already checks its guides against the binary (ADR 0030).

## Decision

- **The skill lives in this repository**, in `claude-plugin/skills/ubc/`: `SKILL.md` (when to use ubc, how to build and report a model, the rules that keep a model honest), `references/language.md` (a compact language and flag reference) and `scripts/ensure-ubc.sh`.
- **The repository is a plugin marketplace** named `uncertain-bc` (`.claude-plugin/marketplace.json` at the root) with one plugin, `ubc`, whose source is `./claude-plugin`. Users install it with `/plugin marketplace add eferro/uncertain-bc` and `/plugin install ubc@uncertain-bc`, or copy the skill folder to `~/.claude/skills/ubc`. The README says how.
- **Claude invokes the skill on its own.** The description names the situations (an answer that rests on guessed or uncertain inputs: estimates, costs, dates, capacity, risk, Fermi problems), even when the user does not mention ubc, and excludes exact arithmetic and fitting statistics to a dataset.
- **The skill checks and installs the binary.** `ensure-ubc.sh` exits 0 and prints `ubc X ready: PATH` when ubc 1.1.1 or later is found (on `PATH` or in the install directory); otherwise it exits 1. With `--install` it first installs the latest release with `install.sh`, fetched from `main`, which verifies the checksum, into `UBC_INSTALL_DIR` (default `~/.local/bin`), so it never needs sudo. The skill tells the user before installing, and its `allowed-tools` pre-approve only the check and `ubc` itself, so Claude Code asks for permission to install. 1.1.1 is the minimum because it is the first release in which an assigned comparison keeps its samples (ADR 0066), which the skill's way of modelling events relies on. POSIX sh like `install.sh` (ADR 0026); Windows users get the PowerShell line instead.
- **The plugin has no `version`.** Without one, Claude Code installs the marketplace's latest commit, so the skill follows `main` without a release step of its own; `claude plugin validate` warns about it, and that warning is accepted. The minimum ubc version in the script protects users from a skill newer than their binary.
- **The test suite keeps it honest:** every console block in `SKILL.md` and `references/language.md` runs against the binary like the guides (ADR 0030); the script is tested end to end against a fake GitHub release, like `install.sh`; and a test checks that the marketplace names the plugin's manifest, that the skill is named after its folder, that its description fits 1024 characters, that it links every reference, and that the skill, the README and the script state the same minimum version. `make check-scripts` lints the script with shellcheck.

## Consequences

- A change to the language, the flags, the output or the scenario names updates the skill in the same commit, or its tests fail. AGENTS.md says so.
- When the skill starts relying on a newer feature, the minimum version in `ensure-ubc.sh`, `SKILL.md` and the README goes up together, once that version is released.
- `claude plugin validate .` (not run in CI, which has no Claude Code) should pass with only the missing-version warning.
- The marketplace and plugin manifests are small files at fixed paths that Claude Code reads from the default branch; renaming `claude-plugin/` or the plugin breaks installed copies.

## Alternatives considered

- **Keep the skill in a separate skills repository.** Easier to bundle with other skills, but nothing would check its examples against the binary, and it would drift.
- **Make the repository root the plugin** (`"source": "./"`, `"strict": false`): one directory fewer, but it mixes the plugin with the Go module and is not the documented layout.
- **Install without asking**, or bundle the binary in the plugin. Installing software is the user's call, and a bundled binary would have to be per platform and kept in step with releases.
- **Version the plugin with the ubc releases.** Users would stay on a version until a bump, at the cost of a step in every release; following `main` is simpler while the skill only documents released behaviour.
