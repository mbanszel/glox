# glox
Crafting Interpreter's lox written in golang.

See [Crafting Interpreters](https://craftinginterpreters.com/).

# TODO
- [x] Write the glox interpreter
- [ ] fix bug when there is assignment (incorrectly)  as the first statement, e.g. x = 5, then
      the parser ends up in an infinit loop. Instead, it should refuse that and return ParserError.
- [ ] Write the glox transpiler to C (the third part of the book)
- [ ] Write a formatter tool (that parses the language and from the AST it generates back the source
      code but nicely formatted, like `black` or `go fmt` do)