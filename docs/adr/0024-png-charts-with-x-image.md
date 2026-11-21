# ADR 0024: PNG charts with golang.org/x/image

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

SPEC §7 asks `:plot` to write SVG and PNG. SVG is written by hand with the standard library. PNG matters for sharing in chat apps and documents that do not render SVG. The standard library can draw rectangles into an image and encode PNG, but it has no font rendering, so a PNG would have no title, summary or tick labels.

## Decision

- Render PNG with `image/png` from the standard library, plus `golang.org/x/image/font` and its built-in bitmap font `basicfont.Face7x13` for text. The font is pure Go, ships inside the module and needs no font files.
- SVG and PNG share one chart layout (title, summary, bars, ticks, clipping note). Each format only renders it.
- The bitmap font covers ASCII only, so PNG text uses the ASCII style (`~`, `...`), the same as `--ascii` (ADR 0022).

## Consequences

- PNG charts carry the same information as SVG charts and can be pasted anywhere.
- One more dependency, maintained by the Go team. Only the font packages are compiled into the binary.
- The bitmap font is small and fixed-size: fine for a 640×360 chart, not for print quality. SVG remains the better format for documents.

## Alternatives considered

- **SVG only:** simplest, but PNG is still needed where SVG does not render.
- **PNG with bars only:** standard library only, but a chart without labels is hard to read.
- **gonum/plot:** full plotting library, much more code and dependencies than one histogram needs.
