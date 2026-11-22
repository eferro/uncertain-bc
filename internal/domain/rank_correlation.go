package domain

import (
	"cmp"
	"math"
	"slices"
)

// RankCorrelation is Spearman's coefficient: the Pearson correlation of the
// ranks, with tied samples sharing their average rank. Ranks make it measure
// any monotonic relation, not only a straight line, and ignore outliers'
// size (ADR 0033). Rows where either sample is not finite are skipped; it is
// undefined when either side does not vary.
func RankCorrelation(x, y SampleSet) (float64, bool) {
	fx, fy := finitePairs(x, y)
	return pearson(ranks(fx), ranks(fy))
}

func finitePairs(x, y SampleSet) ([]float64, []float64) {
	fx := make([]float64, 0, len(x))
	fy := make([]float64, 0, len(y))
	for i := range x {
		if IsFinite(x[i]) && IsFinite(y[i]) {
			fx = append(fx, x[i])
			fy = append(fy, y[i])
		}
	}
	return fx, fy
}

func ranks(values []float64) []float64 {
	order := make([]int, len(values))
	for i := range order {
		order[i] = i
	}
	slices.SortFunc(order, func(a, b int) int { return cmp.Compare(values[a], values[b]) })
	ranked := make([]float64, len(values))
	for start := 0; start < len(order); {
		end := start + 1
		for end < len(order) && values[order[end]] == values[order[start]] {
			end++
		}
		for _, i := range order[start:end] {
			ranked[i] = float64(start+end+1) / 2
		}
		start = end
	}
	return ranked
}

func pearson(x, y []float64) (float64, bool) {
	mx, my := mean(SampleSet(x)), mean(SampleSet(y))
	var sxy, sxx, syy float64
	for i := range x {
		dx, dy := x[i]-mx, y[i]-my
		sxy += dx * dy
		sxx += dx * dx
		syy += dy * dy
	}
	if sxx == 0 || syy == 0 {
		return 0, false
	}
	return sxy / math.Sqrt(sxx*syy), true
}

// Inspect evaluates one expression without binding anything, not even
// last, so a question about the session leaves it as it was.
func Inspect(source string, env Env) (Value, error) {
	expr, err := parseExpression(source)
	if err != nil {
		return nil, err
	}
	return expr.evaluate(env)
}
