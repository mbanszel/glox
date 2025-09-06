package lox

import (
	"fmt"
)

type AstPrinter struct {
	RPN bool
}

func NewAstPrinter() AstPrinter {
	return AstPrinter{RPN: false}
}

func NewAstRPNPrinter() AstPrinter {
	return AstPrinter{RPN: true}
}

func (p AstPrinter) Print(e Expr) string {
	return fmt.Sprintf("%v", e.Accept(p))
}

func (p AstPrinter) VisitBinaryExpr(expr BinaryExpr) any {
	return p.parenthesize(expr.operator.lexeme, expr.left, expr.right)
}

func (p AstPrinter) VisitGroupingExpr(expr GroupingExpr) any {
	return p.parenthesize("group", expr.expression)
}

func (p AstPrinter) VisitLiteralExpr(expr LiteralExpr) any {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.value)
}

func (p AstPrinter) VisitUnaryExpr(expr UnaryExpr) any {
	return p.parenthesize(expr.operator.lexeme, expr.right)
}

func (p AstPrinter) parenthesize(name string, expressions ...Expr) string {
	var out string
	if p.RPN {
		for _, e := range expressions {
			out += fmt.Sprintf("%s", e.Accept(p))
			out += " "
		}
		out += name

		return out
	}
	out = fmt.Sprintf("(%s", name)
	for _, e := range expressions {
		out += " "
		out += fmt.Sprintf("%s", e.Accept(p))
	}
	out += ")"
	return out
}
