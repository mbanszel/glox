package lox

import (
	"fmt"
	"strings"
)

type Interpreter struct {
	environment *Environment
}

// VisitAssignmentExpr implements ExprVisitor.
func (i *Interpreter) VisitAssignmentExpr(expr AssignmentExpr) (any, LoxError) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	return i.environment.assign(expr.name, value)
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(nil),
	}
}

func (i *Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			fmt.Println("Error...")
			runtimeError(err.(RuntimeError))
		}
	}
}

func (i *Interpreter) VisitBinaryExpr(expr BinaryExpr) (any, LoxError) {
	left, err := i.evaluate(expr.left)
	var numbers []float64
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.TokenType {
	case MINUS:
		if numbers, err = validateNumber(expr.operator, left, right); err != nil {
			return nil, err
		}
		return numbers[0] - numbers[1], nil
	case PLUS:
		right_num, right_num_ok := right.(float64)
		left_num, left_num_ok := left.(float64)
		right_str, right_str_ok := right.(string)
		left_str, left_str_ok := left.(string)

		switch {
		case left_num_ok && right_num_ok:
			return left_num + right_num, nil
		case left_str_ok && right_str_ok:
			return left_str + right_str, nil
		default:
			return nil, &RuntimeErrorObj{
				expr.operator,
				"Operands must be numbers or strings",
			}
		}
	case STAR:
		right_num, right_num_ok := right.(float64)
		left_num, left_num_ok := left.(float64)
		right_str, right_str_ok := right.(string)
		left_str, left_str_ok := left.(string)

		switch {
		case left_num_ok && right_num_ok:
			return left_num * right_num, nil
		case left_num_ok && right_str_ok:
			return i.multiplyString(int(left_num), right_str)
		case right_num_ok && left_str_ok:
			return i.multiplyString(int(right_num), left_str)
		default:
			return nil, &RuntimeErrorObj{
				expr.operator,
				"Cannot multiply string by string",
			}
		}
	case SLASH:
		if numbers, err = validateNumber(expr.operator, left, right); err != nil {
			return nil, err
		}
		if numbers[1] == 0 {
			return nil, &RuntimeErrorObj{
				expr.operator,
				"Division by zero.",
			}
		}
		return numbers[0] / numbers[1], nil
	case GREATER:
		if numbers, err = validateNumber(expr.operator, left, right); err != nil {
			return nil, err
		}
		return numbers[0] > numbers[1], nil
	case GREATER_EQUAL:
		if numbers, err = validateNumber(expr.operator, left, right); err != nil {
			return nil, err
		}
		return numbers[0] >= numbers[1], nil
	case LESS:
		if numbers, err = validateNumber(expr.operator, left, right); err != nil {
			return nil, err
		}
		return numbers[0] < numbers[1], nil
	case LESS_EQUAL:
		if numbers, err = validateNumber(expr.operator, left, right); err != nil {
			return nil, err
		}
		return numbers[0] <= numbers[1], nil
	case BANG_EQUAL:
		value, err := i.isEqual(left, right)
		return !value, err
	case EQUAL:
		return i.isEqual(left, right)
	}

	return nil, nil
}
func (i *Interpreter) VisitGroupingExpr(expr GroupingExpr) (any, LoxError) {
	return i.evaluate(expr.expression)
}
func (i *Interpreter) VisitLiteralExpr(expr LiteralExpr) (any, LoxError) {
	return expr.value, nil
}
func (i *Interpreter) VisitLogicalExpr(expr LogicalExpr) (any, LoxError) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}

	left_truthiness, err := i.isTruthy(left)
	if err != nil {
		return nil, err
	}
	if (expr.operator.TokenType == OR) {
		if left_truthiness {
			return left, nil
		}
	} else {
		if !left_truthiness {
			return left, nil
		}
	}
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	return right, nil
}
func (i *Interpreter) VisitUnaryExpr(expr UnaryExpr) (any, LoxError) {
	right, _ := i.evaluate(expr.right)

	switch expr.operator.TokenType {
	case MINUS:
		return -right.(float64), nil
	case BANG:
		value, err := i.isTruthy(right)

		return !value, err
	}

	// unreachable
	return nil, &RuntimeErrorObj{expr.operator, "Unexpected unary operator"}
}
func (i *Interpreter) VisitVariableExpr(expr VariableExpr) (any, LoxError) {
	return i.environment.get(expr.name)
}

// ------------------------------------------------------------------------------------------
func (i *Interpreter) VisitExpressionStmt(stmt ExpressionStmt) (any, LoxError) {
	v, err := i.evaluate(stmt.expression)
	if err != nil {
		return nil, err
	}
	return v, nil
}
func (i *Interpreter) VisitIfStmt(stmt IfStmt) (any, LoxError) {
	val, err := i.evaluate(stmt.condition)
	if err != nil {
		return nil, err
	}

	condition, err := i.isTruthy(val)
	if err != nil {
		return nil, err
	}

	if condition {
		return i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		return i.execute(stmt.elseBranch)
	}
	return nil, nil
}
func (i *Interpreter) VisitWhileStmt(stmt WhileStmt) (any, LoxError) {
	for {
		val, err := i.evaluate(stmt.condition)
		if err != nil {
			return nil, err
		}
		condition, err := i.isTruthy(val)
		if err != nil {
			return nil, err
		}
		if !condition {
			return nil, nil
		}
		_, err = i.execute(stmt.body)
		if err != nil {
			return nil, err
		}
	}
}
func (i *Interpreter) VisitPrintStmt(stmt PrintStmt) (any, LoxError) {
	v, err := i.evaluate(stmt.expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(i.stringify(v))
	return v, nil
}
func (i *Interpreter) VisitVarStmt(stmt VarStmt) (any, LoxError) {
	var value any
	var err LoxError
	if stmt.initializer != nil {
		value, err = i.evaluate(stmt.initializer)
		if err != nil {
			return nil, err
		}
	}
	i.environment.define(stmt.name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) evaluate(expr Expr) (any, LoxError) {
	return expr.Accept(i)
}
func (i *Interpreter) isTruthy(object any) (bool, LoxError) {
	// false and nil is false, everything else is true
	if object == nil {
		return false, nil
	}
	switch v := object.(type) {
	case bool:
		return v, nil
	default:
		return true, nil
	}
}
func (i *Interpreter) isEqual(a, b any) (bool, LoxError) {
	// TODO: would just "return a == b" be enough in Golang?
	if a == nil && b == nil {
		return true, nil
	}
	if a == nil {
		return false, nil
	}

	return a == b, nil
}
func (i *Interpreter) stringify(object any) string {
	if object == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", object)
}
func (i *Interpreter) multiplyString(count int, s string) (string, RuntimeError) {
	var sb strings.Builder
	for range count {
		sb.WriteString(s)
	}
	return sb.String(), nil
}

func validateNumber(operator Token, numbers ...any) ([]float64, RuntimeError) {
	converted := make([]float64, 0, len(numbers))

	for _, aNumber := range numbers {
		f, ok := aNumber.(float64)
		converted = append(converted, f)
		if !ok {
			return []float64{}, &RuntimeErrorObj{operator, "Operands must be numbers."}
		}
	}
	return converted, nil
}

type RuntimeError interface {
	LoxError
	GetToken() Token
	GetMessage() string
}
type RuntimeErrorObj struct {
	token   Token
	message string
}

func (e *RuntimeErrorObj) GetToken() Token {
	return e.token
}
func (e *RuntimeErrorObj) GetMessage() string {
	return e.message
}

func (i *Interpreter) execute(stmt Stmt) (any, LoxError) {
	return stmt.Accept(i)
}

func (i *Interpreter) VisitBlockStmt(stmt BlockStmt) (any, LoxError) {
	return i.executeBlock(stmt.statements, NewEnvironment(i.environment))
}

func (i *Interpreter) descope() {
	// leave the current lexical scope
	if i.environment.Enclosing == nil {
		return
	}

	i.environment = i.environment.Enclosing
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) (any, RuntimeError) {
	i.environment = environment
	defer i.descope()

	for _, stmt := range statements {
		value, err := i.execute(stmt)
		if err != nil {
			return value, err
		}
	}
	return nil, nil
}
