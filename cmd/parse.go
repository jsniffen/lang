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
	l, err := lexer.FromFile("helloworld")
	if err != nil {
		panic(err)
	}
	p := parser.New(l)

	prog, ok := p.ParseProgram()
	if !ok {
		p.PrintErrors()
	}
	fmt.Println(prog.DebugString(0))
}
