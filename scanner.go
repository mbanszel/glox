package main

import (
  "strconv"
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

// ===========================================================================================
type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) Scanner {
	scanner := Scanner{
		source:  source,
		tokens:  make([]Token, 0, 100),
		start:   0,
		current: 0,
		line:    1,
	}
	return scanner
}

func (s *Scanner) ScanTokens() []Token {

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{tokenType: EOF, lexeme: "", literal: "", line: s.line})

	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
  case '!':
    if s.match('=') {
      s.addToken(BANG_EQUAL)
    } else {
      s.addToken(BANG)
    }
  case '<':
    if s.match('=') {
      s.addToken(LESS_EQUAL)
    } else {
      s.addToken(LESS)
    }
  case '>':
    if s.match('=') {
      s.addToken(GREATER_EQUAL)
    } else {
      s.addToken(GREATER)
    }
  case '=':
    if s.match('=') {
      s.addToken(EQUAL_EQUAL)
    } else {
      s.addToken(EQUAL)
    }
  case '/':
    if s.match('/') {
      for s.peek() != '\n' && !s.isAtEnd() {
        s.advance()
      }
    } else {
      s.addToken(SLASH)
    }
  case ' ':
  case '\r':
  case '\t':
  case '\n':
    s.line++
  case '"':
    s.string()
	default:
    if isDigit(c) {
      s.number()
    } else if isAlpha(c) {
      s.identifier()
    } else {
      emit_error(s.line, "Unexpected character")
    }
	}
}

func isDigit(c byte) bool {
  return c >= '0' && c<='9'
}

func isAlpha(c byte) bool {
  return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
  return isAlpha(c) || isDigit(c)
}

func (s *Scanner) advance() byte {
	next := s.source[s.current]
	s.current++
	return next
}

func (s *Scanner) match(expected byte) bool {
  if s.isAtEnd() {
    return false
  }

  if s.source[s.current] != expected {
    return false
  }

  s.current++
  return true
}

func (s *Scanner) peek() byte {
  if s.isAtEnd() {
    return 0 // byte with ordinal 0 should not appear in the source code...?
  } else {
    return s.source[s.current]
  }
}

func (s *Scanner) peekNext() byte {
  nextCurrent := s.current + 1
  if nextCurrent >= len(s.source) {return 0}
  return s.source[nextCurrent]
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{tokenType, text, literal, s.line})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) string() {
  for s.peek() != '"' && !s.isAtEnd() {
    if (s.peek() == '\n') {
      s.line++
    }
    s.advance()
  }

  if s.isAtEnd() {
    emit_error(s.line, "Unterminated string.")
    return
  }

  // the closing "
  s.advance()
  s.addTokenWithLiteral(STRING, s.source[s.start+1:s.current-1])
}

func (s *Scanner) number() {
  for isDigit(s.peek()) {s.advance()}

  if s.peek() == '.' && isDigit(s.peekNext()) {
    s.advance()
    for isDigit(s.peek()) {s.advance()}
  }
  value, error := strconv.ParseFloat(s.source[s.start:s.current], 64)
  if error!=nil {
    panic("Unexpected error")
  }
  s.addTokenWithLiteral(NUMBER, value)
}

func (s *Scanner) identifier() {
  for isAlphaNumeric(s.peek()) {s.advance()} 

  value := s.source[s.start:s.current]
  tokenType, ok := keywords[value]
  if !ok {
    tokenType = IDENTIFIER
  }
  s.addToken(tokenType)
}