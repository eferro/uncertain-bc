package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

var version = "dev"

func main() {
	app := cli.App{
		Version:     version,
		RandomSeed:  clockSeed,
		Interactive: isTerminal(os.Stdin),
		HistoryPath: historyPath(),
		Wait:        pollEvery(300 * time.Millisecond),
	}
	os.Exit(app.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
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
