package lox

// recursive descent parser for (g)lox interpreter

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() []Stmt {
	// expr, err := p.expression()
	// if err != nil {
	// 	return nil
	// }

	statements := []Stmt{}
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			parserError(err)
			p.advance()
		}
		statements = append(statements, stmt)
	}
	return statements
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

func (p *Parser) expression() (Expr, ParserError) {
	return p.assignment();
}

func (p *Parser) assignment() (Expr, ParserError) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		variable_expr, ok := expr.(VariableExpr)
		if ok {
			name := variable_expr.name
			return NewAssignmentExpr(name, value), nil
		}
		Error(equals, "Invalid assignment target")
	}
	return expr, nil
}

func (p *Parser) equality() (Expr, ParserError) {
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

func (p *Parser) comparison() (Expr, ParserError) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) term() (Expr, ParserError) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr, nil

}

func (p *Parser) factor() (Expr, ParserError) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, ParserError) {
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

func (p *Parser) primary() (Expr, ParserError) {
	switch {
	case p.match(FALSE):
		return NewLiteralExpr(false), nil
	case p.match(TRUE):
		return NewLiteralExpr(true), nil
	case p.match(NIL):
		return NewLiteralExpr(nil), nil
	case p.match(NUMBER, STRING):
		return NewLiteralExpr(p.previous().Literal), nil
	case p.match(IDENTIFIER):
		return NewVariableExpr(p.previous()), nil
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
		return nil, p.error("Expect statement.")
	}
}

// -----------------------------------------------------------------

func (p *Parser) consume(tokenType TokenType, message string) (Token, ParserError) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return Token{}, p.error(message)

}

func (p *Parser) match(types ...TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() Token {
	// TODO why this way? What is isAtEnd is true -- advnace then returns the last token?
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		default:
			p.advance()
		}
	}
}

func (p *Parser) error(message string) ParserError {
	token := p.peek()
	return &ParserErrorObj{
		Token: token,
		Message: message,
	}
}

type ParserError interface {
	LoxError
}

type ParserErrorObj struct{
	Token Token
	Message string
}

func (pe ParserErrorObj) GetToken() Token {
	return pe.Token
}

func (pe ParserErrorObj) GetMessage() string {
	return pe.Message
}

func (p *Parser) declaration() (Stmt, ParserError) {
	var stmt Stmt
	var err ParserError
	if p.match(VAR) {
		stmt, err = p.varDeclaration()
	} else {
		stmt, err = p.statement()
	}
	if err != nil {
		p.synchronize()
	}
	return stmt, err
	
}

func (p *Parser) varDeclaration() (Stmt, ParserError) {
	name, err := p.consume(IDENTIFIER, "Expect variable name")
	if err != nil {
		return nil, err
	}
	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after variable declaration")
	if err != nil {
		return nil, err
	}

	return NewVarStmt(name, initializer), nil
}


func (p *Parser) statement() (Stmt, ParserError) {
	if p.match(PRINT) {return p.printStatement()}
	if p.match(LEFT_BRACE) {
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return NewBlockStmt(stmts), nil
	}

	return p.expressionStatement()
}

func (p *Parser) block() ([]Stmt, ParserError) {
	var statements []Stmt

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		dclr, err := p.declaration()
		if err != nil {
			return statements, err
		}
		statements = append(statements, dclr)
	}
	_, err := p.consume(RIGHT_BRACE, "Expect '}' after block")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) printStatement() (Stmt, ParserError) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(SEMICOLON, "Expect ';' after value.")
	return NewPrintStmt(value), nil
}

func (p *Parser) expressionStatement() (Stmt, ParserError) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(SEMICOLON, "Expect ';' after value.")
	return NewExpressionStmt(expr), nil

}