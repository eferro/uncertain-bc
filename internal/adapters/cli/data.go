package cli

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/eferro/uncertain-bc/internal/application"
)

type dataset struct {
	name string
	path string
}

func parseDatasets(values []string) ([]dataset, error) {
	datasets := make([]dataset, 0, len(values))
	for _, value := range values {
		name, path, err := splitNamed("--data", "name=file", value)
		if err != nil {
			return nil, err
		}
		if path == "" {
			return nil, fmt.Errorf("--data %s: no file name", name)
		}
		datasets = append(datasets, dataset{name: name, path: path})
	}
	return datasets, nil
}

func loadDataset(session *application.Session, d dataset) error {
	values, err := readNumbers(d.path)
	if err != nil {
		return err
	}
	if err := session.OverrideData(d.name, values); err != nil {
		return fmt.Errorf("%s: %w", d.path, err)
	}
	return nil
}

// readNumbers reads one plain number per line, the way spreadsheets and
// scripts export a column. Blank lines are skipped; anything else that is
// not a finite number is an error, so a stray header is never data (ADR 0031).
func readNumbers(path string) ([]float64, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var values []float64
	for i, line := range strings.Split(string(content), "\n") {
		field := strings.TrimSpace(line)
		if field == "" {
			continue
		}
		x, err := strconv.ParseFloat(field, 64)
		if err != nil || math.IsNaN(x) || math.IsInf(x, 0) {
			return nil, fmt.Errorf("%s: line %d: '%s' is not a number", path, i+1, field)
		}
		values = append(values, x)
	}
	return values, nil
}
