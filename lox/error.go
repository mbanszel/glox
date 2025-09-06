package lox

import (
	"fmt"
)

var HadError = false

func Emit(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Printf("[line %d ] Error %s: %s\n", line, where, message)
	HadError = true
}

func Error(token Token, message string) {
	if token.TokenType == EOF {
		Report(token.Line, " at end", message)
	} else {
		Report(
			token.Line,
			fmt.Sprintf(" at '%s'", token.Lexeme),
			message,
		)
	}
}

type LoxError interface{}
