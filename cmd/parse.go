package main

import (
	"bytes"
	"fmt"
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

	var b bytes.Buffer
	prog.Codegen(&b)
	if outputFile == "" {
		fmt.Println(b.String())
	} else {
		os.WriteFile(outputFile, b.Bytes(), 0666)
	}
}
