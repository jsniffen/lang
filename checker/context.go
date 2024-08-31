package checker

import "lang/ast"

type Context struct {
	outer *Context
	vars  map[string]*ast.VarDecl
	funcs map[string]*ast.FuncDecl
}

func newContext(outer *Context) *Context {
	return &Context{
		outer: outer,
		vars:  make(map[string]*ast.VarDecl),
		funcs: make(map[string]*ast.FuncDecl),
	}
}

func (c *Context) getVar(name string) (*ast.VarDecl, bool) {
	for ctx := c; ctx != nil; ctx = ctx.outer {
		if vd, ok := ctx.vars[name]; ok {
			return vd, true
		}
	}
	return nil, false
}

func (c *Context) getFuncDecl(name string) (*ast.FuncDecl, bool) {
	for ctx := c; ctx != nil; ctx = ctx.outer {
		if fd, ok := ctx.funcs[name]; ok {
			return fd, true
		}
	}
	return nil, false
}
