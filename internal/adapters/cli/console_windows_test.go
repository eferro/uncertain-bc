//go:build windows

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
)

func TestEnableVirtualTerminal_TurnsOnEscapeProcessing_WhenOutputIsAConsole(t *testing.T) {
	console := openConsole(t)
	before := consoleMode(t, console) &^ windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	setConsoleMode(t, console, before)

	restore := enableVirtualTerminal(console)

	assert.Equal(t, before|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING, consoleMode(t, console))
	restore()
	assert.Equal(t, before, consoleMode(t, console))
}

func TestEnableVirtualTerminal_LeavesTheFileAlone_WhenOutputIsNotAConsole(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.txt")
	file, err := os.Create(path)
	require.NoError(t, err)
	defer func() { _ = file.Close() }()

	enableVirtualTerminal(file)()

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Empty(t, content)
}

func openConsole(t *testing.T) *os.File {
	t.Helper()
	_, _, _ = windows.NewLazySystemDLL("kernel32.dll").NewProc("AllocConsole").Call()
	console, err := os.OpenFile("CONOUT$", os.O_RDWR, 0)
	require.NoError(t, err)
	original := consoleMode(t, console)
	t.Cleanup(func() {
		_ = windows.SetConsoleMode(windows.Handle(console.Fd()), original)
		_ = console.Close()
	})
	return console
}

func consoleMode(t *testing.T, console *os.File) uint32 {
	t.Helper()
	var mode uint32
	require.NoError(t, windows.GetConsoleMode(windows.Handle(console.Fd()), &mode))
	return mode
}

func setConsoleMode(t *testing.T, console *os.File, mode uint32) {
	t.Helper()
	require.NoError(t, windows.SetConsoleMode(windows.Handle(console.Fd()), mode))
}
