package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/twolodzko/kanren/envir"
	"github.com/twolodzko/kanren/eval"
	"github.com/twolodzko/kanren/repl"
	"github.com/twolodzko/kanren/types"
)

const prompt string = "> "

func main() {
	var (
		showHelp bool
		keepRepl bool
	)

	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.BoolVar(&eval.Debug, "debug", false, "run in debug mode")
	flag.BoolVar(&types.Pretty, "pretty", false, "prettify the outputs")
	flag.BoolVar(&keepRepl, "keep", false, "open REPL after evaluating files")
	flag.Parse()

	if showHelp {
		printHelp()
		return
	}

	env := eval.DefaultEnv()
	if flag.NArg() > 0 {
		evalFiles(env, flag.Args())
		if !keepRepl {
			return
		}
	}
	startRepl(env)
}

func evalFiles(env *envir.Env, paths []string) {
	var last any = nil
	for _, path := range paths {
		sexprs, err := eval.LoadEval(path, env)
		if err != nil {
			log.Fatalf("ERROR: %v\n", err)
		}
		if len(sexprs) > 0 {
			last = sexprs[len(sexprs)-1]
		}
	}
	fmt.Printf("%s\n", types.ToString(last))
}

func startRepl(env *envir.Env) {
	repl := repl.NewRepl(os.Stdin, env)

	fmt.Println("Press ^C to exit.")
	fmt.Println()

	for {
		fmt.Printf("%s", prompt)
		objs, err := repl.Repl()
		if err != nil {
			print(fmt.Sprintf("ERROR: %s", err))
			continue
		}
		for _, obj := range objs {
			print(types.ToString(obj))
		}
	}
}

func printHelp() {
	fmt.Printf("%s FLAGS [script]\n", os.Args[0])
	fmt.Println()
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

func print(msg string) {
	_, err := io.WriteString(os.Stdout, fmt.Sprintf("%s\n", msg))
	if err != nil {
		log.Fatal(err)
	}
}
