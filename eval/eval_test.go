package eval

import (
	"reflect"
	"testing"

	"github.com/twolodzko/kanren/parser"
	"github.com/twolodzko/kanren/types"
)

func TestScheme(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{"'a", "a"},
		{"(quote a)", "a"},
		{"(quote (quote a))", "'a"},
		{"'()", "()"},
		{"'(1 2 3)", "(1 2 3)"},
		{"'(+ 1 2)", "(+ 1 2)"},
		{"'()", "()"},
		{"'(1 2 3)", "(1 2 3)"},
		{"'(+ 1 2)", "(+ 1 2)"},
		{"(quasiquote a)", "a"},
		{"(quasiquote (quasiquote a))", "`a"},
		{"`(1 ,(+ 2 2) ,(+ 3 3))", "(1 4 6)"},
		{"`(1 ,(+ 2 2) . ,(+ 3 3))", "(1 4 . 6)"},
		{"`()", "()"},
		{"`(1 2 3)", "(1 2 3)"},
		{"`(+ 1 2)", "(+ 1 2)"},
		{"`,(+ 2 2)", "4"},
		{"`(2 + 2 = ,(+ 2 2))", "(2 + 2 = 4)"},
		{"`(+ 1 , (* 2 3))", "(+ 1 6)"},
		{"`(result = ,(cdr `(1 ,(car '(2 3)) 4)))", "(result = (2 4))"},
		{"`(1 2 (unquote (+ 3 4)))", "(1 2 7)"},
		{"``,,(+ 2 2)", "`,4"},
		{"`(`(,(+ 1 ,(+ 2 3)) ,,(+ 4 5)) ,(+ 6 7))", "(`(,(+ 1 5) ,9) 13)"},
		{"`(`(+ 2 ,,(+ 1 1)) ,(+ 3 3))", "(`(+ 2 ,2) 6)"},
		{"(= '(1 2 3) '(1 . (2 . (3 . ()))))", "#t"},
		{"(= '(1 2 3) `(1 ,(+ 1 1) ,(+ 1 2)))", "#t"},
		{"(= '(1 2 3 . 4) '(1 . (2 . (3 . 4))))", "#t"},
		{"(= (cdr '(1 . (2 3))) '(2 3))", "#t"},
		{"(not (= '(1 (2 (3 (4)))) '(1 2 3 4)))", "#t"},
		{"(list)", "()"},
		{"(list 1 2 3)", "(1 2 3)"},
		{"(list '() '())", "(() ())"},
		{"(car '(1))", "1"},
		{"(car '(1 2 3))", "1"},
		{"(car '((a) b c d))", "(a)"},
		{"(cdr '(1))", "()"},
		{"(cdr '(1 2))", "(2)"},
		{"(cdr '(1 2 3))", "(2 3)"},
		{"(cdr '((a) b c d))", "(b c d)"},
		{"(cdr '(a b))", "(b)"},
		{"(cdr '(a . b))", "b"},
		{"(cdr '(h))", "()"},
		{"(cdr '(h . t))", "t"},
		{"(= (cdr '(h)) '())", "#t"},
		{"(car (cdr '(1 2 3)))", "2"},
		{"(cons 1 '())", "(1)"},
		{"(cons 1 2)", "(1 . 2)"},
		{"(cons 1 '(2 3))", "(1 2 3)"},
		{"(cons '() '())", "(())"},
		{"(cons '(1 2 3) '())", "((1 2 3))"},
		{"(cons '() '(a b c))", "(() a b c)"},
		{"(cons '(a b c) '())", "((a b c))"},
		{"(cons 'hello '(world))", "(hello world)"},
		{"(null? '())", "#t"},
		{"(null? '(h . t))", "#f"},
		{"(null? #t)", "#f"},
		{"(pair? '())", "#f"},
		{"(pair? '(h . t))", "#t"},
		{"(pair? '(1 2 3))", "#t"},
		{"(pair? #t)", "#f"},
		{"(not #t)", "#f"},
		{"(not 3)", "#f"},
		{"(not '(3))", "#f"},
		{"(not #f)", "#t"},
		{"(= 'a 'a)", "#t"},
		{"(= 'a 'b)", "#f"},
		{"(= '() '())", "#t"},
		{"(= '() (cdr '(h)))", "#t"},
		{"(= '(()) '(()))", "#t"},
		{"(= '(1 2 3) '(1 2 3))", "#t"},
		{"(= '(1 2 3) '(1 2 . 3))", "#f"},
		{"(= '(1 2 . 3) '(1 2 . 3))", "#t"},
		{"(and #t #t)", "#t"},
		{"(and #t #f)", "#f"},
		{"(and (< 1 2) (< 2 3))", "#t"},
		{"(or #f #f #t #f)", "#t"},
		{"(or #t #t)", "#t"},
		{"(or #f #f #f #f)", "#f"},
		{"(or (< 10 2) (< 2 3))", "#t"},
		{"(let ((x 1)) (+ x 2))", "3"},
		{"(let ((x 5) (y 4)) (+ x y))", "9"},
		{"(let ((l '(1 2 3)) (y 5)) (/ (+ (car l) y) 2))", "3"},
		{"(let () (+ 2 2))", "4"},
		{"((lambda () 42))", "42"},
		{"((lambda (x) x) 3)", "3"},
		{"((lambda (x) (let ((y 2)) (+ x y))) 3)", "5"},
		{"(= 2 2)", "#t"},
		{"(= 2 2 2)", "#t"},
		{"(= 2 3 2)", "#f"},
		{"(< 2 3)", "#t"},
		{"(< 3 2)", "#f"},
		{"(< 1 2)", "#t"},
		{"(> 2 3)", "#f"},
		{"(> 3 2)", "#t"},
		{"(> 3 1)", "#t"},
		{"(- 7 4)", "3"},
		{"(* 2 2)", "4"},
		{"(/ 6 3)", "2"},
		{"(% 5 2)", "1"},
		{"(- 7 4)", "3"},
		{"(* 2 2)", "4"},
		{"(/ 6 3)", "2"},
		{"(% 5 2)", "1"},
		{"(define x (+ 2 (/ 10 5)))", "4"},
		{"else", "#t"},
		{"(cond ((< 5 2) 'one) ((> 7 2) 'two) (else 'three))", "two"},
		{"(cond ((< 5 2) 'one) (#f 'two) (else 'three))", "three"},
		{"(cond (#f 'one))", "()"},
		{"(cond (else 'one) (#t 'two))", "one"},
		{"(((lambda (x) (lambda (y) (+ x y))) 3 ) 4)", "7"},
		{"(let ((x 72)) ((lambda (y) (+ x y)) -12))", "60"},
		{"(let ((x 5)) (let ((y 4)) (+ x y)))", "9"},
		{"((lambda (x) (let ((y 2)) (+ x y))) 9)", "11"},
		{"((car (list + - * /)) 2 2)", "4"},
	}

	for _, tt := range testCases {
		parser := parser.NewParser(tt.input)
		sexprs, err := parser.Read()
		if err != nil {
			t.Errorf("for %v got an unexpected error: %v", tt.input, err)
			return
		}

		for _, sexpr := range sexprs {
			env := DefaultEnv()
			result, err := Eval(sexpr, env)
			if err != nil {
				t.Errorf("for %v got an unexpected error: %v", types.ToString(sexpr), err)
				return
			}
			if !reflect.DeepEqual(types.ToString(result), tt.expected) {
				t.Errorf("for %v expected %v, got %v", tt.input, tt.expected, types.ToString(result))
			}
		}
	}
}

func TestKanren(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{"(run 1 (q) (== #t #t))", "(_.0)"},
		{"(run 1 (q) (== 4 5))", "()"},
		{"(run 1 (q) succeed)", "(_.0)"},
		{"(run 1 (q) fail)", "()"},
		{"(run 1 (q) (== q q))", "(_.0)"},
		{"(run 1 (q) fail (== #t q))", "()"},
		{"(run 1 (q) succeed (== #t q))", "(#t)"},
		{"(run 1 (q) (== q 'ok))", "(ok)"},
		{"(run 1 (q) (fresh (x) (== q x)))", "(_.0)"},
		{"(run 1 (q) (fresh (x) (== 'ok x) (== x q)))", "(ok)"},
		{"(run 1 (q) (fresh (x y) (== q (list x y))))", "((_.0 _.1))"},
		{"(run 1 (q) (fresh (x y) (== x y) (== y 'ok) (== x q)))", "(ok)"},
		{"(run 1 (q) (let ((x #f)) (== #t x)))", "()"},
		{"(run 1 (q) (let ((x #f)) (== #f x)))", "(_.0)"},
		{"(run 1 (q) (let ((x 'ok)) (== q x)))", "(ok)"},
		{"(run 1 (q) ((lambda (x) x) (== q #f)))", "(#f)"},
		{"(run 1 (q) (fresh (x) (== x #t) (== #t q)))", "(#t)"},
		{"(run 1 (q) (fresh (x y) (== q (cons x (cons y '())))))", "((_.0 _.1))"},
		{"(run 1 (q) (conde (fail succeed) (succeed fail)))", "()"},
		{"(run 1 (q) (conde (fail succeed) (succeed succeed)))", "(_.0)"},
		{"(run 1 (q) (conde (fail succeed) (else succeed)))", "(_.0)"},
		{"(run 2 (q) (conde (succeed (== q 1)) (succeed (== q 2)) (fail (== q 'wrong)) (succeed (== q 3)) ))", "(1 2)"},
		{"(run #f (q) (conde (succeed (== q 1)) (succeed (== q 2)) (fail (== q 'wrong)) (succeed (== q 3)) ))", "(1 2 3)"},
		{"(run* (q) (conde (succeed (== q 1)) (succeed (== q 2)) (fail (== q 'wrong)) (succeed (== q 3)) ))", "(1 2 3)"},
		{"(run* (q) (conde ((== q 'ok))))", "(ok)"},
		{"(run 1 (q) (== `(1 . ,q) '(1)))", "(())"},
		{"(run 1 (q) (== '(1 2 3) `(1 2 3 . ,q)))", "(())"},
		{"(run 1 (q) (fresh (a b) (== `(,a . ,b) '(1 . 2)) (== `(,a . ,b) q)))", "((1 . 2))"},
		{
			`
			(run 5 (q)
				(fresh (x y z)
			    	(conde
						((== 'a x) (== 1 y) (== 'd z))
						((== 2 y) (== 'b x) (== 'e z))
						((== 'f z) (== 'c x) (== 3 y)))
					(== (list x y z) q)))
			`,
			"((a 1 d) (b 2 e) (c 3 f))",
		},
		{"(run* (q) (== q #f))", "(#f)"},
		{"(run* (q) (let ((a (== #t q)) (b (== #f q))) b))", "(#f)"},
		{
			`
			(run 2 (q)
				(fresh (x y z)
					(conde
						((== (list x y z x) q))
						((== (list z y x z) q)))))
			`,
			"((_.0 _.1 _.2 _.0) (_.0 _.1 _.2 _.0))",
		},
		{
			"(run* (q) (conde ((== q #f)) ((== q #t)) ) (== q #t))",
			"(#t)",
		},
		{
			"(run* (q) (let ((a (== #t q )) (b (== #f q))) b))",
			"(#f)",
		},
		{
			`
			(run* (q)
				(fresh (x y)
					(== (list 'a x 'c)
						 (list 'a 'b y))
					(== q (cons x y))))
			`,
			"((b . c))",
		},
		{
			`
			(run* (q)
				(fresh (x)
					(== 5 x)
					(project (x)
						(== (* x x) q))))
			`,
			"(25)",
		},
	}

	for _, tt := range testCases {
		parser := parser.NewParser(tt.input)
		sexprs, err := parser.Read()
		if err != nil {
			t.Errorf("for %v got an unexpected error: %v", tt.input, err)
			return
		}

		env := DefaultEnv()
		result, err := Eval(sexprs[0], env)
		if err != nil {
			t.Errorf("for %v got an unexpected error: %v", types.ToString(sexprs[0]), err)
			return
		}
		result = types.ToString(result)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("for %v expected %s, got %s", tt.input, tt.expected, result)
		}
	}
}

func TestWalk(t *testing.T) {
	memory := NewStream()
	x := types.NewVariable("x")
	memory.extend(x, "ok")
	val := memory.walk(x)
	if val != "ok" {
		t.Error("walk failed")
	}
}

func TestDeepWalk(t *testing.T) {
	x, y, z := types.NewVariable("x"), types.NewVariable("y"), types.NewVariable("z")
	input := types.List(
		5,
		x,
		types.List(true, y, x),
		z,
	)
	expected := types.List(
		5,
		0,
		types.List(true, 1, 0),
		2,
	)
	memory := Stream{[]KeyVal{{x, 0}, {y, 1}, {z, 2}}}
	result := memory.deepWalk(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected: %v, got %v", expected, result)
	}
}

func TestDeepWalkRecursive(t *testing.T) {
	// (test-check "9.46"
	//    (walk* y `((,y . (,w ,z c)) (,v . b) (,x . ,v) (,z . ,x)))
	//    `(,w b c))
	x, y, z, v, w := types.NewVariable("x"), types.NewVariable("y"), types.NewVariable("z"), types.NewVariable("v"), types.NewVariable("w")
	memory := NewStream()
	memory.extend(y, types.List(w, z, types.Symbol("c")))
	memory.extend(v, types.Symbol("b"))
	memory.extend(x, v)
	memory.extend(z, x)
	expected := types.List(w, types.Symbol("b"), types.Symbol("c"))
	result := memory.deepWalk(y)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected: %v, got %v", expected, result)
	}
}

func TestBirthRecords(t *testing.T) {
	x, y := types.NewVariable("x"), types.NewVariable("y")
	memory := NewStream()
	memory.extend(x, "wrong")
	memory.birthRecord(x)
	memory.extend(y, "also wrong")
	result := memory.deepWalk(x)
	if !reflect.DeepEqual(result, x) {
		t.Errorf("expected: %v, got %v", x, result)
	}
}

func TestUnifyConstant(t *testing.T) {
	memory := NewStream()
	ok := memory.unify(2, 2)
	if !ok {
		t.Error("unification failed")
	}
}

func TestUnify(t *testing.T) {
	x := types.NewVariable("x")
	memory := NewStream()
	ok := memory.unify(x, 42)
	if !ok {
		t.Error("unification failed")
	}
	expected := NewStream()
	expected.extend(x, 42)
	if !reflect.DeepEqual(memory, expected) {
		t.Errorf("expected memory: %v, got %v", expected, memory)
	}
}

func TestReify(t *testing.T) {
	x, y, z := types.NewVariable("x"), types.NewVariable("y"), types.NewVariable("z")
	input := types.List(
		5,
		x,
		types.List(true, y, x),
		z,
	)
	expected := types.List(
		5,
		types.Free(0),
		types.List(true, types.Free(1), types.Free(0)),
		types.Free(2),
	)
	memory := NewStream()
	result := memory.reify(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestReifyS(t *testing.T) {
	x := types.NewVariable("x")
	memory := NewStream()
	result := memory.reifyStream(x)
	if !result || memory.len() == 0 {
		t.Errorf("reification failed with memory: %v", memory)
	}
}
