package lox

import (
	"testing"
)

var testcases = []Expr {
	NewBinaryExpr(
		NewUnaryExpr(
			Token{MINUS, "-", nil, 1},
			NewLiteralExpr(123),
		),
		Token{STAR, "*", nil, 1},
		NewGroupingExpr(NewLiteralExpr(45.67)),
	),
	NewBinaryExpr(
		NewBinaryExpr(
			NewLiteralExpr(1),
			Token{PLUS, "+", nil, 1},
			NewLiteralExpr(2),
		),
		Token{STAR, "*", nil, 1},
		NewBinaryExpr(
			NewLiteralExpr(3),
			Token{MINUS, "-", nil, 1},
			NewLiteralExpr(4),
		),
	),
}

var expected = []string {
	"(* (- 123) (group 45.67))",
	"(* (+ 1 2) (- 3 4))",
}

func TestAstPrinter(t *testing.T) {
	printer := NewAstPrinter()

	for i := range len(testcases) {
		result := printer.Print(testcases[i])

		if result != expected[i] {
			t.Errorf("failed: %s != %s", result, expected[i])
		}
	}
}

func TestAstRPNPrinter(t *testing.T) {
	expr := NewBinaryExpr(
		NewBinaryExpr(
			NewLiteralExpr(1),
			Token{PLUS, "+", nil, 1},
			NewLiteralExpr(2),
		),
		Token{STAR, "*", nil, 1},
		NewBinaryExpr(
			NewLiteralExpr(3),
			Token{MINUS, "-", nil, 1},
			NewLiteralExpr(4),
		),
	)
	printer := NewAstRPNPrinter()
	result := printer.Print(expr)

	expected := "1 2 + 3 4 - *"
	if result != expected {
		t.Errorf("failed: %s != %s", result, expected)
	}
}
