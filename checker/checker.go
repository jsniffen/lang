package checker

import (
	"fmt"
	"lang/ast"
	"lang/token"
)

type Checker struct {
	program *ast.Program
	context *Context
	Errors  []string
}

func New(p *ast.Program) *Checker {
	c := &Checker{
		program: p,
		context: newContext(nil),
		Errors:  make([]string, 0),
	}
	return c
}

func (c *Checker) Check() {
	for _, stmt := range c.program.Statements {
		switch v := stmt.(type) {
		case *ast.FuncDecl:
			c.checkFuncDeclDup(v)
			c.context.funcs[v.Token.Value] = v
		case *ast.VarDecl:
			c.checkVarDecl(v)
		default:
			panic("unsupported type")
		}
	}

	c.checkFuncDecls()
}

func (c *Checker) checkReturn(r *ast.Return) {
	c.checkExpression(r.Value)
}

func (c *Checker) checkVar(v *ast.Var) {
	vd, ok := c.context.getVar(v.Token.Value)
	if ok {
		v.VarDecl = vd
	} else {
		c.errorNotFound(v.Token, v.Token.Value)
	}
}

func (c *Checker) checkVarDecl(vd *ast.VarDecl) {
	if dup, ok := c.context.getVar(vd.Token.Value); ok {
		c.errorDuplicate(vd.Token, dup.Token)
	}

	c.checkExpression(vd.Value)

	c.context.vars[vd.Token.Value] = vd
}

func (c *Checker) checkExpression(e ast.Expression) {
	switch v := e.(type) {
	case *ast.Var:
		c.checkVar(v)
	case *ast.InfixExpression:
		c.checkInfixExpression(v)
	case *ast.FuncCall:
		c.checkFuncCall(v)
	case *ast.IntLiteral:
	default:
		panic(fmt.Sprintf("checking unsupported expression: %T", v))
	}
}

func (c *Checker) checkInfixExpression(ie *ast.InfixExpression) {
	c.checkExpression(ie.Left)
	c.checkExpression(ie.Right)
}

func (c *Checker) checkFuncArg(fc *ast.FuncArg) {

}

func (c *Checker) checkFuncCall(fc *ast.FuncCall) {
	fd, ok := c.context.getFuncDecl(fc.Token.Value)
	if !ok {
		c.errorNotFound(fc.Token, fc.Token.Value)
	}

	fc.FuncDecl = fd
}

func (c *Checker) checkFuncDecl(fd *ast.FuncDecl) {
	c.pushContext()
	defer c.popContext()

	// for _, a := range fd.Args {
	// c.checkVarDecl(a)
	// }

	for _, s := range fd.Body {
		switch v := s.(type) {
		case *ast.VarDecl:
			c.checkVarDecl(v)
		case *ast.Return:
			c.checkReturn(v)
		default:
			panic(fmt.Sprintf("cannot check body %T", v))
		}
	}
}

func (c *Checker) checkFuncDecls() {
	it := ast.NewIterator(c.program)
	for n, ok := it.Next(); ok; n, ok = it.Next() {
		switch v := n.(type) {
		case *ast.FuncDecl:
			c.checkFuncDecl(v)
		}
	}
}

func (c *Checker) checkFuncDeclDup(fd *ast.FuncDecl) {
	if dup, ok := c.context.getFuncDecl(fd.Token.Value); ok {
		c.errorDuplicate(fd.Token, dup.Token)
	}
}

func (c *Checker) error(t token.Token, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	err := fmt.Sprintf("%s %s", t.Path(), msg)
	c.Errors = append(c.Errors, err)
}

func (c *Checker) errorDuplicate(t, dup token.Token) {
	c.error(t, "duplicate declaration of '%s', previous declaration at %s", t.Value, dup.Path())
}

func (c *Checker) errorNotFound(t token.Token, name string) {
	c.error(t, "%s not declared", name)
}

func (c *Checker) pushContext() {
	c.context = newContext(c.context)
}

func (c *Checker) popContext() {
	if c.context.outer != nil {
		c.context = c.context.outer
	}
}
