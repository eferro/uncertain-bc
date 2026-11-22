# Uncertain bc (`ubc`)

`bc` for numbers you are not sure about.

```console
$ ubc -e '(8 to 25) * 12'
≈ 170   90% [96 … 296]   98% [75 … 380]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
```

Type ranges instead of exact numbers (`10 to 20` means "90 % sure it is between 10 and 20"). `ubc` runs a Monte Carlo simulation and answers with a central value, two nested intervals (90 % and 98 %) and a small chart, right in the terminal. A single Go binary, inspired by the Unix calculator [`bc`](https://www.gnu.org/software/bc/) (quick calculations in the terminal, from the command line, a pipe or an interactive session) and by [Guesstimate](https://www.getguesstimate.com) and [Squiggle](https://www.squiggle-language.com) (modelling uncertainty with ranges).

```console
$ ubc -e '(8 to 25) * 12'
≈ 170   90% [96 … 296]   98% [75 … 380]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁

$ ubc -e 'yearly = (12 to 30) / 1000 * 24 * 365 * (0.12 to 0.25)' -e 'yearly' -e 'P(yearly > 50)'
≈ 28.7   90% [16.1 … 51]   98% [12 … 65]   ▃▅▇▇█▇▆▅▄▃▃▂▂▂▁▁
0.06

$ ubc --plain -e '1k to 5k'
2240.839064891978

$ ubc examples/homelab.ubc
homelab = ≈ 66.5   90% [47 … 95.7]   98% [40 … 111]   ▂▃▅▆▇█▇▆▅▄▃▃▂▂▂▁
vps = ≈ 169   90% [94.7 … 295]   98% [75 … 380]   ▃▅▇▇█▇▆▅▄▃▃▂▂▂▁▁
saving = ≈ 101   90% [21.5 … 229]   98% [-1.5 … 310]   ▃▅▆▇█▇▇▆▄▃▃▂▂▂▁▁
```

Run `ubc` with no arguments for an interactive session with line editing and history:

```console
$ ubc
ubc> watts = 12 to 30
watts = ≈ 19   90% [12 … 29.7]   98% [9.8 … 37]   ▂▄▆▇██▇▅▅▄▃▂▂▂▁▁
ubc> yearly = watts / 1000 * 24 * 365 * 0.2
yearly = ≈ 33.3   90% [21 … 51.9]   98% [17 … 64]   ▂▄▆▇██▇▅▅▄▃▂▂▂▁▁
ubc> P(yearly > 40)
0.25
ubc> :hist yearly
ubc> :vars
ubc> :plot yearly yearly.png
ubc> :save model.ubc
```

`:help` lists the commands (`:vars`, `:hist`, `:plot`, `:sens`, `:seed`, `:samples`, `:save`, `:load`, `:quit`) and the language. For scripts, `--json` prints one JSON object per result and `--plain` only the median; `--echo` also prints assignments with their names, so `ubc --echo model.ubc` reads like a report. `ubc --watch model.ubc` evaluates the model again every time you save it. `ubc --set team=5 --set 'rate=0.6 to 0.9' model.ubc` runs a scenario: the given values replace the model's own assignments to those names. `ubc --data lead=lead_times.txt model.ubc` does the same with observed data, one number per line, resampled like `[3, 5, 8, 13]`. `ubc --assert 'P(cost > 100k) < 0.1' model.ubc` turns a model into a check for CI: it exits with code 4 unless every assertion certainly holds. `ubc --samples-out samples.csv model.ubc` writes the raw samples of every variable, one column each (or a JSON object with `.json`), for pandas, R or a spreadsheet. `ubc --sens profit model.ubc` (or `:sens profit` in the REPL) ranks the variables by how strongly they drive `profit`'s uncertainty, as a tornado list.

What the language has today:

- Numbers with suffixes (`1.5e3`, `1k`, `2.5M`, `5%`), `+ - * / ^`, parentheses, comments (`#`).
- `a to b`: a 90 % interval, lognormal when both ends are positive, normal otherwise.
- Distributions: `normal`, `lognormal`, `uniform`, `triangular`, `beta`, `pert`, `exponential`, `gamma`, `poisson`, `binomial`, and empirical data `[3, 5, 8, 13]` (resampled).
- Variables that keep their samples (`x - x` is exactly 0), user functions (`cost(n) = n * (2 to 4)`), `last`, `pi`, `e`.
- Comparisons and `P(...)`; summaries `mean`, `median`, `stdev`, `min`, `max`, `p(x, q)`; `sqrt`, `exp`, `log`, `log10`, `abs`, `round`, `floor`, `ceil`.
- Input from `-e` (repeatable), files or stdin; flags `-n/--samples`, `-s/--seed`, `--random`, `--plain`, `--ascii` (see `ubc --help`).

New to `ubc`? Start with the [tutorial](docs/TUTORIAL.md): from plain arithmetic to project estimates, Fermi problems and cloud costs, with every example in [examples/](examples/). Have a question rather than time for a tutorial? The [scenario catalogue](scenarios/README.md) has ready-to-use models, such as the cost of replacing a person, hiring, attrition, a team's real capacity, a recurring meeting, the bus factor, a delivery date from your own throughput, Little's law, scope creep, the cost of delay of a feature, CD3 prioritisation, a business case, LTV and CAC payback, buy vs build, the cost of technical debt, a slow CI pipeline, flaky tests, a migration budget, the cost of incidents, composite availability, an SLO error budget, capacity planning, the cost of a bug by where it is caught, an LLM feature billed by the token, log ingestion, growing storage, reserved capacity vs on-demand, the infrastructure cost per customer, per-seat SaaS, ephemeral environments per pull request, the annual loss expectancy of a risk against its mitigation or the cost of a data breach: change the assumptions and run. Full behaviour: [docs/SPEC.md](docs/SPEC.md). Goals and roadmap: [docs/VISION.md](docs/VISION.md).

## Install

Pick the line for your system. None of them needs Go.

**macOS and Linux**, with the install script (detects your system, verifies the checksum, installs to `/usr/local/bin` or `~/.local/bin`):

```console
$ curl -fsSL https://raw.githubusercontent.com/eferro/uncertain-bc/main/install.sh | sh
```

**macOS and Linux**, with [Homebrew](https://brew.sh) (also installs shell completion):

```console
$ brew install eferro/tap/ubc
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

**With Go**:

```console
$ go install github.com/eferro/uncertain-bc/cmd/ubc@latest
```

**By hand**: every [release](https://github.com/eferro/uncertain-bc/releases) has archives for Linux, macOS and Windows (amd64 and arm64) and a `checksums.txt`.

Check it works with `ubc --version`. To pin a version with the scripts, set `UBC_VERSION=0.7.0`; to choose the directory, `UBC_INSTALL_DIR=~/bin`.

The binaries are not signed. The script, Homebrew and Scoop installs avoid the warnings; if you download an archive with a browser, macOS needs `xattr -d com.apple.quarantine ubc` and Windows SmartScreen needs "More info → Run anyway".

Shell completion (bash, zsh, fish): `source <(ubc --completion bash)`, or the same with `zsh`/`fish` in your shell's startup file. Homebrew installs it for you.

## Development

```console
$ make local-setup   # once: enable the versioned git hooks
$ make help          # list targets
$ make validate      # everything CI checks; the pre-commit hook runs it
```

Rules and decisions: [AGENTS.md](AGENTS.md), [docs/adr/](docs/adr/README.md).

## License

[MIT](LICENSE)
