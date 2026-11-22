// Command releasenotes prints the section of one version from CHANGELOG.md,
// so the GitHub release says the same as the changelog (ADR 0060).
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) != 2 {
		fmt.Fprintln(stderr, "usage: releasenotes CHANGELOG VERSION")
		return 2
	}
	changelog, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "releasenotes: %v\n", err)
		return 1
	}
	notes, err := notesFor(string(changelog), args[1])
	if err != nil {
		fmt.Fprintf(stderr, "releasenotes: %s: %v\n", args[0], err)
		return 1
	}
	fmt.Fprint(stdout, notes)
	return 0
}
