package eval

import (
	"github.com/twolodzko/kanren/envir"
	"github.com/twolodzko/kanren/types"
)

func quote(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	if p.Next != nil {
		return nil, ArityError
	}
	return p.This, nil
}

// `unquote` procedure
func unquote(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	if p.Next != nil {
		return nil, ArityError
	}
	return Eval(p.This, env)
}

// `quasiQuote` procedure
func quasiQuote(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	if p.Next != nil {
		return nil, ArityError
	}
	return unquoteRecursively(p.This, 1, env)
}

func unquoteRecursively(val any, numQuotes int, env *envir.Env) (any, error) {
	p, ok := val.(types.Pair)
	if !ok {
		return val, nil
	}
	if sym, ok := p.This.(types.Symbol); ok {
		switch sym {
		case "quasiquote":
			numQuotes++
		case "unquote":
			numQuotes--
			if numQuotes == 0 {
				return unquote(p.Next, env)
			}
		}
	}
	head, err := unquoteRecursively(p.This, numQuotes, env)
	if err != nil {
		return nil, err
	}
	tail, err := unquoteRecursively(p.Next, numQuotes, env)
	if err != nil {
		return nil, err
	}
	return types.Cons(head, tail), err
}
