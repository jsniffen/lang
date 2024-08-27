package ast

import (
	"lang/ir"
	"lang/token"
	"lang/types"
	"strconv"
)

type Node interface {
	Codegen(w *ir.Writer)
}

type Statement interface {
	Node
	isStatement()
}

type Expression interface {
	Node
	isExpression()
	GetType() *Type
	GetLocation() *Location
}

type Program struct {
	Statements []Statement
}

func (p *Program) Codegen(w *ir.Writer) {
	for _, stmt := range p.Statements {
		stmt.Codegen(w)
		w.Write("\n")
	}
}

type FuncArg struct {
	Name     string
	Type     *Type
	Location *Location
}

func (fa *FuncArg) isStatement() {}

func (fa *FuncArg) Codegen(w *ir.Writer) {
	fa.Type.Codegen(w)
	if fa.Location != nil {
		w.Write(" ")
		fa.Location.Codegen(w)
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

func (fd *FuncDecl) Codegen(w *ir.Writer) {
	if fd.Extern {
		w.Write("declare ")
	} else {
		w.Write("define ")
	}
	fd.ReturnType.Codegen(w)
	w.Write(" @")
	w.Write(fd.Name)
	w.Write("(")
	for i, fa := range fd.Args {
		fa.Codegen(w)
		if i < len(fd.Args)-1 {
			w.Write(", ")
		}
	}
	w.Write(")")

	if !fd.Extern {
		w.Write(" {")
		for i, stmt := range fd.Body {
			w.Indent()
			w.NewLine()
			stmt.Codegen(w)

			if i == len(fd.Body)-1 {
				w.DeIndent()
				w.NewLine()
			}
		}
		w.Write("}")
	}
	w.Write(" ")
}

type IntLiteral struct {
	Token token.Token
	Value int
}

func (il *IntLiteral) isExpression() {}

func (il *IntLiteral) GetLocation() *Location {
	return &Location{strconv.Itoa(il.Value)}
}

func (il *IntLiteral) GetType() *Type {
	return &Type{
		Name: types.Int32,
	}
}

func (il *IntLiteral) Codegen(w *ir.Writer) {}

type Return struct {
	Token    token.Token
	HasValue bool
	Value    Expression
}

func (r *Return) isStatement() {}

func (r *Return) Codegen(w *ir.Writer) {
	if r.HasValue {
		r.Value.Codegen(w)
	}

	w.Write("ret")
	if r.HasValue {
		w.Write(" ")
		r.Value.GetType().Codegen(w)
		w.Write(" ")
		r.Value.GetLocation().Codegen(w)
	}
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Right    Expression
	Location *Location
	Type     *Type
}

func (ie *InfixExpression) isExpression() {}

func (ie *InfixExpression) GetType() *Type { return ie.Type }

func (ie *InfixExpression) GetLocation() *Location { return ie.Location }

func (ie *InfixExpression) Codegen(w *ir.Writer) {
	ie.Left.Codegen(w)
	ie.Right.Codegen(w)

	ie.Location.Codegen(w)
	w.Write(" = add ")
	ie.GetType().Codegen(w)
	w.Write(" ")
	ie.Left.GetLocation().Codegen(w)
	w.Write(", ")
	ie.Right.GetLocation().Codegen(w)

	w.NewLine()
}

type Location struct {
	Name string
}

func (l *Location) Codegen(w *ir.Writer) {
	w.Write(l.Name)
}

type Type struct {
	Name string
}

func (t *Type) Codegen(w *ir.Writer) {
	w.Write(t.Name)
}
