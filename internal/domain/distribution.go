package domain

import (
	"math"
	"math/rand/v2"
)

// z95 is the standard normal quantile at 95%, so [mean - z95·σ, mean + z95·σ]
// holds 90% of the mass (SPEC §4, ADR 0004).
const z95 = 1.6448536269514722

func ninetyPercentRange(env Env, from, to float64) SampleSet {
	if from > 0 && to > 0 {
		return lognormal(env, (math.Log(from)+math.Log(to))/2, (math.Log(to)-math.Log(from))/(2*z95))
	}
	return normal(env, (from+to)/2, (to-from)/(2*z95))
}

func normal(env Env, mean, sd float64) SampleSet {
	return env.draw(func(r *rand.Rand) float64 { return mean + sd*r.NormFloat64() })
}

func lognormal(env Env, mu, sigma float64) SampleSet {
	return env.draw(func(r *rand.Rand) float64 { return math.Exp(mu + sigma*r.NormFloat64()) })
}
