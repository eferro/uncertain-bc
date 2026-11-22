package cli

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/eferro/uncertain-bc/internal/domain"
)

type style struct {
	approx    string
	ellipsis  string
	blocks    []rune
	axis      string
	bar       string
	plain     bool
	json      bool
	echo      bool
	hist      bool
	quantiles []float64
	interval  float64
	// jsonQuantiles are the percentiles --json adds to the standard fields.
	jsonQuantiles []float64
}

// defaultInterval is the width of the summary interval unless --interval says
// otherwise: the 90 % that `a to b` also means (ADR 0004, ADR 0006).
const defaultInterval = 90

var (
	unicodeStyle = style{approx: "≈", ellipsis: "…", blocks: unicodeBlocks, axis: "┤", bar: "█", interval: defaultInterval}
	asciiStyle   = style{approx: "~", ellipsis: "...", blocks: asciiBlocks, axis: "|", bar: "#", interval: defaultInterval}
)

func styleFor(opts options) style {
	st := unicodeStyle
	if opts.ascii {
		st = asciiStyle
	}
	st.plain = opts.plain
	st.json = opts.json
	st.echo = opts.echo
	st.hist = opts.hist
	st.quantiles = opts.percentiles
	st.interval = opts.width
	st.jsonQuantiles = opts.percentiles
	if opts.interval != "" {
		st.jsonQuantiles = intervalEnds(opts.width)
	}
	return st
}

// display is a result's line, followed by its histogram with --hist.
func display(result domain.Result, st style) string {
	line := formatResult(result, st)
	samples, uncertain := result.Value.(domain.SampleSet)
	if !st.hist || !uncertain || samples.Dropped() == len(samples) || samples.Quantile(0) == samples.Quantile(1) {
		return line
	}
	return line + "\n" + strings.TrimSuffix(histogramChart(samples, st), "\n")
}

// formatResult labels an assignment with its variable name.
func formatResult(result domain.Result, st style) string {
	switch {
	case result.Name == "":
		return format(result.Value, st)
	case st.json:
		return jsonLine(result.Value, result.Name, st.jsonQuantiles)
	}
	return result.Name + " = " + format(result.Value, st)
}

func format(value domain.Value, st style) string {
	if st.plain {
		return centralValue(value)
	}
	if st.json {
		return jsonLine(value, "", st.jsonQuantiles)
	}
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

// centralValue prints the median with full precision, for $(ubc --plain ...) (SPEC §6).
func centralValue(value domain.Value) string {
	if s, ok := value.(domain.SampleSet); ok {
		return shortest(s.Quantile(0.5))
	}
	if p, ok := value.(domain.Probability); ok {
		return shortest(float64(p))
	}
	return shortest(float64(value.(domain.Scalar)))
}

// summarise leaves out the sparkline in a style without blocks, as error
// messages on stderr use.
func summarise(v domain.SampleSet, st style) string {
	switch {
	case v.Quantile(0) == v.Quantile(1):
		return plain(v.Quantile(0))
	case st.blocks == nil:
		return summaryNumbers(v, st)
	}
	return summaryNumbers(v, st) + "   " + sparkline(v, st.blocks)
}

func summaryNumbers(v domain.SampleSet, st style) string {
	width := innerWidth(v, st.interval)
	median := st.approx + " " + suffixed(roundForSummary(v.Quantile(0.5), width))
	if len(st.quantiles) > 0 {
		return median + quantileList(v, st.quantiles, width)
	}
	return median + "   " + intervals(v, st, width)
}

// innerWidth is the width of the inner interval even when --quantiles hides
// it, so one rounding rule covers every summary (ADR 0052).
func innerWidth(v domain.SampleSet, interval float64) float64 {
	ends := intervalEnds(interval)
	return v.Quantile(ends[1]/100) - v.Quantile(ends[0]/100)
}

func quantileList(v domain.SampleSet, quantiles []float64, width float64) string {
	var list strings.Builder
	for _, q := range quantiles {
		fmt.Fprintf(&list, "   %s %s", quantileLabel(q), suffixed(roundForSummary(v.Quantile(q/100), width)))
	}
	return list.String()
}

// intervals nests the tail interval around the inner one only while it is
// the wider of the two (ADR 0040).
func intervals(v domain.SampleSet, st style, width float64) string {
	ends := intervalEnds(st.interval)
	lo, hi := roundForSummary(v.Quantile(ends[0]/100), width), roundForSummary(v.Quantile(ends[1]/100), width)
	inner := interval(st.interval, suffixed(lo), suffixed(hi), st)
	if st.interval >= tailInterval {
		return inner
	}
	outerLo, outerHi := tailEnds(v.Quantile(0.01), v.Quantile(0.99), lo, hi)
	return inner + "   " + interval(tailInterval, suffixed(outerLo), suffixed(outerHi), st)
}

const tailInterval = 98

func interval(width float64, lo, hi string, st style) string {
	return fmt.Sprintf("%s%% [%s %s %s]", shortest(width), lo, st.ellipsis, hi)
}

// tailEnds rounds P1 down and P99 up to two significant figures of their
// distance, since their noise grows with the spread, but never coarser than
// two of their own, and never inside the printed inner interval, so the 98%
// interval always contains the true tails and the inner interval (ADR 0051).
func tailEnds(p1, p99, innerLo, innerHi float64) (float64, float64) {
	if p1 == p99 {
		return roundForSummary(p1, 0), roundForSummary(p99, 0)
	}
	spread := figuresExponent(p99-p1, tailFigures)
	lo := roundOutwards(p1, min(spread, figuresExponent(p1, tailFigures)), -1)
	hi := roundOutwards(p99, min(spread, figuresExponent(p99, tailFigures)), 1)
	return math.Min(lo, innerLo), math.Max(hi, innerHi)
}

// roundOutwards rounds to the nearest step first and then moves one step in
// direction if that crossed x, so a value already on a step stays there
// instead of picking up the float error of x·10^k (ADR 0039).
func roundOutwards(x float64, exponent int, direction float64) float64 {
	steps := math.Round(toSteps(x, exponent))
	if (fromSteps(steps, exponent)-x)*direction < 0 {
		steps += direction
	}
	return fromSteps(steps, exponent)
}

func intervalEnds(width float64) []float64 {
	return []float64{50 - width/2, 50 + width/2}
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
	return suffixed(roundForSummary(x, 0))
}

// summaryFigures and widthFigures set the precision of every summary number
// but the 98% ends, which use tailFigures; the 98% clamp relies on rounding
// the inner ends exactly as they are printed (ADR 0006, ADR 0051, ADR 0052).
const (
	summaryFigures = 3
	widthFigures   = 2
	tailFigures    = 2
)

// roundForSummary keeps the finer of three significant figures of x and two
// of the inner interval's width, so the ends of a narrow interval never
// collapse into one number (ADR 0052). A zero width leaves the first rule.
func roundForSummary(x, width float64) float64 {
	if x == 0 {
		return 0
	}
	exponent := figuresExponent(x, summaryFigures)
	if width > 0 {
		exponent = min(exponent, figuresExponent(width, widthFigures))
	}
	return roundToPowerOfTen(x, exponent)
}

func suffixed(rounded float64) string {
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
	return roundToPowerOfTen(x, figuresExponent(x, digits))
}

// figuresExponent is the power of ten of the last of the first figures
// significant digits of x, the step every rounding rule is stated in.
func figuresExponent(x float64, figures int) int {
	return magnitude(x) - figures + 1
}

func magnitude(x float64) int {
	return int(math.Floor(math.Log10(math.Abs(x))))
}

func roundToPowerOfTen(x float64, exponent int) float64 {
	return fromSteps(math.Round(toSteps(x, exponent)), exponent)
}

// toSteps and fromSteps scale by an exact power of ten, dividing for negative
// exponents, because 10^-k has no exact float (ADR 0039).
func toSteps(x float64, exponent int) float64 {
	if exponent > 0 {
		return x / math.Pow(10, float64(exponent))
	}
	return x * math.Pow(10, float64(-exponent))
}

func fromSteps(steps float64, exponent int) float64 {
	if exponent > 0 {
		return steps * math.Pow(10, float64(exponent))
	}
	return steps / math.Pow(10, float64(-exponent))
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

// parseQuantiles reads a comma-separated list of percentages, such as 10,50,90.
func parseQuantiles(list string) ([]float64, error) {
	if list == "" {
		return nil, nil
	}
	var quantiles []float64
	for _, item := range strings.Split(list, ",") {
		q, err := strconv.ParseFloat(item, 64)
		if err != nil || q < 0 || q > 100 {
			return nil, fmt.Errorf("--quantiles expects percentages between 0 and 100, got '%s'", list)
		}
		quantiles = append(quantiles, q)
	}
	return quantiles, nil
}

// parseInterval reads the width of the printed interval, the default when not given.
func parseInterval(width string) (float64, error) {
	if width == "" {
		return defaultInterval, nil
	}
	w, err := strconv.ParseFloat(width, 64)
	if err != nil || !(w > 0 && w < 100) {
		return 0, fmt.Errorf("--interval expects a percentage between 0 and 100, got '%s'", width)
	}
	return w, nil
}

func quantileLabel(q float64) string {
	return "p" + strconv.FormatFloat(q, 'f', -1, 64)
}
