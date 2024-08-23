package parser

import (
	"lang/ast"
	"strconv"
)

func (p *Parser) parseIntegerLiteral() (*ast.IntegerLiteral, bool) {
	i := &ast.IntegerLiteral{Token: p.curr}
	p.advance()
	val, err := strconv.ParseInt(i.Token.Value, 10, 64)
	if err != nil {
		p.error(p.curr, "Error parsing integer literal: %v", err)
		return i, false
	}
	i.Value = val
	return i, true
}
