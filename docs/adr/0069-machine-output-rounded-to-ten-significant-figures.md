# ADR 0069: `--plain` and `--json` round to 10 significant figures (amends 0005)

- Status: Accepted
- Effective date: 2026-11-25
- Deciders: eduardo.ferro

## Context

ADR 0005 promises "same program, same flags, same output, on every machine". `--plain` and `--json` printed the full float64, and its last bits are not the same on every machine: Go's `math.Exp` and `math.Log` have separate assembly per architecture, and on amd64 even choose an FMA path by CPU. On darwin/arm64 `ubc --plain examples/homelab.ubc` printed `66.49795720184639` where amd64 printed `66.4979572018464`, and the guide tests failed there (#43). CI ran only on amd64, while releases ship arm64 binaries.

The human output already rounds (3 significant figures for summaries, 10 for plain numbers), so it never showed the difference.

## Decision

- Every number `--plain` and `--json` print is rounded to **10 significant figures**, the rule plain numbers already follow in the normal output (SPEC §6).
- CI runs `make test` on `macos-latest` (arm64) as well, so a platform difference fails the build.
- This is a patch fix: ADR 0062 keeps the exact digits of a Monte Carlo result out of the stability promise, and the JSON fields and their meaning do not change.

## Consequences

- The same program prints the same `--plain` and `--json` output on every platform we ship, as ADR 0005 intended.
- Scripts get 10 significant figures instead of ~17. Monte Carlo error with 5000 samples is around 1%, so no meaningful information is lost.
- A value that lands exactly on a rounding boundary could still differ in the 10th figure between platforms; it is very unlikely, and the arm64 CI job would show it.
- `--samples-out` keeps the shortest exact form of each sample: it exports raw data, not a result.

## Alternatives considered

- **Round `--plain` only**: fixes the failing guide, but leaves `--json` different across platforms, against ADR 0005.
- **Keep full precision and tolerate last-digit differences in tests**: hides the problem from tests while users still see different output per machine.
- **Pure-Go `math.Exp`/`math.Log` everywhere**: not available without forking the standard library math, and slower.
