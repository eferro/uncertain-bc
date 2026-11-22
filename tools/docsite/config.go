package main

import (
	"errors"
	"regexp"
	"strings"
)

// The three things the site needs from mkdocs.yml, read line by line so the
// tool needs no YAML library: the config is ours and keeps this simple shape.
var (
	repoURLPattern = regexp.MustCompile(`^repo_url:\s*(\S+)$`)
	docsDirPattern = regexp.MustCompile(`^docs_dir:\s*(\S+)$`)
	navPagePattern = regexp.MustCompile(`^\s+- [^:]+:\s*(\S+\.md)$`)
)

type config struct {
	repoURL string
	docsDir string
	pages   []string
}

func parseConfig(text string) (config, error) {
	var c config
	for _, line := range strings.Split(text, "\n") {
		if m := repoURLPattern.FindStringSubmatch(line); m != nil {
			c.repoURL = m[1]
		} else if m := docsDirPattern.FindStringSubmatch(line); m != nil {
			c.docsDir = m[1]
		} else if m := navPagePattern.FindStringSubmatch(line); m != nil {
			c.pages = append(c.pages, m[1])
		}
	}
	return c, c.validate()
}

func (c config) validate() error {
	switch {
	case c.repoURL == "":
		return errors.New("no repo_url")
	case c.docsDir == "":
		return errors.New("no docs_dir")
	case len(c.pages) == 0:
		return errors.New("no pages in nav")
	}
	return nil
}

func (c config) publishes(page string) bool {
	for _, p := range c.pages {
		if p == page {
			return true
		}
	}
	return false
}
