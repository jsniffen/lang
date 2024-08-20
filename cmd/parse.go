package main

import (
	"fmt"
	"lang/lexer"
	"lang/parser"
	"os"
)

const TestInput = `
	func main() {
		s = "hello world"
		x = 1+2/34
	}
`

func main() {
	inputFile := "examples/testfile"
	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	}
	outputFile := ""
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}
	l, err := lexer.FromFile(inputFile)
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
	fmt.Println("---")

	code := prog.CodeGen()
	if outputFile == "" {
		fmt.Println(code)
	} else {
		os.WriteFile(outputFile, []byte(code), 0666)
	}
}
