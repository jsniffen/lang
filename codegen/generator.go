package codegen

import (
	"bytes"
	"lang/ast"
	"lang/token"
	"strconv"
)

type Generator struct {
	varMap  map[string]*ast.VarDecl
	funcMap map[string]*ast.FuncDecl
	prog    *ast.Program
	buf     bytes.Buffer
}

func New(prog *ast.Program) *Generator {
	g := &Generator{prog: prog}
	g.buildMaps()
	return g
}

func (g *Generator) Generate() string {
	for _, s := range g.prog.Statements {
		switch v := s.(type) {
		case *ast.FuncDecl:
			g.genFuncDecl(v)
		case *ast.VarDecl:
			g.genVarDecl(v)
		}
		g.write("\n")
	}
	return g.buf.String()
}

func (g *Generator) genFuncDecl(fd *ast.FuncDecl) {
	if len(fd.Body) > 0 {
		g.write("define ")
	} else {
		g.write("declare ")
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

func (g *Generator) genFuncCall(fc *ast.FuncCall) {
	g.write("call ")

	name := fc.Token.Value
	fd, ok := g.funcMap[name]
	if !ok {
		return
	}

	if fd.ReturnType.Value == "" {
		g.write("void")
	} else {
		g.write(fd.ReturnType.Value)
	}
	g.write(" @")
	g.write(fc.Token.Value)
	g.write("(")

	if len(fc.Args) > 0 {
		for i, a := range fc.Args {
			g.genExpression(a)

			if i < len(fc.Args)-1 {
				g.write(", ")
			}
		}
	}

	g.write(")")
}

func (g *Generator) genVarDecl(vd *ast.VarDecl) {
	if vd.Global {
		g.write("@")
		g.write(vd.Name.Value)
		g.write(" = global ")
		if vd.Type.Value != "" {
			g.write(vd.Type.Value)
			if vd.Pointer {
				g.write("*")
			}
			g.write(" ")
		}
		g.genExpression(vd.Value)
	}
}

func (g *Generator) genStatement(s ast.Statement) {
	switch v := s.(type) {
	case *ast.FuncCall:
		g.genFuncCall(v)
	case *ast.Return:
		g.genReturn(v)
	}
}

func (g *Generator) genExpression(exp ast.Expression) {
	switch v := exp.(type) {
	case *ast.StringLiteral:
		g.genStringLiteral(v)
	case *ast.Var:
		g.genVar(v)
	case *ast.IntegerLiteral:
		g.genIntLiteral(v)
	case *ast.InfixExpression:
		g.genInfixExpression(v)
	}
}

func (g *Generator) genInfixExpression(exp *ast.InfixExpression) {
	g.genExpression(exp.Left)
	switch exp.Token.Type {
	case token.PLUS:
		g.write(" + ")
	}
	g.genExpression(exp.Right)
}

func (g *Generator) genStringLiteral(sl *ast.StringLiteral) {
	g.write("[")
	g.write(strconv.Itoa(len(sl.Token.Value) - 2))
	g.write(" x ")
	g.write("i8")
	g.write("] ")
	g.write("c")
	g.write(sl.Token.Value)
}

func (g *Generator) genIntLiteral(il *ast.IntegerLiteral) {
	g.write(il.Token.Value)
}

func (g *Generator) genVar(v *ast.Var) {
	name := v.Name.Value

	vd, ok := g.varMap[name]
	if !ok {
		return
	}

	g.write(vd.Type.Value)
	if vd.Pointer {
		g.write("*")
	}
	g.write(" ")
	if vd.Global {
		g.write("@")
	}
	g.write(name)
}

func (g *Generator) genReturn(r *ast.Return) {
	g.write("ret ")
	g.genExpression(r.Value)
}

func (g *Generator) write(s string) {
	g.buf.WriteString(s)
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
