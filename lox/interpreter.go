package lox

type Interpreter struct{}

func NewIterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expr Expr) any {
	result := expr.Accept(i)
	return result
}

func (i *Interpreter) VisitBinaryExpr(expr BinaryExpr) any {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		return left.(float64) - right.(float64)
	case PLUS:
		right_num, right_num_ok := right.(float64)
		left_num, left_num_ok := left.(float64)
		right_str, right_str_ok := right.(string)
		left_str, left_str_ok := left.(string)

		switch {
		case left_num_ok && right_num_ok:
			return left_num + right_num
		case left_str_ok && right_str_ok:
			return left_str + right_str
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
			return left_num * right_num
		case left_num_ok && right_str_ok:
			return i.multiplyString(int(left_num), right_str)
		case right_num_ok && left_str_ok:
			return i.multiplyString(int(right_num), left_str)
		default:
			panic("Cannot multiply string by string")
		}
	case SLASH:
		return left.(float64) / right.(float64)
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL:
		return i.isEqual(left, right)
	}

	return nil
}
func (i *Interpreter) VisitGroupingExpr(expr GroupingExpr) any {
	return i.evaluate(expr.expression)
}
func (i *Interpreter) VisitLiteralExpr(expr LiteralExpr) any {
	return expr.value
}
func (i *Interpreter) VisitUnaryExpr(expr UnaryExpr) any {
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		return -right.(float64)
	case BANG:
		return !i.isTruthy(right)
	}

	// unreachable
	return nil
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}
func (i *Interpreter) isTruthy(object any) bool {
	// false and nil is false, everything else is true
	if object == nil {
		return false
	}
	switch v := object.(type) {
	case bool:
		return v
	default:
		return true
	}
}
func (i *Interpreter) isEqual(a, b any) bool {
	// TODO: would just "return a == b" be enough in Golang?
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return a == b
}
func (i *Interpreter) multiplyString(count int, s string) string {
	result := ""
	for range count {
		result += s
	}
	return result
}

type InterpreterError struct{}
