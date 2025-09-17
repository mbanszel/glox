package lox

import (
	"fmt"
)

type AstPrinter struct {
	RPN bool
}

// VisitAssignmentExpr implements ExprVisitor.
func (p AstPrinter) VisitAssignmentExpr(expr AssignmentExpr) (any, LoxError) {
	panic("unimplemented")
}
func (p AstPrinter) VisitBlockStmt(expr BlockStmt) (any, LoxError) {
	panic("unimplemented")
}
func (p AstPrinter) VisitIfStmt(expr BlockStmt) (any, LoxError) {
	panic("unimplemented")
}
func (p AstPrinter) VisitLogicalExpr(expr LogicalExpr) (any, LoxError) {
	panic("unimplemented")
}
func (p AstPrinter) VisitCallExpr(expr CallExpr) (any, LoxError) {
	panic("unimplemented")
}


func NewAstPrinter() AstPrinter {
	return AstPrinter{RPN: false}
}

func NewAstRPNPrinter() AstPrinter {
	return AstPrinter{RPN: true}
}

func (p AstPrinter) Print(e Expr) string {
	result, err := e.Accept(p)
	if err != nil {
		fmt.Println("Error while printing AST")
	}
	return fmt.Sprintf("%v", result)
}



func (p AstPrinter) VisitBinaryExpr(expr BinaryExpr) (any, LoxError) {
	return p.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (p AstPrinter) VisitGroupingExpr(expr GroupingExpr) (any, LoxError) {
	return p.parenthesize("group", expr.expression)
}

func (p AstPrinter) VisitLiteralExpr(expr LiteralExpr) (any, LoxError) {
	if expr.value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.value), nil
}

func (p AstPrinter) VisitUnaryExpr(expr UnaryExpr) (any, LoxError) {
	return p.parenthesize(expr.operator.Lexeme, expr.right)
}

func (p AstPrinter) VisitVariableExpr(expr VariableExpr) (any, LoxError) {
	return expr.name.Lexeme, nil
}

func (p AstPrinter) VisitVarStmt(stmt VarStmt) (any, LoxError) {
	return p.parenthesize(fmt.Sprintf("(var %s)", stmt.initializer))
}

func (p AstPrinter) parenthesize(name string, expressions ...Expr) (string, LoxError) {
	var out string
	if p.RPN {
		for _, e := range expressions {
			res, _ := e.Accept(p)
			out += fmt.Sprintf("%s", res)
			out += " "
		}
		out += name

		return out, nil
	}
	out = fmt.Sprintf("(%s", name)
	for _, e := range expressions {
		out += " "
		res, _ := e.Accept(p)
		out += fmt.Sprintf("%s", res)
	}
	out += ")"
	return out, nil
}
