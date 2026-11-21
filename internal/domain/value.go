package domain

import (
	"math/rand/v2"
	"slices"
)

type Value interface {
	isValue()
}

type Scalar float64

func (Scalar) isValue() {}

type SampleSet []float64

func (SampleSet) isValue() {}

// Quantile interpolates linearly between the two nearest sorted samples
// (Hyndman & Fan type 7, the default in R and NumPy).
func (s SampleSet) Quantile(q float64) float64 {
	sorted := slices.Sorted(slices.Values(s))
	h := q * float64(len(sorted)-1)
	lower := int(h)
	if lower == len(sorted)-1 {
		return sorted[lower]
	}
	return sorted[lower] + (h-float64(lower))*(sorted[lower+1]-sorted[lower])
}

type Env struct {
	rng       *rand.Rand
	samples   int
	variables map[string]Value
}

func NewEnv(seed uint64, samples int) Env {
	return Env{rng: rand.New(rand.NewPCG(seed, seed)), samples: samples, variables: map[string]Value{}}
}

func (e Env) draw(sample func(*rand.Rand) float64) SampleSet {
	set := make(SampleSet, e.samples)
	for i := range set {
		set[i] = sample(e.rng)
	}
	return set
}

// combine applies op pairwise (SPEC §4 rule 1): sample i with sample i,
// with a scalar standing for the same value in every sample.
func combine(op func(a, b float64) float64, left, right Value) Value {
	l, leftIsScalar := left.(Scalar)
	r, rightIsScalar := right.(Scalar)
	if leftIsScalar && rightIsScalar {
		return Scalar(op(float64(l), float64(r)))
	}
	out := make(SampleSet, max(size(left), size(right)))
	for i := range out {
		out[i] = op(at(left, i), at(right, i))
	}
	return out
}

func size(v Value) int {
	if s, ok := v.(SampleSet); ok {
		return len(s)
	}
	return 1
}

func at(v Value, i int) float64 {
	if s, ok := v.(SampleSet); ok {
		return s[i]
	}
	return float64(v.(Scalar))
}
