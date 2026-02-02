package eval

import (
	"github.com/twolodzko/kanren/envir"
	"github.com/twolodzko/kanren/types"
)

func DefaultEnv() *envir.Env {
	env := envir.NewEnv()
	// scheme
	env.Set("quote", quote)
	env.Set("unquote", unquote)
	env.Set("quasiquote", quasiQuote)
	env.Set("lambda", newLambda)
	env.Set("let", let)
	env.Set("define", define)
	env.Set("car", car)
	env.Set("cdr", cdr)
	env.Set("cons", cons)
	env.Set("else", true)
	env.Set("load", load)
	env.Set("null?", isNull)
	env.Set("pair?", isPair)
	env.Set("and", and)
	env.Set("or", or)
	env.Set("not", not)
	env.Set("cond", cond)
	env.Set("list", list)
	env.Set("=", func(args any, env *envir.Env) (any, error) {
		return cmp(args, env, func(a, b any) (bool, error) {
			return a == b, nil
		})
	})
	env.Set(">", func(args any, env *envir.Env) (any, error) {
		return cmp(args, env, func(a, b any) (bool, error) {
			ai, ok := a.(int)
			if !ok {
				return false, NaN{a}
			}
			bi, ok := b.(int)
			if !ok {
				return false, NaN{b}
			}
			return ai > bi, nil
		})
	})
	env.Set("<", func(args any, env *envir.Env) (any, error) {
		return cmp(args, env, func(a, b any) (bool, error) {
			ai, ok := a.(int)
			if !ok {
				return false, NaN{a}
			}
			bi, ok := b.(int)
			if !ok {
				return false, NaN{b}
			}
			return ai < bi, nil
		})
	})
	env.Set("+", op(func(a, b int) int { return a + b }))
	env.Set("-", func(args any, env *envir.Env) (any, error) {
		p, ok := args.(types.Pair)
		if !ok {
			return nil, SyntaxError
		}
		if p.Next == nil {
			val, err := Eval(p.This, env)
			if err != nil {
				return nil, err
			}
			num, ok := val.(int)
			if !ok {
				return nil, NaN{val}
			}
			return -num, nil
		}
		return foldLeft(p, env, func(a, b int) (int, error) {
			return a - b, nil
		})
	})
	env.Set("*", op(func(a, b int) int { return a * b }))
	env.Set("/", op(func(a, b int) int { return a / b }))
	env.Set("%", op(func(a, b int) int { return a % b }))
	// extras
	env.Set("test-check", testCheck)
	// kanren
	env.Set("run", run)
	env.Set("run*", runAll)
	env.Set("succeed", ConstGoal{"succeed", true})
	env.Set("fail", ConstGoal{"fail", false})
	env.Set("==", newUnify)
	env.Set("fresh", newFresh)
	env.Set("conde", newConde)
	env.Set("project", newProject)
	return envir.NewEnvFrom(env)
}

func op(fn func(a, b int) int) func(any, *envir.Env) (any, error) {
	return func(args any, env *envir.Env) (any, error) {
		return foldLeft(args, env, func(a, b int) (int, error) {
			return fn(a, b), nil
		})
	}
}
