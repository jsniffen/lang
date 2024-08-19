package ast

import (
	"bytes"
	"lang/token"
)

type FuncDecl struct {
	Token      token.Token
	Params     []*VarDecl
	Body       []Statement
	ReturnType token.Token
}

func (f *FuncDecl) isStatement() {}

func (f *FuncDecl) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString("FuncDecl ")
	out.WriteString(f.Token.Value)
	out.WriteString("(")
	for _, p := range f.Params {
		out.WriteString(p.String())
	}
	out.WriteString(")")
	if f.ReturnType.Value != "" {
		out.WriteString(" -> ")
		out.WriteString(f.ReturnType.Value)
	}
	for _, s := range f.Body {
		out.WriteString(s.DebugString(i + 1))
	}
	return out.String()
}
func (f *FuncDecl) String() string {
	var out bytes.Buffer
	out.WriteString("func ")
	out.WriteString(f.Token.Value)
	out.WriteString("(")
	for i, p := range f.Params {
		out.WriteString(p.String())
		if i < len(f.Params)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")
	if len(f.Body) > 0 {
		out.WriteString(" {")
		for _, s := range f.Body {
			out.WriteString("\n\t")
			out.WriteString(s.String())
		}
		out.WriteString("\n")
		out.WriteString("}")
	}
	return out.String()
}

type FuncCall struct {
	Token token.Token
	Args  []Expression
}

func (f *FuncCall) isExpression() {}
func (f *FuncCall) isStatement()  {}

func (f *FuncCall) CodeGen() string {
	var out bytes.Buffer
	out.WriteString("call ")
	out.WriteString(f.Token.Value)
	out.WriteString("(")
	out.WriteString(")")
	return out.String()
}

func (f *FuncCall) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString("FuncCall ")
	out.WriteString(f.Token.Value)
	for _, arg := range f.Args {
		out.WriteString(arg.DebugString(i + 1))
	}
	return out.String()
}

func (f *FuncCall) String() string {
	var out bytes.Buffer
	out.WriteString(f.Token.Value)
	out.WriteString("(")
	for i, arg := range f.Args {
		out.WriteString(arg.String())
		if i < len(f.Args)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")
	return out.String()
}
