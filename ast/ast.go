package ast

import (
	"io"
	"lang/token"
	"lang/types"
	"strconv"
)

type Node interface {
	Codegen(w io.StringWriter)
}

type Statement interface {
	Node
	isStatement()
}

type Expression interface {
	Node
	isExpression()
	GetType() *Type
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
		w.WriteString("declare ")
	} else {
		w.WriteString("define ")
	}
	fd.ReturnType.Codegen(w)
	w.WriteString(" @")
	w.WriteString(fd.Name)
	w.WriteString("(")
	for i, fa := range fd.Args {
		fa.Codegen(w)
		if i < len(fd.Args)-1 {
			w.WriteString(", ")
		}
	}
	w.WriteString(")")

	if !fd.Extern {
		w.WriteString(" {")
		for i, stmt := range fd.Body {
			w.WriteString("\n\t")
			stmt.Codegen(w)

			if i == len(fd.Body)-1 {
				w.WriteString("\n")
			}
		}
		w.WriteString("}")
	}
	w.WriteString(" ")
}

type IntLiteral struct {
	Token token.Token
	Value int
}

func (il *IntLiteral) isExpression() {}

func (il *IntLiteral) GetType() *Type {
	return &Type{
		Name: types.Int32,
	}
}

func (il *IntLiteral) Codegen(w io.StringWriter) {
	w.WriteString(strconv.Itoa(il.Value))
}

type Return struct {
	Token    token.Token
	HasValue bool
	Value    Expression
}

func (r *Return) isStatement() {}

func (r *Return) Codegen(w io.StringWriter) {
	w.WriteString("ret")
	if r.HasValue {
		w.WriteString(" ")
		r.Value.GetType().Codegen(w)
		w.WriteString(" ")
		r.Value.Codegen(w)
	}
}

type Type struct {
	Name string
}

func (t *Type) Codegen(w io.StringWriter) {
	w.WriteString(t.Name)
}
