package lox

type Interpreter struct{}

func NewIterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expr Expr) (any, LoxError) {
	res, err := expr.Accept(i)
	return res, err
}

func (i *Interpreter) VisitBinaryExpr(expr BinaryExpr) (any, LoxError) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.TokenType {
	case MINUS:
		if validateNumber(left, right) {
			return left.(float64) - right.(float64), nil
		}
		return nil, &RuntimeErrorObj{expr.operator, "has invalid operands"}
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
			panic("Incompatible types for +: both must be either string or number")
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
			panic("Cannot multiply string by string")
		}
	case SLASH:
		return left.(float64) / right.(float64), nil
	case GREATER:
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64), nil
	case LESS:
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		return left.(float64) <= right.(float64), nil
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
func (i *Interpreter) multiplyString(count int, s string) (string, *RuntimeErrorObj) {
	result := ""
	for range count {
		result += s
	}
	return result, nil
}

func validateNumber(numbers ...any) bool {
	for _, aNumber := range numbers {
		_, ok := aNumber.(float64)
		if !ok {
			return false
		}
	}
	return true
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
