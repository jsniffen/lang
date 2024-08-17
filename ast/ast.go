package ast

type Node interface {
	String() string
}

type Expression interface {
	Node
	isExpression()
}

type Statement interface {
	Node
	isStatement()
}

type Program struct {
	Statements []Statement
}

type IntegerLiteral struct {
	Value string
}

func (i *IntegerLiteral) String() string {
	return "IntegerLiteral{" + i.Value + "}"
}
func (i *IntegerLiteral) isExpression() {}

type VariableDeclaration struct {
	Identifier string
	Expression Expression
}

func (v *VariableDeclaration) String() string {
	return "VariableDeclaration{" + v.Identifier + " = " + v.Expression.String() + "}"
}
func (v *VariableDeclaration) isStatement() {}
