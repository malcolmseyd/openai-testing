# OpenAI Codex: Lisp Interpreter

My prompt was:

```go
// Write a Lisp interpreter in Go
```

Its output included an entire Lisp interpreter in ~350 lines of Go! There were a few bugs that I had to resolve but I'm really impressed at the solution that it came up with.

[Here's a link to the changes that I've made.](https://github.com/malcolmseyd/openai-testing/compare/6ac745dad68d1ad911dd81d803ad3aa28dfdcb12...38ce4f9940d81ae0c701383b99e0e0cc83bbed98#diff-a04fe97432699e7cc309b5a6656a9ba7ca72e860cf25f53e85ab2944c7124c8b)

## What it got wrong

The generated code was surprisingly good, and the mistakes that it made were something that I could see a human accidentally doing. The bugs were fixed in order of most-to-least severe, where the first bug made compilation fail.

Well, the commit history says it all, but I'll make a short summary:

* [It tried to store an `int` and a `bool` in the same variable for number comparison functions](https://github.com/malcolmseyd/openai-testing/commit/60568018c4c65159218768ce3959ee94950c4e14)
* [It thought that typing asserting on a defined type was equivalent to asserting on the underlying type.](https://github.com/malcolmseyd/openai-testing/commit/f6f59282d4c7f6cfc41dab2e85f28e0ced3c81a6) Honestly, I didn't know that this wasn't possible and I learned something new from this.
* [The parser didn't consume input after parsing it, using a weird hack with the symbol type so that only a flat list of symbols would parse](https://github.com/malcolmseyd/openai-testing/commit/f42b8a64bdba3969a2c7515e7f3016a4f9eb57fe)
* [The parser would consume closing parentheses if they were immediately after a number or symbol](https://github.com/malcolmseyd/openai-testing/commit/e38deaffa993fe4ebc166111b6fd0ea7b4256f5e)
* [If the first element of a list wasn't a symbol, it treated it as a list literal.](https://github.com/malcolmseyd/openai-testing/commit/354e239d6e5ef79832e7c99eb66efc7b06cc903b) I could see this being a valid implementation behaviour, and I only changed it because I thought it was weird because it didn't allow immediate lambda application, breaking expressions like the [Y combinator](https://rosettacode.org/wiki/Y_combinator#Scheme) for example. With this change, the example I linked should run! (You need to replace `define` with `def` though, and `display` with `print`, and remove `newline`).
* [The quote reader macro created a weird "quote" expression type that couldn't be used, when it should have simply expanded to `(quote <expr>)`](https://github.com/malcolmseyd/openai-testing/commit/38ce4f9940d81ae0c701383b99e0e0cc83bbed98)
* It's missing an `eq?` function which means it can't compare symbols.

## What it got right

Although there were bugs, it did successfully design a metacircular evaluator in shockingly few lines of code. The lisp that it wrote has:

* A recursive descent parser
* Variables
* Functions with lexical scope
* Conditionals and comparison operators (`null?`, `=`, `>`, etc)
* Numbers
* Strings
* A quote reader macro
* `car`, `cdr`, `cons`, and `list`
* An interactive REPL
* Of course, an `eval` function that can evaluate forms

As far as design, here's a few things that made me think:

* Using Go slices to represent Lisp lists is very simple and very efficient. It allows Go's `nil` to represent `'()` which is very cool. However, this is only possible because this Lisp only supports immutable values, so you couldn't splice values into the middle of a list with `set-cdr!` if that existed.
* String interning is not necessary for Lisp symbols, although it does improve performance.
* Rigid, error-prone Go code looks very nice, while error checking often makes it look quite noisy. I know that it was a concious decision for Go's error handling to be very explicit, but it rubs me the wrong way that bad code looks so good and vice-versa.