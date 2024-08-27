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
	outputFile := ""
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

	var w ir.Writer
	prog.Codegen(&w)
	if outputFile == "" {
		fmt.Println(w.String())
	} else {
		os.WriteFile(outputFile, w.Bytes(), 0666)
	}
}
