package main

import (
	"fmt"
	"lang/lexer"
	"lang/parser"
)

func main() {
	l := lexer.New("x = 1")
	p := parser.New(l)

	prog, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}

	for _, s := range prog.Statements {
		fmt.Println(s)
	}
}
