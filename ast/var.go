package ast

import (
	"bytes"
	"lang/token"
)

type Var struct {
	Token token.Token
	Name  string
	Type  string
}

func (v *Var) isExpression()  {}
func (v *Var) String() string { return v.Name }
func (v *Var) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString("Var ")
	out.WriteString(v.Name)
	if v.Type != "" {
		out.WriteString(" ")
		out.WriteString(v.Type)
	}
	return out.String()
}

type VarDecl struct {
	Name  token.Token
	Type  token.Token
	Value Expression
}

func (vd *VarDecl) isStatement()   {}
func (vd *VarDecl) String() string { return vd.Name.Value }
func (vd *VarDecl) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString("VarDecl ")
	out.WriteString(vd.String())
	out.WriteString(" = ")
	if vd.Type.Value != "" {
		out.WriteString(" ")
		out.WriteString(vd.Type.Value)
	}
	out.WriteString(vd.Value.DebugString(i + 1))
	return out.String()
}
