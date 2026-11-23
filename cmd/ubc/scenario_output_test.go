package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateScenarioOutputs = flag.Bool("update-scenario-outputs", false,
	"rewrite the output block at the end of every scenario (make scenario-outputs)")

func TestScenarios_EveryScenarioEndsWithItsOutput(t *testing.T) {
	files, err := modelFiles(scenariosDir)
	require.NoError(t, err)
	for _, file := range scenarioFiles(files) {
		t.Run(file, func(t *testing.T) {
			path := filepath.Join(scenariosDir, filepath.FromSlash(file))
			source, err := os.ReadFile(path)
			require.NoError(t, err)

			want := withOutputBlock(string(source), scenarioOutput(t, file))

			if *updateScenarioOutputs {
				require.NoError(t, os.WriteFile(path, []byte(want), 0o600))
				source = []byte(want)
			}
			assert.Equal(t, want, string(source), "stale output block: run make scenario-outputs")
		})
	}
}

func TestWithOutputBlock_AppendsTheOutputAsCommentsAfterABlankLine(t *testing.T) {
	source := "x = 1 to 2   # why\n\n# Result\nx\nP(x > 1)\n"

	got := withOutputBlock(source, "x = ≈ 1.5\nP(x > 1) = 0.95\n")

	assert.Equal(t, source+"\n# Output (seed 1, 5000 samples):\n#   x = ≈ 1.5\n#   P(x > 1) = 0.95\n", got)
}

func TestWithOutputBlock_ReplacesAStaleBlock(t *testing.T) {
	source := "x\n\n# Output (seed 1, 5000 samples):\n#   x = ≈ 9\n#   old line\n"

	got := withOutputBlock(source, "x = ≈ 1.5\n")

	assert.Equal(t, "x\n\n# Output (seed 1, 5000 samples):\n#   x = ≈ 1.5\n", got)
}

func TestWithOutputBlock_LeavesOneBlankLineBeforeTheBlock_WhenTheSourceEndsWithSeveral(t *testing.T) {
	got := withOutputBlock("x\n\n\n", "x = 1\n")

	assert.Equal(t, "x\n\n# Output (seed 1, 5000 samples):\n#   x = 1\n", got)
}

func TestWithOutputBlock_KeepsALineThatOnlyQuotesTheHeader(t *testing.T) {
	source := "x   # see # Output (seed 1, 5000 samples): below\n"

	got := withOutputBlock(source, "x = 1\n")

	assert.Equal(t, source+"\n# Output (seed 1, 5000 samples):\n#   x = 1\n", got)
}

func TestWithOutputBlock_KeepsAnUpToDateBlockAsItIs(t *testing.T) {
	source := "x\n\n# Output (seed 1, 5000 samples):\n#   x = 1\n"

	assert.Equal(t, source, withOutputBlock(source, "x = 1\n"))
}

const outputHeader = "# Output (seed 1, 5000 samples):"

func withOutputBlock(source, output string) string {
	body, _, _ := strings.Cut(source, "\n"+outputHeader)
	var block strings.Builder
	block.WriteString(strings.TrimRight(body, "\n") + "\n\n" + outputHeader + "\n")
	for line := range strings.Lines(output) {
		block.WriteString("#   " + line)
	}
	return block.String()
}

func scenarioOutput(t *testing.T, file string) string {
	t.Helper()
	code, stdout, stderr := runScenario(file)
	require.Equal(t, 0, code, stderr)
	return stdout
}
