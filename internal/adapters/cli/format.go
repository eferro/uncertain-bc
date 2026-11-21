package cli

import (
	"fmt"
	"math"
	"strconv"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func format(value domain.Value) string {
	switch v := value.(type) {
	case domain.SampleSet:
		if v.Quantile(0) == v.Quantile(1) {
			return format(domain.Scalar(v[0]))
		}
		return fmt.Sprintf("≈ %s   90%% [%s … %s]",
			significant(v.Quantile(0.5)), significant(v.Quantile(0.05)), significant(v.Quantile(0.95)))
	case domain.Scalar:
		return strconv.FormatFloat(float64(v), 'g', -1, 64)
	case domain.Probability:
		return probability(float64(v))
	}
	panic(fmt.Sprintf("cli: unknown value type %T", value))
}

var suffixes = []struct {
	name  string
	scale float64
	from  float64
}{
	{"T", 1e12, 1e12},
	{"B", 1e9, 1e9},
	{"M", 1e6, 1e6},
	{"k", 1e3, 1e4},
}

// significant keeps three significant figures, because Monte Carlo noise makes
// further digits meaningless (ADR 0006), and uses suffixes from 10k up (SPEC §6).
func significant(x float64) string {
	rounded := roundSignificant(x)
	for _, s := range suffixes {
		if math.Abs(rounded) >= s.from {
			return shortest(rounded/s.scale) + s.name
		}
	}
	return shortest(rounded)
}

func roundSignificant(x float64) float64 {
	if x == 0 {
		return 0
	}
	decimals := 2 - int(math.Floor(math.Log10(math.Abs(x))))
	scale := math.Pow(10, float64(decimals))
	return math.Round(x*scale) / scale
}

func shortest(x float64) string {
	return strconv.FormatFloat(x, 'f', -1, 64)
}

// probability shows two decimals, but never rounds an uncertain outcome to a
// certain-looking 0 or 1 (SPEC §6).
func probability(p float64) string {
	switch {
	case p == 0 || p == 1:
		return shortest(p)
	case p < 0.01:
		return "<0.01"
	case p > 0.99:
		return ">0.99"
	}
	return strconv.FormatFloat(p, 'f', 2, 64)
}
