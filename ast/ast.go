package ast

import "io"

type Node interface {
	Codegen(w io.StringWriter)
}

type Statement interface {
	Node
	isStatement()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Codegen(w io.StringWriter) {
	for _, stmt := range p.Statements {
		stmt.Codegen(w)
		w.WriteString("\n")
	}
}

type FuncArg struct {
	Name     string
	Type     *Type
	Location string
}

func (fa *FuncArg) isStatement() {}

func (fa *FuncArg) Codegen(w io.StringWriter) {
	fa.Type.Codegen(w)
	if fa.Location != "" {
		w.WriteString(" ")
		w.WriteString(fa.Location)
	}
}

type FuncDecl struct {
	Args       []*FuncArg
	Body       []Statement
	Extern     bool
	Name       string
	ReturnType *Type
}

func (fd *FuncDecl) isStatement() {}

func (fd *FuncDecl) Codegen(w io.StringWriter) {
	if fd.Extern {
		w.WriteString("define ")
	} else {
		w.WriteString("declare ")
	}
	fd.ReturnType.Codegen(w)
	w.WriteString(" ")
	w.WriteString(fd.Name)
	w.WriteString("(")
	for i, fa := range fd.Args {
		fa.Codegen(w)
		if i < len(fd.Args)-1 {
			w.WriteString(", ")
		}
	}
	w.WriteString(")")
	w.WriteString(" ")
}

type Type struct {
	Name string
}

func (t *Type) Codegen(w io.StringWriter) {
	w.WriteString(t.Name)
}
