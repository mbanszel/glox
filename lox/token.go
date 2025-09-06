package lox

import (
	"fmt"
)

type TokenType int

//go:generate stringer -type=TokenType
const (
	// single character tokens
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// literals
	IDENTIFIER
	STRING
	NUMBER

	// keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOL
	EOF
)

func (tt TokenType) String() string {
	return [...]string{
		"LEFT_PAREN",
		"RIGHT_PAREN",
		"LEFT_BRACE",
		"RIGHT_BRACE",
		"COMMA",
		"DOT",
		"MINUS",
		"PLUS",
		"SEMICOLON",
		"SLASH",
		"STAR",
		"BANG",
		"BANG_EQUAL",
		"EQUAL",
		"EQUAL_EQUAL",
		"GREATER",
		"GREATER_EQUAL",
		"LESS",
		"LESS_EQUAL",
		"IDENTIFIER",
		"STRING",
		"NUMBER",
		"AND",
		"CLASS",
		"ELSE",
		"FALSE",
		"FUN",
		"FOR",
		"IF",
		"NIL",
		"OR",
		"PRINT",
		"RETURN",
		"SUPER",
		"THIS",
		"TRUE",
		"VAR",
		"WHILE",
		"EOL",
		"EOF",
	}[tt]
}

// ===========================================================================================
var keywords = map[string]TokenType{
  "and": AND,
  "class": CLASS,
  "else": ELSE,
  "false": FALSE,
  "for": FOR,
  "fun": FUN,
  "if": IF,
  "nil": NIL,
  "or": OR,
  "print": PRINT,
  "return": RETURN,
  "super": SUPER,
  "this": THIS,
  "true": TRUE,
  "var": VAR,
  "while": WHILE,
}
// ===========================================================================================
type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func (t Token) String() string {
  return fmt.Sprintf("%s %s %v", t.tokenType, t.lexeme, t.literal)
}