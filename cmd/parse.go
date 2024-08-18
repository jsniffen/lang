package main

import (
	"fmt"
	"lang/codegen"
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
	l, err := lexer.FromFile("examples/testfile")
	if err != nil {
		panic(err)
	}
	p := parser.New(l)

	prog, ok := p.ParseProgram()
	if !ok {
		p.PrintErrors()
	}
	fmt.Println(prog.DebugString(0))
	fmt.Println("---")
	fmt.Println(prog.String())

	g := codegen.New(prog)
	fmt.Println(g.Generate())
}
