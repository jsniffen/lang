package nasm

import (
	"bytes"
	"fmt"
	"lang/ast"
)

type Assembler struct {
	data    *Section
	global  *Section
	text    *Section
	program *ast.Program
}

func New(p *ast.Program) *Assembler {
	return &Assembler{
		data:    newSection("data"),
		global:  newSection(""),
		text:    newSection("text"),
		program: p,
	}
}

func (a *Assembler) Generate() string {
	for _, n := range a.program.Statements {
		a.generateNode(n)
	}

	var out bytes.Buffer
	out.WriteString(a.global.String())
	out.WriteString(a.data.String())
	out.WriteString(a.text.String())
	return out.String()
}

func (a *Assembler) generateNode(n ast.Node) {
	switch v := n.(type) {
	case *ast.IntLiteral:
		a.generateIntLiteral(v)
	case *ast.FuncDecl:
		a.generateFuncDecl(v)
	case *ast.FuncCall:
		a.generateFuncCall(v)
	default:
		panic(fmt.Sprintf("unsupported generation of node: %v", v))
	}
}

func (a *Assembler) generateFuncCall(fc *ast.FuncCall) {
	for _, arg := range fc.Args {
		a.text.write("mov ecx, ")
		a.generateNode(arg)
		a.text.newLine()
	}
	a.text.write("call ")
	a.text.write(fc.Token.Value)
}

func (a *Assembler) generateFuncDecl(fd *ast.FuncDecl) {
	if fd.Extern {
		a.global.write("extern ")
	} else {
		a.global.write("global ")
	}
	a.global.write(fd.Token.Value)
	a.global.write("\n")

	if !fd.Extern {
		a.text.write(fd.Token.Value)
		a.text.write(":")
		a.text.indent()
		a.text.newLine()

		for _, n := range fd.Body {
			a.generateNode(n)
		}
	}
}

func (a *Assembler) generateIntLiteral(il *ast.IntLiteral) {
	a.text.write(il.Token.Value)
}
