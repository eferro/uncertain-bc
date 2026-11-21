package domain

import "math/rand/v2"

type list struct {
	elements []node
	at       position
}

// evaluate resamples the listed values with replacement, so empirical data
// mixes with any distribution (SPEC §5).
func (l list) evaluate(env Env) (Value, error) {
	if len(l.elements) == 0 {
		return nil, evalErrorAt(l.at, "an empty list has no values to resample")
	}
	values := make([]float64, len(l.elements))
	for i, element := range l.elements {
		value, err := element.evaluate(env)
		if err != nil {
			return nil, err
		}
		x, ok := scalarOf(value)
		if !ok {
			return nil, evalErrorAt(l.at, "list elements must be plain numbers")
		}
		values[i] = x
	}
	return env.resample(values), nil
}

func (e Env) resample(values []float64) SampleSet {
	return e.draw(func(r *rand.Rand) float64 { return values[r.IntN(len(values))] })
}
