package main

import (
	"fmt"
	"lang/lexer"
	"lang/parser"
)

func main() {
	input := `
	x = -1*2+3-1
	`
	l := lexer.New(input)
	p := parser.New(l)

	prog, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}

	fmt.Println(prog)
}
