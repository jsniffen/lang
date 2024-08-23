package parser

import (
	"lang/ast"
	"lang/token"
)

func (p *Parser) parseVar() (*ast.Var, bool) {
	id := &ast.Var{Name: p.curr}
	p.advance()
	return id, true
}

func (p *Parser) parseVarDecl() (*ast.VarDecl, bool) {
	if !p.assertCurrIs(token.VAR) {
		return nil, false
	}
	p.advance()

	if !p.assertCurrIs(token.IDENT) {
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
