package domain

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// Below these limits the direct methods are fast: their cost grows with lambda
// and n, so larger parameters are first reduced exactly (ADR 0035).
const (
	poissonDirectLimit  = 16
	binomialDirectLimit = 16
)

func poisson(lambda float64) sampler {
	return func(r *rand.Rand) float64 { return poissonDraw(r, lambda) }
}

func binomial(n, p float64) sampler {
	return func(r *rand.Rand) float64 { return binomialDraw(r, n, p) }
}

// poissonDraw counts the arrivals of a rate-1 Poisson process in [0, lambda]
// (Knuth, TAOCP vol. 2, 3.4.1 F). The m-th arrival time is Gamma(m): before
// lambda, m arrivals are counted and the rest of the interval recurses; after
// it, the first m-1 arrivals are uniform on [0, x], so their count in
// [0, lambda] is binomial. Every step is exact.
func poissonDraw(r *rand.Rand, lambda float64) float64 {
	if lambda < poissonDirectLimit {
		return poissonByProduct(r, lambda)
	}
	m := math.Floor(lambda * 7 / 8)
	x := gamma(r, m)
	if x < lambda {
		return m + poissonDraw(r, lambda-x)
	}
	return binomialDraw(r, m-1, lambda/x)
}

// poissonByProduct multiplies uniforms until the product falls below
// e^-lambda; the number of factors beyond the first is Poisson(lambda).
func poissonByProduct(r *rand.Rand, lambda float64) float64 {
	limit := math.Exp(-lambda)
	count := 0.0
	for product := r.Float64(); product > limit; product *= r.Float64() {
		count++
	}
	return count
}

// binomialDraw counts which of n uniforms fall below p (Knuth, TAOCP vol. 2,
// 3.4.1 F). The a-th smallest of them is Beta(a, n+1-a): if it is at or above
// p, the successes are among the a-1 uniforms below it, uniform on [0, x];
// otherwise those a are all successes and the n-a above x, uniform on [x, 1],
// recurse. Each step halves n and is exact.
func binomialDraw(r *rand.Rand, n, p float64) float64 {
	if n < binomialDirectLimit {
		return bernoulliTrials(r, n, p)
	}
	a := math.Floor(n/2) + 1
	x := beta(a, n+1-a)(r)
	if x >= p {
		return binomialDraw(r, a-1, p/x)
	}
	return a + binomialDraw(r, n-a, (p-x)/(1-x))
}

func bernoulliTrials(r *rand.Rand, n, p float64) float64 {
	successes := 0.0
	for range int(n) {
		if r.Float64() < p {
			successes++
		}
	}
	return successes
}

func finite(i int, name string) func(p []float64) string {
	return func(p []float64) string {
		if math.IsInf(p[i], 0) || math.IsNaN(p[i]) {
			return fmt.Sprintf("%s must be finite, got %g", name, p[i])
		}
		return ""
	}
}

func positiveWhole(i int, name string) func(p []float64) string {
	return func(p []float64) string {
		if p[i] < 1 || p[i] != math.Trunc(p[i]) || math.IsInf(p[i], 0) {
			return fmt.Sprintf("%s must be a positive whole number, got %g", name, p[i])
		}
		return ""
	}
}

func probability(i int, name string) func(p []float64) string {
	return func(p []float64) string {
		if !(p[i] >= 0 && p[i] <= 1) {
			return fmt.Sprintf("%s must be between 0 and 1, got %g", name, p[i])
		}
		return ""
	}
}
