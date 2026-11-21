# ADR 0022: Unicode output everywhere, ASCII on request

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro
- Amends: [ADR 0006](0006-bc-like-cli-with-honest-summary-output.md) ("Colour and Unicode only on a TTY")

## Context

ADR 0006 limited Unicode (`≈`, `…`, sparkline blocks) to a TTY. In practice that makes every piped or redirected result (`echo … | ubc`, `ubc model.ubc > out.txt`, `ubc … | less`) look different from what the user sees in the terminal, including most README examples. UTF-8 is the norm in terminals, pagers, files and CI logs, and tools such as ripgrep and fd print it unconditionally.

## Decision

- `ubc` always prints UTF-8: `≈`, `…` and the block sparkline `▁▂▃▄▅▆▇█`.
- `--ascii` switches to ASCII (`~`, `...` and the sparkline ` .:-=+*#`) for terminals or consumers that cannot show UTF-8.
- Colour, when it arrives, stays TTY-only, as ADR 0006 says.

## Consequences

- Output is identical in the terminal, in pipes and in files. Golden tests show exactly what users see.
- A user on a legacy console has to pass `--ascii`. It is not detected automatically.

## Alternatives considered

- **Keep TTY detection (ADR 0006):** piped output silently changes shape, which is surprising and makes documentation examples inaccurate.
