package parser

import (
	"reflect"
	"testing"

	"github.com/twolodzko/kanren/types"
)

func TestParse(t *testing.T) {
	var testCases = []struct {
		input    string
		expected any
	}{
		{"a", types.Symbol("a")},
		{"42", 42},
		{"-100", -100},
		{"#t", true},
		{"#f", false},
		{"()", nil},
		{"(a)", types.List(types.Symbol("a"))},
		{"(())", types.List(nil)},
		{"(1 2 3)", types.List(1, 2, 3)},
		{"((1 2) 3)", types.List(types.List(1, 2), 3)},
		{"(1 (2 3))", types.List(1, types.List(2, 3))},
		{"'a", types.List(types.Symbol("quote"), types.Symbol("a"))},
		{"'(a)", types.List(types.Symbol("quote"), types.List(types.Symbol("a")))},
		{"('a)", types.List(quote(types.Symbol("a")))},
		{"'''a", types.List(types.Symbol("quote"), types.List(types.Symbol("quote"), types.List(types.Symbol("quote"), types.Symbol("a"))))},
		{"'()", types.List(types.Symbol("quote"), nil)},
		{"''()", types.List(types.Symbol("quote"), types.List(types.Symbol("quote"), nil))},
		{"  \n\ta", types.Symbol("a")},
		{"\n  \t\n(\n   a\t\n)  ", types.List(types.Symbol("a"))},
		{"(list 1 2 ;; a comment\n3)", types.List(types.Symbol("list"), 1, 2, 3)},
	}

	for _, tt := range testCases {
		parser := NewParser(tt.input)
		result, err := parser.Read()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(result[0], tt.expected) {
			t.Errorf("for %q expected %v, got: %v", tt.input, tt.expected, result[0])
		}
	}
}

func TestParseAndPrint(t *testing.T) {
	var testCases = []string{
		"(1 2 3)",
		"(1 (2 3))",
		"((1 2) 3)",
		"((1) (((2)) 3))",
	}

	for _, input := range testCases {
		parser := NewParser(input)
		result, err := parser.Read()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if types.ToString(result[0]) != input {
			t.Errorf("%q is printable as %v", input, result[0])
		}
	}
}

func TestReadAtomValue(t *testing.T) {
	var testCases = []struct {
		input    string
		expected types.Symbol
	}{
		{"a", "a"},
		{"a   ", "a"},
		{"a)   ", "a"},
		{"a(b)", "a"},
	}

	for _, tt := range testCases {
		parser := NewParser(tt.input)
		result, err := parser.readAtom()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != tt.expected {
			t.Errorf("for %q expected %v, got: %v", tt.input, tt.expected, result)
		}
	}
}

func TestParseExpectError(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{"(", "list was not closed with closing bracket"},
		{"(a", "list was not closed with closing bracket"},
		{"(lorem ipsum", "list was not closed with closing bracket"},
		{"lorem ipsum)", "unexpected closing bracket"},
	}
	for _, tt := range testCases {
		parser := NewParser(tt.input)
		if _, err := parser.Read(); err.Error() != tt.expected {
			t.Errorf("for %q expected an error %q, got: %q", tt.input, tt.expected, err)
		}
	}
}
