package opt

import (
	"fmt"
	"lang/ast"
	"lang/lexer"
	"lang/parser"
	"testing"
)

func TestOptGetExpressions(t *testing.T) {
	input := `
	var x i32 = 1
	var y i32 = 2 + x
	`
	prog := parse(t, input)
	o := New(prog)
	exps := o.getExpressions()

	if len(exps) != 4 {
		s := ""
		for _, e := range exps {
			s += fmt.Sprintf("%T, ", e)
		}
		t.Fatalf("want 4, got %d:\n%v", len(exps), s)
	}
}

func parse(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	prog, ok := p.ParseProgram()
	if !ok {
		for _, err := range p.Errors {
			t.Fatalf(err)
		}
		return nil
	}
	return prog
}
