package lox

import (
	"fmt"
	"os"
)

var HadError = false
var HadRuntimeError = false

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

func parserError(err ParserError) {
	token := err.GetToken()
	message := err.GetMessage()
	fmt.Fprintf(os.Stderr, "[line %v] %s at %s\n", token.Line, message, token.Lexeme)
	HadError = true
}

type LoxError interface{
	GetToken() Token
	GetMessage() string
}

func runtimeError(err RuntimeError) {
	fmt.Fprintf(os.Stderr, "[line %v] %s\n", err.GetToken().Line, err.GetMessage())
	HadRuntimeError = true
}
