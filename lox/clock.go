package lox

import (
	"time"
)

type ClockNativeFunction struct{}

func (c ClockNativeFunction) Arity() int {
	return 0
}

func (c ClockNativeFunction) Call(i *Interpreter, arguments []any) (any, LoxError) {
	return time.Now().UnixMilli(), nil
}
