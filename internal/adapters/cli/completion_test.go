package cli

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func registeredFlags() []string {
	fs := flag.NewFlagSet("ubc", flag.ContinueOnError)
	registerFlags(fs, &options{})
	var names []string
	fs.VisitAll(func(f *flag.Flag) { names = append(names, f.Name) })
	return names
}

func TestCompletion_OffersEveryFlagInEveryShell(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish"} {
		script, ok := completionScript(shell, nil)
		require.True(t, ok, shell)
		for _, name := range registeredFlags() {
			assert.Contains(t, script, name, "%s completion misses %s", shell, name)
		}
	}
}

func TestCompletion_RegistersItselfForUbc(t *testing.T) {
	cases := map[string]string{
		"bash": "complete -o filenames -F _ubc ubc",
		"zsh":  "#compdef ubc",
		"fish": "complete -c ubc -l samples -d 'Monte Carlo sample count'",
	}
	for shell, line := range cases {
		script, _ := completionScript(shell, nil)
		assert.Contains(t, script, line, shell)
	}
}

func TestCompletion_ProducesValidScripts(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish"} {
		t.Run(shell, func(t *testing.T) {
			check := shellBesideCompletionScript(t, shell, t.TempDir(), "-n", "ubc."+shell)

			out, err := check.CombinedOutput()

			assert.NoError(t, err, string(out))
		})
	}
}

func TestCompletion_KnowsOnlyThreeShells(t *testing.T) {
	_, ok := completionScript("powershell", nil)

	assert.False(t, ok)
}

var completedScenarios = []string{"team/hire", "team/replace", "delivery/date"}

func shellBesideCompletionScript(t *testing.T, shell, dir string, args ...string) *exec.Cmd {
	t.Helper()
	binary, err := exec.LookPath(shell)
	if err != nil {
		t.Skip(shell + " is not installed")
	}
	script, _ := completionScript(shell, completedScenarios)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "ubc."+shell), []byte(script), 0o600))
	cmd := exec.Command(binary, args...)
	cmd.Dir = dir
	return cmd
}

func bashCompletes(t *testing.T, dir string, words ...string) string {
	t.Helper()
	probe := `source ubc.bash; COMP_WORDS=(ubc "$@"); COMP_CWORD=$#; _ubc; echo "${COMPREPLY[*]}"`
	complete := shellBesideCompletionScript(t, "bash", dir, append([]string{"-c", probe, "bash"}, words...)...)

	out, err := complete.Output()

	require.NoError(t, err)
	return strings.TrimSpace(string(out))
}

func TestBashCompletion_CompletesLongFlags(t *testing.T) {
	assert.Equal(t, "--samples --samples-out", bashCompletes(t, t.TempDir(), "--sa"))
}

func TestBashCompletion_CompletesFileNames(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "model.ubc"), nil, 0o600))

	assert.Equal(t, "model.ubc", bashCompletes(t, dir, "mo"))
}

func TestFishCompletion_MarksShortFlags(t *testing.T) {
	script, _ := completionScript("fish", nil)

	assert.Contains(t, script, "complete -c ubc -s n -d 'Monte Carlo sample count'")
}

func TestBashCompletion_CompletesScenarioNames(t *testing.T) {
	for _, flag := range []string{"--scenario", "--scenario-source"} {
		assert.Equal(t, "team/hire team/replace", bashCompletes(t, t.TempDir(), flag, "te"), flag)
	}
}

func TestBashCompletion_CompletesFileNames_AfterAScenario(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "model.ubc"), nil, 0o600))

	assert.Equal(t, "model.ubc", bashCompletes(t, dir, "--scenario", "formulas", "mo"))
}

func TestZshCompletion_OffersScenarioNames(t *testing.T) {
	script, _ := completionScript("zsh", completedScenarios)

	assert.Contains(t, script, "  '--scenario[run a built-in scenario after the shared formulas]:scenario:(formulas team/hire team/replace delivery/date)' \\\n")
	assert.Contains(t, script, "  '--scenario-source[print the source of a built-in scenario]:scenario:(formulas team/hire team/replace delivery/date)' \\\n")
}

func TestFishCompletion_OffersScenarioNames(t *testing.T) {
	script, _ := completionScript("fish", completedScenarios)

	assert.Contains(t, script, "complete -c ubc -l scenario -d 'run a built-in scenario after the shared formulas' -x -a 'formulas team/hire team/replace delivery/date'\n")
	assert.Contains(t, script, "complete -c ubc -l scenario-source -d 'print the source of a built-in scenario' -x -a 'formulas team/hire team/replace delivery/date'\n")
}
