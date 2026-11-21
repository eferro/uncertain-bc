# Uncertain bc (`ubc`)

`bc` for numbers you are not sure about.

```console
$ ubc -e '(8 to 25) * 12'
≈ 170   90% [96 … 300]   ▁▃▇█▅▃▂▁
```

Type ranges instead of exact numbers (`10 to 20` means "90 % sure it is between 10 and 20"). `ubc` runs a Monte Carlo simulation and answers with a central value, an interval and a small chart, right in the terminal. A single Go binary, inspired by [Guesstimate](https://www.getguesstimate.com) and [Squiggle](https://www.squiggle-language.com).

Status: phase 1 in progress. Today `ubc` has numbers with suffixes (`1k`, `5%`), `+ - * / ^`, comparisons and `P(...)`, ranges and distributions, summaries and math functions, variables and comments, and reads `-e` expressions, files or stdin:

```console
$ ubc -e '(8 to 25) * 12'
≈ 170   90% [96 … 296]   ▃▅▆██▇▆▅▄▃▂▂▂▁▁▁
$ ubc -e 'x = 1 to 5' -e 'x - x'
0
$ ubc examples/homelab.ubc
L15: ≈ 66.5   90% [47 … 95.7]   ▂▃▅▆▇█▇▆▅▄▃▃▂▂▂▁
L16: ≈ 169   90% [94.7 … 295]   ▃▅▇▇█▇▆▅▄▃▃▂▂▂▁▁
L18: ≈ 101   90% [21.5 … 229]   ▃▅▆▇█▇▇▆▄▃▃▂▂▂▁▁
```

See [examples/homelab.ubc](examples/homelab.ubc). Flags for samples, seed and ASCII output, and handling of non-finite samples are next ([docs/SPEC.md](docs/SPEC.md)); the interactive REPL comes in phase 2. Goals and roadmap: [docs/VISION.md](docs/VISION.md).

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
