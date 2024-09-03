package ast

import (
	"lang/token"
)

type Node interface {
	isNode()
}

type Statement interface {
	Node
	isStatement()
}

type Expression interface {
	Node
	isExpression()
	Type() *Type
	Location() string
}

type Program struct {
	Statements []Statement
}

func (p *Program) isNode() {}

type FuncArg struct {
	Name string
	Type *Type
}

func (fa *FuncArg) isNode()      {}
func (fa *FuncArg) isStatement() {}

type FuncCall struct {
	Token    token.Token
	Args     []Expression
	FuncDecl *FuncDecl
	Register string
}

func (fc *FuncCall) isNode()          {}
func (fc *FuncCall) isStatement()     {}
func (fc *FuncCall) isExpression()    {}
func (fc *FuncCall) Type() *Type      { return fc.FuncDecl.ReturnType }
func (fc *FuncCall) Location() string { return fc.Register }

type FuncDecl struct {
	Params     []*VarDecl
	Body       []Statement
	Extern     bool
	Token      token.Token
	ReturnType *Type
}

func (fd *FuncDecl) isNode()      {}
func (fd *FuncDecl) isStatement() {}

type IntLiteral struct {
	Token token.Token
	Value int
}

func (il *IntLiteral) isNode()          {}
func (il *IntLiteral) isExpression()    {}
func (il *IntLiteral) Type() *Type      { return Int32 }
func (il *IntLiteral) Location() string { return il.Token.Value }

type Return struct {
	Token    token.Token
	HasValue bool
	Value    Expression
}

func (r *Return) isNode()      {}
func (r *Return) isStatement() {}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Right    Expression
	Register string
}

func (ie *InfixExpression) isNode()          {}
func (ie *InfixExpression) isExpression()    {}
func (ie *InfixExpression) Type() *Type      { return ie.Left.Type() }
func (ie *InfixExpression) Location() string { return ie.Register }

type Var struct {
	Token token.Token

	// set by checker
	VarDecl *VarDecl

	// set by assembler
	Register string
}

func (v *Var) isNode()          {}
func (v *Var) isExpression()    {}
func (v *Var) Type() *Type      { return v.VarDecl.Type }
func (v *Var) Location() string { return v.Register }

type VarDecl struct {
	Token token.Token
	Type  *Type
	Value Expression
}

func (vd *VarDecl) isNode()      {}
func (vd *VarDecl) isStatement() {}

type Type struct {
	Name string
}

func (t *Type) isNode() {
}

type EmptyExpression struct{}

func (ee *EmptyExpression) isNode()          {}
func (ee *EmptyExpression) isExpression()    {}
func (ee *EmptyExpression) Type() *Type      { return Empty }
func (ee *EmptyExpression) Location() string { return "" }
