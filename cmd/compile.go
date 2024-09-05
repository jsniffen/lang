package main

import (
	"fmt"
	"lang/checker"
	"lang/lexer"
	"lang/llvm"
	"lang/parser"
	"os"
)

func main() {
	inputFile := "examples/testfile"
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
	if len(p.Errors) > 0 {
		for _, err := range p.Errors {
			fmt.Println(err)
		}
		return
	}

	ch := checker.New(prog)
	ch.Check()

	if len(ch.Errors) > 0 {
		for _, err := range ch.Errors {
			fmt.Println(err)
		}
		return
	}

	gen := llvm.NewGenerator()
	code := gen.Generate(prog)
	// code := a.Generate()
	fmt.Println(code)
	os.WriteFile(outputFile, []byte(code), 0666)
}
