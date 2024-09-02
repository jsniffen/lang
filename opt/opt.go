package opt

import "lang/ast"

type Optimizer struct {
	program *ast.Program
}

func New(p *ast.Program) *Optimizer {
	o := &Optimizer{
		program: p,
	}
	return o
}

func (o *Optimizer) Optimizer() {
	exps := o.getExpressions()
	o.allocateRegisters(exps)
}

func (o *Optimizer) allocateRegisters(exps []ast.Expression) {
}

func (o *Optimizer) getExpressions() []ast.Expression {
	it := ast.NewIterator(o.program)

	exps := make([]ast.Expression, 0)
	for n, ok := it.Next(); ok; n, ok = it.Next() {
		switch v := n.(type) {
		case ast.Expression:
			exps = append(exps, v)
		}
	}

	return exps
}
