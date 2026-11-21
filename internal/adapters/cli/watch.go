package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/eferro/uncertain-bc/internal/application"
)

type stamp struct {
	modified time.Time
	size     int64
}

// watch evaluates the sources, then evaluates them again whenever a watched
// file changes. Wait sleeps between checks and returns false to stop.
func (a App) watch(opts options, settings application.Settings, st style, stdout, stderr io.Writer) int {
	code := evaluateSources(opts, settings, st, nil, stdout, stderr)
	stamps := stampsOf(opts.files)
	for a.Wait != nil && a.Wait() {
		current := stampsOf(opts.files)
		if changed := changedFiles(opts.files, stamps, current); len(changed) > 0 {
			fmt.Fprintf(stdout, "--- %s changed ---\n", strings.Join(changed, ", "))
			code = evaluateSources(opts, settings, st, nil, stdout, stderr)
			stamps = current
		}
	}
	return code
}

func stampsOf(files []string) map[string]stamp {
	stamps := make(map[string]stamp, len(files))
	for _, path := range files {
		if info, err := os.Stat(path); err == nil {
			stamps[path] = stamp{modified: info.ModTime(), size: info.Size()}
		}
	}
	return stamps
}

func changedFiles(files []string, before, after map[string]stamp) []string {
	var changed []string
	for _, path := range files {
		if before[path] != after[path] {
			changed = append(changed, path)
		}
	}
	return changed
}
