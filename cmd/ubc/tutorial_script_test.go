package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rogpeppe/go-internal/txtar"
)

var errSourceNotRun = errors.New("ubc block is not followed by a console block that runs it")

type fencedBlock struct {
	info  string
	lines []string
}

type tutorialScript struct {
	commands strings.Builder
	files    []txtar.File
	source   *fencedBlock
	outputs  int
	sources  int
}

func tutorialArchive(markdown string) (*txtar.Archive, error) {
	script := &tutorialScript{}
	for _, block := range fencedBlocks(markdown) {
		if err := script.add(block); err != nil {
			return nil, err
		}
	}
	if script.source != nil {
		return nil, errSourceNotRun
	}
	return &txtar.Archive{Comment: []byte(script.commands.String()), Files: script.files}, nil
}

func fencedBlocks(markdown string) []fencedBlock {
	var blocks []fencedBlock
	var current *fencedBlock
	for _, line := range strings.Split(markdown, "\n") {
		switch {
		case current == nil && strings.HasPrefix(line, "```"):
			current = &fencedBlock{info: strings.TrimPrefix(line, "```")}
		case current != nil && line == "```":
			blocks = append(blocks, *current)
			current = nil
		case current != nil:
			current.lines = append(current.lines, line)
		}
	}
	return blocks
}

func (s *tutorialScript) add(block fencedBlock) error {
	switch {
	case s.source != nil && block.info != "console":
		return errSourceNotRun
	case block.info == "ubc":
		s.source = &block
	case block.info == "console":
		return s.addConsole(block.lines)
	}
	return nil
}

func (s *tutorialScript) addConsole(lines []string) error {
	if len(lines) == 0 {
		return errors.New("console block is empty")
	}
	if !strings.HasPrefix(lines[0], "$ ") {
		return fmt.Errorf("console block starts with %q instead of a $ command", lines[0])
	}
	start := 0
	for i := 1; i <= len(lines); i++ {
		if i == len(lines) || strings.HasPrefix(lines[i], "$ ") {
			if err := s.addCommand(strings.TrimPrefix(lines[start], "$ "), lines[start+1:i]); err != nil {
				return err
			}
			start = i
		}
	}
	return nil
}

func (s *tutorialScript) addCommand(command string, output []string) error {
	if !strings.HasPrefix(command, "ubc ") {
		return fmt.Errorf("console command %q does not run ubc", command)
	}
	if err := s.checkSource(command); err != nil {
		return err
	}
	fmt.Fprintf(&s.commands, "exec %s\n", command)
	if len(output) == 0 {
		s.commands.WriteString("! stdout .\n")
		return nil
	}
	s.outputs++
	name := fmt.Sprintf("output%d", s.outputs)
	fmt.Fprintf(&s.commands, "cmp stdout %s\n", name)
	s.files = append(s.files, txtar.File{Name: name, Data: joinLines(output)})
	return nil
}

func (s *tutorialScript) checkSource(command string) error {
	if s.source == nil {
		return nil
	}
	file := lastUbcFile(command)
	if file == "" {
		return fmt.Errorf("console command %q does not run the ubc file shown before it", command)
	}
	s.sources++
	name := fmt.Sprintf("source%d.ubc", s.sources)
	fmt.Fprintf(&s.commands, "cmp %s %s\n", file, name)
	s.files = append(s.files, txtar.File{Name: name, Data: joinLines(s.source.lines)})
	s.source = nil
	return nil
}

func lastUbcFile(command string) string {
	file := ""
	for _, field := range strings.Fields(command) {
		if strings.HasSuffix(field, ".ubc") {
			file = field
		}
	}
	return file
}

func joinLines(lines []string) []byte {
	return []byte(strings.Join(lines, "\n") + "\n")
}
