package install_test

import (
	"fmt"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ensureScript = "claude-plugin/skills/ubc/scripts/ensure-ubc.sh"

type ensureRun struct {
	path       string
	installDir string
	server     *httptest.Server
}

func newEnsureRun(t *testing.T) *ensureRun {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("ensure-ubc.sh is for Linux and macOS")
	}
	for _, tool := range []string{"sh", "curl", "tar", "sed", "sort"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skip(tool + " is not installed")
		}
	}
	for _, system := range []string{"/usr/bin/ubc", "/bin/ubc"} {
		if _, err := os.Stat(system); err == nil {
			t.Skip("a system-wide ubc at " + system + " hides the fake one")
		}
	}
	return &ensureRun{
		path:       filepath.Join(t.TempDir(), "path"),
		installDir: filepath.Join(t.TempDir(), "bin"),
	}
}

func fakeUbc(version string) string {
	return fmt.Sprintf("#!/bin/sh\necho '%s'\n", version)
}

func (r *ensureRun) withUbcOnPath(t *testing.T, version string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(r.path, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(r.path, "ubc"), []byte(fakeUbc(version)), 0o755)) //nolint:gosec // the fake binary must be executable
}

func (r *ensureRun) run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command("sh", append([]string{ensureScript}, args...)...)
	installer, err := filepath.Abs("install.sh")
	require.NoError(t, err)
	cmd.Env = append(os.Environ(),
		"PATH="+r.path+":/usr/bin:/bin",
		"UBC_INSTALL_DIR="+r.installDir,
		"UBC_INSTALL_SCRIPT=file://"+installer,
	)
	if r.server != nil {
		cmd.Env = append(cmd.Env, "UBC_GITHUB="+r.server.URL)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func exitCode(err error) int {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return 0
}

func TestEnsureUbc_ReportsWhereUbcIs_WhenItIsRecentEnough(t *testing.T) {
	for _, version := range []string{"ubc 1.1.1", "ubc 1.2.0", "ubc 2.0.0", "ubc v1.1.1-3-gabc1234"} {
		r := newEnsureRun(t)
		r.withUbcOnPath(t, version)

		out, err := r.run(t)

		require.NoError(t, err, out)
		assert.Contains(t, out, "ready: "+filepath.Join(r.path, "ubc"), version)
	}
}

func TestEnsureUbc_FindsUbcInTheInstallDirectory_WhenItIsNotOnPath(t *testing.T) {
	r := newEnsureRun(t)
	require.NoError(t, os.MkdirAll(r.installDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(r.installDir, "ubc"), []byte(fakeUbc("ubc 1.1.1")), 0o755)) //nolint:gosec // the fake binary must be executable

	out, err := r.run(t)

	require.NoError(t, err, out)
	assert.Contains(t, out, "ubc 1.1.1 ready: "+filepath.Join(r.installDir, "ubc"))
}

func TestEnsureUbc_FailsWithoutInstalling_WhenUbcIsMissing(t *testing.T) {
	r := newEnsureRun(t)

	out, err := r.run(t)

	assert.Equal(t, 1, exitCode(err), out)
	assert.Contains(t, out, "ubc is not installed")
	assert.Contains(t, out, "--install")
	assert.NoFileExists(t, filepath.Join(r.installDir, "ubc"))
}

func TestEnsureUbc_Fails_WhenUbcIsOlderThanTheMinimum(t *testing.T) {
	for _, version := range []string{"ubc 0.7.2", "ubc 1.1.0", "ubc v1.0.0-6-gc298a33", "something else"} {
		r := newEnsureRun(t)
		r.withUbcOnPath(t, version)

		out, err := r.run(t)

		assert.Equal(t, 1, exitCode(err), version+": "+out)
		assert.Contains(t, out, "older than 1.1.1", version)
	}
}

func TestEnsureUbc_InstallsTheLatestRelease_WhenUbcIsMissing(t *testing.T) {
	r := newEnsureRun(t)
	r.server = newFakeRelease(t, "1.2.3", fakeUbc("ubc 1.2.3"), false)

	out, err := r.run(t, "--install")

	require.NoError(t, err, out)
	assert.FileExists(t, filepath.Join(r.installDir, "ubc"))
	assert.Contains(t, out, "ubc 1.2.3 ready: "+filepath.Join(r.installDir, "ubc"))
}

func TestEnsureUbc_InstallsTheLatestRelease_WhenUbcIsTooOld(t *testing.T) {
	r := newEnsureRun(t)
	r.withUbcOnPath(t, "ubc 0.7.2")
	r.server = newFakeRelease(t, "1.2.3", fakeUbc("ubc 1.2.3"), false)

	out, err := r.run(t, "--install")

	require.NoError(t, err, out)
	assert.Contains(t, out, "ubc 1.2.3 ready: "+filepath.Join(r.installDir, "ubc"))
	assert.Contains(t, out, "another ubc comes first on PATH")
}

func TestEnsureUbc_DoesNotDownloadAnything_WhenUbcIsAlreadyRecentEnough(t *testing.T) {
	r := newEnsureRun(t)
	r.withUbcOnPath(t, "ubc 1.1.1")

	out, err := r.run(t, "--install")

	require.NoError(t, err, out)
	assert.NoFileExists(t, filepath.Join(r.installDir, "ubc"))
}

func TestEnsureUbc_Fails_WhenTheDownloadedArchiveIsCorrupt(t *testing.T) {
	r := newEnsureRun(t)
	r.server = newFakeRelease(t, "1.2.3", fakeUbc("ubc 1.2.3"), true)

	out, err := r.run(t, "--install")

	assert.Equal(t, 1, exitCode(err), out)
	assert.Contains(t, out, "checksum mismatch")
	assert.NoFileExists(t, filepath.Join(r.installDir, "ubc"))
}

func TestEnsureUbc_RejectsArgumentsOtherThanInstall(t *testing.T) {
	for _, args := range [][]string{{"--force"}, {"--install", "--install"}} {
		r := newEnsureRun(t)
		r.withUbcOnPath(t, "ubc 1.1.1")

		out, err := r.run(t, args...)

		assert.Equal(t, 2, exitCode(err), out)
		assert.Contains(t, out, "usage: ensure-ubc.sh [--install]")
	}
}

func TestEnsureUbc_UsesTheInstallDirectory_WhenTheUbcOnPathIsTooOld(t *testing.T) {
	r := newEnsureRun(t)
	r.withUbcOnPath(t, "ubc 0.7.2")
	require.NoError(t, os.MkdirAll(r.installDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(r.installDir, "ubc"), []byte(fakeUbc("ubc 1.1.1")), 0o755)) //nolint:gosec // the fake binary must be executable

	out, err := r.run(t)

	require.NoError(t, err, out)
	assert.Contains(t, out, "ubc 1.1.1 ready: "+filepath.Join(r.installDir, "ubc"))
}

func TestEnsureUbc_InstallsWithoutAPathHint_WhenTheInstallDirectoryIsFirstOnPath(t *testing.T) {
	r := newEnsureRun(t)
	r.path = r.installDir
	r.server = newFakeRelease(t, "1.2.3", fakeUbc("ubc 1.2.3"), false)

	out, err := r.run(t, "--install")

	require.NoError(t, err, out)
	assert.True(t, strings.HasSuffix(out, "ubc 1.2.3 ready: "+filepath.Join(r.installDir, "ubc")+"\n"), out)
}
