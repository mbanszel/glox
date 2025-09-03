package main

type Token struct {}

type Scanner struct {
  source string
  tokens []Token
}

func NewScanner(source string) Scanner {
  scanner := Scanner{source: source}
  return scanner
}

func (s *Scanner) ScanTokens() []Token {
  tokens := make([]Token, 10, 10)
  return tokens
}
