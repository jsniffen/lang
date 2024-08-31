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

	Filename string
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
			c.checkFuncDecl(v)
			c.context.funcs[v.Token.Value] = v
		case *ast.VarDecl:
			c.checkVarDecl(v)
			c.context.vars[v.Token.Value] = v
		default:
			panic("unsupported type")
		}
	}
}

func (c *Checker) checkFuncDecl(fd *ast.FuncDecl) {
	if dup, ok := c.context.getFuncDecl(fd.Token.Value); ok {
		c.errorDuplicate(fd.Token, dup.Token)
	}
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
}

func (c *Checker) checkExpression(e ast.Expression) {
	switch v := e.(type) {
	case *ast.Var:
		c.checkVar(v)
	case *ast.IntLiteral:
		return
	default:
		panic(v)
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
