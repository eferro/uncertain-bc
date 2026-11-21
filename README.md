# Uncertain bc (`ubc`)

`bc` for numbers you are not sure about.

```console
$ ubc -e '(8 to 25) * 12'
≈ 170   90% [96 … 300]   ▁▃▇█▅▃▂▁
```

Type ranges instead of exact numbers (`10 to 20` means "90 % sure it is between 10 and 20"). `ubc` runs a Monte Carlo simulation and answers with a central value, an interval and a small chart, right in the terminal. A single Go binary, inspired by [Guesstimate](https://www.getguesstimate.com) and [Squiggle](https://www.squiggle-language.com).

Status: idea, nothing to install yet. See [docs/VISION.md](docs/VISION.md) for goals, language sketch, architecture and roadmap.
