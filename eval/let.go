package eval

import (
	"github.com/twolodzko/kanren/envir"
	"github.com/twolodzko/kanren/types"
)

// `let` procedure
//
//	(let ((key1 value1) (key2 value2) ...) expr1 expr2 ...)
func let(args any, env *envir.Env) (any, *envir.Env, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, nil, SyntaxError
	}
	local := envir.NewEnvFrom(env)
	switch b := p.This.(type) {
	case types.Pair:
		if err := setBindings(b, local, env); err != nil {
			return nil, nil, err
		}
	case nil:
	default:
		return nil, nil, NonList{p.This}
	}
	return partialEval(p.Next, local)
}

// Iterate through the bindings ((key1 value1) (key2 value2) ...) and set them to an environment
func setBindings(bindings types.Pair, local, parent *envir.Env) error {
	return bindings.ForEach(func(val any) error {
		p, ok := val.(types.Pair)
		if !ok {
			return &NonList{val}
		}
		return bind(p, local, parent)
	})
}

// Bind value to the name in the local env
func bind(binding types.Pair, local, parent *envir.Env) error {
	name, sexpr, err := extractBinding(binding)
	if err != nil {
		return err
	}
	// arguments are evaluated in env enclosing let
	val, err := Eval(sexpr, parent)
	if err != nil {
		return err
	}
	local.Set(name, val)
	return nil
}

// Extract name and value for the binding
func extractBinding(arg types.Pair) (types.Symbol, any, error) {
	switch name := arg.This.(type) {
	case types.Symbol:
		p, ok := arg.Next.(types.Pair)
		if !ok {
			return "", nil, SyntaxError
		}
		return name, p.This, nil
	default:
		return "", nil, InvalidName{arg.This}
	}
}
