package eval

import (
	"github.com/twolodzko/kanren/envir"
	"github.com/twolodzko/kanren/types"
)

// Create `lambda` function
//
//	(lambda (args ...) body ...)
func newLambda(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	var (
		body = p.Next
		vars []types.Symbol
		err  error
	)
	switch a := p.This.(type) {
	case types.Pair:
		vars, err = extractSymbols(a)
		if err != nil {
			return nil, err
		}
	case nil:
		vars = nil
	default:
		return nil, NonList{p.This}
	}
	return func(args any, runEnv *envir.Env) (any, *envir.Env, error) {
		local, err := createClosure(args, vars, env, runEnv)
		if err != nil {
			return nil, local, err
		}
		// the body of the function is evaluated in the local env of the lambda
		return partialEval(body, local)
	}, nil
}

// Transform pair to slice
func extractSymbols(args types.Pair) ([]types.Symbol, error) {
	var (
		vars []types.Symbol
		head any = args
	)
	for head != nil {
		p, ok := head.(types.Pair)
		if !ok {
			return nil, NonList{head}
		}
		if name, ok := p.This.(types.Symbol); ok {
			vars = append(vars, name)
		} else {
			return vars, InvalidName{args.This}
		}
		head = p.Next
	}
	return vars, nil
}

func createClosure(args any, vars []types.Symbol, parentEnv, env *envir.Env) (*envir.Env, error) {
	// local env inherits from the env where the lambda was defined
	local := envir.NewEnvFrom(parentEnv)
	switch p := args.(type) {
	case types.Pair:
		i := 0
		err := p.ForEach(func(val any) error {
			if i >= len(vars) {
				return ArityError
			}
			// arguments are evaluated in the env enclosing the lambda call
			val, err := Eval(val, env)
			local.Set(vars[i], val)
			i++
			return err
		})
		if err != nil {
			return nil, err
		}
		if i != len(vars) {
			return nil, ArityError
		}
		return local, err
	default:
		return local, nil
	}
}
