package codegen

import (
	"bytes"
	"lang/ast"
)

type Generator struct {
	funcMap map[string]*ast.FuncDecl
	prog    *ast.Program
	buf     bytes.Buffer
}

func New(prog *ast.Program) *Generator {
	g := &Generator{prog: prog}
	g.buildFuncMap()
	return g
}

func (g *Generator) Generate() string {
	for _, s := range g.prog.Statements {
		switch v := s.(type) {
		case *ast.FuncDecl:
			g.genFuncDecl(v)
		}
		g.buf.WriteString("\n")
	}
	return g.buf.String()
}

func (g *Generator) genFuncDecl(fd *ast.FuncDecl) {
	if len(fd.Body) > 0 {
		g.buf.WriteString("define ")
	} else {
		g.buf.WriteString("declare ")
	}
	if fd.ReturnType.Value != "" {
		g.buf.WriteString(fd.ReturnType.Value)
		g.buf.WriteString(" ")
	}
	g.buf.WriteString("@")
	g.buf.WriteString(fd.Token.Value)
	g.buf.WriteString("(")
	for _, p := range fd.Params {
		g.buf.WriteString(p.CodeGen())
	}
	g.buf.WriteString(")")
}

func (g *Generator) buildFuncMap() {}
