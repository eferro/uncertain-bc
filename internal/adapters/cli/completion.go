package cli

import (
	"flag"
	"fmt"
	"strings"
)

type flagHelp struct {
	name  string
	usage string
}

// flagsForCompletion reads the real flag set, so completion never drifts
// from the flags ubc accepts.
func flagsForCompletion() []flagHelp {
	fs := flag.NewFlagSet("ubc", flag.ContinueOnError)
	registerFlags(fs, &options{})
	var flags []flagHelp
	fs.VisitAll(func(f *flag.Flag) { flags = append(flags, flagHelp{f.Name, strings.ReplaceAll(f.Usage, "'", "")}) })
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

func completionScript(shell string) (string, bool) {
	script, ok := completionScripts[shell]
	if !ok {
		return "", false
	}
	return script(flagsForCompletion()), true
}

func bashCompletion(flags []flagHelp) string {
	words := make([]string, len(flags))
	for i, f := range flags {
		words[i] = dashed(f.name)
	}
	return fmt.Sprintf(`_ubc() {
  local cur="${COMP_WORDS[COMP_CWORD]}"
  if [[ $cur == -* ]]; then
    COMPREPLY=($(compgen -W "%s" -- "$cur"))
  else
    COMPREPLY=($(compgen -f -- "$cur"))
  fi
}
complete -o filenames -F _ubc ubc
`, strings.Join(words, " "))
}

func zshCompletion(flags []flagHelp) string {
	var script strings.Builder
	script.WriteString("#compdef ubc\n_arguments \\\n")
	for _, f := range flags {
		fmt.Fprintf(&script, "  '%s[%s]' \\\n", dashed(f.name), f.usage)
	}
	script.WriteString("  '*:model file:_files'\n")
	return script.String()
}

func fishCompletion(flags []flagHelp) string {
	var script strings.Builder
	for _, f := range flags {
		option := "-l"
		if len(f.name) == 1 {
			option = "-s"
		}
		fmt.Fprintf(&script, "complete -c ubc %s %s -d '%s'\n", option, f.name, f.usage)
	}
	return script.String()
}
