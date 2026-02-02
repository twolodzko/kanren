# kanren (関連)

Kanren is a superset of Scheme adding relational/logic programming features to it. It is described by
Friedman, Byrd, and Kiselyov in *The Reasoned Schemer* book, in the [thesis by Byrd (2009)][byrd09],
and multiple papers including [Byrd and Friedman (2006)][byrd06]. More details about it
can be found on a dedicated <https://minikanren.org/> website.

This implementation is based on minimal Scheme interpreter written in Go, with the dedicated kanren
primitives also implemented from scratch in Go (unlike original implementation using pure Scheme
and making use of it's macros). It is a minimal implementation, created as a learning exercise.

## Scheme methods

The Scheme interpreter is based on my [gosch] implementation, with some improvements.
The original implementation was thinned and improved, including support for proper dotted pairs.

A small subset of Scheme build-in methods is available, e.g. `define`, `lambda`, `let`,
`quote` (`'x`), `quasiquote` (``` `x ```), `unquote` (`,x`),
`cons`, `car`, `cdr`, `null?`, `pair?`, `=`, `and`, `or`, `not`, `cond`,
and basic arithmetic operations. `(load "path")` can be used for running another Scheme script.

The supported atomic data types are integers and booleans (`#t` and `#f`), strings
can be used as constants, but there is no string manipulation procedures implemented.

## miniKanren methods

Only the minimal subset of kanren was implemented. This includes the following functions:

* `(== e1 e2)` [unify] the results of the two expressions.
* `(fresh (x ...) g1 g2 ...)` initialize the fresh variables `x ...`. It works in a similar way as `let` in Scheme.
* `(conde (g1a g2a ...) (g1b g2b ...) ... )` on subsequent calls, returns the result of succeeding branches.
  It works in a similar way as `cond` in Scheme, but checks the following branches on subsequent queries.
* `(run* (x) g1 g2 ...)` run the `g1 g2 ...` goals and collect the results for the `x` target variable.
  Repeat until failure.
* `(run n (x) g1 g2 ...)` run the `g1 g2 ...` goals and collect the results for the `x` target variable.
  Repeat at least `n` times.
* `succeed` is a goal that always succeeds.
* `fail` is a goal that always fails.

The language is fully specified and explained in the great *The Reasoned Schemer* book. The code is tested using 
an integration test that runs all the relevant examples from the book.

[gosch]: https://github.com/twolodzko/gosch
[byrd09]: https://scholarworks.iu.edu/iuswrrest/api/core/bitstreams/27f1ebb8-5114-4fa5-b598-dcfaddfd6af5/content
[byrd06]: http://scheme2006.cs.uchicago.edu/12-byrd.pdf
[unify]: https://www.cs.bu.edu/fac/snyder/publications/UnifChapter.pdf
