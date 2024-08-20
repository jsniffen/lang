package parser

import (
	"fmt"
	"lang/ast"
	"lang/lexer"
	"lang/token"
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

func (p *Parser) ParseProgram() (*ast.Program, bool) {
	prog := &ast.Program{make([]ast.Statement, 0)}

	for !p.currTokenIs(token.EOF) {
		var stmt ast.Statement
		var ok bool

		switch p.curr.Type {
		case token.FUNC:
			stmt, ok = p.parseFuncDecl()
		default:
			p.Error(p.curr, "invalid token: '%s'", p.curr.Value)
			ok = false
		}

		if !ok {
			return prog, false
		}

		prog.Statements = append(prog.Statements, stmt)
	}

	return prog, true
}

func (p *Parser) parseFuncStatement() (ast.Statement, bool) {
	switch p.curr.Type {
	case token.VAR:
		return p.parseVarDecl()
	case token.RETURN:
		return p.parseReturn()
	case token.IDENT:
		switch p.next.Type {
		case token.LPAREN:
			return p.parseFuncCall()
		case token.ASSIGN:
			fallthrough
		case token.IDENT:
			fallthrough
		case token.ASTERISK:
			return p.parseVarDecl()
		default:
			p.Error(p.next, "invalid token: '%s'", p.next.Value)
			return nil, false
		}
	default:
		p.Error(p.curr, "invalid token: '%s'", p.curr.Value)
		return nil, false
	}
}

func (p *Parser) parseReturn() (*ast.Return, bool) {
	if !p.assertCurrIs(token.RETURN) {
		return nil, false
	}
	ret := &ast.Return{Token: p.curr}
	p.advance()

	e, ok := p.parseExpression(LOWEST)
	if !ok {
		return nil, false
	}
	ret.Value = e

	return ret, true
}

func (p *Parser) parseFuncDecl() (*ast.FuncDecl, bool) {
	var ok bool
	ok = p.assertCurrIs(token.FUNC)
	if !ok {
		return nil, false
	}
	p.advance()

	f := &ast.FuncDecl{}

	if !p.assertCurrIs(token.IDENT) {
		return nil, false
	}

	f.Token = p.curr
	p.advance()

	if !p.assertCurrIs(token.LPAREN) {
		return nil, false
	}
	p.advance()

	f.Params = make([]*ast.VarDecl, 0)
	for !p.currTokenIs(token.RPAREN) {
		param, ok := p.parseVarDecl()
		if !ok {
			return nil, false
		}
		f.Params = append(f.Params, param)

		if p.currTokenIs(token.COMMA) {
			p.advance()
		}
	}

	if !p.assertCurrIs(token.RPAREN) {
		return nil, false
	}
	p.advance()

	if p.currTokenIs(token.IDENT) {
		t, ok := p.parseType()
		fmt.Println(t)
		if !ok {
			return nil, false
		}
		f.ReturnType = t
	} else {
		f.ReturnType = &ast.Type{Type: ast.Void}
	}

	if !p.assertCurrIs(token.LBRACE) {
		return nil, false
	}
	p.advance()

	f.Body = make([]ast.Statement, 0)
	for !p.currTokenIs(token.RBRACE) {
		s, ok := p.parseFuncStatement()
		if !ok {
			return nil, false
		}
		f.Body = append(f.Body, s)
	}

	if !p.assertCurrIs(token.RBRACE) {
		return nil, false
	}
	p.advance()

	return f, true
}

func (p *Parser) parseFuncCall() (*ast.FuncCall, bool) {
	fc := &ast.FuncCall{Token: p.curr, Args: make([]ast.Expression, 0)}

	if !p.assertCurrIs(token.IDENT) {
		return nil, false
	}
	p.advance()

	if !p.assertCurrIs(token.LPAREN) {
		return nil, false
	}
	p.advance()

	for !p.currTokenIs(token.RPAREN) {
		exp, ok := p.parseExpression(LOWEST)
		if !ok {
			return nil, false
		}
		fc.Args = append(fc.Args, exp)

		if p.currTokenIs(token.COMMA) {
			p.advance()
		}
	}

	if !p.assertCurrIs(token.RPAREN) {
		return nil, false
	}
	p.advance()

	return fc, true
}

func (p *Parser) parseType() (*ast.Type, bool) {
	t := &ast.Type{}

	if p.currTokenIs(token.ASTERISK) {
		t.Pointer = true
		p.advance()
	}

	if !p.assertCurrIs(token.IDENT) {
		return nil, false
	}
	t.Type = p.curr.Value
	p.advance()

	return t, true
}

func (p *Parser) advance() {
	p.curr = p.next
	p.next = p.l.NextToken()
}

func (p *Parser) assertCurrIs(t token.TokenType) bool {
	if p.curr.Type == t {
		return true
	} else {
		p.Error(p.curr, fmt.Sprintf("expected %v, got %v", t, p.curr.Type))
		return false
	}
}

func (p *Parser) assertSameType(a, b ast.Type) bool {
	if a.Type == b.Type {
		return true
	} else {
		p.Error(p.curr, fmt.Sprintf("TypeError: %v != %v", a, b))
		return false
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

func (p *Parser) PrintErrors() {
	for _, err := range p.errors {
		fmt.Println(err)
	}
}

func (p *Parser) Error(t token.Token, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	err := fmt.Sprintf("%s:%d:%d: %s", p.l.Filename, t.Line, t.Column, msg)
	p.errors = append(p.errors, err)
}
