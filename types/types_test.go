package types

import (
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	var testCases = []struct {
		input    Pair
		expected string
	}{
		{
			Cons(1, 2),
			"(1 . 2)",
		},
		{
			Cons(1, Cons(2, 3)),
			"(1 2 . 3)",
		},
		{
			Cons(1, 2, 3, 4),
			"(1 2 3 . 4)",
		},
		{
			Cons(1, 2, 3, nil),
			"(1 2 3)",
		},
		{
			List(1, 2, 3).(Pair),
			"(1 2 3)",
		},
		{
			List(1, 2, 3, nil).(Pair),
			"(1 2 3 ())",
		},
		{
			List(true).(Pair),
			"(#t)",
		},
	}
	for _, tt := range testCases {
		result := tt.input.String()
		if result != tt.expected {
			t.Errorf("expected '%s', got '%s'", tt.expected, result)
		}
	}
}

func TestRepack(t *testing.T) {
	input := List(1, 2, 3).(Pair)

	result := Cons(input.Map(func(val any) any { return val })...)
	if !reflect.DeepEqual(input, result) {
		t.Errorf("cons after map failed: %v", result)
	}

	cons := Cons(1, 2, 3, nil)
	if !reflect.DeepEqual(cons, result) {
		t.Errorf("cons does not return same result as list: %v", cons)
	}
}

func TestVariables(t *testing.T) {
	a := NewVariable("x")
	b := NewVariable("x")
	if a == b {
		t.Error("the variables should differ ragerdless of same names")
	}
}
