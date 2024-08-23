package parser

import (
	"fmt"
	"lang/ast"
	"lang/lexer"
	"lang/token"
	"lang/types"
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
			p.error(p.curr, "invalid token: '%s'", p.curr.Value)
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
	if !p.assertCurrIs(token.FUNC) {
		return nil, false
	}
	p.advance()

	name := p.curr.Value
	p.advance()

	if !p.assertCurrIs(token.LPAREN) {
		return nil, false
	}
	p.advance()

	args := make([]*ast.FuncArg, 0)
	for !p.currIsOrEOF(token.RPAREN) {
		fa, ok := p.parseFuncArg()
		if !ok {
			return nil, false
		}

		args = append(args, fa)

		if p.currIs(token.COMMA) {
			p.advance()
		}
	}

	if !p.assertCurrIs(token.RPAREN) {
		return nil, false
	}
	p.advance()

	returnType := &ast.Type{Name: types.Void}
	if p.currIs(token.IDENT) {
		t, ok := p.parseType()
		if !ok {
			return nil, false
		}
		returnType = t
	}

	extern := true

	if p.currIs(token.LBRACE) {
		extern = false
		p.advance()

		if !p.assertCurrIs(token.RBRACE) {
			return nil, false
		}
		p.advance()
	}

	fd := &ast.FuncDecl{
		Args:       args,
		Extern:     extern,
		Name:       name,
		ReturnType: returnType,
	}

	return fd, true
}

func (p *Parser) parseFuncArg() (*ast.FuncArg, bool) {
	if !p.assertCurrIs(token.IDENT) {
		return nil, false
	}
	name := p.curr.Value
	p.advance()

	t, ok := p.parseType()
	if !ok {
		return nil, false
	}

	fa := &ast.FuncArg{
		Name:     name,
		Type:     t,
		Location: "%" + name,
	}

	return fa, true
}

func (p *Parser) parseType() (*ast.Type, bool) {
	if !p.assertCurrIs(token.IDENT) {
		return nil, false
	}
	name := p.curr.Value
	p.advance()

	return &ast.Type{Name: name}, true
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
