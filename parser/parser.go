package parser

import (
	"fmt"
	"lang/ast"
	"lang/lexer"
	"lang/token"
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

	for p.curr.Type != token.EOF {
		var statement ast.Statement
		var err error
		if p.curr.Type == token.IDENT && p.next.Type == token.ASSIGN {
			statement, err = p.parseVariableDeclaration()
			if err != nil {
				return prog, err
			}
			prog.Statements = append(prog.Statements, statement)
		}

		p.advance()
	}

	return prog, nil
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	if p.curr.Type == token.INT {
		return p.parseIntegerLiteral()
	}

	return nil, fmt.Errorf("Error parsing expression: No valid token found")
}

func (p *Parser) parseIntegerLiteral() (*ast.IntegerLiteral, error) {
	node := &ast.IntegerLiteral{}
	node.Value = p.curr.Value
	p.advance()
	return node, nil
}

func (p *Parser) parseVariableDeclaration() (*ast.VariableDeclaration, error) {
	node := &ast.VariableDeclaration{}

	node.Identifier = p.curr.Value

	p.advance()
	p.advance()

	expression, err := p.parseExpression()
	if err != nil {
		return node, err
	}
	node.Expression = expression

	return node, nil
}

func (p *Parser) advance() {
	p.curr = p.next
	p.next = p.l.NextToken()
}
