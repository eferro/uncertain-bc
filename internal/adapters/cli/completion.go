package cli

import (
	"flag"
	"fmt"
	"slices"
	"strings"
)

type flagHelp struct {
	name   string
	usage  string
	values []string
}

// scenarioFlags take the name of a built-in scenario, which completion
// offers from the catalogue (ADR 0054).
var scenarioFlags = []string{"scenario", "scenario-source"}

// flagsForCompletion reads the real flag set, so completion never drifts
// from the flags ubc accepts.
func flagsForCompletion(scenarios []string) []flagHelp {
	fs := flag.NewFlagSet("ubc", flag.ContinueOnError)
	registerFlags(fs, &options{})
	names := append([]string{formulas}, scenarios...)
	var flags []flagHelp
	fs.VisitAll(func(f *flag.Flag) {
		help := flagHelp{name: f.Name, usage: strings.ReplaceAll(f.Usage, "'", "")}
		if slices.Contains(scenarioFlags, f.Name) {
			help.values = names
		}
		flags = append(flags, help)
	})
	return flags
}

func dashed(name string) string {
	if len(name) == 1 {
		return "-" + name
	}
	return "--" + name
}

var completionScripts = map[string]func(flags []flagHelp) string{
	"bash": bashCompletion,
	"zsh":  zshCompletion,
	"fish": fishCompletion,
}

func completionScript(shell string, scenarios []string) (string, bool) {
	script, ok := completionScripts[shell]
	if !ok {
		return "", false
	}
	return script(flagsForCompletion(scenarios)), true
}

func bashCompletion(flags []flagHelp) string {
	words := make([]string, len(flags))
	var valued, values []string
	for i, f := range flags {
		words[i] = dashed(f.name)
		if f.values != nil {
			valued, values = append(valued, words[i]), f.values
		}
	}
	return fmt.Sprintf(`_ubc() {
  local cur="${COMP_WORDS[COMP_CWORD]}" prev="${COMP_WORDS[COMP_CWORD-1]}"
  if [[ " %s " == *" $prev "* ]]; then
    COMPREPLY=($(compgen -W "%s" -- "$cur"))
  elif [[ $cur == -* ]]; then
    COMPREPLY=($(compgen -W "%s" -- "$cur"))
  else
    COMPREPLY=($(compgen -f -- "$cur"))
  fi
}
complete -o filenames -F _ubc ubc
`, strings.Join(valued, " "), strings.Join(values, " "), strings.Join(words, " "))
}

func zshCompletion(flags []flagHelp) string {
	var script strings.Builder
	script.WriteString("#compdef ubc\n_arguments \\\n")
	for _, f := range flags {
		fmt.Fprintf(&script, "  '%s[%s]%s' \\\n", dashed(f.name), f.usage, zshValues(f.values))
	}
	script.WriteString("  '*:model file:_files'\n")
	return script.String()
}

func zshValues(values []string) string {
	if values == nil {
		return ""
	}
	return ":scenario:(" + strings.Join(values, " ") + ")"
}

func fishCompletion(flags []flagHelp) string {
	var script strings.Builder
	for _, f := range flags {
		option := "-l"
		if len(f.name) == 1 {
			option = "-s"
		}
		fmt.Fprintf(&script, "complete -c ubc %s %s -d '%s'%s\n", option, f.name, f.usage, fishValues(f.values))
	}
	return script.String()
}

func fishValues(values []string) string {
	if values == nil {
		return ""
	}
	return " -x -a '" + strings.Join(values, " ") + "'"
}
