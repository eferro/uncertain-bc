# Uncertain bc (`ubc`)

`bc` for numbers you are not sure about.

```console
$ ubc -e '(8 to 25) * 12'
≈ 170   90% [96 … 296]   98% [75 … 380]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
```

A VPS costs between 8 and 25 a month; what does it cost a year? Type a range instead of a number (`8 to 25` means "90 % sure it is between 8 and 25"), and `ubc` runs a Monte Carlo simulation and answers with a central value, two intervals and a small chart, right in the terminal.

## Why

Most estimates are ranges dressed up as numbers. A spreadsheet multiplies the guesses and hides how unsure they were; `ubc` keeps the uncertainty all the way to the answer, so you can ask how likely the bad case is instead of hoping it will not happen.

It is inspired by the Unix calculator [`bc`](https://www.gnu.org/software/bc/), from which it takes the habit and the name: quick calculations in the terminal, from the command line, a pipe or an interactive session. The idea of modelling with ranges comes from [Guesstimate](https://www.getguesstimate.com) and [Squiggle](https://www.squiggle-language.com). It is a single binary with no dependencies, for Linux, macOS and Windows.

## Install

**macOS and Linux**, with the install script (detects your system, verifies the checksum, installs to `/usr/local/bin` or `~/.local/bin`):

```sh
curl -fsSL https://raw.githubusercontent.com/eferro/uncertain-bc/main/install.sh | sh
```

**macOS and Linux**, with [Homebrew](https://brew.sh) (also installs shell completion):

```sh
brew install eferro/tap/ubc
```

**Windows**, in PowerShell (installs to `%LOCALAPPDATA%\Programs\ubc` and adds it to your `PATH`):

```powershell
irm https://raw.githubusercontent.com/eferro/uncertain-bc/main/install.ps1 | iex
```

**Windows**, with [Scoop](https://scoop.sh):

```powershell
scoop bucket add eferro https://github.com/eferro/scoop-bucket
scoop install ubc
```

**Debian, Ubuntu, Fedora, Alpine**: download the `.deb`, `.rpm` or `.apk` from the [latest release](https://github.com/eferro/uncertain-bc/releases/latest) and install it with `sudo dpkg -i ./ubc_*.deb`, `sudo dnf install ./ubc-*.rpm` or `sudo apk add --allow-untrusted ./ubc_*.apk` (the packages are not signed, so plain `rpm -i` refuses them on recent Fedora).

**With Go**: `go install github.com/eferro/uncertain-bc/cmd/ubc@latest`. **By hand**: every [release](https://github.com/eferro/uncertain-bc/releases) has archives for Linux, macOS and Windows (amd64 and arm64) and a `checksums.txt`.

Check it works with `ubc --version`. To pin a version with the scripts, set `UBC_VERSION=0.7.0`; to choose the directory, `UBC_INSTALL_DIR=~/bin`. The binaries are not signed: the script, Homebrew and Scoop installs avoid the warnings; if you download an archive with a browser, macOS needs `xattr -d com.apple.quarantine ubc` and Windows SmartScreen needs "More info → Run anyway". Shell completion (bash, zsh, fish): `source <(ubc --completion bash)` in your shell's startup file; Homebrew installs it for you.

## Use it from Claude Code

This repository is also a [Claude Code](https://code.claude.com) plugin marketplace with a `ubc` skill. With it, Claude uses `ubc` on its own whenever an answer rests on guessed numbers: a cost, an effort or a delivery date, a risk, a cloud bill, a Fermi problem. It writes the model as a `.ubc` file you can keep and change, runs it, and answers with the probability that matters to you instead of a single number built on guesses. The skill checks that ubc 1.1.1 or later is installed; when it is not, it offers to install the latest release into `~/.local/bin` with the install script above (Linux and macOS, checksum verified, no sudo).

Install the plugin from a Claude Code session:

```text
/plugin marketplace add eferro/uncertain-bc
/plugin install ubc@uncertain-bc
```

The same from the shell is `claude plugin marketplace add eferro/uncertain-bc` and then `claude plugin install ubc@uncertain-bc`; on Claude Code 2.1.275 or later one step is enough: `/plugin install ubc --marketplace eferro/uncertain-bc`. Add `--scope project` to the shell command to share it with everyone who works in a repository. Plugins from marketplaces you add yourself do not update on their own: run `claude plugin marketplace update uncertain-bc`, or turn on auto-update for the marketplace under Marketplaces in `/plugin`.

Without the plugin system, copy [`claude-plugin/skills/ubc`](claude-plugin/skills/ubc/SKILL.md) to `~/.claude/skills/ubc` for every project, or to a project's `.claude/skills/ubc`.

## Get started in three commands

Ask a question with `-e`. Variables keep their uncertainty, and `P(...)` says how likely something is:

```console
$ ubc -e 'hours = 20 to 60' -e 'cost = hours * (50 to 80)' -e 'P(cost > 3k)'
0.19
```

Run a model from a file. The [example](examples/homelab.ubc) compares the yearly cost of a home server with renting a VPS:

```console
$ ubc examples/homelab.ubc
homelab = ≈ 66.5   90% [47 … 95.7]   98% [40 … 111]   ▂▃▅▆▇█▇▆▅▄▃▃▂▂▂▁
vps = ≈ 169   90% [94.7 … 295]   98% [75 … 380]   ▃▅▇▇█▇▆▅▄▃▃▂▂▂▁▁
saving = ≈ 101   90% [21.5 … 229]   98% [-1.5 … 310]   ▃▅▆▇█▇▇▆▄▃▃▂▂▂▁▁
```

Run a ready-made model from the built-in catalogue (`ubc --scenarios` lists them all) and change its assumptions with `--set`. What does it cost to replace an engineer who leaves, if hiring took one to three months?

```console
$ ubc --set 'vacancy_months = 1 to 3' --scenario team/replace-a-person
total = ≈ 47.2k   90% [29.4k … 76.8k]   98% [24k … 95k]   ▃▅▇██▇▇▆▄▄▃▂▂▂▂▁
P(total > 50k) = 0.43
```

Run `ubc` with no arguments for an interactive session with line editing and history; `:help` lists its commands.

## Reading the output

```
≈ 170   90% [96 … 296]   98% [75 … 380]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
  │       │                │                └ the shape, from the 1st to the 99th percentile
  │       │                └ 98% sure it is in this range
  │       └ 90% sure it is in this range
  └ the median: as likely to be above as below
```

A probability prints with two decimals, a plain number as it is. When a program prints more than one result, each line starts with what it is, usually the variable name. The output is the same every time, because the seed is fixed; `--random` changes it.

## Flags you will use most

| Flag | What it does |
|---|---|
| `-e EXPR` | Evaluate an expression (repeatable); without `-e` or files, `ubc` reads stdin |
| `--echo` | Also print assignments, so a model reads like a report |
| `--set NAME=EXPR` | Replace a model's value with yours, a number or a range |
| `--data NAME=FILE` | Use your own measurements, one number per line |
| `--sens EXPR` | Rank the assumptions that drive a result's uncertainty |
| `--assert EXPR` | Exit with code 4 unless the condition certainly holds, for CI |
| `--scenarios`, `--scenario NAME` | List and run the built-in scenarios |
| `--plain`, `--json` | Output for scripts: only the median, or one JSON object per result |

`ubc --help` lists every flag; the [reference](docs/REFERENCE.md) explains them.

## Learn more

The guides below are also published, with search, at <https://eferro.github.io/uncertain-bc/>.

- [Tutorial](docs/TUTORIAL.md): the language step by step, from plain arithmetic to project estimates, a Fermi problem and a cloud bill.
- [Modelling guide](docs/MODELLING.md): how to estimate well, what a 90 % interval claims, choosing a shape, correlation, typical mistakes and Monte Carlo noise.
- [How-to guides](docs/HOWTO.md): recipes for common tasks, such as using your own data, comparing two options, a check in CI or a delivery forecast.
- [Scenario catalogue](scenarios/README.md): ready-to-use models for teams, delivery, product, engineering, reliability, cloud costs and risk.
- [Reference](docs/REFERENCE.md): every function, flag, REPL command, output format, exit code and common error.
- [Changelog](CHANGELOG.md): what changed in each release.

For contributors: the [specification](docs/SPEC.md), the [vision and roadmap](docs/VISION.md), the rules in [AGENTS.md](AGENTS.md) and the decisions in [docs/adr/](docs/adr/README.md). `make local-setup` enables the git hooks, `make help` lists the targets and `make validate` runs everything CI checks.

## License

[MIT](LICENSE)
