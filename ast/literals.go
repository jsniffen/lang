package ast

import (
	"bytes"
	"lang/token"
)

type IntegerLiteral struct {
	Token token.Token
	Value int64
	Type  Type
}

func (i *IntegerLiteral) CodeGen() string  { return i.Token.Value }
func (i *IntegerLiteral) ReturnType() Type { return i.Type }
func (i *IntegerLiteral) isExpression()    {}
func (i *IntegerLiteral) DebugString(indent int) string {
	var out bytes.Buffer
	printIndentLine(indent, &out)
	out.WriteString("IntegerLiteral ")
	out.WriteString(i.Token.Value)
	out.WriteString(" ")
	out.WriteString(i.Type.Type)
	return out.String()
}
func (i *IntegerLiteral) String() string { return i.Token.Value }
