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

func (p *Parser) ParseProgram() (*ast.Program, bool) {
	prog := &ast.Program{make([]ast.Statement, 0)}

	for !p.currTokenIs(token.EOF) {
		s, ok := p.parseProgramStatement()
		if !ok {
			return prog, false
		}
		prog.Statements = append(prog.Statements, s)
	}

	return prog, true
}

func (p *Parser) PrintErrors() {
	for _, err := range p.errors {
		fmt.Println(err)
	}
}

func (p *Parser) parseProgramStatement() (ast.Statement, bool) {
	switch p.curr.Type {
	case token.FUNC:
		return p.parseFuncDecl()
	case token.VAR:
		return p.parseVarDecl()
	default:
		p.Error(p.curr, "invalid token: '%s'", p.curr.Value)
		return nil, false
	}
}

func (p *Parser) parseFuncStatement() (ast.Statement, bool) {
	switch p.curr.Type {
	case token.VAR:
		return p.parseVarDecl()
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

func (p *Parser) parseExpression(precedence int) (ast.Expression, bool) {
	var left ast.Expression
	var ok bool

	switch p.curr.Type {
	case token.LPAREN:
		left, ok = p.parseGroupedExpression()
	case token.MINUS:
		left, ok = p.parsePrefixExpression()
	case token.INT:
		left, ok = p.parseIntegerLiteral()
	case token.IDENT:
		left, ok = p.parseVar()
	case token.STRING:
		left, ok = p.parseStringLiteral()
	default:
		p.Error(p.curr, "invalid token: '%s'", p.curr.Value)
		return left, false
	}

	if !ok {
		return left, false
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
			left, ok = p.parseInfixExpression(left)
		default:
			return left, true
		}

		if !ok {
			return left, false
		}
	}

	return left, true
}

func (p *Parser) parseGroupedExpression() (ast.Expression, bool) {
	p.advance()
	exp, ok := p.parseExpression(LOWEST)
	if !ok {
		return nil, false
	}

	if ok, msg := p.assertCurrIs(token.LPAREN); !ok {
		p.Error(p.curr, msg)
		return nil, false
	}
	p.advance()
	return exp, true
}

func (p *Parser) parsePrefixExpression() (ast.Expression, bool) {
	var ok bool
	exp := &ast.PrefixExpression{Token: p.curr}
	p.advance()
	exp.Right, ok = p.parseExpression(PREFIX)
	if !ok {
		return nil, false
	}
	return exp, true
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, bool) {
	var ok bool
	exp := &ast.InfixExpression{Token: p.curr, Left: left}
	precedence := p.currTokenPrecedence()
	p.advance()
	exp.Right, ok = p.parseExpression(precedence)
	if !ok {
		return nil, ok
	}
	return exp, true
}

func (p *Parser) parseVar() (*ast.Var, bool) {
	id := &ast.Var{p.curr}
	p.advance()
	return id, true
}

func (p *Parser) parseIntegerLiteral() (*ast.IntegerLiteral, bool) {
	i := &ast.IntegerLiteral{Token: p.curr}
	p.advance()
	val, err := strconv.ParseInt(i.Token.Value, 10, 64)
	if err != nil {
		p.Error(p.curr, "Error parsing integer literal: %v", err)
		return i, false
	}
	i.Value = val
	return i, true
}

func (p *Parser) parseStringLiteral() (*ast.StringLiteral, bool) {
	s := &ast.StringLiteral{Token: p.curr}
	p.advance()
	return s, true
}

func (p *Parser) parseFuncDecl() (*ast.FuncDecl, bool) {
	var ok bool
	ok = p.checkCurrIs(token.FUNC)
	if !ok {
		return nil, false
	}
	p.advance()

	f := &ast.FuncDecl{}

	if ok, msg := p.assertCurrIs(token.IDENT); !ok {
		p.Error(p.curr, msg)
		return nil, false
	}

	f.Token = p.curr
	p.advance()

	if ok, msg := p.assertCurrIs(token.LPAREN); !ok {
		p.Error(p.curr, msg)
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

	if ok, msg := p.assertCurrIs(token.RPAREN); !ok {
		p.Error(p.curr, msg)
		return nil, false
	}
	p.advance()

	if p.currTokenIs(token.IDENT) {
		f.ReturnType = p.curr
		p.advance()
	}

	if !p.currTokenIs(token.LBRACE) {
		return f, true
	}

	if ok, msg := p.assertCurrIs(token.LBRACE); !ok {
		p.Error(p.curr, msg)
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

	if ok, msg := p.assertCurrIs(token.RBRACE); !ok {
		p.Error(p.curr, msg)
		return nil, false
	}
	p.advance()

	return f, true
}

func (p *Parser) parseFuncCall() (*ast.FuncCall, bool) {
	fc := &ast.FuncCall{Token: p.curr, Args: make([]ast.Expression, 0)}

	if ok, msg := p.assertCurrIs(token.IDENT); !ok {
		p.Error(p.curr, msg)
		return nil, false
	}
	p.advance()

	if ok, msg := p.assertCurrIs(token.LPAREN); !ok {
		p.Error(p.curr, msg)
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

	if ok, msg := p.assertCurrIs(token.RPAREN); !ok {
		p.Error(p.curr, msg)
		return nil, false
	}
	p.advance()

	return fc, true
}

func (p *Parser) parseVarDecl() (*ast.VarDecl, bool) {
	if p.currTokenIs(token.VAR) {
		p.advance()
	}

	if !p.checkCurrIs(token.IDENT) {
		return nil, false
	}

	vd := &ast.VarDecl{Name: p.curr}
	p.advance()

	if p.currTokenIs(token.ASTERISK) {
		vd.Pointer = true
		p.advance()
	}

	if p.currTokenIs(token.IDENT) {
		vd.Type = p.curr
		p.advance()
	}

	if p.currTokenIs(token.ASSIGN) {
		p.advance()

		exp, ok := p.parseExpression(LOWEST)
		if !ok {
			return nil, false
		}
		vd.Value = exp
	}

	return vd, true
}

func (p *Parser) advance() {
	p.curr = p.next
	p.next = p.l.NextToken()
}

func (p *Parser) checkCurrIs(t token.TokenType) bool {
	if p.curr.Type == t {
		return true
	} else {
		p.Error(p.curr, fmt.Sprintf("expected %v, got %v", t, p.curr.Type))
		return false
	}
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
	err := fmt.Sprintf("%s:%d:%d: %s", p.l.Filename, t.Line, t.Column, msg)
	p.errors = append(p.errors, err)
}
