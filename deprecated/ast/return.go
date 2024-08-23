package ast

import (
	"bytes"
	"lang/token"
)

type Return struct {
	Token token.Token
	Value Expression
}

func (r *Return) CodeGen() string {
	var out bytes.Buffer
	out.WriteString("ret ")
	out.WriteString(r.Value.GetResult().Type.Name)
	return out.String()
}

func (r *Return) isStatement() {}

func (r *Return) String() string {
	var out bytes.Buffer
	out.WriteString(r.Token.Value)
	out.WriteString(" ")
	out.WriteString(r.Value.String())
	return out.String()
}

func (r *Return) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString(r.String())
	out.WriteString(" ")
	out.WriteString(r.Value.DebugString(i + 1))
	return out.String()
}
