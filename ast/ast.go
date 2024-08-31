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
	GetType() *Type
}

type Program struct {
	Statements []Statement
}

type FuncArg struct {
	Name     string
	Type     *Type
	Location *Location
}

func (fa *FuncArg) isNode()      {}
func (fa *FuncArg) isStatement() {}

type FuncDecl struct {
	Args       []*FuncArg
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

func (il *IntLiteral) isNode()       {}
func (il *IntLiteral) isExpression() {}

func (il *IntLiteral) GetType() *Type {
	return &Type{
		Name: types.Int32,
	}
}

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
	Location *Location
	Type     *Type
}

func (ie *InfixExpression) isNode()       {}
func (ie *InfixExpression) isExpression() {}

func (ie *InfixExpression) GetType() *Type { return ie.Type }

type Var struct {
	Token token.Token

	// lazily evaluated
	VarDecl *VarDecl
}

func (v *Var) isNode()       {}
func (v *Var) isExpression() {}

func (v *Var) GetType() *Type { return v.VarDecl.Type }

type VarDecl struct {
	Location *Location
	Token    token.Token
	Type     *Type
	Value    Expression
}

func (vd *VarDecl) isNode()      {}
func (vd *VarDecl) isStatement() {}

type Location struct {
	Name string
}

func (l *Location) isNode() {}

type Type struct {
	Name string
}

func (t *Type) isNode() {}
