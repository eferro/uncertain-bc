//go:build !windows

package cli

import "os"

// enableVirtualTerminal has nothing to do outside Windows, where terminals
// always interpret escape sequences (ADR 0037).
func enableVirtualTerminal(*os.File) (restore func()) { return func() {} }
