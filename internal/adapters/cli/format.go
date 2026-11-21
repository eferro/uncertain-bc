package cli

import (
	"fmt"
	"math"
	"strconv"

	"github.com/eferro/uncertain-bc/internal/domain"
)

type style struct {
	approx   string
	ellipsis string
	blocks   []rune
}

var (
	unicodeStyle = style{approx: "≈", ellipsis: "…", blocks: unicodeBlocks}
	asciiStyle   = style{approx: "~", ellipsis: "...", blocks: asciiBlocks}
)

func styleFor(opts options) style {
	if opts.ascii {
		return asciiStyle
	}
	return unicodeStyle
}

func format(value domain.Value, st style) string {
	switch v := value.(type) {
	case domain.SampleSet:
		if v.Dropped() == len(v) {
			return fmt.Sprintf("no finite samples (%d dropped)", len(v))
		}
		return withDropped(summarise(v, st), v.Dropped())
	case domain.Scalar:
		return plain(float64(v))
	case domain.Probability:
		return probability(float64(v))
	}
	panic(fmt.Sprintf("cli: unknown value type %T", value))
}

func summarise(v domain.SampleSet, st style) string {
	if v.Quantile(0) == v.Quantile(1) {
		return plain(v.Quantile(0))
	}
	return fmt.Sprintf("%s %s   90%% [%s %s %s]   %s", st.approx,
		significant(v.Quantile(0.5)), significant(v.Quantile(0.05)), st.ellipsis, significant(v.Quantile(0.95)),
		sparkline(v, st.blocks))
}

// withDropped says how many non-finite samples the summary ignored (SPEC §4 rule 6).
func withDropped(summary string, dropped int) string {
	if dropped == 0 {
		return summary
	}
	return fmt.Sprintf("%s   (%d non-finite dropped)", summary, dropped)
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
	rounded := roundSignificant(x, 3)
	for _, s := range suffixes {
		if math.Abs(rounded) >= s.from {
			return shortest(rounded/s.scale) + s.name
		}
	}
	return shortest(rounded)
}

// plain shows a scalar with at most ten significant figures (SPEC §6).
func plain(x float64) string {
	return shortest(roundSignificant(x, 10))
}

func roundSignificant(x float64, digits int) float64 {
	if x == 0 {
		return 0
	}
	decimals := digits - 1 - int(math.Floor(math.Log10(math.Abs(x))))
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
