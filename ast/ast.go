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
	ReturnType() Type
}

type Statement interface {
	Node
	isStatement()
}

func printIndentLine(i int, b *bytes.Buffer) {
	b.WriteString("\n")
	for j := 0; j < i; j += 1 {
		b.WriteString("\t")
	}
}
