package parser

import (
	"fmt"
	"lang/ast"
	"lang/lexer"
	"lang/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	COMPARISON
	SUM
	PRODUCT
	PREFIX
	FUNCCALL
)

type Parser struct {
	curr token.Token
	next token.Token
	l    *lexer.Lexer
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.advance()
	p.advance()
	return p
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	prog := &ast.Program{make([]ast.Statement, 0)}

	for !p.currTokenIs(token.EOF) {
		s, err := p.parseStatement()
		if err != nil {
			return prog, fmt.Errorf("Error parsing program: %v", err)
		}
		prog.Statements = append(prog.Statements, s)
	}

	return prog, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curr.Type {
	case token.IDENT:
		return p.parseVariableDecl()
	case token.FUNC:
		return p.parseFuncDecl()
	default:
		return nil, fmt.Errorf("Error parsing statement: invalid token: %v", p.curr)
	}
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	var left ast.Expression
	var err error

	switch p.curr.Type {
	case token.LPAREN:
		left, err = p.parseGroupedExpression()
	case token.MINUS:
		left, err = p.parsePrefixExpression()
	case token.INT:
		left, err = p.parseIntegerLiteral()
	case token.IDENT:
		left, err = p.parseIdentifier()
	case token.STRING:
		left, err = p.parseStringLiteral()
	default:
		err = fmt.Errorf("invalid token: %v", p.curr)
	}

	if err != nil {
		return left, fmt.Errorf("Error parsing expression: %v", err)
	}

	for precedence < p.currTokenPrecedence() {
		switch p.curr.Type {
		case token.SLASH:
			fallthrough
		case token.ASTERISK:
			fallthrough
		case token.MINUS:
			fallthrough
		case token.PLUS:
			left, err = p.parseInfixExpression(left)
		default:
			return left, nil
		}

		if err != nil {
			return left, fmt.Errorf("Error parsing expression: %v", err)
		}
	}

	return left, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	p.advance()
	exp, err := p.parseExpression(LOWEST)
	if err == nil && !p.currTokenIs(token.RPAREN) {
		return exp, fmt.Errorf("Error parsing grouped expression: missing right parenthesis")
	}
	p.advance()
	return exp, err
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	var err error
	exp := &ast.PrefixExpression{Token: p.curr}

	p.advance()
	exp.Right, err = p.parseExpression(PREFIX)
	if err != nil {
		return exp, fmt.Errorf("Error parsing infix expression: %v", err)
	}

	return exp, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	var err error

	exp := &ast.InfixExpression{Token: p.curr, Left: left}

	precedence := p.currTokenPrecedence()
	p.advance()
	exp.Right, err = p.parseExpression(precedence)
	if err != nil {
		return exp, fmt.Errorf("Error parsing infix expression: %v", err)
	}

	return exp, nil
}

func (p *Parser) parseIdentifier() (*ast.Identifier, error) {
	id := &ast.Identifier{Token: p.curr}
	p.advance()
	return id, nil
}

func (p *Parser) parseIntegerLiteral() (*ast.IntegerLiteral, error) {
	i := &ast.IntegerLiteral{Token: p.curr}
	p.advance()

	val, err := strconv.ParseInt(i.Token.Value, 10, 64)
	if err != nil {
		return i, fmt.Errorf("Error parsing integer literal: %v", err)
	}
	i.Value = val
	return i, nil
}

func (p *Parser) parseStringLiteral() (*ast.StringLiteral, error) {
	s := &ast.StringLiteral{Token: p.curr}
	p.advance()
	return s, nil
}

func (p *Parser) parseFuncDecl() (*ast.FuncDecl, error) {
	var err error
	f := &ast.FuncDecl{}
	p.advance()

	if err = p.assertCurrIs(token.IDENT); err != nil {
		return f, fmt.Errorf("Error parsing function declaration: %v", err)
	}

	f.Token = p.curr
	p.advance()

	if err = p.assertCurrIs(token.LPAREN); err != nil {
		return f, fmt.Errorf("Error parsing function declaration: %v", err)
	}
	p.advance()

	if err = p.assertCurrIs(token.RPAREN); err != nil {
		return f, fmt.Errorf("Error parsing function declaration: %v", err)
	}
	p.advance()

	if err = p.assertCurrIs(token.LBRACE); err != nil {
		return f, fmt.Errorf("Error parsing function declaration: %v", err)
	}
	p.advance()

	f.Body = make([]ast.Statement, 0)
	for !p.currTokenIs(token.RBRACE) {
		s, err := p.parseStatement()
		if err != nil {
			return f, fmt.Errorf("Error parsing function declaration: %v", err)
		}
		f.Body = append(f.Body, s)
	}

	if err = p.assertCurrIs(token.RBRACE); err != nil {
		return f, fmt.Errorf("Error parsing function declaration: %v", err)
	}
	p.advance()

	return f, nil
}

func (p *Parser) parseVariableDecl() (*ast.VariableDecl, error) {
	var err error
	s := &ast.VariableDecl{Name: p.curr}
	p.advance()

	if !p.currTokenIs(token.ASSIGN) {
		return s, fmt.Errorf("Error parsing variable declaration: expected '=', got: %v", p.next)
	}
	p.advance()

	s.Value, err = p.parseExpression(LOWEST)
	if err != nil {
		return s, fmt.Errorf("Error parsing variable declaration: %v", err)
	}

	if p.currTokenIs(token.SEMICOLON) {
		p.advance()
	}

	return s, nil
}

func (p *Parser) advance() {
	p.curr = p.next
	p.next = p.l.NextToken()
}

func (p *Parser) assertCurrIs(t token.TokenType) error {
	if p.curr.Type == t {
		return nil
	} else {
		return fmt.Errorf("expected %v, got %v", t, p.curr.Type)
	}
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.curr.Type == t
}

func (p *Parser) nextTokenIs(t token.TokenType) bool {
	return p.next.Type == t
}

func (p *Parser) currTokenPrecedence() int {
	switch p.curr.Type {
	case token.ASTERISK:
		fallthrough
	case token.SLASH:
		return PRODUCT
	case token.MINUS:
		fallthrough
	case token.PLUS:
		return SUM
	default:
		return LOWEST
	}
}
