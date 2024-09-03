package llvm

import (
	"bytes"
	"fmt"
	"lang/ast"
	"lang/token"
	"strconv"
)

type Assembler struct {
	ind     int
	reg     int
	out     bytes.Buffer
	program *ast.Program
}

func New(p *ast.Program) *Assembler {
	return &Assembler{
		program: p,
	}
}

func (a *Assembler) Generate() string {
	for i, n := range a.program.Statements {
		a.generateNode(n)
		if i < len(a.program.Statements)-1 {
			a.newLine()
		}
	}

	return a.out.String()
}

func (a *Assembler) generateNode(n ast.Node) int {
	start := a.out.Len()
	switch v := n.(type) {
	case *ast.FuncDecl:
		a.generateFuncDecl(v)
	case *ast.FuncCall:
		a.generateFuncCall(v)
	case *ast.InfixExpression:
		a.generateInfixExpression(v)
	case *ast.Return:
		a.generateReturn(v)
	case *ast.Type:
		a.generateType(v)
	case *ast.Var:
		a.generateVar(v)
	case *ast.VarDecl:
		a.generateVarDecl(v)
	case *ast.EmptyExpression:
	case *ast.IntLiteral:
	default:
		panic(fmt.Sprintf("cannot generate node: %T", v))
	}
	return a.out.Len() - start
}

func (a *Assembler) generateFuncCall(fc *ast.FuncCall) {
	for i := range fc.FuncDecl.Params {
		// val := fc.FuncDecl.Params[i].Value

		fc.FuncDecl.Params[i].Value = fc.Args[i]
		if a.generateNode(fc.FuncDecl.Params[i]) > 0 {
			a.newLine()
		}

		// fc.FuncDecl.Params[i].Value = val
	}

	// for _, param := range fc.FuncDecl.Params {
	// param.Value =
	// }
	//
	// for _, arg := range fc.Args {
	// if a.generateNode(arg) > 0 {
	// a.newLine()
	// }
	// }

	fc.Register = a.getRegister()
	a.writef("%s = call ", fc.Register)
	a.generateType(fc.FuncDecl.ReturnType)
	a.writef(" @%s(", fc.Token.Value)

	for i, p := range fc.FuncDecl.Params {
		a.writef("ptr %%%s", p.Token.Value)
		if i < len(fc.FuncDecl.Params)-1 {
			a.write(", ")
		}
	}
	a.write(")")
}

func (a *Assembler) generateFuncDecl(fd *ast.FuncDecl) {
	if fd.Extern {
		a.write("declare ")
	} else {
		a.write("define ")
	}
	a.generateNode(fd.ReturnType)
	a.write(" @")
	a.write(fd.Token.Value)
	a.write("(")
	for i, vd := range fd.Params {
		a.writef("ptr %%%s", vd.Token.Value)

		if i < len(fd.Params)-1 {
			a.write(", ")
		}
	}
	a.write(")")

	if !fd.Extern {
		a.write(" {")
		a.indent()
		for _, n := range fd.Body {
			a.newLine()
			a.generateNode(n)
		}
		a.unindent()
		a.newLine()
		a.write("}")
	}
}

func (a *Assembler) generateInfixExpression(ie *ast.InfixExpression) {
	if a.generateNode(ie.Left) > 0 {
		a.newLine()
	}
	if a.generateNode(ie.Right) > 0 {
		a.newLine()
	}

	ie.Register = a.getRegister()

	a.write(ie.Register)
	switch ie.Token.Value {
	case token.PLUS:
		a.write(" = add ")
	case token.ASTERISK:
		a.write(" = mul ")
	case token.MINUS:
		a.write(" = sub ")
	default:
		panic(fmt.Sprintf("cannot generate InfixExpression: %s", ie.Token.Value))
	}
	a.generateNode(ie.Type())
	a.writef(" %s, %s", ie.Left.Location(), ie.Right.Location())
}

func (a *Assembler) generateReturn(r *ast.Return) {
	a.generateNode(r.Value)
	a.newLine()
	a.write("ret ")
	a.generateNode(r.Value.Type())
	a.writef(" %s", r.Value.Location())
}

func (a *Assembler) generateType(t *ast.Type) {
	a.write(t.Name)
}

func (a *Assembler) generateVar(v *ast.Var) {
	v.Register = a.getRegister()
	a.writef("%s = load ", v.Register)
	a.generateNode(v.VarDecl.Type)
	a.writef(", ptr %%%s", v.VarDecl.Token.Value)
}

func (a *Assembler) generateVarDecl(vd *ast.VarDecl) {
	a.writef("%%%s = alloca ", vd.Token.Value)
	a.generateNode(vd.Type)
	a.newLine()

	if a.generateNode(vd.Value) > 0 {
		a.newLine()
	}

	a.write("store ")
	a.generateNode(vd.Value.Type())
	a.writef(" %s, ptr %%%s", vd.Value.Location(), vd.Token.Value)
}

func (a *Assembler) indent() {
	a.ind += 1
}

func (a *Assembler) unindent() {
	a.ind -= 1
}

func (a *Assembler) newLine() {
	curr := a.out.String()
	if curr[len(curr)-1] != '\n' {
		a.out.WriteString("\n")
	}
	for i := 0; i < a.ind; i++ {
		a.out.WriteString("\t")
	}
}

func (a *Assembler) getRegister() string {
	a.reg++
	return "%" + strconv.Itoa(a.reg)
}

func (a *Assembler) writef(str string, args ...interface{}) {
	a.write(fmt.Sprintf(str, args...))
}

func (a *Assembler) write(s string) {
	a.out.WriteString(s)
}
