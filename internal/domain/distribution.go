package domain

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type sampler func(*rand.Rand) float64

// z95 is the standard normal quantile at 95%, so [mean - z95·σ, mean + z95·σ]
// holds 90% of the mass (SPEC §4, ADR 0004).
const z95 = 1.6448536269514722

// ninetyPercentRange expects from < to, so a positive start means both ends are positive.
func ninetyPercentRange(env Env, from, to float64) SampleSet {
	if from > 0 {
		return env.draw(lognormal(
			(math.Log(from)+math.Log(to))/2, (math.Log(to)-math.Log(from))/(2*z95)))
	}
	return env.draw(normal((from+to)/2, (to-from)/(2*z95)))
}

func normal(mean, sd float64) sampler {
	return func(r *rand.Rand) float64 { return mean + sd*r.NormFloat64() }
}

func lognormal(mu, sigma float64) sampler {
	return func(r *rand.Rand) float64 { return math.Exp(mu + sigma*r.NormFloat64()) }
}

func uniform(lo, hi float64) sampler {
	return func(r *rand.Rand) float64 { return lo + (hi-lo)*r.Float64() }
}

// triangular samples by inverting the CDF, which has a closed form on each
// side of the mode.
func triangular(lo, mode, hi float64) sampler {
	split := (mode - lo) / (hi - lo)
	return func(r *rand.Rand) float64 {
		u := r.Float64()
		if u < split {
			return lo + math.Sqrt(u*(hi-lo)*(mode-lo))
		}
		return hi - math.Sqrt((1-u)*(hi-lo)*(hi-mode))
	}
}

func beta(a, b float64) sampler {
	return func(r *rand.Rand) float64 {
		x := gamma(r, a)
		return x / (x + gamma(r, b))
	}
}

// gamma uses Marsaglia & Tsang (2000) for shape >= 1 and boosts smaller shapes
// with Gamma(shape+1)·U^(1/shape).
func gamma(r *rand.Rand, shape float64) float64 {
	if shape < 1 {
		return gamma(r, shape+1) * math.Pow(r.Float64(), 1/shape)
	}
	d := shape - 1.0/3
	c := 1 / math.Sqrt(9*d)
	for {
		x := r.NormFloat64()
		v := 1 + c*x
		if v <= 0 {
			continue
		}
		v = v * v * v
		if u := r.Float64(); math.Log(u) < 0.5*x*x+d*(1-v+math.Log(v)) {
			return d * v
		}
	}
}

func exponential(rate float64) sampler {
	return func(r *rand.Rand) float64 { return r.ExpFloat64() / rate }
}

func gammaScaled(shape, scale float64) sampler {
	return func(r *rand.Rand) float64 { return gamma(r, shape) * scale }
}

// pert is the Beta-PERT distribution used in project estimation: a beta
// rescaled to [lo, hi] whose mean is (lo + 4·mode + hi) / 6.
func pert(lo, mode, hi float64) sampler {
	shape := beta(1+4*(mode-lo)/(hi-lo), 1+4*(hi-mode)/(hi-lo))
	return func(r *rand.Rand) float64 { return lo + shape(r)*(hi-lo) }
}

type distribution struct {
	arity int
	check func(p []float64) string
	build func(p []float64) sampler
}

var distributions = map[string]distribution{
	"normal":      {2, positive(1, "sd"), func(p []float64) sampler { return normal(p[0], p[1]) }},
	"lognormal":   {2, positive(1, "sigma"), func(p []float64) sampler { return lognormal(p[0], p[1]) }},
	"uniform":     {2, ordered(0, 1), func(p []float64) sampler { return uniform(p[0], p[1]) }},
	"triangular":  {3, triangle, func(p []float64) sampler { return triangular(p[0], p[1], p[2]) }},
	"beta":        {2, all(positive(0, "a"), positive(1, "b")), func(p []float64) sampler { return beta(p[0], p[1]) }},
	"exponential": {1, positive(0, "rate"), func(p []float64) sampler { return exponential(p[0]) }},
	"gamma":       {2, all(positive(0, "shape"), positive(1, "scale")), func(p []float64) sampler { return gammaScaled(p[0], p[1]) }},
	"pert":        {3, triangle, func(p []float64) sampler { return pert(p[0], p[1], p[2]) }},
}

func (d distribution) call(c call, env Env, args []Value) (Value, error) {
	params, err := c.plainNumbers(args, d.arity)
	if err != nil {
		return nil, err
	}
	if problem := d.check(params); problem != "" {
		return nil, evalErrorAt(c.at, "%s: %s", c.name, problem)
	}
	return env.draw(d.build(params)), nil
}

func positive(i int, name string) func(p []float64) string {
	return func(p []float64) string {
		if p[i] <= 0 {
			return fmt.Sprintf("%s must be positive, got %g", name, p[i])
		}
		return ""
	}
}

func all(checks ...func(p []float64) string) func(p []float64) string {
	return func(p []float64) string {
		for _, check := range checks {
			if problem := check(p); problem != "" {
				return problem
			}
		}
		return ""
	}
}

func ordered(lo, hi int) func(p []float64) string {
	return func(p []float64) string {
		if p[lo] >= p[hi] {
			return fmt.Sprintf("lo must be below hi, got %g and %g", p[lo], p[hi])
		}
		return ""
	}
}

var triangle = all(ordered(0, 2), modeInside)

func modeInside(p []float64) string {
	if p[1] < p[0] || p[1] > p[2] {
		return fmt.Sprintf("mode must be between lo and hi, got %g, %g and %g", p[0], p[1], p[2])
	}
	return ""
}
