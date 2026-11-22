// Command docsite assembles the pages of the documentation website from the
// repository's Markdown files, for MkDocs to build (ADR 0061).
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stderr))
}

func run(args []string, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: docsite MKDOCS_CONFIG")
		return 2
	}
	if err := buildSite(args[0]); err != nil {
		fmt.Fprintf(stderr, "docsite: %v\n", err)
		return 1
	}
	return 0
}

func buildSite(configPath string) error {
	text, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	site, err := parseConfig(string(text))
	if err != nil {
		return fmt.Errorf("%s: %w", configPath, err)
	}
	root := filepath.Dir(configPath)
	return assemble(root, site, filepath.Join(root, filepath.FromSlash(site.docsDir)))
}
