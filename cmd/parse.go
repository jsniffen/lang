package main

import (
	"fmt"
	"lang/ir"
	"lang/lexer"
	"lang/parser"
	"os"
)

func main() {
	inputFile := "deprecated/examples/testfile"
	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	}
	outputFile := "out.ll"
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}
	l, err := lexer.FromFile(inputFile)
	if err != nil {
		panic(err)
	}
	p := parser.New(l)

	prog, _ := p.ParseProgram()
	for _, err := range p.Errors {
		fmt.Println(err)

	}

	g := ir.New(prog)
	code := g.Generate()
	fmt.Println(code)
	os.WriteFile(outputFile, []byte(code), 0666)
}
