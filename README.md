# Uncertain bc (`ubc`)

`bc` for numbers you are not sure about.

```console
$ ubc -e '(8 to 25) * 12'
≈ 170   90% [96 … 300]   ▁▃▇█▅▃▂▁
```

Type ranges instead of exact numbers (`10 to 20` means "90 % sure it is between 10 and 20"). `ubc` runs a Monte Carlo simulation and answers with a central value, an interval and a small chart, right in the terminal. A single Go binary, inspired by [Guesstimate](https://www.getguesstimate.com) and [Squiggle](https://www.squiggle-language.com).

Status: calculator MVP with an interactive REPL, JSON output and chart export (phase 3, v0.4.0).

```console
$ ubc -e '(8 to 25) * 12'
≈ 170   90% [96 … 296]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁

$ ubc -e 'yearly = (12 to 30) / 1000 * 24 * 365 * (0.12 to 0.25)' -e 'yearly' -e 'P(yearly > 50)'
≈ 28.7   90% [16.1 … 51]   ▃▅▇▇█▇▆▅▄▃▃▂▂▂▁▁
0.06

$ ubc --plain -e '1k to 5k'
2240.839064891978

$ ubc examples/homelab.ubc
L15: ≈ 66.5   90% [47 … 95.7]   ▂▃▅▆▇█▇▆▅▄▃▃▂▂▂▁
L16: ≈ 169   90% [94.7 … 295]   ▃▅▇▇█▇▆▅▄▃▃▂▂▂▁▁
L18: ≈ 101   90% [21.5 … 229]   ▃▅▆▇█▇▇▆▄▃▃▂▂▂▁▁
```

Run `ubc` with no arguments for an interactive session with line editing and history:

```console
$ ubc
ubc> watts = 12 to 30
watts = ≈ 19   90% [12 … 29.7]   ▂▄▆▇██▇▅▅▄▃▂▂▂▁▁
ubc> yearly = watts / 1000 * 24 * 365 * 0.2
yearly = ≈ 33.3   90% [21 … 51.9]   ▂▄▆▇██▇▅▅▄▃▂▂▂▁▁
ubc> P(yearly > 40)
0.25
ubc> :hist yearly
ubc> :vars
ubc> :plot yearly yearly.png
ubc> :save model.ubc
```

`:help` lists the commands (`:vars`, `:hist`, `:plot`, `:seed`, `:samples`, `:save`, `:load`, `:quit`) and the language. For scripts, `--json` prints one JSON object per result and `--plain` only the median; `--echo` also prints assignments with their names, so `ubc --echo model.ubc` reads like a report. `ubc --watch model.ubc` evaluates the model again every time you save it.

What the language has today:

- Numbers with suffixes (`1.5e3`, `1k`, `2.5M`, `5%`), `+ - * / ^`, parentheses, comments (`#`).
- `a to b`: a 90 % interval, lognormal when both ends are positive, normal otherwise.
- Distributions: `normal`, `lognormal`, `uniform`, `triangular`, `beta`, `pert`, `exponential`, `gamma`, and empirical data `[3, 5, 8, 13]` (resampled).
- Variables that keep their samples (`x - x` is exactly 0), user functions (`cost(n) = n * (2 to 4)`), `last`, `pi`, `e`.
- Comparisons and `P(...)`; summaries `mean`, `median`, `stdev`, `min`, `max`, `p(x, q)`; `sqrt`, `exp`, `log`, `log10`, `abs`, `round`, `floor`, `ceil`.
- Input from `-e` (repeatable), files or stdin; flags `-n/--samples`, `-s/--seed`, `--random`, `--plain`, `--ascii` (see `ubc --help`).

Full behaviour: [docs/SPEC.md](docs/SPEC.md). Goals and roadmap: [docs/VISION.md](docs/VISION.md).

## Install

```console
$ go install github.com/eferro/uncertain-bc/cmd/ubc@latest
```

Or download a binary for Linux, macOS or Windows from [Releases](https://github.com/eferro/uncertain-bc/releases).

## Development

```console
$ make local-setup   # once: enable the versioned git hooks
$ make help          # list targets
$ make validate      # everything CI checks; the pre-commit hook runs it
```

Rules and decisions: [AGENTS.md](AGENTS.md), [docs/adr/](docs/adr/README.md).

## License

[MIT](LICENSE)
