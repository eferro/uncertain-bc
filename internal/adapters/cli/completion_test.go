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
		script, ok := completionScript(shell)
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
		script, _ := completionScript(shell)
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
	_, ok := completionScript("powershell")

	assert.False(t, ok)
}

func shellBesideCompletionScript(t *testing.T, shell, dir string, args ...string) *exec.Cmd {
	t.Helper()
	binary, err := exec.LookPath(shell)
	if err != nil {
		t.Skip(shell + " is not installed")
	}
	script, _ := completionScript(shell)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "ubc."+shell), []byte(script), 0o600))
	cmd := exec.Command(binary, args...)
	cmd.Dir = dir
	return cmd
}

func bashCompletes(t *testing.T, dir, word string) string {
	t.Helper()
	probe := `source ubc.bash; COMP_WORDS=(ubc "$1"); COMP_CWORD=1; _ubc; echo "${COMPREPLY[*]}"`
	complete := shellBesideCompletionScript(t, "bash", dir, "-c", probe, "bash", word)

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
	script, _ := completionScript("fish")

	assert.Contains(t, script, "complete -c ubc -s n -d 'Monte Carlo sample count'")
}
