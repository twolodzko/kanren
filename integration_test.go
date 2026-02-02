package main

import (
	"testing"

	"github.com/twolodzko/kanren/eval"
	"github.com/twolodzko/kanren/types"
)

func TestIntegration(t *testing.T) {
	types.Pretty = true
	var files = []string{
		"examples/other.scm",
		"examples/peano.scm",
		"examples/mktests.scm",
	}
	for _, file := range files {
		env := eval.DefaultEnv()
		_, err := eval.LoadEval(file, env)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
	}
}
