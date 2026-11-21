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
	}
	panic(fmt.Sprintf("cli: unknown value type %T", value))
}

// significant keeps three significant figures: Monte Carlo noise makes
// further digits meaningless (ADR 0006).
func significant(x float64) string {
	if x == 0 {
		return "0"
	}
	decimals := 2 - int(math.Floor(math.Log10(math.Abs(x))))
	scale := math.Pow(10, float64(decimals))
	return strconv.FormatFloat(math.Round(x*scale)/scale, 'f', -1, 64)
}
