package main

import (
	"fmt"
	"lang/lexer"
	"lang/parser"
)

func main() {
	input := `
	func main() {
		s = "hello world"
		x = 1+2/34
	}
	`
	l := lexer.New(input)
	p := parser.New(l)

	prog, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}

	fmt.Println(prog.DebugString(0))
}
