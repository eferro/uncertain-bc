package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type dirRepository string

func (root dirRepository) kind(path string) (targetKind, bool) {
	info, err := os.Stat(filepath.Join(string(root), filepath.FromSlash(path)))
	if err != nil {
		return 0, false
	}
	if info.IsDir() {
		return directory, true
	}
	return file, true
}

// assemble writes every page of the nav into out at the same relative path it
// has in the repository, so the links between published pages keep working.
func assemble(root string, site config, out string) error {
	for _, page := range site.pages {
		if err := assemblePage(root, site, page, out); err != nil {
			return err
		}
	}
	return nil
}

func assemblePage(root string, site config, page, out string) error {
	text, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(page)))
	if err != nil {
		return err
	}
	rewritten, err := rewriteLinks(page, string(text), site, dirRepository(root))
	if err != nil {
		return fmt.Errorf("%s: %w", page, err)
	}
	dest := filepath.Join(out, filepath.FromSlash(page))
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dest, []byte(rewritten), 0o644) //nolint:gosec // G306: a website page, meant to be world-readable
}
