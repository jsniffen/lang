package main

import (
	"fmt"
	"lang/codegen"
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
	inputFile := ""
	if len(os.Args) < 2 {
		fmt.Printf("input file required")
		return
	}
	inputFile = os.Args[1]
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

	g := codegen.New(prog)
	code := g.Generate()
	if outputFile == "" {
		fmt.Println(code)
	} else {
		os.WriteFile(outputFile, []byte(code), 0666)
	}
}
