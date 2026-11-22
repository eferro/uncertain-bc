package main

import (
	"fmt"
	"regexp"
	"strings"
)

// A Keep a Changelog heading: "## [1.2.0] - 2026-01-31" or "## [Unreleased]".
var headingPattern = regexp.MustCompile(`^## \[([^\]]+)\](?: - (\S+))?$`)

// Link reference definitions are invisible in rendered Markdown, so they never
// belong to a version's notes.
var linkDefinitionPattern = regexp.MustCompile(`^\[[^\]]+\]: `)

type section struct {
	version string
	date    string
	body    string
}

func sections(changelog string) []section {
	var all []section
	var body []string
	for _, line := range strings.Split(changelog, "\n") {
		if match := headingPattern.FindStringSubmatch(line); match != nil {
			all = closeLast(all, body)
			all = append(all, section{version: match[1], date: match[2]})
			body = nil
		} else if !linkDefinitionPattern.MatchString(line) {
			body = append(body, line)
		}
	}
	return closeLast(all, body)
}

func closeLast(all []section, body []string) []section {
	if len(all) > 0 {
		all[len(all)-1].body = strings.TrimSpace(strings.Join(body, "\n"))
	}
	return all
}

func notesFor(changelog, version string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	for _, s := range sections(changelog) {
		if s.version != version {
			continue
		}
		if s.body == "" {
			return "", fmt.Errorf("the %s section is empty", version)
		}
		return s.body + "\n", nil
	}
	return "", fmt.Errorf("no section for %s", version)
}
