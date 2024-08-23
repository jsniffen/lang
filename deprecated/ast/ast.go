package ast

import (
	"bytes"
)

type Node interface {
	CodeGen() string
	DebugString(int) string
	String() string
}

type Expression interface {
	Node
	isExpression()
	GetResult() Result
}

type Statement interface {
	Node
	isStatement()
}

type Result struct {
	Type Type
}

func printIndentLine(i int, b *bytes.Buffer) {
	b.WriteString("\n")
	for j := 0; j < i; j += 1 {
		b.WriteString("\t")
	}
}
