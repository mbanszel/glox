package lox

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []any) (any, LoxError)
	Arity() int
}