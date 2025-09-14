package lox

type Environment struct {
	Values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
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

	return nil, &RuntimeErrorObj{
		name,
		"Undefined variable '" + name.Lexeme + "'",
	}
}

func (e *Environment) get(name Token) (any, RuntimeError) {
	value, ok := e.Values[name.Lexeme]
	if !ok {
		return "", &RuntimeErrorObj{name, "Undefined variable '" + name.Lexeme}
	}
	return value, nil
}