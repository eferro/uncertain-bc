package cli

import (
	"encoding/json"
	"math"

	"github.com/eferro/uncertain-bc/internal/domain"
)

type numberJSON struct {
	Kind  string   `json:"kind"`
	Value *float64 `json:"value"`
}

type samplesJSON struct {
	Kind    string   `json:"kind"`
	Median  *float64 `json:"median"`
	Mean    *float64 `json:"mean"`
	Stdev   *float64 `json:"stdev"`
	P5      *float64 `json:"p5"`
	P95     *float64 `json:"p95"`
	N       int      `json:"n"`
	Dropped int      `json:"dropped"`
}

// jsonLine describes a result as one JSON object (SPEC §6).
func jsonLine(value domain.Value) string {
	var object any
	switch v := value.(type) {
	case domain.SampleSet:
		object = samplesJSON{
			Kind: "samples", Median: finite(v.Quantile(0.5)), Mean: finite(v.Mean()), Stdev: finite(v.Stdev()),
			P5: finite(v.Quantile(0.05)), P95: finite(v.Quantile(0.95)), N: len(v), Dropped: v.Dropped(),
		}
	case domain.Probability:
		object = numberJSON{Kind: "probability", Value: finite(float64(v))}
	case domain.Scalar:
		object = numberJSON{Kind: "scalar", Value: finite(float64(v))}
	}
	line, _ := json.Marshal(object)
	return string(line)
}

// finite turns NaN and ±Inf into null, which JSON can represent.
func finite(x float64) *float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return nil
	}
	return &x
}
