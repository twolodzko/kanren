package envir

import (
	"fmt"
	"strings"

	"github.com/twolodzko/kanren/types"
)

type Env struct {
	Vars   map[types.Symbol]any
	Parent *Env
}

func NewEnv() *Env {
	vars := make(map[types.Symbol]any)
	return &Env{vars, nil}
}

func NewEnvFrom(parent *Env) *Env {
	new := NewEnv()
	new.Parent = parent
	return new
}

func (e *Env) Set(name types.Symbol, value any) {
	e.Vars[name] = value
}

// Find an enclosing environment for the variable
func (e Env) FindEnv(name types.Symbol) (*Env, bool) {
	current := &e
	for current != nil {
		if _, ok := current.Vars[name]; ok {
			return current, true
		}
		current = current.Parent
	}
	return nil, false
}

func (e Env) Get(name types.Symbol) (any, bool) {
	if env, ok := e.FindEnv(name); ok {
		return env.Vars[name], true
	}
	return nil, false
}

func (e Env) String() string {
	var acc []string
	for k, v := range e.Vars {
		acc = append(acc, fmt.Sprintf("%s:%s", k, types.ToString(v)))
	}
	return fmt.Sprintf("{[%s] %p}", strings.Join(acc, " "), e.Parent)
}
