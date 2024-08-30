package checker

import "lang/ast"

type Checker struct {
	prog *ast.Program
}

func New(p *ast.Program) *Checker {
	return &Checker{p}
}

func (c *Checker) Check() bool {
	for _, stmt := range c.prog.Statements {
		if !c.checkNode(stmt) {
			return false
		}
	}
	return true
}

func (c *Checker) checkNode(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.FuncDecl:
		return true
	default:
		panic(v)
		return false
	}
}
