package eval

import (
	"fmt"
	"reflect"

	"github.com/twolodzko/kanren/envir"
	"github.com/twolodzko/kanren/types"
)

func define(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	key, ok := p.This.(types.Symbol)
	if !ok {
		return nil, InvalidName{p.This}
	}
	lhs, ok := p.Next.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	val, err := Eval(lhs.This, env)
	if err != nil {
		return nil, err
	}
	env.Set(key, val)
	return val, nil
}

func car(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	v, err := Eval(p.This, env)
	if err != nil {
		return nil, err
	}
	switch p := v.(type) {
	case types.Pair:
		return p.This, nil
	default:
		return nil, NonList{v}
	}
}

func cdr(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	v, err := Eval(p.This, env)
	if err != nil {
		return nil, err
	}
	switch p := v.(type) {
	case types.Pair:
		return p.Next, nil
	default:
		return nil, NonList{v}
	}
}

func cons(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	a, b, err := evalTwo(p, env)
	if err != nil {
		return nil, err
	}
	return types.Cons(a, b), nil
}

func list(args any, env *envir.Env) (any, error) {
	if args == nil {
		return nil, nil
	}
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	s, err := p.TryMap(func(val any) (any, error) {
		return Eval(val, env)
	})
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return nil, nil
	}
	return types.Cons(s...), nil
}

func isNull(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	v, err := Eval(p.This, env)
	if err != nil {
		return nil, err
	}
	return v == nil, nil
}

func isPair(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	v, err := Eval(p.This, env)
	if err != nil {
		return nil, err
	}
	_, ok = v.(types.Pair)
	return ok, nil
}

func and(args any, env *envir.Env) (any, error) {
	var (
		last any = true
		err  error
	)
	head := args
	for head != nil {
		p, ok := head.(types.Pair)
		if !ok {
			return nil, SyntaxError
		}
		last, err = Eval(p.This, env)
		if err != nil {
			return nil, err
		}
		if !types.IsTrue(last) {
			return false, nil
		}
		head = p.Next
	}
	return last, nil
}

func or(args any, env *envir.Env) (any, error) {
	head := args
	for head != nil {
		p, ok := head.(types.Pair)
		if !ok {
			return nil, SyntaxError
		}
		this, err := Eval(p.This, env)
		if err != nil {
			return nil, err
		}
		if types.IsTrue(this) {
			return this, nil
		}
		head = p.Next
	}
	return false, nil
}

func not(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	val, err := Eval(p.This, env)
	if err != nil {
		return nil, err
	}
	return !types.IsTrue(val), nil
}

// `cond` procedure
//
//	(cond (test1 expr1) (test2 expr2)...)
func cond(args any, env *envir.Env) (any, *envir.Env, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, nil, SyntaxError
	}
	var head any = p
	for head != nil {
		p, ok := head.(types.Pair)
		if !ok {
			return nil, nil, SyntaxError
		}
		b, ok := p.This.(types.Pair)
		if !ok {
			return nil, nil, NonList{p.This}
		}
		test, err := Eval(b.This, env)
		if err != nil {
			return nil, env, err
		}
		if types.IsTrue(test) {
			v, ok := b.Next.(types.Pair)
			if !ok {
				return nil, nil, SyntaxError
			}
			return v.This, env, nil
		}
		head = p.Next
	}
	return nil, env, nil
}

func testCheck(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	tag, ok := p.This.(string)
	if !ok {
		return nil, SyntaxError
	}
	p, ok = p.Next.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	a, b, err := evalTwo(p, env)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(a, b) {
		return nil, fmt.Errorf("test %s failed:\n        %v\n is not %v", tag, types.ToString(a), types.ToString(b))
	}
	return nil, nil
}

func load(args any, env *envir.Env) (any, error) {
	var head any = args
	for head != nil {
		p, ok := head.(types.Pair)
		if !ok {
			return nil, SyntaxError
		}
		val, err := Eval(p.This, env)
		if err != nil {
			return nil, err
		}
		path, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("%v is not a valid filename", val)
		}
		if _, err := LoadEval(path, env); err != nil {
			return nil, err
		}
		head = p.Next
	}
	return nil, nil
}

func cmp(args any, env *envir.Env, cmp func(a, b any) (bool, error)) (bool, error) {
	var (
		prev any
		this any
		err  error
	)
	p, ok := args.(types.Pair)
	if !ok {
		return false, SyntaxError
	}
	prev, err = Eval(p.This, env)
	if err != nil {
		return false, err
	}
	head := p.Next
	for head != nil {
		p, ok := head.(types.Pair)
		if !ok {
			return false, SyntaxError
		}
		this, err = Eval(p.This, env)
		if err != nil {
			return false, err
		}
		ok, err = cmp(prev, this)
		if !ok || err != nil {
			return ok, nil
		}
		head = p.Next
	}
	return true, nil
}

func foldLeft(args any, env *envir.Env, fn func(acc, val int) (int, error)) (int, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return 0, SyntaxError
	}
	val, err := Eval(p.This, env)
	if err != nil {
		return 0, err
	}
	acc, ok := val.(int)
	if !ok {
		return 0, NaN{val}
	}
	head := p.Next
	for head != nil {
		p, ok := head.(types.Pair)
		if !ok {
			return 0, SyntaxError
		}
		val, err := Eval(p.This, env)
		if err != nil {
			return 0, err
		}
		this, ok := val.(int)
		if !ok {
			return 0, NaN{val}
		}
		acc, err = fn(acc, this)
		if err != nil {
			return 0, nil
		}
		head = p.Next
	}
	return acc, nil
}
