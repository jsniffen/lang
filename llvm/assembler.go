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
	case *ast.Var:
		a.generateVar(v)
	case *ast.VarDecl:
		a.generateVarDecl(v)
	case *ast.EmptyExpression:
	case *ast.IntLiteral:
	case *ast.RegisterExpression:
	default:
		panic(fmt.Sprintf("cannot generate node: %T", v))
	}
	return a.out.Len() - start
}

func (a *Assembler) generateFuncCall(fc *ast.FuncCall) {
	for _, arg := range fc.Args {
		if a.generateNode(arg) > 0 {
			a.newLine()
		}
	}

	fc.Register = a.getRegister()
	a.writef("%s = call %s @%s(", fc.Register, fc.FuncDecl.ReturnType.Type.Name(), fc.Token.Value)
	for i, arg := range fc.Args {
		a.writef("%s %s", arg.Type().Name(), arg.Location())
		if i < len(fc.FuncDecl.Params)-1 {
			a.write(", ")
		}
	}
	a.write(")")
}

func (a *Assembler) generateFuncDecl(fd *ast.FuncDecl) {
	a.resetRegister()

	if fd.Extern {
		a.write("declare ")
	} else {
		a.write("define ")
	}
	a.writef("%s @%s(", fd.ReturnType.Type.Name(), fd.Token.Value)

	for i, vd := range fd.Params {
		reg := a.getRegister()
		a.writef("%s %s", vd.Type.Type.Name(), reg)

		if i < len(fd.Params)-1 {
			a.write(", ")
		}

		vd.Value = &ast.RegisterExpression{
			Register:     reg,
			RegisterType: vd.Type.Type,
		}
	}
	a.write(")")
	a.getRegister()

	if !fd.Extern {
		a.write(" {")
		a.indent()
		a.newLine()
		for _, n := range fd.Params {
			if a.generateNode(n) > 0 {
				a.newLine()
			}
		}
		for i, n := range fd.Body {
			if a.generateNode(n) > 0 && i < len(fd.Body)-1 {
				a.newLine()
			}
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
	a.write(ie.Type().Name())
	a.writef(" %s, %s", ie.Left.Location(), ie.Right.Location())
}

func (a *Assembler) generateReturn(r *ast.Return) {
	if a.generateNode(r.Value) > 0 {
		a.newLine()
	}
	a.writef("ret %s %s", r.Value.Type().Name(), r.Value.Location())
}

func (a *Assembler) generateType(t *ast.Type) {
	a.write(t.Type.Name())
}

func (a *Assembler) generateVar(v *ast.Var) {
	v.Register = a.getRegister()
	a.writef("%s = load %s, ptr %%%s", v.Register, v.VarDecl.Type.Type.Name(), v.VarDecl.Token.Value)
}

func (a *Assembler) generateVarDecl(vd *ast.VarDecl) {
	a.writef("%%%s = alloca %s", vd.Token.Value, vd.Type.Type.Name())
	a.newLine()

	if a.generateNode(vd.Value) > 0 {
		a.newLine()
	}

	_, ok := vd.Value.(*ast.EmptyExpression)
	if !ok {
		a.writef("store %s %s, ptr %%%s", vd.Type.Type.Name(), vd.Value.Location(), vd.Token.Value)
	}
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

func (a *Assembler) resetRegister() {
	a.reg = 0
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
