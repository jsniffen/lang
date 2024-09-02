package ast

import (
	"fmt"
)

type Iterator struct {
	stack []Node
}

func NewIterator(n Node) *Iterator {
	return &Iterator{[]Node{n}}
}

func (it *Iterator) push(n Node) {
	it.stack = append(it.stack, n)
}

func (it *Iterator) pop() (Node, bool) {
	if len(it.stack) == 0 {
		return nil, false
	}

	n := it.stack[len(it.stack)-1]
	it.stack = it.stack[:len(it.stack)-1]

	return n, true
}

func (it *Iterator) Next() (Node, bool) {
	n, ok := it.pop()
	if !ok {
		return nil, false
	}

	switch v := n.(type) {
	case *Program:
		for i := len(v.Statements) - 1; i >= 0; i-- {
			it.push(v.Statements[i])
		}
	case *IntLiteral:
	case *InfixExpression:
		it.push(v.Right)
		it.push(v.Left)
	case *FuncArg:
	case *FuncCall:
		for i := len(v.Args) - 1; i >= 0; i-- {
			it.push(v.Args[i])
		}
	case *FuncDecl:
		for i := len(v.Args) - 1; i >= 0; i-- {
			it.push(v.Args[i])
		}
		for i := len(v.Body) - 1; i >= 0; i-- {
			it.push(v.Body[i])
		}
	case *Return:
		if v.HasValue {
			it.push(v.Value)
		}
	case *VarDecl:
		it.push(v.Value)
	case *Var:
	default:
		panic(fmt.Sprintf("Iterator error: unsupported Node: %T", v))
	}

	return n, true
}
