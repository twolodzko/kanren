package eval

import (
	"fmt"
	"strings"

	"github.com/twolodzko/kanren/envir"
	"github.com/twolodzko/kanren/types"
)

type Goal interface {
	Query(*Stream) (bool, error)
	Next() bool
	Reset()
}

func run(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	if p.This == false {
		// (run #f (x) ... ) -> (run* (x) ...)
		return runAll(p.Next, env)
	} else if p.Next == nil {
		return nil, ArityError
	}
	local := envir.NewEnvFrom(env)

	reps, ok := p.This.(int)
	if !ok {
		return nil, NaN{p.This}
	}

	p, ok = p.Next.(types.Pair)
	if !ok {
		return nil, WrongArg{p.Next}
	}
	b, ok := p.This.(types.Pair)
	if !ok {
		return nil, WrongArg{p.Next}
	}
	name, ok := b.This.(types.Symbol)
	if !ok {
		return nil, InvalidName{b.This}
	}
	target := types.NewVariable(string(name))
	local.Set(name, target)

	body, ok := p.Next.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	goals, err := extractGoals(body, local)
	if err != nil {
		return nil, err
	}

	var acc []any
	i := 0
	for i < reps {
		s := NewStream()
		s.birthRecord(target)
		ok, err := queryAll(goals, s)
		if err != nil {
			return nil, err
		}
		if ok {
			r := s.reify(target)
			if Debug {
				fmt.Printf("  result: %v\n", types.ToString(r))
			}
			acc = append(acc, r)
			i += 1
		}
		if !next(goals) {
			if Debug {
				fmt.Println("       ∎  final goal")
			}
			break
		}
	}
	return types.List(acc...), nil
}

func runAll(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	local := envir.NewEnvFrom(env)

	binding, ok := p.This.(types.Pair)
	if !ok || p.Next == nil {
		return nil, WrongArg{p.This}
	}
	name, ok := binding.This.(types.Symbol)
	if !ok {
		return nil, InvalidName{binding.This}
	}
	target := types.NewVariable(string(name))
	local.Set(name, target)

	body, ok := p.Next.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	goals, err := extractGoals(body, local)
	if err != nil {
		return nil, err
	}

	var acc []any
	for {
		s := NewStream()
		s.birthRecord(target)
		ok, err := queryAll(goals, s)
		if err != nil {
			return nil, err
		}
		if ok {
			r := s.reify(target)
			if Debug {
				fmt.Printf("  result: %v\n", types.ToString(r))
			}
			acc = append(acc, r)
		}
		if !next(goals) {
			if Debug {
				fmt.Println("       ∎  final goal")
			}
			break
		}
	}
	return types.List(acc...), nil
}

type ConstGoal struct {
	name  string
	value bool
}

func (g ConstGoal) Query(_ *Stream) (bool, error) {
	return g.value, nil
}

func (g ConstGoal) Next() bool {
	return false
}

func (g ConstGoal) Reset() {}

func (g ConstGoal) String() string {
	return g.name
}

type Unify struct {
	u, v any
	env  *envir.Env
}

func (g Unify) Query(s *Stream) (bool, error) {
	u, err := Eval(g.u, g.env)
	if err != nil {
		return false, err
	}
	v, err := Eval(g.v, g.env)
	if err != nil {
		return false, err
	}
	ok := s.unify(u, v)
	return ok, nil
}

func (g Unify) Next() bool {
	return false
}

func (g Unify) Reset() {}

func (g Unify) String() string {
	return fmt.Sprintf("(== %v %v)", types.ToString(g.u), types.ToString(g.v))
}

func newUnify(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	first := p.This
	p, ok = p.Next.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	if p.Next != nil {
		return nil, ArityError
	}
	return Unify{first, p.This, env}, nil
}

type Fresh struct {
	vars  []types.Variable
	goals []Goal
}

func (g *Fresh) Query(s *Stream) (bool, error) {
	for _, v := range g.vars {
		s.birthRecord(v)
	}
	return queryAll(g.goals, s)
}

func (g *Fresh) Next() bool {
	return next(g.goals)
}

func (g *Fresh) Reset() {
	for _, g := range g.goals {
		g.Reset()
	}
}

func (g Fresh) String() string {
	var vars, goals []string
	for _, v := range g.vars {
		vars = append(vars, types.ToString(v))
	}
	for _, g := range g.goals {
		goals = append(goals, fmt.Sprintf("%v", g))
	}
	return fmt.Sprintf("(fresh (%s) %s)", strings.Join(vars, " "), strings.Join(goals, " "))
}

func newFresh(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	local := envir.NewEnvFrom(env)
	pair, ok := p.This.(types.Pair)
	if !ok {
		return nil, NonList{p.This}
	}
	names, err := extractSymbols(pair)
	if err != nil {
		return nil, err
	}
	var vars []types.Variable
	for _, name := range names {
		v := types.NewVariable(string(name))
		local.Set(name, v)
		vars = append(vars, v)
	}
	body, ok := p.Next.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	goals, err := extractGoals(body, local)
	if err != nil {
		return nil, err
	}
	return &Fresh{vars, goals}, nil
}

type Conde struct {
	branches []any
	current  int
	branch   []Goal
	env      *envir.Env
}

func (g *Conde) Query(s *Stream) (bool, error) {
	start := s.len()
	for g.current < len(g.branches) {
		if err := g.EnsureBranch(); err != nil {
			return false, err
		}
		ok, err := queryAll(g.branch, s)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
		if !g.Next() {
			break
		}
		// backtrack
		s.keep(start)
	}
	if Debug {
		fmt.Println("       ∎  final goal")
	}
	return false, nil
}

func (g *Conde) Next() bool {
	if next(g.branch) {
		return true
	}
	g.current += 1
	g.branch = nil
	return g.current < len(g.branches)
}

func (g *Conde) Reset() {
	g.current = 0
	g.branch = nil
}

func (g Conde) String() string {
	var branches []string
	for _, b := range g.branches {
		branches = append(branches, fmt.Sprintf("%v", b))
	}
	return fmt.Sprintf("(conde %s)", strings.Join(branches, " "))
}

func (g *Conde) EnsureBranch() error {
	if g.branch == nil {
		return g.InitBranch()
	}
	return nil
}

func (g *Conde) InitBranch() error {
	p, ok := g.branches[g.current].(types.Pair)
	if !ok {
		return NonList{g.branches[g.current]}
	}
	if p.This == types.Symbol("else") {
		// no-op: this is a syntactic sugar
		p, ok = p.Next.(types.Pair)
		if !ok {
			return SyntaxError
		}
	}
	var err error
	g.branch, err = extractGoals(p, g.env)
	return err
}

func newConde(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	var (
		branches []any
		head     any = p
	)
	for head != nil {
		p, ok := head.(types.Pair)
		if !ok {
			return nil, SyntaxError
		}
		branches = append(branches, p.This)
		head = p.Next
	}
	return &Conde{branches, 0, nil, env}, nil
}

type Project struct {
	vars  []types.Symbol
	goals []Goal
	env   *envir.Env
}

func (g *Project) Query(s *Stream) (bool, error) {
	for _, name := range g.vars {
		val, ok := g.env.Get(name)
		if !ok {
			return false, fmt.Errorf("unbound variable %v", name)
		}
		val = s.deepWalk(val)
		g.env.Set(name, val)
	}
	return queryAll(g.goals, s)
}

func (g *Project) Next() bool {
	return next(g.goals)
}

func (g *Project) Reset() {
	for _, g := range g.goals {
		g.Reset()
	}
}

func (g Project) String() string {
	var vars, goals []string
	for _, v := range g.vars {
		vars = append(vars, string(v))
	}
	for _, g := range g.goals {
		goals = append(goals, fmt.Sprintf("%v", g))
	}
	return fmt.Sprintf("(project (%s) %s)", strings.Join(vars, " "), strings.Join(goals, " "))
}

func newProject(args any, env *envir.Env) (any, error) {
	p, ok := args.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	local := envir.NewEnvFrom(env)
	pair, ok := p.This.(types.Pair)
	if !ok {
		return nil, NonList{p.This}
	}
	vars, err := extractSymbols(pair)
	if err != nil {
		return nil, err
	}
	body, ok := p.Next.(types.Pair)
	if !ok {
		return nil, SyntaxError
	}
	goals, err := extractGoals(body, local)
	if err != nil {
		return nil, err
	}
	return &Project{vars, goals, local}, nil
}

func queryAll(goals []Goal, s *Stream) (bool, error) {
	for _, g := range goals {
		if Debug {
			fmt.Printf(" ↪ query: %v\n", g)
			fmt.Printf("   subst: %v\n", s)
		}
		ok, err := g.Query(s)
		if Debug {
			if ok {
				fmt.Println("       ✔  success")
			} else {
				fmt.Println("       ✘  failure")
			}
		}
		if !ok || err != nil {
			return false, err
		}
	}
	return true, nil
}

func next(goals []Goal) bool {
	// advance the deepest goals first
	for i := len(goals) - 1; i >= 0; i-- {
		g := goals[i]
		if g.Next() {
			return true
		} else {
			g.Reset()
		}
	}
	return false
}

func extractGoals(pair types.Pair, env *envir.Env) ([]Goal, error) {
	var (
		goals []Goal
		head  any = pair
	)
	for head != nil {
		p, ok := head.(types.Pair)
		if !ok {
			return nil, SyntaxError
		}
		v, err := Eval(p.This, env)
		if err != nil {
			return nil, err
		}
		g, ok := v.(Goal)
		if !ok {
			return nil, WrongArg{v}
		}
		goals = append(goals, g)
		head = p.Next
	}
	return goals, nil
}
