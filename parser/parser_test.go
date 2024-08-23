package parser

import (
	"lang/ast"
	"lang/lexer"
	"lang/types"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestParseFuncDecl(t *testing.T) {
	input := `
	func a(a i32, b i32)
	func b() i32 {}
	`
	stmts := []ast.Statement{
		&ast.FuncDecl{
			Args: []*ast.FuncArg{
				&ast.FuncArg{
					Name:     "a",
					Type:     &ast.Type{types.Int32},
					Location: "%a",
				},
				&ast.FuncArg{
					Name:     "b",
					Type:     &ast.Type{types.Int32},
					Location: "%b",
				},
			},
			Name:       "a",
			ReturnType: &ast.Type{types.Void},
			Extern:     true,
		},
		&ast.FuncDecl{
			Args:       []*ast.FuncArg{},
			Name:       "b",
			ReturnType: &ast.Type{types.Int32},
			Extern:     false,
		},
	}

	test(t, input, stmts)
}

func test(t *testing.T, input string, want []ast.Statement) {
	l := lexer.New(input)
	p := New(l)

	prog, _ := p.ParseProgram()

	for _, err := range p.Errors {
		t.Fatalf(err)
	}

	got := prog.Statements
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %s, want %s", spew.Sdump(got), spew.Sdump(want))
	}
}
