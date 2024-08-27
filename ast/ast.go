package ast

import (
	"lang/old_ir"
	"lang/token"
	"lang/types"
	"strconv"
)

type Node interface {
	Codegen(w *old_ir.Writer)
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

func (p *Program) Codegen(w *old_ir.Writer) {
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

func (fa *FuncArg) Codegen(w *old_ir.Writer) {
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

func (fd *FuncDecl) Codegen(w *old_ir.Writer) {
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
		w.Indent()
		for i, stmt := range fd.Body {
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

func (il *IntLiteral) Codegen(w *old_ir.Writer) {}

type Return struct {
	Token    token.Token
	HasValue bool
	Value    Expression
}

func (r *Return) isStatement() {}

func (r *Return) Codegen(w *old_ir.Writer) {
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

func (ie *InfixExpression) Codegen(w *old_ir.Writer) {
	ie.Left.Codegen(w)
	ie.Right.Codegen(w)

	ie.Location.Codegen(w)
	w.Write(" = ")

	switch ie.Token.Type {
	case token.SLASH:
		w.Write("sdiv")
	case token.ASTERISK:
		w.Write("mul")
	case token.MINUS:
		w.Write("sub")
	case token.PLUS:
		w.Write("add")
	}

	w.Write(" ")
	ie.GetType().Codegen(w)
	w.Write(" ")
	ie.Left.GetLocation().Codegen(w)
	w.Write(", ")
	ie.Right.GetLocation().Codegen(w)

	w.NewLine()
}

type Var struct {
	Token    token.Token
	Name     string
	Type     *Type
	Location *Location
}

func (v *Var) isExpression() {}

func (v *Var) Codegen(w *old_ir.Writer) {}
func (v *Var) GetLocation() *Location   { return v.Location }
func (v *Var) GetType() *Type           { return v.Type }

type VarDecl struct {
	Location *Location
	Name     string
	Token    token.Token
	Type     *Type
	Value    Expression
}

func (vd *VarDecl) isStatement() {}

func (vd *VarDecl) Codegen(w *old_ir.Writer) {
	vd.Location.Codegen(w)
	w.Write(" = alloca ")
	vd.Type.Codegen(w)
	w.NewLine()
	vd.Value.Codegen(w)
	w.NewLine()
	w.Write("store ")
	vd.Value.GetType().Codegen(w)
	w.Write(" ")
	vd.Value.GetLocation().Codegen(w)
	w.Write(", ptr ")
	vd.Location.Codegen(w)
}

type Location struct {
	Name string
}

func (l *Location) Codegen(w *old_ir.Writer) {
	w.Write(l.Name)
}

type Type struct {
	Name string
}

func (t *Type) Codegen(w *old_ir.Writer) {
	w.Write(t.Name)
}
