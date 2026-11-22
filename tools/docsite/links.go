package main

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// An inline Markdown link's target: "[text](target)".
var linkPattern = regexp.MustCompile(`\]\(([^)\s]+)\)`)

type targetKind int

const (
	file targetKind = iota + 1
	directory
)

type repository interface {
	kind(path string) (targetKind, bool)
}

// rewriteLinks keeps the links between published pages, which MkDocs resolves,
// and points every other relative link at the file on GitHub, so the site has
// no dead links to files it does not publish.
func rewriteLinks(page, text string, site config, repo repository) (string, error) {
	var broken error
	rewritten := linkPattern.ReplaceAllStringFunc(text, func(link string) string {
		target := linkPattern.FindStringSubmatch(link)[1]
		url, err := resolve(page, target, site, repo)
		if err != nil && broken == nil {
			broken = err
		}
		return "](" + url + ")"
	})
	return rewritten, broken
}

func resolve(page, target string, site config, repo repository) (string, error) {
	if strings.Contains(target, ":") || strings.HasPrefix(target, "#") {
		return target, nil
	}
	file, anchor, _ := strings.Cut(target, "#")
	repoPath := path.Join(path.Dir(page), file)
	if site.publishes(repoPath) {
		return target, nil
	}
	kind, ok := repo.kind(repoPath)
	if !ok || !filepath.IsLocal(repoPath) {
		return target, fmt.Errorf("broken link to %s", repoPath)
	}
	return githubURL(site.repoURL, kind, file, repoPath) + withHash(anchor), nil
}

func githubURL(repoURL string, kind targetKind, target, repoPath string) string {
	view := "blob"
	if kind == directory {
		view = "tree"
	}
	if strings.HasSuffix(target, "/") {
		repoPath += "/"
	}
	return fmt.Sprintf("%s/%s/main/%s", repoURL, view, repoPath)
}

func withHash(anchor string) string {
	if anchor == "" {
		return ""
	}
	return "#" + anchor
}
