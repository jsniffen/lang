package ast

import (
	"bytes"
	"lang/token"
)

type Node interface {
	TokenValue() string
	String() string
}

type Expression interface {
	Node
	isExpression()
}

type Statement interface {
	Node
	isStatement()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenValue() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenValue()
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	return out.String()
}

type VariableDeclaration struct {
	Name  token.Token
	Value Expression
}

func (v *VariableDeclaration) isStatement()       {}
func (v *VariableDeclaration) TokenValue() string { return v.Name.Value }
func (v *VariableDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString(v.TokenValue())
	out.WriteString(" = ")
	out.WriteString(v.Value.String())
	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) isExpression()      {}
func (i *IntegerLiteral) TokenValue() string { return i.Token.Value }
func (i *IntegerLiteral) String() string     { return i.Token.Value }

type InfixExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (i *InfixExpression) isExpression()      {}
func (i *InfixExpression) TokenValue() string { return i.Token.Value }
func (i *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Token.Value + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

type PrefixExpression struct {
	Token token.Token
	Right Expression
}

func (p *PrefixExpression) isExpression()      {}
func (p *PrefixExpression) TokenValue() string { return p.Token.Value }
func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(" " + p.Token.Value + " ")
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

type Identifier struct {
	Token token.Token
}

func (i *Identifier) isExpression()      {}
func (i *Identifier) TokenValue() string { return i.Token.Value }
func (i *Identifier) String() string     { return i.Token.Value }
