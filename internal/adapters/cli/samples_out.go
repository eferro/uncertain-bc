package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eferro/uncertain-bc/internal/application"
	"github.com/eferro/uncertain-bc/internal/domain"
)

// A samplesOut writes every variable's samples after the program, so they
// can be analysed with other tools (ADR 0032).
type samplesOut struct {
	path   string
	encode func(columns []column) string
}

type column struct {
	name    string
	samples []float64
}

var tableEncoders = map[string]func(columns []column) string{
	".csv":  csvTable,
	".json": jsonTable,
}

func parseSamplesOut(path string) (*samplesOut, error) {
	if path == "" {
		return nil, nil
	}
	encode, ok := tableEncoders[filepath.Ext(path)]
	if !ok {
		return nil, fmt.Errorf("--samples-out writes .csv or .json files, got '%s'", path)
	}
	return &samplesOut{path: path, encode: encode}, nil
}

func (o samplesOut) write(session *application.Session, samples int) error {
	columns := columnsOf(session.VariablesAndLast(), samples)
	if len(columns) == 0 {
		return errors.New("nothing to export")
	}
	return os.WriteFile(o.path, []byte(o.encode(columns)), 0o644)
}

// columnsOf lines the values up sample by sample, repeating plain numbers,
// so that row i holds sample i of every variable and keeps their correlation.
func columnsOf(variables []application.Variable, samples int) []column {
	columns := make([]column, 0, len(variables))
	for _, v := range variables {
		columns = append(columns, column{name: v.Name, samples: samplesOf(v.Value, samples)})
	}
	return columns
}

func samplesOf(value domain.Value, samples int) []float64 {
	switch v := value.(type) {
	case domain.SampleSet:
		return v
	case domain.Probability:
		return repeated(float64(v), samples)
	}
	return repeated(float64(value.(domain.Scalar)), samples)
}

func repeated(x float64, samples int) []float64 {
	values := make([]float64, samples)
	for i := range values {
		values[i] = x
	}
	return values
}

// csvTable leaves non-finite samples empty, the usual missing value in CSV.
func csvTable(columns []column) string {
	var b strings.Builder
	cells := make([]string, len(columns))
	for i, c := range columns {
		cells[i] = c.name
	}
	b.WriteString(strings.Join(cells, ",") + "\n")
	for row := range columns[0].samples {
		for i, c := range columns {
			cells[i] = csvNumber(c.samples[row])
		}
		b.WriteString(strings.Join(cells, ",") + "\n")
	}
	return b.String()
}

func csvNumber(x float64) string {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return ""
	}
	return strconv.FormatFloat(x, 'g', -1, 64)
}

// jsonTable keeps the variables in order, which a Go map would not.
func jsonTable(columns []column) string {
	fields := make([]string, len(columns))
	for i, c := range columns {
		name, _ := json.Marshal(c.name)
		values := make([]*float64, len(c.samples))
		for j, x := range c.samples {
			values[j] = finite(x)
		}
		array, _ := json.Marshal(values)
		fields[i] = string(name) + ":" + string(array)
	}
	return "{" + strings.Join(fields, ",") + "}\n"
}

func exportSamples(session *application.Session, out *samplesOut, samples int, stderr io.Writer) int {
	if out == nil {
		return exitOK
	}
	if err := out.write(session, samples); err != nil {
		reportError(source{name: "--samples-out " + out.path}, err, stderr)
		return exitEvalError
	}
	return exitOK
}
