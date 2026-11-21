# ADR 0001: Ship a single static Go binary

- Status: Accepted
- Effective date: 2026-11-21
- Deciders: eduardo.ferro

## Context

The tool is meant to replace `bc` for quick estimates and to be shared with other people. Anything that needs a runtime (Node, Python) or a server adds friction for both. Startup time matters because the tool is used in pipes and shell scripts.

## Decision

Write `ubc` in Go and distribute it as one statically linked binary per platform (Linux, macOS, Windows; amd64 and arm64), with no runtime dependencies, no network access and no cgo.

## Consequences

- Install is `go install` or downloading a release asset (GoReleaser).
- Instant startup, comparable to `bc`.
- Rules out embedding a JavaScript engine through cgo (v8go) and depending on Node at runtime.
- Charts must be rendered in Go (terminal text, SVG, PNG).

## Alternatives considered

- **Web app or HTML page** (like Guesstimate): rejected, the target habit is the terminal.
- **Node/TypeScript CLI on top of squiggle-lang:** rejected, needs Node installed.
