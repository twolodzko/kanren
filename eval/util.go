package eval

import (
	"fmt"
	"os"

	"github.com/twolodzko/kanren/envir"
	"github.com/twolodzko/kanren/parser"
	"github.com/twolodzko/kanren/types"
)

func EvalString(code string, env *envir.Env) ([]any, *envir.Env, error) {
	var out []any
	parser := parser.NewParser(code)
	sexprs, err := parser.Read()
	if err != nil {
		return nil, env, err
	}
	for _, sexpr := range sexprs {
		result, err := Eval(sexpr, env)
		if err != nil {
			return nil, env, err
		}
		out = append(out, result)
	}
	return out, env, err
}

func LoadEval(path string, env *envir.Env) ([]any, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	sexprs, _, err := EvalString(string(content), env)
	return sexprs, err
}

func debugPrintUnify(u, v any, ok bool, s *Stream) string {
	if val, ok := u.(types.Variable); ok {
		u = s.reify(val)
	}
	if val, ok := v.(types.Variable); ok {
		v = s.reify(val)
	}
	u = types.ToString(u)
	v = types.ToString(v)
	var op string
	if ok {
		op = "≡"
	} else {
		op = "≢"
	}
	return fmt.Sprintf("%s %s %s", u, op, v)
}
