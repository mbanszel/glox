package lox

import (
	"fmt"
)

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []any) (any, LoxError)
	Arity() int
}

type LoxFunction struct {
	declaration FunctionStmt
}

func NewLoxFunction(declaration FunctionStmt) *LoxFunction {
	return &LoxFunction{declaration: declaration}
}

func (lf *LoxFunction) Arity() int {
	return len(lf.declaration.params)
}

func (lf *LoxFunction) Call(i *Interpreter, arguments []any) (any, LoxError) {
	environment := NewEnvironment(i.environment)
	for j := 0; j < lf.Arity(); j++ {
		environment.define(lf.declaration.params[j].Lexeme, arguments[j])
	}
	_, err := i.executeBlock(lf.declaration.body, environment)
	if err!=nil {
		return_value, ok := err.(*ReturnObj)
		if ok {
			return return_value.GetValue(), nil
		}
		return nil, err
	}
	return nil, nil
}

func (lf *LoxFunction) Stringer() string {
	return fmt.Sprintf("<fn %s >", lf.declaration.name.Lexeme)
}
