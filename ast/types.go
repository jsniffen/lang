package ast

import (
	"bytes"
)

const (
	Int32 = "i32"
	Void  = "void"
)

type Type struct {
	Pointer bool
	Type    string
}

func (t *Type) CodeGen() string {
	var out bytes.Buffer
	if t.Pointer {
		out.WriteString("*")
	}
	out.WriteString(t.Type)
	return out.String()
}
func (t *Type) DebugString(int) string {
	return "Type"
}
func (t *Type) String() string {
	var out bytes.Buffer
	if t.Pointer {
		out.WriteString("*")
	}
	out.WriteString(t.Type)
	return out.String()
}
