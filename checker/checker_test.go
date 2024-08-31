package checker

import (
	"lang/ast"
	"lang/lexer"
	"lang/parser"
	"testing"
)

func TestDupFuncDecl(t *testing.T) {
	input := `
func x() {}
func x() {}
	`
	p := parse(t, input)
	checker := New(p)
	checker.Check()

	if len(checker.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(checker.Errors))
	}
}

func TestDuplicateVarDecl(t *testing.T) {
	input := `
var x i21 = 1
var x i21 = 2
	`
	p := parse(t, input)
	checker := New(p)
	checker.Check()

	if len(checker.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(checker.Errors))
	}
}

func TestVarResolution(t *testing.T) {
	input := `
var x i32 = 1
var y i32 = x
	`
	p := parse(t, input)

	vd := p.Statements[0].(*ast.VarDecl)
	v := p.Statements[1].(*ast.VarDecl).Value.(*ast.Var)

	if v.VarDecl != nil {
		t.Fatalf("expected nil, got %v", v.VarDecl)
	}

	checker := New(p)
	checker.Check()

	if v.VarDecl != vd {
		t.Fatalf("expected %v, got %v", vd, v.VarDecl)
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
