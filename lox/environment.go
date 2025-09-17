package lox

import (
	"fmt"
)

type Environment struct {
	Enclosing *Environment
	Values map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		Values: make(map[string]any),
	}
}

func (e *Environment) define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) assign(name Token, value any) (any, RuntimeError) {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return value, nil
	}
	if e.Enclosing != nil {
		return e.Enclosing.assign(name, value)
	}
	return nil, &RuntimeErrorObj{
		name,
		"Undefined variable '" + name.Lexeme + "'",
	}
}

func (e *Environment) get(name Token) (any, RuntimeError) {
	if value, ok := e.Values[name.Lexeme]; ok {
		if value == nil {
			return "", &RuntimeErrorObj{
				name,
				fmt.Sprintf("Uninitialized variable '%s'", name.Lexeme),
			}
		}
		return value, nil
	}
	if e.Enclosing != nil {
		return e.Enclosing.get(name)
	}
	return "", &RuntimeErrorObj{name, "Undefined variable '" + name.Lexeme + "'"}
}