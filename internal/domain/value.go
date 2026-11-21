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
	rng     *rand.Rand
	samples int
}

func NewEnv(seed uint64, samples int) Env {
	return Env{rng: rand.New(rand.NewPCG(seed, seed)), samples: samples}
}

func (e Env) draw(sample func(*rand.Rand) float64) SampleSet {
	set := make(SampleSet, e.samples)
	for i := range set {
		set[i] = sample(e.rng)
	}
	return set
}
