package codegen

import (
	"bytes"
	"fmt"
	"lang/ast"
	"lang/token"
	"lang/types"
	"strconv"
)

type Generator struct {
	varMap  map[string]*ast.VarDecl
	funcMap map[string]*ast.FuncDecl
	prog    *ast.Program
	buf     bytes.Buffer
	reg     int
}

func New(prog *ast.Program) *Generator {
	g := &Generator{prog: prog, reg: 1}
	g.buildMaps()
	return g
}

func (g *Generator) Generate() string {
	for _, s := range g.prog.Statements {
		switch v := s.(type) {
		case *ast.FuncDecl:
			g.genFuncDecl(v)
		}
		g.write("\n")
	}
	return g.buf.String()
}

func (g *Generator) genFuncDecl(fd *ast.FuncDecl) {
	extern := len(fd.Body) == 0

	if extern {
		g.write("declare ")
	} else {
		g.write("define ")
	}
	if fd.ReturnType.Value != "" {
		g.write(fd.ReturnType.Value)
		g.write(" ")
	} else {
		g.write("void ")
	}
	g.write("@")
	g.write(fd.Token.Value)
	g.write("(")
	for i, vd := range fd.Params {
		if !extern {
			g.write(vd.Name.Value)
			g.write(" ")
		}
		g.write(vd.Type.Value)
		if vd.Pointer {
			g.write("*")
		}
		if i < len(fd.Params)-1 {
			g.write(", ")
		}
	}
	g.write(")")

	if len(fd.Body) > 0 {
		g.write(" {")
		for _, s := range fd.Body {
			g.write("\n\t")
			g.genStatement(s)
		}
		g.write("\n}")
	}
}

func (g *Generator) genStatement(s ast.Statement) {
	switch v := s.(type) {
	case *ast.Return:
		g.genReturn(v)
	}
}

func (g *Generator) genExpression(exp ast.Expression) (string, types.Type) {
	switch v := exp.(type) {
	case *ast.IntegerLiteral:
		return g.genIntLiteral(v)
	case *ast.InfixExpression:
		return g.genInfixExpression(v)
	}
	return "", ""
}

func (g *Generator) genInfixExpression(exp *ast.InfixExpression) (string, types.Type) {
	lReg, lType := g.genExpression(exp.Left)
	rReg, rType := g.genExpression(exp.Right)

	if lType != rType {
		fmt.Println("TypeError: %s != %s", lType, rType)
		return "", ""
	}

	v := g.getNextRegister()
	g.write(v)
	g.write(" = ")
	switch exp.Token.Type {
	case token.ASTERISK:
		g.write("mul ")

	case token.PLUS:
		g.write("add ")
	}
	g.write(string(lType))
	g.write(" ")
	g.write(lReg)
	g.write(", ")
	g.write(rReg)
	g.newLine()
	return v, lType
}

func (g *Generator) genIntLiteral(il *ast.IntegerLiteral) (string, types.Type) {
	return il.Token.Value, il.Type
}

func (g *Generator) genReturn(r *ast.Return) {
	v, t := g.genExpression(r.Value)
	g.write("ret ")
	g.write(string(t))
	g.write(" ")
	g.write(v)
}

func (g *Generator) write(s string) {
	g.buf.WriteString(s)
}

func (g *Generator) getNextRegister() string {
	reg := strconv.Itoa(g.reg)
	g.reg += 1
	return "%" + reg
}

func (g *Generator) newLine() {
	g.write("\n\t")
}

func (g *Generator) buildMaps() {
	g.funcMap = make(map[string]*ast.FuncDecl, 0)
	g.varMap = make(map[string]*ast.VarDecl, 0)

	for _, s := range g.prog.Statements {
		switch v := s.(type) {
		case *ast.FuncDecl:
			g.funcMap[v.Token.Value] = v
		case *ast.VarDecl:
			g.varMap[v.Name.Value] = v
		}
	}
}
