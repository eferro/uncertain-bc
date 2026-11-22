package domain

import (
	"math"
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

type Probability float64

func (Probability) isValue() {}

// Quantile interpolates linearly between the two nearest sorted finite
// samples (Hyndman & Fan type 7, the default in R and NumPy). Non-finite
// samples are ignored (SPEC §4 rule 6).
func (s SampleSet) Quantile(q float64) float64 {
	sorted := slices.Sorted(slices.Values(s.finite()))
	if len(sorted) == 0 {
		return math.NaN()
	}
	h := q * float64(len(sorted)-1)
	lower := int(h)
	if lower == len(sorted)-1 {
		return sorted[lower]
	}
	return sorted[lower] + (h-float64(lower))*(sorted[lower+1]-sorted[lower])
}

func (s SampleSet) Mean() float64 { return mean(s) }

func (s SampleSet) Stdev() float64 { return stdev(s) }

func (s SampleSet) Dropped() int {
	return len(s) - len(s.finite())
}

func (s SampleSet) finite() []float64 {
	finite := make([]float64, 0, len(s))
	for _, x := range s {
		if IsFinite(x) {
			finite = append(finite, x)
		}
	}
	return finite
}

// IsFinite tells the samples that summaries and correlations count from the
// NaN and infinities that division by zero leaves (SPEC §4 rule 6).
func IsFinite(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0)
}

func finiteValues(v Value) []float64 {
	if s, ok := v.(SampleSet); ok {
		return s.finite()
	}
	x, _ := scalarOf(v)
	return []float64{x}
}

type Env struct {
	source     *rand.PCG
	rng        *rand.Rand
	samples    int
	variables  map[string]Value
	functions  map[string]definition
	overridden map[string]bool
	depth      int
}

func NewEnv(seed uint64, samples int) Env {
	source := rand.NewPCG(seed, seed)
	return Env{
		source: source, rng: rand.New(source), samples: samples,
		variables: predefined(), functions: map[string]definition{}, overridden: map[string]bool{},
	}
}

func predefined() map[string]Value {
	return map[string]Value{"pi": Scalar(math.Pi), "e": Scalar(math.E)}
}

func (e Env) Variable(name string) (Value, bool) {
	value, ok := e.variables[name]
	return value, ok
}

func (e Env) Reseed(seed uint64) {
	e.source.Seed(seed, seed)
}

// WithSamples keeps the random stream but forgets variables, whose sample
// sets would no longer match the new count, and their overrides with them.
func (e Env) WithSamples(samples int) Env {
	return Env{
		source: e.source, rng: e.rng, samples: samples,
		variables: predefined(), functions: e.functions, overridden: map[string]bool{},
	}
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
	l, leftIsScalar := scalarOf(left)
	r, rightIsScalar := scalarOf(right)
	if leftIsScalar && rightIsScalar {
		return Scalar(op(l, r))
	}
	out := make(SampleSet, max(size(left), size(right)))
	for i := range out {
		out[i] = op(at(left, i), at(right, i))
	}
	return out
}

// AsSamples gives v as n samples, a plain number standing for the same value
// in every sample, as in arithmetic (SPEC §4 rule 1).
func AsSamples(v Value, n int) SampleSet {
	samples := make(SampleSet, n)
	for i := range samples {
		samples[i] = at(v, i)
	}
	return samples
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
	x, _ := scalarOf(v)
	return x
}

func scalarOf(v Value) (float64, bool) {
	switch v := v.(type) {
	case Scalar:
		return float64(v), true
	case Probability:
		return float64(v), true
	}
	return 0, false
}

func mean(v Value) float64 {
	values := finiteValues(v)
	total := 0.0
	for _, x := range values {
		total += x
	}
	return total / float64(len(values))
}
