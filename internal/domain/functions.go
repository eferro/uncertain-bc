package domain

import "fmt"

type call struct {
	name string
	args []node
	at   position
}

type builtin func(c call, args []Value) (Value, error)

var builtins = map[string]builtin{
	"P": probabilityOf,
}

func (c call) evaluate(env Env) (Value, error) {
	fn, ok := builtins[c.name]
	if !ok {
		return nil, evalErrorAt(c.at, "unknown function '%s'", c.name)
	}
	args := make([]Value, len(c.args))
	for i, arg := range c.args {
		value, err := arg.evaluate(env)
		if err != nil {
			return nil, err
		}
		args[i] = value
	}
	return fn(c, args)
}

func (c call) expectArgs(args []Value, want int) error {
	if len(args) != want {
		return evalErrorAt(c.at, "%s expects %s, got %d", c.name, plural(want, "argument"), len(args))
	}
	return nil
}

func plural(n int, noun string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", noun)
	}
	return fmt.Sprintf("%d %ss", n, noun)
}

func probabilityOf(c call, args []Value) (Value, error) {
	if err := c.expectArgs(args, 1); err != nil {
		return nil, err
	}
	if !isIndicator(args[0]) {
		return nil, evalErrorAt(c.at, "P expects a comparison")
	}
	return Probability(mean(args[0])), nil
}
