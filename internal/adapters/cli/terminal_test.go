package cli

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/application"
)

const arrowUp = "\x1b[A"

type screen struct {
	keys    io.Reader
	display bytes.Buffer
}

func (s *screen) Read(p []byte) (int, error)  { return s.keys.Read(p) }
func (s *screen) Write(p []byte) (int, error) { return s.display.Write(p) }

func typeInto(t *testing.T, historyPath, keys string) string {
	t.Helper()
	terminal := &screen{keys: strings.NewReader(keys)}
	editor := newLineEditor(terminal, loadHistory(historyPath))
	code := repl(application.NewSession(application.DefaultSettings()), unicodeStyle, editor, editor, editor)
	assert.Equal(t, 0, code)
	return terminal.display.String()
}

func TestLineEditor_EvaluatesWhatIsTyped(t *testing.T) {
	assert.Contains(t, typeInto(t, "", "1 + 2\r"), "3\r\n")
}

func TestLineEditor_RecallsThePreviousLineWithTheUpArrow(t *testing.T) {
	screen := typeInto(t, "", "x = 2\rx * 3\r"+arrowUp+"\r")

	assert.Equal(t, 2, strings.Count(screen, "6\r\n"))
}

func TestLineEditor_RecallsLinesFromAPreviousSession(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history")
	typeInto(t, path, "40 + 2\r")

	assert.Contains(t, typeInto(t, path, arrowUp+"\r"), "42\r\n")
}
