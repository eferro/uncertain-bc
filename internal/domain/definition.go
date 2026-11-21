package domain

import "maps"

// maxDepth stops runaway recursion: without conditionals, a function that
// calls itself can never finish.
const maxDepth = 64

type definition struct {
	name   string
	params []token
	body   node
	at     position
}

func (e Env) define(d definition) error {
	if _, builtIn := lookup(d.name); builtIn {
		return evalErrorAt(d.at, "cannot redefine built-in function '%s'", d.name)
	}
	seen := map[string]bool{}
	for _, param := range d.params {
		if seen[param.text] {
			return evalErrorAt(param.position, "parameter '%s' appears twice", param.text)
		}
		seen[param.text] = true
	}
	e.functions[d.name] = d
	return nil
}

func (e Env) lookup(name string) (builtin, bool) {
	if d, ok := e.functions[name]; ok {
		return d.apply, true
	}
	return lookup(name)
}

// apply evaluates the body in a scope that sees the globals as they are now,
// with the parameters bound to the arguments.
func (d definition) apply(c call, env Env, args []Value) (Value, error) {
	if err := c.expectArgs(args, len(d.params)); err != nil {
		return nil, err
	}
	if env.depth >= maxDepth {
		return nil, evalErrorAt(c.at, "function calls nest too deeply (does %s call itself?)", c.name)
	}
	scope := env
	scope.variables = maps.Clone(env.variables)
	scope.depth++
	for i, param := range d.params {
		scope.variables[param.text] = args[i]
	}
	return d.body.evaluate(scope)
}
