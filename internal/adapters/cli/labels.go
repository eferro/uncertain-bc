package cli

import (
	"regexp"

	"github.com/eferro/uncertain-bc/internal/domain"
)

// maxLabelLength keeps a label short enough to read as a name; a longer
// expression is better pointed at by its line (ADR 0053).
const maxLabelLength = 30

var variableName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// labelled names an expression after its source text, as an assignment is
// named after its variable, when that text reads as a name.
func labelled(result domain.Result) (domain.Result, bool) {
	if len(result.Source) > maxLabelLength && !variableName.MatchString(result.Source) {
		return result, false
	}
	result.Name = result.Source
	return result, true
}
