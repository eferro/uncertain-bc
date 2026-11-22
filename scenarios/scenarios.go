// Package scenarios embeds the scenario catalogue in the binary, so
// `ubc --scenario` works wherever ubc is installed, not only from a clone
// (ADR 0054).
package scenarios

import "embed"

// Files holds formulas.ubc and every <area>/<name>.ubc.
//
//go:embed *.ubc */*.ubc
var Files embed.FS
