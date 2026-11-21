# Uncertain bc (`ubc`)

`bc` for numbers you are not sure about.

```console
$ ubc -e '(8 to 25) * 12'
≈ 170   90% [96 … 300]   ▁▃▇█▅▃▂▁
```

Type ranges instead of exact numbers (`10 to 20` means "90 % sure it is between 10 and 20"). `ubc` runs a Monte Carlo simulation and answers with a central value, an interval and a small chart, right in the terminal. A single Go binary, inspired by [Guesstimate](https://www.getguesstimate.com) and [Squiggle](https://www.squiggle-language.com).

Status: walking skeleton (phase 0). Today `ubc` evaluates one `-e` expression with numbers, `+` and `a to b` ranges:

```console
$ ubc -e '1 + 2'
3
$ ubc -e '10 to 20'
≈ 14.2   90% [10 … 19.8]
```

The rest of the language (see [docs/SPEC.md](docs/SPEC.md)) arrives in phase 1. Goals, language sketch and roadmap: [docs/VISION.md](docs/VISION.md).

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
