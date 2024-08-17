package ast

import (
	"bytes"
	"lang/token"
)

type Node interface {
	DebugString(int) string
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

func (p *Program) DebugString(i int) string {
	var out bytes.Buffer
	out.WriteString("Program")
	for _, s := range p.Statements {
		out.WriteString(s.DebugString(i + 1))
	}
	return out.String()
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

func (v *VariableDeclaration) isStatement() {}
func (v *VariableDeclaration) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString("VariableDeclation ")
	out.WriteString(v.Name.Value)
	out.WriteString(" =")
	out.WriteString(v.Value.DebugString(i + 1))
	return out.String()
}
func (v *VariableDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString(v.Name.Value)
	out.WriteString(" = ")
	out.WriteString(v.Value.String())
	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) isExpression() {}
func (i *IntegerLiteral) DebugString(indent int) string {
	var out bytes.Buffer
	printIndentLine(indent, &out)
	out.WriteString("IntegerLiteral(")
	out.WriteString(i.Token.Value)
	out.WriteString(")")
	return out.String()
}
func (i *IntegerLiteral) String() string { return i.Token.Value }

type InfixExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (i *InfixExpression) isExpression() {}
func (i *InfixExpression) DebugString(indent int) string {
	var out bytes.Buffer
	printIndentLine(indent, &out)
	out.WriteString("InfixExpression")
	out.WriteString(i.Left.DebugString(indent + 1))
	printIndentLine(indent+1, &out)
	out.WriteString(i.Token.Value)
	out.WriteString(i.Right.DebugString(indent + 1))
	printIndentLine(indent, &out)
	return out.String()
}
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

func (p *PrefixExpression) isExpression() {}
func (p *PrefixExpression) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString("PrefixExpression")
	printIndentLine(i, &out)
	out.WriteString(p.Token.Value)
	out.WriteString(p.Right.DebugString(i + 1))
	return out.String()
}
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

func (id *Identifier) isExpression() {}
func (id *Identifier) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString("Identifier(")
	out.WriteString(id.Token.Value)
	out.WriteString(")")
	return out.String()
}
func (i *Identifier) String() string { return i.Token.Value }

func printIndentLine(i int, b *bytes.Buffer) {
	b.WriteString("\n")
	for j := 0; j < i; j += 1 {
		b.WriteString("\t")
	}
}