package ast

import (
	"bytes"
	"lang/token"
)

type FuncParam struct {
	Name token.Token
	Type token.Token
}

func (fp *FuncParam) isStatement() {}
func (fp *FuncParam) String() string {
	var out bytes.Buffer
	out.WriteString(fp.Name.Value)
	out.WriteString(" ")
	out.WriteString(fp.Type.Value)
	return out.String()
}
func (fp *FuncParam) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString(fp.String())
	return out.String()
}

type FuncDecl struct {
	Token  token.Token
	Params []*FuncParam
	Body   []Statement
}

func (f *FuncDecl) isStatement() {}
func (f *FuncDecl) DebugString(i int) string {
	var out bytes.Buffer
	printIndentLine(i, &out)
	out.WriteString("FuncDecl ")
	out.WriteString(f.Token.Value)
	printIndentLine(i+1, &out)
	out.WriteString("Params")
	for _, p := range f.Params {
		out.WriteString(p.DebugString(i + 2))
	}
	printIndentLine(i+1, &out)
	out.WriteString("Body")
	for _, s := range f.Body {
		out.WriteString(s.DebugString(i + 2))
	}
	return out.String()
}
func (f *FuncDecl) String() string {
	var out bytes.Buffer
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
