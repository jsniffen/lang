package parser

import (
	"fmt"
	"lang/ast"
	"lang/lexer"
	"lang/token"
	"lang/types"
	"strconv"
)

type Parser struct {
	l      *lexer.Lexer
	curr   token.Token
	next   token.Token
	Errors []string
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

func (p *Parser) parseFuncDecl() (*ast.FuncDecl, bool) {
	fd := &ast.FuncDecl{
		Extern:     true,
		ReturnType: &ast.Type{Name: types.Void},
	}
	var ok bool

	if !p.assertCurrIs(token.FUNC) {
		return nil, false
	}
	p.advance()

	fd.Name = p.curr.Value
	p.advance()

	if !p.assertCurrIs(token.LPAREN) {
		return nil, false
	}
	p.advance()

	if !p.currIs(token.RPAREN) {
		fd.Args, ok = p.parseFuncArgs()
		if !ok {
			return nil, false
		}
	}
	p.advance()

	if p.currIs(token.IDENT) {
		fd.ReturnType, ok = p.parseType()
		if !ok {
			return nil, false
		}
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
		case token.RETURN:
			stmt, ok = p.parseReturn()
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

func (p *Parser) parseReturn() (*ast.Return, bool) {
	if !p.assertCurrIs(token.RETURN) {
		return nil, false
	}
	r := &ast.Return{Token: p.curr}
	p.advance()

	if !p.currIs("}") {
		e, ok := p.parseExpression()
		if !ok {
			return nil, false
		}
		r.Value = e
		r.HasValue = true
	}

	return r, true
}

func (p *Parser) parseExpression() (ast.Expression, bool) {
	var left ast.Expression
	var ok bool

	switch p.curr.Type {
	case token.INT:
		left, ok = p.parseIntLiteral()
	default:
		p.errorInvalidToken()
		ok = false
	}

	if !ok {
		return left, false
	}

	return left, true
}

func (p *Parser) parseFuncArgs() ([]*ast.FuncArg, bool) {
	args := make([]*ast.FuncArg, 0)

	for !p.currIsOrEOF(token.RPAREN) {
		fa := &ast.FuncArg{}
		var ok bool

		if !p.assertCurrIs(token.IDENT) {
			return nil, false
		}
		fa.Name = p.curr.Value
		fa.Location = "%" + p.curr.Value
		p.advance()

		fa.Type, ok = p.parseType()
		if !ok {
			return nil, false
		}

		if p.currIs(token.COMMA) {
			p.advance()
		}

		args = append(args, fa)
	}

	return args, true
}

func (p *Parser) parseType() (*ast.Type, bool) {
	if !p.assertCurrIs(token.IDENT) {
		return nil, false
	}
	name := p.curr.Value
	p.advance()
	return &ast.Type{Name: name}, true
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

func (p *Parser) error(t token.Token, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	err := fmt.Sprintf("%s:%d:%d: %s", p.l.Filename, t.Line, t.Column, msg)
	p.Errors = append(p.Errors, err)
}

func (p *Parser) errorInvalidToken() {
	p.error(p.curr, "invalid token: '%s'", p.curr.Value)
}

func (p *Parser) errorParse(err error) {
	p.error(p.curr, "%s", err)
}
