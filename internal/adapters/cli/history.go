package cli

import (
	"os"
	"path/filepath"
	"strings"
)

// fileHistory implements term.History, appending every entry to a file so
// the next session can recall it (SPEC §7, ADR 0023).
type fileHistory struct {
	path    string
	entries []string
}

func loadHistory(path string) *fileHistory {
	h := &fileHistory{path: path}
	if content, err := os.ReadFile(path); err == nil {
		for _, line := range strings.Split(string(content), "\n") {
			if strings.TrimSpace(line) != "" {
				h.entries = append(h.entries, line)
			}
		}
	}
	return h
}

func (h *fileHistory) Add(entry string) {
	if strings.TrimSpace(entry) == "" {
		return
	}
	h.entries = append(h.entries, entry)
	h.persist(entry)
}

func (h *fileHistory) Len() int { return len(h.entries) }

func (h *fileHistory) At(idx int) string { return h.entries[len(h.entries)-1-idx] }

// persist ignores errors on purpose: a history file that cannot be written
// must never break the session (ADR 0023).
func (h *fileHistory) persist(entry string) {
	if h.path == "" || os.MkdirAll(filepath.Dir(h.path), 0o700) != nil {
		return
	}
	file, err := os.OpenFile(h.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	_, _ = file.WriteString(entry + "\n")
	_ = file.Close()
}
