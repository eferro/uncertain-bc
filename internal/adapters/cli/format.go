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
	if len(st.quantiles) > 0 {
		var line strings.Builder
		fmt.Fprintf(&line, "%s %s", st.approx, significant(v.Quantile(0.5)))
		for _, q := range st.quantiles {
			fmt.Fprintf(&line, "   %s %s", quantileLabel(q), significant(v.Quantile(q/100)))
		}
		return line.String()
	}
	ends := intervalEnds(st.interval)
	return fmt.Sprintf("%s %s   %s%% [%s %s %s]", st.approx, significant(v.Quantile(0.5)), shortest(st.interval),
		significant(v.Quantile(ends[0]/100)), st.ellipsis, significant(v.Quantile(ends[1]/100)))
}

// intervalEnds are the percentiles that leave the same share of samples out
// on each side of a central interval of the given width.
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
	// Only whole powers of ten are exact in binary, so rounding to tens or more
	// divides and multiplies by one instead of using an inexact 10^-n (ADR 0039).
	if decimals < 0 {
		step := math.Pow(10, float64(-decimals))
		return math.Round(x/step) * step
	}
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
