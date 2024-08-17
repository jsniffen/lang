package ast

import (
	"bytes"
	"lang/token"
)

type Var struct {
	Name token.Token
}

func (v *Var) isExpression()  {}
func (v *Var) String() string { return v.Name.Value }
func (v *Var) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString("Var ")
	out.WriteString(v.String())
	return out.String()
}

type VarDecl struct {
	Name    token.Token
	Type    token.Token
	Pointer bool
	Value   Expression
}

func (vd *VarDecl) isStatement() {}
func (vd *VarDecl) String() string {
	var out bytes.Buffer
	out.WriteString(vd.Name.Value)
	if vd.Type.Value != "" {
		out.WriteString(" ")
		if vd.Pointer {
			out.WriteString("*")
		}
		out.WriteString(vd.Type.Value)
	}
	return out.String()
}
func (vd *VarDecl) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString("VarDecl ")
	out.WriteString(vd.String())
	if vd.Type.Value != "" {
		out.WriteString(" ")
		if vd.Pointer {
			out.WriteString("*")
		}
		out.WriteString(vd.Type.Value)
	}
	if vd.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vd.Value.DebugString(i + 1))
	}
	return out.String()
}
