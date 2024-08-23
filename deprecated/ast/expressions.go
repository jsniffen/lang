package ast

import (
	"bytes"
	"lang/token"
)

type InfixExpression struct {
	Token  token.Token
	Left   Expression
	Right  Expression
	Result Result
	Type   Type
}

func (i *InfixExpression) CodeGen() string   { return "" }
func (i *InfixExpression) GetResult() Result { return i.Result }
func (i *InfixExpression) isExpression()     {}
func (i *InfixExpression) DebugString(indent int) string {
	var out bytes.Buffer
	printIndentLine(indent, &out)
	out.WriteString("InfixExpression ")
	out.WriteString(i.Type.Name)
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
	Token  token.Token
	Right  Expression
	Type   Type
	Result Result
}

func (p *PrefixExpression) CodeGen() string   { return "" }
func (p *PrefixExpression) GetResult() Result { return p.Result }
func (p *PrefixExpression) isExpression()     {}
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
