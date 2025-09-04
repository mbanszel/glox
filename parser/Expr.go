package parser

import (
  "github.com/mbanszel/glox/scanner"
)

type Expr interface {
}

type Binary struct {
  left Expr
  operator scanner.Token
  right Expr
}

func NewBinary(left Expr, operator scanner.Token, right Expr) Binary {
  return Binary{
    left:left,
    operator:operator,
    right:right,
  }
}

type Grouping struct {
  expression Expr
}

func NewGrouping(expression Expr) Grouping {
  return Grouping{
    expression:expression,
  }
}

type Literal struct {
  value any
}

func NewLiteral(value any) Literal {
  return Literal{
    value:value,
  }
}

type Unary struct {
  operator scanner.TokenType
  right Expr
}

func NewUnary(operator scanner.TokenType, right Expr) Unary {
  return Unary{
    operator:operator,
    right:right,
  }
}

