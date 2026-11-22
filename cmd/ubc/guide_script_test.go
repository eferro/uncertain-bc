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

type guideScript struct {
	commands strings.Builder
	files    []txtar.File
	source   *fencedBlock
	outputs  int
	messages int
	sources  int
}

func guideArchive(markdown string) (*txtar.Archive, error) {
	script := &guideScript{}
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

func (s *guideScript) add(block fencedBlock) error {
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

type consoleCommand struct {
	command string
	output  []string
}

func (s *guideScript) addConsole(lines []string) error {
	commands, err := consoleCommands(lines)
	if err != nil {
		return err
	}
	for _, c := range commands {
		if err := s.addCommand(c.command, c.output); err != nil {
			return err
		}
	}
	return nil
}

func consoleCommands(lines []string) ([]consoleCommand, error) {
	if len(lines) == 0 {
		return nil, errors.New("console block is empty")
	}
	if !strings.HasPrefix(lines[0], "$ ") {
		return nil, fmt.Errorf("console block starts with %q instead of a $ command", lines[0])
	}
	var commands []consoleCommand
	for _, line := range lines {
		if command, ok := strings.CutPrefix(line, "$ "); ok {
			commands = append(commands, consoleCommand{command: command})
			continue
		}
		last := &commands[len(commands)-1]
		last.output = append(last.output, line)
	}
	return commands, nil
}

func (s *guideScript) addCommand(command string, output []string) error {
	if !strings.HasPrefix(command, "ubc ") {
		return fmt.Errorf("console command %q does not run ubc", command)
	}
	if err := s.checkSource(command); err != nil {
		return err
	}
	stdout, stderr := splitErrors(output)
	s.exec(command, len(stderr) > 0)
	s.expectStdout(stdout)
	s.expectStderr(stderr)
	return nil
}

func (s *guideScript) exec(command string, fails bool) {
	if fails {
		s.commands.WriteString("! ")
	}
	fmt.Fprintf(&s.commands, "exec %s\n", command)
}

func (s *guideScript) expectStdout(lines []string) {
	if len(lines) == 0 {
		s.commands.WriteString("! stdout .\n")
		return
	}
	s.outputs++
	s.expect("stdout", fmt.Sprintf("output%d", s.outputs), lines)
}

func (s *guideScript) expectStderr(lines []string) {
	if len(lines) == 0 {
		return
	}
	s.messages++
	s.expect("stderr", fmt.Sprintf("error%d", s.messages), lines)
}

// splitErrors separates the lines ubc writes to stderr, which all start with
// "ubc: ", from its results (ADR 0057).
func splitErrors(output []string) (stdout, stderr []string) {
	for _, line := range output {
		if strings.HasPrefix(line, "ubc: ") {
			stderr = append(stderr, line)
		} else {
			stdout = append(stdout, line)
		}
	}
	return stdout, stderr
}

func (s *guideScript) expect(stream, name string, lines []string) {
	fmt.Fprintf(&s.commands, "cmp %s %s\n", stream, name)
	s.files = append(s.files, txtar.File{Name: name, Data: joinLines(lines)})
}

func (s *guideScript) checkSource(command string) error {
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
