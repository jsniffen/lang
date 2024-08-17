package main

import (
	"fmt"
	"lang/lexer"
	"lang/parser"
)

const TestInput = `
	func main() {
		s = "hello world"
		x = 1+2/34
	}
`

func main() {
	l, err := lexer.FromFile("testfile")
	if err != nil {
		panic(err)
	}
	p := parser.New(l)

	prog := p.ParseProgram()
	if p.HasErrors() {
		p.PrintErrors()
	} else {
		fmt.Println(prog.DebugString(0))
	}
}
