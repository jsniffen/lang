package ir

import (
	"bytes"
	"fmt"
	"lang/ast"
	"lang/token"
	"strconv"
)

type Generator struct {
	w          bytes.Buffer
	currIndent int
	register   int
	program    *ast.Program
	varMap     map[string]*ast.VarDecl
	funcMap    map[string]*ast.FuncDecl
}

func New(p *ast.Program) *Generator {
	g := &Generator{
		program: p,
		varMap:  make(map[string]*ast.VarDecl, 0),
		funcMap: make(map[string]*ast.FuncDecl, 0),
	}
	g.buildFuncMap()
	g.buildVarMap()
	return g
}

func (g *Generator) Generate() string {
	for _, s := range g.program.Statements {
		g.generateStatement(s)
	}
	return g.w.String()
}

func (g *Generator) generateStatement(s ast.Statement) {
	switch v := s.(type) {
	case *ast.FuncDecl:
		g.generateFuncDecl(v)
	case *ast.VarDecl:
		g.generateVarDecl(v)
	case *ast.Return:
		g.generateReturn(v)
	}
}

func (g *Generator) generateFuncDecl(fd *ast.FuncDecl) {
	if fd.Extern {
		g.write("declare ")
	} else {
		g.write("define ")
	}
	g.generateType(fd.ReturnType)
	g.write(" @")
	g.write(fd.Name)
	g.write("(")
	for i, fa := range fd.Args {
		g.generateFuncArg(fa)
		if i < len(fd.Args)-1 {
			g.write(", ")
		}
	}
	g.write(")")

	if !fd.Extern {
		g.write(" {")
		g.indent()

		for i, stmt := range fd.Body {
			g.newLine()
			g.generateStatement(stmt)

			if i == len(fd.Body)-1 {
				g.deIndent()
				g.newLine()
			}
		}
		g.write("}")
	}
	g.write(" ")
}

func (g *Generator) generateExpression(e ast.Expression) {
	switch v := e.(type) {
	case *ast.InfixExpression:
		g.generateInfixExpression(v)
	case *ast.Var:
		g.generateVar(v)
	}
}

func (g *Generator) generateVar(v *ast.Var) {
	vd, ok := g.varMap[v.Name]
	if !ok {
		fmt.Println(v.Name, "not found")
		return
	}

	v.Type = vd.Type
	v.Location = g.getRegisterLocation()

	g.generateLocation(v.Location)
	g.write(" = load ")
	g.generateType(vd.Type)
	g.write(", ptr ")
	g.generateLocation(vd.Location)
}

func (g *Generator) generateInfixExpression(ie *ast.InfixExpression) {
	g.generateExpression(ie.Left)
	g.newLine()
	g.generateExpression(ie.Right)
	g.newLine()

	if !g.assertSameType(ie.Left.GetType(), ie.Right.GetType()) {
		return
	}
	ie.Type = ie.Left.GetType()

	ie.Location = g.getRegisterLocation()
	g.generateLocation(ie.Location)
	g.write(" = ")

	switch ie.Token.Type {
	case token.SLASH:
		g.write("sdiv")
	case token.ASTERISK:
		g.write("mul")
	case token.MINUS:
		g.write("sub")
	case token.PLUS:
		g.write("add")
	}

	g.write(" ")
	g.generateType(ie.Type)
	g.write(" ")
	g.generateLocation(ie.Left.GetLocation())
	g.write(", ")
	g.generateLocation(ie.Right.GetLocation())
}

func (g *Generator) generateFuncArg(fa *ast.FuncArg) {
	fa.Location = g.getRegisterLocation()

	g.generateType(fa.Type)
	g.write(" ")
	g.generateLocation(fa.Location)
}

func (g *Generator) generateVarDecl(vd *ast.VarDecl) {
	g.generateLocation(vd.Location)
	g.write(" = alloca ")
	g.generateType(vd.Type)
	g.generateExpression(vd.Value)
	g.newLine()
	g.write("store ")
	g.generateType(vd.Value.GetType())
	g.write(" ")
	g.generateLocation(vd.Value.GetLocation())
	g.write(", ptr ")
	g.generateLocation(vd.Location)

	g.varMap[vd.Name] = vd
}

func (g *Generator) generateReturn(r *ast.Return) {
	if r.HasValue {
		g.generateExpression(r.Value)
		g.newLine()
	}

	g.write("ret")
	if r.HasValue {
		g.write(" ")
		g.generateType(r.Value.GetType())
		g.write(" ")
		g.generateLocation(r.Value.GetLocation())
	}
}

func (g *Generator) generateLocation(l *ast.Location) {
	g.write(l.Name)
}

func (g *Generator) generateType(t *ast.Type) {
	g.write(t.Name)
}

func (g *Generator) write(s string) {
	g.w.WriteString(s)
}

func (g *Generator) newLine() {
	g.write("\n")
	for i := 0; i < g.currIndent; i += 1 {
		g.write("\t")
	}
}

func (g *Generator) indent() {
	g.currIndent += 1
}

func (g *Generator) deIndent() {
	g.currIndent -= 1
}

func (g *Generator) getRegisterLocation() *ast.Location {
	g.register += 1
	name := "%" + strconv.Itoa(g.register)
	return &ast.Location{
		Name: name,
	}
}

func (g *Generator) getStackLocation(name string) *ast.Location {
	return &ast.Location{
		Name: "%" + name,
	}
}

func (g *Generator) buildFuncMap() {}
func (g *Generator) buildVarMap()  {}

func (g *Generator) assertSameType(a, b *ast.Type) bool {
	if a.Name != b.Name {
		fmt.Println(fmt.Sprintf("TypeError: %v != %v", a, b))
		return false
	}
	return true
}
