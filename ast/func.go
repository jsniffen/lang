package ast

import (
	"bytes"
	"lang/token"
)

type FuncDecl struct {
	Token      token.Token
	Params     []*VarDecl
	Body       []Statement
	ReturnType *Type
}

func (fd *FuncDecl) isStatement() {}

func (fd *FuncDecl) CodeGen() string {
	var out bytes.Buffer

	extern := len(fd.Body) == 0

	if extern {
		out.WriteString("declare ")
	} else {
		out.WriteString("define ")
	}
	out.WriteString(fd.ReturnType.CodeGen())
	out.WriteString(" ")
	out.WriteString("@")
	out.WriteString(fd.Token.Value)
	out.WriteString("(")
	out.WriteString(")")

	if !extern {
		out.WriteString(" {")
		for _, s := range fd.Body {
			out.WriteString("\n\t")
			out.WriteString(s.CodeGen())
		}
		out.WriteString("\n}")
	}

	return out.String()
}

func (fd *FuncDecl) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString("FuncDecl ")
	out.WriteString(fd.Token.Value)
	out.WriteString("(")
	for i, p := range fd.Params {
		out.WriteString(p.String())
		if i < len(fd.Params)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")
	if fd.ReturnType.Type != "" {
		out.WriteString(" -> ")
		out.WriteString(fd.ReturnType.Type)
	}
	for _, s := range fd.Body {
		out.WriteString(s.DebugString(i + 1))
	}
	return out.String()
}
func (fd *FuncDecl) String() string {
	var out bytes.Buffer
	out.WriteString("func ")
	out.WriteString(fd.Token.Value)
	out.WriteString("(")
	for i, p := range fd.Params {
		out.WriteString(p.String())
		if i < len(fd.Params)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")
	if fd.ReturnType.Type != "" {
		out.WriteString(" ")
		out.WriteString(fd.ReturnType.Type)
	}
	if len(fd.Body) > 0 {
		out.WriteString(" {")
		for _, s := range fd.Body {
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

func (f *FuncCall) CodeGen() string { return "" }

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
