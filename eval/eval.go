package eval

import (
	"errors"
	"fmt"

	"github.com/twolodzko/kanren/envir"
	"github.com/twolodzko/kanren/types"
)

var Debug = false

type (
	tco  = func(any, *envir.Env) (any, *envir.Env, error)
	proc = func(any, *envir.Env) (any, error)
)

func Eval(sexpr any, env *envir.Env) (any, error) {
	for {
		if Debug {
			fmt.Printf(" ↪ eval:  %v\n", types.ToString(sexpr))
			fmt.Printf("   env:   %v\n", env)
		}

		switch val := sexpr.(type) {
		case types.Symbol:
			return getSymbol(val, env)
		case types.Pair:
			name := val.This
			args := val.Next

			callable, err := Eval(name, env)
			if err != nil {
				return nil, err
			}

			switch fn := callable.(type) {
			case tco:
				sexpr, env, err = fn(args, env)
				if err != nil {
					return nil, err
				}
			case proc:
				return fn(args, env)
			default:
				return nil, fmt.Errorf("%v is not callable", types.ToString(fn))
			}
		case types.Free, types.Variable:
			return nil, errors.New("kanren variable was used outside of its context")
		default:
			return sexpr, nil
		}
	}
}

func getSymbol(sexpr any, env *envir.Env) (any, error) {
	switch val := sexpr.(type) {
	case types.Symbol:
		if val, ok := env.Get(val); ok {
			return val, nil
		}
		return nil, fmt.Errorf("unbound variable %v", val)
	default:
		return val, nil
	}
}

// Evaluate all args but last, return the last arg and the enclosing environment
func partialEval(args any, env *envir.Env) (any, *envir.Env, error) {
	head := args
	for {
		switch p := head.(type) {
		case types.Pair:
			if p.Next == nil {
				return p.This, env, nil
			}
			if _, err := Eval(p.This, env); err != nil {
				return nil, nil, err
			}
			head = p.Next
		default:
			return head, env, nil
		}
	}
}

// Evaluate two expressions
func evalTwo(args types.Pair, env *envir.Env) (any, any, error) {
	a, err := Eval(args.This, env)
	if err != nil {
		return nil, nil, err
	}
	p, ok := args.Next.(types.Pair)
	if !ok || p.Next != nil {
		return nil, nil, ArityError
	}
	b, err := Eval(p.This, env)
	return a, b, err
}
