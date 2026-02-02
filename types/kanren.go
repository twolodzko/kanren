package types

import (
	"fmt"
)

var Pretty = false

// kanren variable used in unification processes
type (
	// Variables are named, but are always passed as references,
	// so that when comparing them we compare their addrses, what makes
	// them unique (variables with same names are not necessary the same).
	// In the original implementation they are vectors, and Scheme's eq?
	// for vectors compares their addresses.
	variableName string
	Variable     = *variableName
	Free         int
)

func NewVariable(name string) Variable {
	v := variableName(name)
	return &v
}

func (v variableName) String() string {
	if Pretty {
		// see: https://stackoverflow.com/questions/13559276/can-i-write-italics-to-the-python-shell/13559470#13559470
		return fmt.Sprintf("\x1B[3m%s\x1B[0m", string(v))
	}
	return string(v)
}

func (f Free) String() string {
	if Pretty {
		// see: https://stackoverflow.com/questions/60064647/how-do-i-use-subscript-digits-in-my-c-program
		return fmt.Sprintf("\u208B%s", lowerDigits(int(f)))
	}
	return fmt.Sprintf("_.%d", f)
}

func lowerDigits(num int) string {
	var acc []rune
	n := digits(num)
	for i := len(n) - 1; i >= 0; i-- {
		acc = append(acc, '\u2080'+rune(n[i]))
	}
	return string(acc)
}

func digits(n int) []int {
	if n == 0 {
		return []int{0}
	}
	var acc []int
	for n > 0 {
		acc = append(acc, n%10)
		n = n / 10
	}
	return acc
}
