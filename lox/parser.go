package lox
// recursive descent parser for (g)lox interpreter

type Parser struct {
	tokens []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() Expr {
	expr, err := p.expression()
	if err != nil {
		return nil
	}
	return expr
}

// The lox grammar:
// ----------------
// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;

func (p *Parser) expression() (Expr, *ParserError) {
	return p.equality()
}

func (p *Parser) equality() (Expr, *ParserError) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, *ParserError) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for (p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL)) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) term() (Expr, *ParserError) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for (p.match(MINUS, PLUS)) {
		operator := p.previous()
		right, err := p.factor()
		if err!=nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr, nil

}

func (p *Parser) factor() (Expr, *ParserError) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for (p.match(SLASH, STAR)) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, *ParserError) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnaryExpr(operator, right), nil
	} 
	return p.primary()
}

func (p *Parser) primary() (Expr, *ParserError) {
	switch {
	case p.match(FALSE):
		return NewLiteralExpr(false), nil
	case p.match(TRUE):
		return NewLiteralExpr(true), nil
	case p.match(NIL):
		return NewLiteralExpr(nil), nil
	case p.match(NUMBER, STRING):
		return NewLiteralExpr(p.previous().literal), nil
	case p.match(LEFT_PAREN):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return NewGroupingExpr(expr), nil
	default:
		return nil, p.error(p.peek(), "Expect statement.")
	}
}

// -----------------------------------------------------------------

func (p *Parser) consume(tokenType TokenType, message string) (Token, *ParserError) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return Token{}, p.error(p.peek(), message)

}

func (p *Parser) match(types ...TokenType) bool {
	for _, tokenType := range types {
		if (p.check(tokenType)) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if (p.isAtEnd()) {return false}
	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() Token {
	// TODO why this way? What is isAtEnd is true -- advnace then returns the last token?
	if !p.isAtEnd() {p.current++}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current - 1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {return}

		switch p.peek().tokenType {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		default:
			p.advance()
		}
	}
}

func (p *Parser) error(token Token, message string) *ParserError {
	Error(token, message)
	return &ParserError{}
}

type ParserError struct{}