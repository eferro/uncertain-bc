package main

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
	"github.com/eferro/uncertain-bc/scenarios"
)

var version = "dev"

func main() {
	info, ok := debug.ReadBuildInfo()
	app := cli.App{
		Version:     resolveVersion(version, info, ok),
		RandomSeed:  clockSeed,
		Interactive: isTerminal(os.Stdin),
		HistoryPath: historyPath(),
		Wait:        pollEvery(300 * time.Millisecond),
		Scenarios:   scenarios.Files,
	}
	os.Exit(app.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// resolveVersion prefers the version goreleaser injects; a binary built with
// go install ...@vX.Y.Z only carries the module version in its build info.
func resolveVersion(injected string, info *debug.BuildInfo, ok bool) string {
	if injected != "dev" || !ok || info == nil {
		return injected
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		return strings.TrimPrefix(v, "v")
	}
	return injected
}

func clockSeed() uint64 {
	return uint64(time.Now().UnixNano())
}

func pollEvery(interval time.Duration) func() bool {
	return func() bool {
		time.Sleep(interval)
		return true
	}
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

// historyPath follows the XDG base directory spec (SPEC §7).
func historyPath() string {
	if state := os.Getenv("XDG_STATE_HOME"); state != "" {
		return filepath.Join(state, "ubc", "history")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "state", "ubc", "history")
}
