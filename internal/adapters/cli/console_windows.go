package cli

import (
	"os"

	"golang.org/x/sys/windows"
)

// enableVirtualTerminal turns on ANSI escape processing for a Windows console
// while the line editor runs: x/term only puts the input in raw mode, and
// consoles that start without it (conhost, Windows PowerShell 5.1) would print
// the editor's cursor movements as text (ADR 0037). Errors leave the console
// as it was, as with a redirected stdout that is not a console.
func enableVirtualTerminal(out *os.File) (restore func()) {
	handle := windows.Handle(out.Fd())
	var mode uint32
	if windows.GetConsoleMode(handle, &mode) != nil {
		return func() {}
	}
	if windows.SetConsoleMode(handle, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING) != nil {
		return func() {}
	}
	return func() { _ = windows.SetConsoleMode(handle, mode) }
}
