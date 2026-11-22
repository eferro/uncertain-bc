package install_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fakeBinary = "#!/bin/sh\necho fake ubc\n"

func archiveName(version string) string {
	return fmt.Sprintf("ubc_%s_%s_%s.tar.gz", version, runtime.GOOS, runtime.GOARCH)
}

func fakeArchive(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	require.NoError(t, tw.WriteHeader(&tar.Header{Name: "ubc", Mode: 0o755, Size: int64(len(fakeBinary))}))
	_, err := tw.Write([]byte(fakeBinary))
	require.NoError(t, err)
	require.NoError(t, tw.Close())
	require.NoError(t, gz.Close())
	return buf.Bytes()
}

type fakeGitHub struct {
	archive   []byte
	checksums string
}

func newFakeGitHub(t *testing.T, version string, corrupt bool) *httptest.Server {
	t.Helper()
	archive := fakeArchive(t)
	sum := sha256.Sum256(archive)
	if corrupt {
		sum[0] ^= 0xff
	}
	release := fakeGitHub{archive: archive, checksums: fmt.Sprintf("%x  %s\n", sum, archiveName(version))}
	mux := http.NewServeMux()
	mux.HandleFunc("/eferro/uncertain-bc/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/eferro/uncertain-bc/releases/tag/v"+version, http.StatusFound)
	})
	mux.HandleFunc("/eferro/uncertain-bc/releases/tag/", func(http.ResponseWriter, *http.Request) {})
	mux.HandleFunc("/eferro/uncertain-bc/releases/download/v"+version+"/"+archiveName(version), func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(release.archive)
	})
	mux.HandleFunc("/eferro/uncertain-bc/releases/download/v"+version+"/checksums.txt", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(release.checksums))
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

func runInstaller(t *testing.T, server *httptest.Server, env ...string) (string, string, error) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("install.sh is for Linux and macOS; Windows uses install.ps1")
	}
	for _, tool := range []string{"sh", "curl", "tar"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skip(tool + " is not installed")
		}
	}
	dir := filepath.Join(t.TempDir(), "bin")
	cmd := exec.Command("sh", "install.sh")
	cmd.Env = append(os.Environ(), append(env, "UBC_GITHUB="+server.URL, "UBC_INSTALL_DIR="+dir)...)
	out, err := cmd.CombinedOutput()
	return dir, string(out), err
}

func TestInstaller_InstallsTheLatestRelease(t *testing.T) {
	dir, out, err := runInstaller(t, newFakeGitHub(t, "1.2.3", false))

	require.NoError(t, err, out)
	installed, err := os.ReadFile(filepath.Join(dir, "ubc"))
	require.NoError(t, err)
	assert.Equal(t, fakeBinary, string(installed))
	assert.Contains(t, out, "ubc 1.2.3 installed to "+filepath.Join(dir, "ubc"))
}

func TestInstaller_InstallsTheRequestedVersion(t *testing.T) {
	for _, requested := range []string{"4.5.6", "v4.5.6"} {
		dir, out, err := runInstaller(t, newFakeGitHub(t, "4.5.6", false), "UBC_VERSION="+requested)

		require.NoError(t, err, out)
		assert.FileExists(t, filepath.Join(dir, "ubc"), requested)
	}
}

func TestInstaller_MakesTheBinaryExecutable(t *testing.T) {
	dir, out, err := runInstaller(t, newFakeGitHub(t, "1.2.3", false))
	require.NoError(t, err, out)

	ran, err := exec.Command(filepath.Join(dir, "ubc")).Output()

	require.NoError(t, err)
	assert.Equal(t, "fake ubc\n", string(ran))
}

func TestInstaller_RefusesAnArchiveWithTheWrongChecksum(t *testing.T) {
	dir, out, err := runInstaller(t, newFakeGitHub(t, "1.2.3", true))

	assert.Error(t, err)
	assert.Contains(t, out, "checksum mismatch")
	assert.NoFileExists(t, filepath.Join(dir, "ubc"))
}

func TestInstaller_ExplainsHowToReachTheBinary_WhenItsDirectoryIsNotOnPath(t *testing.T) {
	dir, out, err := runInstaller(t, newFakeGitHub(t, "1.2.3", false))

	require.NoError(t, err, out)
	assert.Contains(t, out, `export PATH="`+dir+`:$PATH"`)
}
