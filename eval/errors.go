package eval

import (
	"errors"
	"fmt"

	"github.com/twolodzko/kanren/types"
)

var ArityError = errors.New("wrong number of arguments")
var SyntaxError = errors.New("invalid syntax")
var TypeError = errors.New("invalid type")

type WrongArg struct {
	Val any
}

func (e WrongArg) Error() string {
	return fmt.Sprintf("invalid argument: %s", types.ToString(e.Val))
}

type NonList struct {
	Val any
}

func (e NonList) Error() string {
	return fmt.Sprintf("%s is not a list", types.ToString(e.Val))
}

type NaN struct {
	Val any
}

func (e NaN) Error() string {
	return fmt.Sprintf("%s is not a number", types.ToString(e.Val))
}

type InvalidName struct {
	Val any
}

func (e InvalidName) Error() string {
	return fmt.Sprintf("%s is not a valid name", types.ToString(e.Val))
}
