package main

import (
	"fmt"
	"lang/asm"
	"lang/checker"
	"lang/lexer"
	"lang/parser"
	"os"
)

func main() {
	inputFile := "examples/testfile"
	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	}
	outputFile := "out.asm"
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

	ch := checker.New(prog)
	ch.Check()
	for _, err := range ch.Errors {
		fmt.Println(err)
	}

	a := asm.New(prog)
	code := a.Generate()
	fmt.Println(code)
	os.WriteFile(outputFile, []byte(code), 0666)
}
