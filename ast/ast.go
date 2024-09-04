package ast

import (
	"lang/token"
	"lang/types"
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
	Type() types.Type
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
func (fc *FuncCall) Type() types.Type { return fc.FuncDecl.ReturnType.Type }
func (fc *FuncCall) Location() string { return fc.Register }

type FuncDecl struct {
	Params     []*VarDecl
	Body       []Statement
	Extern     bool
	Token      token.Token
	HasReturn  bool
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
func (il *IntLiteral) Type() types.Type { return types.Int32 }
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
func (ie *InfixExpression) Type() types.Type { return ie.Left.Type() }
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
func (v *Var) Type() types.Type { return v.VarDecl.Type.Type }
func (v *Var) Location() string { return v.Register }

type VarDecl struct {
	Token    token.Token
	Type     *Type
	Value    Expression
	Register string
}

func (vd *VarDecl) isNode()      {}
func (vd *VarDecl) isStatement() {}

type Type struct {
	Token token.Token
	Type  types.Type
}

func (t *Type) isNode() {}

type EmptyExpression struct{}

func (ee *EmptyExpression) isNode()          {}
func (ee *EmptyExpression) isExpression()    {}
func (ee *EmptyExpression) Type() types.Type { return types.Nil }
func (ee *EmptyExpression) Location() string { return "" }

type RegisterExpression struct {
	RegisterType types.Type
	Register     string
}

func (re *RegisterExpression) isNode()          {}
func (re *RegisterExpression) isExpression()    {}
func (re *RegisterExpression) Type() types.Type { return re.RegisterType }
func (re *RegisterExpression) Location() string { return re.Register }
