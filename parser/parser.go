package parser

import (
	"fmt"
	"lang/ast"
	"lang/lexer"
	"lang/token"
	"lang/types"
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
	l        *lexer.Lexer
	curr     token.Token
	next     token.Token
	register int
	Errors   []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		Errors: make([]string, 0),
	}
	p.advance()
	p.advance()
	return p
}

func (p *Parser) ParseProgram() (*ast.Program, bool) {
	prog := &ast.Program{make([]ast.Statement, 0)}

	for !p.currIs(token.EOF) {
		var stmt ast.Statement
		var ok bool

		switch p.curr.Type {
		case token.FUNC:
			stmt, ok = p.parseFuncDecl()
		case token.VAR:
			stmt, ok = p.parseVarDecl()
		default:
			p.errorInvalidToken()
			ok = false
		}

		if !ok {
			return prog, false
		}

		prog.Statements = append(prog.Statements, stmt)
	}

	return prog, true
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, bool) {
	var left ast.Expression
	var ok bool

	switch p.curr.Type {
	case token.INT:
		left, ok = p.parseIntLiteral()
	case token.IDENT:
		switch p.next.Type {
		case token.LPAREN:
			left, ok = p.parseFuncCall()
		default:
			left, ok = p.parseVar()
		}
	default:
		p.errorInvalidToken()
		ok = false
	}

	if !ok {
		return left, false
	}

	for precedence < p.currPrecedence() {
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

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, bool) {
	exp := &ast.InfixExpression{Token: p.curr, Left: left}

	precedence := p.currPrecedence()
	p.advance()

	right, ok := p.parseExpression(precedence)
	if !ok {
		return nil, false
	}
	exp.Right = right

	return exp, true
}

func (p *Parser) parseFuncDecl() (*ast.FuncDecl, bool) {
	fd := &ast.FuncDecl{
		Extern: true,
		Params: make([]*ast.VarDecl, 0),
	}
	var ok bool

	if !p.assertCurrIs(token.FUNC) {
		return nil, false
	}
	p.advance()

	fd.Token = p.curr
	p.advance()

	if !p.assertCurrIs(token.LPAREN) {
		return nil, false
	}
	p.advance()

	for !p.currIs(token.RPAREN) {
		vd, ok := p.parseFuncParam()
		if !ok {
			return nil, false
		}
		fd.Params = append(fd.Params, vd)

		if p.currIs(token.COMMA) {
			p.advance()
		}
	}
	p.advance()

	if p.currIs(token.IDENT) {
		fd.ReturnType, ok = p.parseType()
		if !ok {
			return nil, false
		}
		fd.HasReturn = true
	}

	if p.currIs(token.LBRACE) {
		fd.Extern = false
		p.advance()

		fd.Body, ok = p.parseFuncBody()
		if !ok {
			return nil, false
		}
		p.advance()
	}

	return fd, true
}

func (p *Parser) parseFuncBody() ([]ast.Statement, bool) {
	body := make([]ast.Statement, 0)

	for !p.currIsOrEOF(token.RBRACE) {
		var stmt ast.Statement
		var ok bool

		switch p.curr.Type {
		case token.IDENT:
			switch p.next.Type {
			case token.LPAREN:
				stmt, ok = p.parseFuncCall()
			default:
				p.errorInvalidToken()
				ok = false
			}
		case token.RETURN:
			stmt, ok = p.parseReturn()
		case token.VAR:
			stmt, ok = p.parseVarDecl()
		default:
			p.errorInvalidToken()
			ok = false
		}

		if !ok {
			return body, false
		}

		body = append(body, stmt)
	}

	return body, true
}

func (p *Parser) parseFuncCall() (*ast.FuncCall, bool) {
	if !p.assertCurrIs(token.IDENT) {
		return nil, false
	}
	fc := &ast.FuncCall{Token: p.curr, Args: make([]ast.Expression, 0)}
	p.advance()

	if !p.assertCurrIs(token.LPAREN) {
		return nil, false
	}
	p.advance()

	for !p.currIsOrEOF(token.RPAREN) {
		e, ok := p.parseExpression(LOWEST)
		if !ok {
			return nil, false
		}
		fc.Args = append(fc.Args, e)

		if p.currIs(token.COMMA) {
			p.advance()
		}
	}
	p.advance()

	return fc, true
}

func (p *Parser) parseVarDecl() (*ast.VarDecl, bool) {
	if !p.assertCurrIs(token.VAR) {
		return nil, false
	}
	p.advance()

	if !p.assertCurrIs(token.IDENT) {
		return nil, false
	}
	vd := &ast.VarDecl{Token: p.curr}
	p.advance()

	t, ok := p.parseType()
	if !ok {
		return nil, false
	}
	vd.Type = t

	if !p.assertCurrIs(token.ASSIGN) {
		return nil, false
	}
	p.advance()

	e, ok := p.parseExpression(LOWEST)
	if !ok {
		return nil, false
	}
	vd.Value = e

	return vd, true
}

func (p *Parser) parseVar() (*ast.Var, bool) {
	if !p.assertCurrIs(token.IDENT) {
		return nil, false
	}
	v := &ast.Var{Token: p.curr}
	p.advance()

	return v, true
}

func (p *Parser) parseReturn() (*ast.Return, bool) {
	if !p.assertCurrIs(token.RETURN) {
		return nil, false
	}
	r := &ast.Return{Token: p.curr}
	p.advance()

	if !p.currIs("}") {
		e, ok := p.parseExpression(LOWEST)
		if !ok {
			return nil, false
		}
		r.Value = e
		r.HasValue = true
	}

	return r, true
}

func (p *Parser) parseFuncParam() (*ast.VarDecl, bool) {
	vd := &ast.VarDecl{Value: &ast.EmptyExpression{}}

	if !p.assertCurrIs(token.IDENT) {
		return nil, false
	}
	vd.Token = p.curr
	p.advance()

	t, ok := p.parseType()
	if !ok {
		return nil, false
	}
	vd.Type = t

	return vd, true
}

func (p *Parser) parseType() (*ast.Type, bool) {
	if !p.assertCurrIs(token.IDENT) {
		return nil, false
	}
	t := &ast.Type{Token: p.curr, Type: types.FromToken(p.curr)}
	p.advance()
	return t, true
}

func (p *Parser) parseIntLiteral() (*ast.IntLiteral, bool) {
	if !p.assertCurrIs(token.INT) {
		return nil, false
	}
	il := &ast.IntLiteral{Token: p.curr}

	n, err := strconv.Atoi(p.curr.Value)
	if err != nil {
		p.errorParse(err)
		return nil, false
	}
	p.advance()
	il.Value = n

	return il, true
}

func (p *Parser) advance() {
	p.curr = p.next
	p.next = p.l.NextToken()
}

func (p *Parser) assertCurrIs(t token.TokenType) bool {
	if !p.currIs(t) {
		p.error(p.curr, fmt.Sprintf("expected %v, got %v", t, p.curr.Type))
		return false
	}
	return true
}

func (p *Parser) currIsOrEOF(t token.TokenType) bool {
	return t == p.curr.Type || p.curr.Type == token.EOF
}

func (p *Parser) currIs(t token.TokenType) bool {
	return t == p.curr.Type
}

func (p *Parser) currPrecedence() int {
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

func (p *Parser) error(t token.Token, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	err := fmt.Sprintf("%s ParseError: %s", t.Path(), msg)
	p.Errors = append(p.Errors, err)
}

func (p *Parser) errorInvalidToken() {
	p.error(p.curr, "invalid token: '%s'", p.curr.Value)
}

func (p *Parser) errorParse(err error) {
	p.error(p.curr, "%s", err)
}
