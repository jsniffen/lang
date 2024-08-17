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
	curr   token.Token
	next   token.Token
	l      *lexer.Lexer
	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: make([]string, 0)}
	p.advance()
	p.advance()
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{make([]ast.Statement, 0)}

	for !p.currTokenIs(token.EOF) {
		s := p.parseStatement()
		if s == nil {
			break
		}
		prog.Statements = append(prog.Statements, s)
	}

	return prog
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}

func (p *Parser) PrintErrors() {
	for _, err := range p.errors {
		fmt.Println(err)
	}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curr.Type {
	case token.IDENT:
		return p.parseVariableDecl()
	case token.FUNC:
		return p.parseFuncDecl()
	default:
		p.Error(p.curr, "invalid token: '%s'", p.curr.Value)
		return nil
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	var left ast.Expression

	switch p.curr.Type {
	case token.LPAREN:
		left = p.parseGroupedExpression()
	case token.MINUS:
		left = p.parsePrefixExpression()
	case token.INT:
		left = p.parseIntegerLiteral()
	case token.IDENT:
		left = p.parseIdentifier()
	case token.STRING:
		left = p.parseStringLiteral()
	default:
		p.Error(p.curr, "invalid token: '%s'", p.curr.Value)
	}

	if left == nil {
		return left
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
			left = p.parseInfixExpression(left)
		default:
			return left
		}

		if left == nil {
			return left
		}
	}

	return left
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advance()
	exp := p.parseExpression(LOWEST)
	if ok, msg := p.assertCurrIs(token.LPAREN); !ok {
		p.Error(p.curr, msg)
		return exp
	}
	p.advance()
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{Token: p.curr}
	p.advance()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{Token: p.curr, Left: left}
	precedence := p.currTokenPrecedence()
	p.advance()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	id := &ast.Identifier{Token: p.curr}
	p.advance()
	return id
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	i := &ast.IntegerLiteral{Token: p.curr}
	p.advance()
	val, err := strconv.ParseInt(i.Token.Value, 10, 64)
	if err != nil {
		p.Error(p.curr, "Error parsing integer literal: %v", err)
		return i
	}
	i.Value = val
	return i
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	s := &ast.StringLiteral{Token: p.curr}
	p.advance()
	return s
}

func (p *Parser) parseFuncDecl() *ast.FuncDecl {
	f := &ast.FuncDecl{}
	p.advance()

	if ok, msg := p.assertCurrIs(token.IDENT); !ok {
		p.Error(p.curr, msg)
		return nil
	}

	f.Token = p.curr
	p.advance()

	if ok, msg := p.assertCurrIs(token.LPAREN); !ok {
		p.Error(p.curr, msg)
		return nil
	}
	p.advance()

	if ok, msg := p.assertCurrIs(token.RPAREN); !ok {
		p.Error(p.curr, msg)
		return nil
	}
	p.advance()

	if ok, msg := p.assertCurrIs(token.LBRACE); !ok {
		p.Error(p.curr, msg)
		return nil
	}
	p.advance()

	f.Body = make([]ast.Statement, 0)
	for !p.currTokenIs(token.RBRACE) {
		s := p.parseStatement()
		if s == nil {
			return nil
		}
		f.Body = append(f.Body, s)
	}

	if ok, msg := p.assertCurrIs(token.RBRACE); !ok {
		p.Error(p.curr, msg)
		return nil
	}
	p.advance()

	return f
}

func (p *Parser) parseVariableDecl() *ast.VariableDecl {
	s := &ast.VariableDecl{Name: p.curr}
	p.advance()

	if !p.currTokenIs(token.ASSIGN) {
		p.Error(p.curr, "expected '='")
		return nil
	}
	p.advance()

	s.Value = p.parseExpression(LOWEST)

	return s
}

func (p *Parser) advance() {
	p.curr = p.next
	p.next = p.l.NextToken()
}

func (p *Parser) assertCurrIs(t token.TokenType) (bool, string) {
	if p.curr.Type == t {
		return true, ""
	} else {
		return false, fmt.Sprintf("expected %v, got %v", t, p.curr.Type)
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

func (p *Parser) Error(t token.Token, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	err := fmt.Sprintf("%s:%d:%d: %s\n", p.l.Filename, t.Line, t.Column, msg)
	p.errors = append(p.errors, err)
}
