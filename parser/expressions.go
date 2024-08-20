package parser

import (
	"lang/ast"
	"lang/token"
)

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
	if !p.assertCurrIs(token.RPAREN) {
		return nil, false
	}
	p.advance()

	exp, ok := p.parseExpression(LOWEST)
	if !ok {
		return nil, false
	}

	if !p.assertCurrIs(token.LPAREN) {
		return nil, false
	}
	p.advance()

	return exp, true
}

func (p *Parser) parsePrefixExpression() (ast.Expression, bool) {
	exp := &ast.PrefixExpression{Token: p.curr}
	p.advance()

	var ok bool
	right, ok := p.parseExpression(PREFIX)
	if !ok {
		return nil, false
	}
	exp.Right = right
	exp.Type = right.ReturnType()

	return exp, true
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, bool) {
	exp := &ast.InfixExpression{Token: p.curr, Left: left}

	precedence := p.currTokenPrecedence()
	p.advance()

	right, ok := p.parseExpression(precedence)
	if !ok {
		return nil, false
	}
	exp.Right = right

	ok = p.assertSameType(left.ReturnType(), right.ReturnType())
	if !ok {
		return nil, false
	}

	exp.Type = left.ReturnType()

	return exp, true
}
