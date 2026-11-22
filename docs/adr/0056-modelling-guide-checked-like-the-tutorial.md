# ADR 0056: A modelling guide checked like the tutorial

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro
- Amends: [ADR 0030](0030-tutorial-checked-by-the-test-suite.md)

## Context

Issue #34 asks for a guide on how to model well: what a 90 % interval claims and how to calibrate it, choosing a distribution, independence and correlation, typical mistakes and Monte Carlo noise. It must be linked from the README, the tutorial and the scenario catalogue, and every example must be checked by the test suite like the tutorial's. The tutorial teaches the language; this guide is about judgement, so mixing them would make the tutorial longer for readers who only want to learn the syntax.

## Decision

- The guide is `docs/MODELLING.md`, separate from the tutorial, in five sections: the claim of a 90 % interval with Hubbard's calibration techniques, choosing a shape, ranges are independent and variables are not, typical mistakes (precedence of `to`, double counting, mixing units, mean vs median, heavy tails) and Monte Carlo noise (seeds and `--samples`).
- **The same converter checks both guides.** `TestGuides_EveryConsoleBlockMatchesTheBinary` runs the ADR 0030 conversion on every guide in its list (`TUTORIAL.md`, `MODELLING.md`), one subtest each; the conventions of ADR 0030 apply unchanged. A new guide is checked by adding its name to the list.
- Most examples are one-line `-e` commands, because each makes one point; the correlation example needs a model, `examples/shared-rate.ubc`, which follows the tutorial's `ubc` block convention.
- Noise is shown by running the same expression with `--seed 1`, `2` and `3`, at the default and at 100000 samples, using a 99th percentile, because probabilities print with two decimals and hide the noise.
- A test checks that the README, the tutorial and the catalogue index link to the guide, so the links do not get lost in later edits.

## Consequences

- The guide cannot drift from the binary: a change in the RNG, the output or the language breaks it until it is updated, as with the tutorial.
- The guide's prose about the numbers (for example "the mean is eight times the median") is not checked; only the outputs are. Whoever updates the outputs must re-read the sentences around them.
- The converter is renamed from `tutorialArchive` to `guideArchive`, in `cmd/ubc/guide_script_test.go`, since it now serves any guide.

## Alternatives considered

- **A new section of the tutorial:** one file fewer, but it mixes learning the syntax with learning to estimate, and the tutorial already ends with three complete models.
- **A hand-written `.txtar` for the guide's examples:** checks the binary but not the Markdown, the drift ADR 0030 rejected.
- **Show noise with probabilities (`P(x > limit)`):** they print with two decimals, so the noise is invisible at the default sample count.
