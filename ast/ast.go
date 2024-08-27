package ast

import (
	"lang/token"
	"lang/types"
	"strconv"
)

type Statement interface {
	isStatement()
}

type Expression interface {
	isExpression()
	GetType() *Type
	GetLocation() *Location
}

type Program struct {
	Statements []Statement
}

type FuncArg struct {
	Name     string
	Type     *Type
	Location *Location
}

func (fa *FuncArg) isStatement() {}

type FuncDecl struct {
	Args       []*FuncArg
	Body       []Statement
	Extern     bool
	Name       string
	ReturnType *Type
}

func (fd *FuncDecl) isStatement() {}

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

type Return struct {
	Token    token.Token
	HasValue bool
	Value    Expression
}

func (r *Return) isStatement() {}

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

type Var struct {
	Token    token.Token
	Name     string
	Type     *Type
	Location *Location
}

func (v *Var) isExpression() {}

func (v *Var) GetLocation() *Location { return v.Location }
func (v *Var) GetType() *Type         { return v.Type }

type VarDecl struct {
	Location *Location
	Name     string
	Token    token.Token
	Type     *Type
	Value    Expression
}

func (vd *VarDecl) isStatement() {}

type Location struct {
	Name string
}

type Type struct {
	Name string
}
